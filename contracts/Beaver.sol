pragma solidity ^0.4.17;

contract Beaver {
  struct Product {
    address seller;
    bytes32 name;
    uint price;
    bool deleted;
  }

  mapping (address => uint) pendingWithdrawals;

  Product[] products;

  event Purchased(uint productId, address buyer);

  function add(bytes32 name, uint price) public returns (uint) {
    return products.push(Product(msg.sender, name, price, false)) - 1;
  }

  function setPrice(uint productId, uint price) public {
    require(msg.sender == products[productId].seller);
    products[productId].price = price;
  }

  function remove(uint productId) public {
    require(msg.sender == products[productId].seller);
    products[productId].deleted = true;
  }

  function query(uint productId) public view returns (address, bytes32, uint, bool) {
    Product storage p = products[productId];

    return (p.seller, p.name, p.price, p.deleted);
  }

  function buy(uint productId) public payable returns (bool) {
    Product storage p = products[productId];

    if (msg.value >= p.price) {
      pendingWithdrawals[p.seller] = p.price;
      pendingWithdrawals[msg.sender] = msg.value - p.price;
      emit Purchased(productId, msg.sender);
      return true;
    } else {
      return false;
    }
  }

  function withdraw() public {
    uint amount = pendingWithdrawals[msg.sender];
    pendingWithdrawals[msg.sender] = 0;
    msg.sender.transfer(amount);
  }
}
