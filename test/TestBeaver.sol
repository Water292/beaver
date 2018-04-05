pragma solidity ^0.4.17;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/Beaver.sol";

contract TestBeaver {
  Beaver beaver = Beaver(DeployedAddresses.Beaver());

  function testCanAddProduct() public {
    uint productId = beaver.add("Test Product", 2 finney);
    Assert.equal(productId, 0, "The first product should be at index 0.");
  }

  function testProductIsCorrect() public {
    address seller;
    bytes32 name;
    uint price;
    bool deleted;
    (seller, name, price,) = beaver.query(0);
    Assert.equal(seller, this, "Seller should be this contract.");
    Assert.equal(name, "Test Product", "Product name should be as expected.");
    Assert.equal(price, 2 finney, "Product price should be as expected.");
    Assert.isFalse(deleted, "Product should not be delted.");
  }

  function testSetPrice() public {
    beaver.setPrice(0, 100);
    uint price;
    (,,price,) = beaver.query(0);
    Assert.equal(price, 100, "Product price should have changed.");
  }

  function testRemoveProduct() public {
    beaver.remove(0);
    bool deleted;
    (,,,deleted) = beaver.query(0);
    Assert.isTrue(deleted, "Product should be deleted.");
  }
}
