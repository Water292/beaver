App = {
  web3Provider: null,
  contracts: {},

  init: function() {
    return App.initWeb3();
  },

  initWeb3: function() {
    if (typeof web3 !== 'undefined') {
      App.web3Provider = web3.currentProvider;
    } else {
      App.web3Provider = new Web3.providers.HttpProvider('http://localhost:7545');
    }
    web3 = new Web3(App.web3Provider);

    return App.initContract();
  },

  initContract: function() {
    $.getJSON('Beaver.json', function(data) {
      var BeaverArtifact = data;
      App.contracts.Beaver = TruffleContract(BeaverArtifact);

      // Set the provider for our contract
      App.contracts.Beaver.setProvider(App.web3Provider);

      return App.loadProducts();
    });

    return App.bindEvents();
  },

  bindEvents: function() {
    $(document).on('click', '.btn-primary', App.handleAdd);
    $(document).on('click', '.btn-success', App.handleBuy);
  },

  loadProducts: function(adopters, account) {
    var beaverInstance;

    App.contracts.Beaver.deployed().then(function(instance) {
      beaverInstance = instance;

      return beaverInstance.getProductsCount.call();
    }).then(async(count) => {
      $('#listing').html('');
      for (i = 0; i < count; i++) {
        var q = await beaverInstance.query.call(i);
        let owner = q[0].substr(0,10);
        let name = web3.toAscii(q[1]);
        let price = web3.fromWei(q[2], "ether");
        let deleted = q[3];
        if (!deleted) {
          let u = await beaverInstance.getUnreviewedCount(i);
          let newRow = $('<tr><td>' + (i+1) + '</td><td>' + owner + '...</td><td>' + name + '</td><td>' +
            price + '</td><td>0</td><td><button data-id="' + i + '" data-price="' + q[2] + '" class="btn btn-success">Buy</button>&nbsp;<button class="btn" id="review' + i + '">Review</button></td>'
          )
          $('#listing').append(newRow);
          if (u == 0) {
            $('#review' + i).prop('disabled', true);
          }
        }
      }
    });
  },

  handleBuy: function(event) {
    event.preventDefault();

    var productId = parseInt($(event.target).data('id'));
    var price = parseInt($(event.target).data('price'));

    web3.eth.getAccounts(function(error, accounts) {
      if (error) {
        console.log(error);
      }

      var account = accounts[0];

      App.contracts.Beaver.deployed().then(function(instance) {
        return instance.buy(productId, {from: account, value: price});
      }).then(function(result) {
        return App.loadProducts();
      }).catch(function(err) {
        console.log(err.message);
      });
    });
  },

  handleAdd: function(event) {
    event.preventDefault();

    var name = $('#name').val();
    var price = web3.toWei($('#price').val(), "ether");

    var beaverInstance;

    web3.eth.getAccounts(function(error, accounts) {
      if (error) {
        console.log(error);
      }

      var account = accounts[0];

      App.contracts.Beaver.deployed().then(function(instance) {
        beaverInstance = instance;

        return beaverInstance.add(name, price, {from: account});
      }).then(function(result) {
        return App.loadProducts();
      }).catch(function(err) {
        console.log(err.message);
      });
    });
  }

};

$(function() {
  $(window).on('load', (function() {
    App.init();
  }));
});
