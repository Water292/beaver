pragma solidity ^0.4.17;

contract Beaver {
  struct ReviewGroup {
    address[] reviewers;
    int score;
    mapping (address => bool) endorsements;
  }

  struct Product {
    address seller;
    bytes32 name;
    uint price;
    bool deleted;

    // Buyers - the number of reviews an address can leave
    mapping (address => uint) buyers;

    address[] pendingReviewers;
    ReviewGroup[] reviewGroups;
  }

  mapping (address => uint) pendingWithdrawals;

  Product[] products;

  event Purchased(uint productId, address buyer);
  event GroupReady(uint productId, uint groupId);

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

  function review(uint productId) public returns (uint) {
    Product storage p = products[productId];
    require(p.buyers[msg.sender] > 0);
    p.buyers[msg.sender]--;
    uint reviewerId = p.pendingReviewers.push(msg.sender) - 1;
    if (reviewerId >= 2) { // We have three reviewers
      uint groupId = p.reviewGroups.length;
      p.reviewGroups.length++;
      ReviewGroup storage g = p.reviewGroups[groupId];
      g.reviewers = p.pendingReviewers;
      p.pendingReviewers.length = 0;
      emit GroupReady(productId, groupId);
    }
    return reviewerId;
  }

  function score(uint productId, uint groupId, uint reviewerId, int _score) public {
    Product storage p = products[productId];
    ReviewGroup storage g = p.reviewGroups[groupId];
    require(g.reviewers[reviewerId] == msg.sender);
    // need to deal with malicious changes of score
    g.score = _score;
    g.endorsements[msg.sender] = true;
  }

  function endorse(uint productId, uint groupId, uint reviewerId) public {
    Product storage p = products[productId];
    ReviewGroup storage g = p.reviewGroups[groupId];

    // Make sure the caller is part of the MPC
    require(g.reviewers[reviewerId] == msg.sender);

    // Endorse
    g.endorsements[msg.sender] = true;
  }

  // Allows withdrawal of ether held by the contract
  // on behalf of users
  function withdraw() public {
    uint amount = pendingWithdrawals[msg.sender];
    pendingWithdrawals[msg.sender] = 0;
    msg.sender.transfer(amount);
  }
}
