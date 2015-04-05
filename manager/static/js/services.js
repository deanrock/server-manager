var srv = angular.module('managerServices', []);
srv.factory('managerServices', function($http) {
    var managerServices = {
        getShells: function(id) {
            var resp = $http.get('/api/v1/shells').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getContainers: function(id) {
            var resp = $http.get('/api/v1.0/containers/').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getAccounts: function(id) {
            var resp = $http.get('/api/v1/accounts').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getAccountByName: function(name) {
            var resp = $http.get('/api/v1/accounts/'+name).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getProfile: function(name) {
            var resp = $http.get('/api/v1/profile').
            then(function(response) {
                return response.data;
            });
            return resp;
        }
    }
    return managerServices;
});
