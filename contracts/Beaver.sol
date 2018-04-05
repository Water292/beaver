pragma solidity ^0.4.17;

contract Beaver {
  struct Product {
    address seller;
    bytes32 name;
    uint price;
    bool deleted;

    // Buyers - the number of reviews an address can leave
    mapping (address => uint) buyers;

    uint[] reviews;
  }

  mapping (address => uint) pendingWithdrawals;

  Product[] products;

  event Purchased(uint productId, address buyer);

  // Product Management
  function add(bytes32 name, uint price) public returns (uint) {
    uint productId = products.length;
    products.length++;
    Product storage p = products[productId];
    p.seller = msg.sender;
    p.name = name;
    p.price = price;
    return productId;
  }

  function setPrice(uint productId, uint price) public {
    require(msg.sender == products[productId].seller);
    products[productId].price = price;
  }

  function remove(uint productId) public {
    require(msg.sender == products[productId].seller);
    products[productId].deleted = true;
  }

  // Query
  function query(uint productId) public view returns (address, bytes32, uint, bool) {
    Product storage p = products[productId];

    return (p.seller, p.name, p.price, p.deleted);
  }

  function reviewCount(uint productId) public view returns (uint) {
    return products[productId].reviews.length;
  }

  function getReview(uint productId, uint i) public view returns (uint) {
    return products[productId].reviews[i];
  }

  // Buyer Actions
  function buy(uint productId) public payable {
    Product storage p = products[productId];
    require(msg.value >= p.price);
    require(!p.deleted);

    pendingWithdrawals[p.seller] += p.price;
    pendingWithdrawals[msg.sender] += msg.value - p.price;
    p.buyers[msg.sender]++;
    emit Purchased(productId, msg.sender);
  }

  function review(uint productId, uint score) public returns (uint) {
    Product storage p = products[productId];
    require(p.buyers[msg.sender] > 0);
    require(score <= 5); // score unsigned, so cannot be less than 0
    p.buyers[msg.sender]--;
    return p.reviews.push(score) - 1;
  }

  // Allows withdrawal of ether held by the contract
  // on behalf of users
  function withdraw() public {
    uint amount = pendingWithdrawals[msg.sender];
    pendingWithdrawals[msg.sender] = 0;
    msg.sender.transfer(amount);
  }
}
