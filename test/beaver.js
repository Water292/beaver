const Beaver = artifacts.require("Beaver");

contract("Beaver test", async (accounts) => {
  it("should be possible for many users to buy and review a product", async () => {
    var total = 0;
    function gas(comment, tx) {
      total += tx.receipt.gasUsed;
      console.log(comment + ": " + tx.receipt.gasUsed);
    }
    let inst = await Beaver.deployed();

    // Add the product
    let productId = await inst.add.call("Test Product", 2000000000000000);
    assert.equal(productId, 0);

    gas("Add", await inst.add("Test Product", 2000000000000000));
    
    // Customers buy
    gas("Buy 1", await inst.buy(0, {from: accounts[1], value: 2000000000000000}));
    gas("Buy 2", await inst.buy(0, {from: accounts[2], value: 2000000000000000}));
    gas("Buy 3", await inst.buy(0, {from: accounts[3], value: 2000000000000000}));

    // Customers express interest in reviewing
    gas("Review Interest 1", await inst.review(0, {from: accounts[1]}));
    gas("Review Interest 2", await inst.review(0, {from: accounts[2]}));
    await inst.review(0, {from: accounts[3]}).then(async(result) => {
      gas("Review Interest 3", result);
      let e = result.logs[0];
      assert.equal(e.event, "GroupReady");
      assert.equal(e.args.productId, 0);
      assert.equal(e.args.groupId, 0);

      gas("Post Score", await inst.score(0, 0, 0, 3, {from: accounts[1]}));
      gas("Endorse Score 1", await inst.endorse(0, 0, 1, {from: accounts[2]}));
      gas("Endorse Score 2", await inst.endorse(0, 0, 2, {from: accounts[3]}));

      console.log("Total: " + total);
    });
  })
})
