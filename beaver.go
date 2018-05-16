package main

import (
	"fmt"
	"math/big"
	"github.com/dedis/kyber"
	"github.com/dedis/kyber/group/edwards25519"
	"github.com/dedis/kyber/xof/blake2xb"
	"github.com/dedis/kyber/share/vss/pedersen"
	"github.com/dedis/kyber/share"
)

type Reviewer struct {
	dealer 		*vss.Dealer
	verifiers 	[]*vss.Verifier
}

func genPair() (kyber.Scalar, kyber.Point) {
	secret := suite.Scalar().Pick(suite.RandomStream())
	public := suite.Point().Mul(secret, nil)
	return secret, public
}

func genCommits(n int) ([]kyber.Scalar, []kyber.Point, []kyber.Scalar, [][]kyber.Scalar) {
	var secrets = make([]kyber.Scalar, n)
	var publics = make([]kyber.Point, n)
	var vals = make([]kyber.Scalar, n)
	var randVals = make([][]kyber.Scalar, 500)
	for k := 0; k < 500; k++ {
		randVals[k] = make([]kyber.Scalar, n)
	}
	for i := 0; i < n; i++ {
		secrets[i], publics[i] = genPair()
		vals[i] = suite.Scalar().SetInt64(1) // dummy review scores
		for j := 0; j < 500; j++ {
			randVals[j][i], _ = genPair()
		}
	}
	return secrets, publics, vals, randVals
}

func genDealer(i int, secretVals []kyber.Scalar) *vss.Dealer {
	d, _ := vss.NewDealer(suite, verifiersSec[i], secretVals[i], verifiersPub, vssThreshold)
	return d
}

func genVerifier(i, j int) *vss.Verifier {
	v, _ := vss.NewVerifier(suite, verifiersSec[i], verifiersPub[j], verifiersPub)
	return v
}

func genAll(secretVals []kyber.Scalar) ([]*Reviewer) {
	var reviewers = make([]*Reviewer, nbVerifiers)
	for i := 0; i < nbVerifiers; i++ {
		d := genDealer(i, secretVals)
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

var nbVerifiers = 4

var vssThreshold int

var verifiersPub []kyber.Point
var verifiersSec []kyber.Scalar

var secretVals []kyber.Scalar
var randSecretVals [][]kyber.Scalar

func main() {
	verifiersSec, verifiersPub, secretVals, randSecretVals = genCommits(nbVerifiers)
	vssThreshold = vss.MinimumT(nbVerifiers)
	// sum of reviews
   	reviewers := genAll(secretVals)
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

	total := suite.Scalar().SetInt64(0)
	for i := 0; i < 500; i++ {
		v := randomVal(i)
		total.Add(total, v)
		fmt.Println(total)
	}
	if count >= vssThreshold {
		sec, _ := vss.RecoverSecret(suite, deals, nbVerifiers, vssThreshold*2)
		val := suite.Scalar().Add(sec, total)
		fmt.Printf("%v\n%v\n", sec, val)
	}
		
}

func randomVal(n int) kyber.Scalar {

	reviewers := genAll(randSecretVals[n])
   	for j := 0; j < nbVerifiers; j++ {
   		resps := make([]*vss.Response, nbVerifiers)
		encDeals, _ := reviewers[j].dealer.EncryptedDeals()
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
	squareShares := make([]*share.PriShare, nbVerifiers)
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
			squareShares[j] = &share.PriShare{deals[j].SecShare.I, deals[j].SecShare.V.Clone()}
			squareShares[j].V.Mul(squareShares[j].V, squareShares[j].V)
		}
	}


	if count >= vssThreshold {
		randSq, _ := share.RecoverSecret(suite, squareShares, vssThreshold*2, nbVerifiers)
		randSq.Mul(randSq, randSq) // to avoid dealing with quartic residues
		// e = (p + 3)/8 - hardcoded value
		e, _ := new(big.Int).SetString("904625697166532776746648320380374280107139544922488450750243867285681781374", 10)
		sqrt1 := fastModExp(randSq, e) //assuming quartic residue, formula for modular square root when p % 8 = 5
		inv1 := suite.Scalar().Inv(sqrt1)
		for i := 0; i < nbVerifiers; i++ {
			squareShares[i].V.Mul(squareShares[i].V, inv1)
		}
	}
	// val = +- 1
	val, _ := share.RecoverSecret(suite, squareShares, vssThreshold*2, nbVerifiers)
	return val
}

func fastModExp(x kyber.Scalar, e *big.Int) kyber.Scalar {
	y := suite.Scalar().SetInt64(1)
	z := big.NewInt(0)
	for e.Cmp(big.NewInt(0)) == 1 {
        z.Mod(e, big.NewInt(2))
        if z.Cmp(big.NewInt(0)) == 1 {
            y.Mul(x, y)
            e.Sub(e, big.NewInt(1))
        }
        e.Div(e, big.NewInt(2))
        x.Mul(x, x)
    }
    return y
}
