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
      var AdoptionArtifact = data;
      App.contracts.Adoption = TruffleContract(AdoptionArtifact);

      // Set the provider for our contract
      App.contracts.Adoption.setProvider(App.web3Provider);

      return App.loadProducts();
    });

    return App.bindEvents();
  },

  bindEvents: function() {
    $(document).on('click', '.btn-primary', App.handleAdd);
  },

  loadProducts: function(adopters, account) {
    var beaverInstance;

    App.contracts.Adoption.deployed().then(function(instance) {
      beaverInstance = instance;

      return beaverInstance.products.call();
    }).then(function(products) {
      for (i = 0; i < products.length; i++) {
        console.log(products[i]);
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

    var petId = parseInt($(event.target).data('id'));

    var adoptionInstance;

    web3.eth.getAccounts(function(error, accounts) {
      if (error) {
        console.log(error);
      }

      var account = accounts[0];

      App.contracts.Adoption.deployed().then(function(instance) {
        adoptionInstance = instance;

        return adoptionInstance.adopt(petId, {from: account});
      }).then(function(result) {
        return App.markAdopted();
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
