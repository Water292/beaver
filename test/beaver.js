const Beaver = artifacts.require("Beaver");

contract("Beaver test", async (accounts) => {
  it("should be possible to buy and review a product", async () => {
    let inst = await Beaver.deployed();
    let productId = await inst.add("Test Product", 2000000000000000);
    
    await inst.buy(0, {from: accounts[1], value: 2000000000000000});
    await inst.review(0, 3, {from: accounts[1]});
    let review = await inst.getReview.call(0, 0);
    assert.equal(review, 3);
  })
})
