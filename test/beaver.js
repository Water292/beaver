const Beaver = artifacts.require("Beaver");

contract("Beaver test", async (accounts) => {
  it("should be possible for many users to buy and review a product", async () => {
    let inst = await Beaver.deployed();

    // Add the product
    let productId = await inst.add.call("Test Product", 2000000000000000);
    assert.equal(productId, 0);

    await inst.add("Test Product", 2000000000000000);
    
    // Customers buy
    await inst.buy(0, {from: accounts[1], value: 2000000000000000});
    await inst.buy(0, {from: accounts[2], value: 2000000000000000});
    await inst.buy(0, {from: accounts[3], value: 2000000000000000});

    // Customers express interest in reviewing
    let cust1 = await inst.review(0, {from: accounts[1]});
    let cust2 = await inst.review(0, {from: accounts[2]});
    inst.review(0, {from: accounts[3]}).then(function(result) {
      let e = result.logs[0];
      assert.equal(e.event, "GroupReady");
      assert.equal(e.args.productId, 0);
      assert.equal(e.args.groupId, 0);
    });
  })
})
