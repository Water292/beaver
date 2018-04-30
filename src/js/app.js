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
  },

  loadProducts: function(adopters, account) {
    var beaverInstance;

    App.contracts.Beaver.deployed().then(function(instance) {
      beaverInstance = instance;

      return beaverInstance.getProductsCount.call();
    }).then(function(count) {
      console.log(count.toString());
      for (i = 0; i < count; i++) {
        console.log("hey");
        if (!products[i].deleted) {
          console.log("add this");
          //$('.panel-pet').eq(i).find('button').text('Success').attr('disabled', true);
        }
      }
    })/*.catch(function(err) {
      console.log(err.message);
    })*/;
  },

  handleAdd: function(event) {
    event.preventDefault();

    var name = $('#name').val();
    var price = $('#price').val();

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
