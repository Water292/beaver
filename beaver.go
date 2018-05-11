package main

import (
  "fmt"
  "encoding/binary"

  "github.com/dedis/kyber"
  "github.com/dedis/kyber/group/edwards25519"
  "github.com/dedis/kyber/share/vss/pedersen"
  "github.com/dedis/kyber/xof/blake2xb"
)

type Reviewer struct {
  dealer    *vss.Dealer
  verifiers []*vss.Verifier
}

func genPair() (kyber.Scalar, kyber.Point) {
  secret := suite.Scalar().Pick(suite.RandomStream())
  public := suite.Point().Mul(secret, nil)
  return secret, public
}

func genCommits(n int) ([]kyber.Scalar, []kyber.Point, []kyber.Scalar) {
  var secrets = make([]kyber.Scalar, n)
  var publics = make([]kyber.Point, n)
  var vals = make([]kyber.Scalar, n)
  for i := 0; i < n; i++ {
    secrets[i], publics[i] = genPair()
    vals[i] = suite.Scalar().SetInt64(int64(i % 5)) // dummy review scores
  }
  return secrets, publics, vals
}

func genDealer(i int) *vss.Dealer {
  d, _ := vss.NewDealer(suite, verifiersSec[i], secretVals[i], verifiersPub, vssThreshold)
  return d
}

func genVerifier(i, j int) *vss.Verifier {
  v, _ := vss.NewVerifier(suite, verifiersSec[i], verifiersPub[j], verifiersPub)
  return v
}

func genAll() []*Reviewer {
  var reviewers = make([]*Reviewer, nbVerifiers)
  for i := 0; i < nbVerifiers; i++ {
    d := genDealer(i)
    var verifiers = make([]*vss.Verifier, nbVerifiers)
    for j := 0; j < nbVerifiers; j++ {
      verifiers[j] = genVerifier(i, j)
    }
    reviewers[i] = &Reviewer{d, verifiers}
  }
  return reviewers
}

var rng = blake2xb.New(nil)

var suite = edwards25519.NewBlakeSHA256Ed25519WithRand(rng)

var nbVerifiers = 6

var vssThreshold int

var verifiersPub []kyber.Point
var verifiersSec []kyber.Scalar

var secretVals []kyber.Scalar

func printArray(vals []kyber.Scalar) {
  for i := 0; i < nbVerifiers; i++ {
    b, _ := vals[i].MarshalBinary()
    fmt.Printf("  %v\n", binary.LittleEndian.Uint32(b))
  }
}

func main() {
  verifiersSec, verifiersPub, secretVals = genCommits(nbVerifiers)
  vssThreshold = vss.MinimumT(nbVerifiers)
  reviewers := genAll()
  for j := 0; j < nbVerifiers; j++ {
    resps := make([]*vss.Response, nbVerifiers)
    encDeals, _ := reviewers[j].dealer.EncryptedDeals() //encDeals on Server
    for i, d := range encDeals {
      resp, err := reviewers[i].verifiers[j].ProcessEncryptedDeal(d)
      if err == nil {
        resps[i] = resp
      }
    }

    for _, resp := range resps {
      for _, r := range reviewers {
        r.verifiers[j].ProcessResponse(resp)
      }
      reviewers[j].dealer.ProcessResponse(resp)
    }
  }

  count := 0
  deals := make([]*vss.Deal, nbVerifiers)
  for j, r := range reviewers {
    certified := true
    for _, v := range r.verifiers {
      if !v.DealCertified() {
        certified = false
      }
    }
    if certified {
      count += 1
      deals[j] = r.verifiers[0].Deal()
      for i := 1; i < nbVerifiers; i++ {
        deals[j].SecShare.V.Add(deals[j].SecShare.V,
          r.verifiers[i].Deal().SecShare.V)
      }
    }
  }

  if count >= vssThreshold {
    sec, _ := vss.RecoverSecret(suite, deals, nbVerifiers, vssThreshold)
    b, _ := sec.MarshalBinary()
    printArray(secretVals)
    fmt.Printf("Total: %v\n", binary.LittleEndian.Uint32(b))
  }
}
