pragma solidity ^0.4.17;

contract Beaver {
  struct Product {
    address seller;
    string name;
    uint price;
    bool deleted;
  }

  Product[] products;

  function add(string name, uint price) public returns (uint) {
    return products.push(Product(msg.sender, name, price, false)) - 1;
  }

  function remove(uint productId) public {
    if (products[productId].seller == msg.sender) {
      products[productId].deleted = true;
    }
  }

  function query(uint productId) public view returns (address, string, uint, bool) {
    Product storage p = products[productId];

    return (p.seller, p.name, p.price, p.deleted);
  }
}
