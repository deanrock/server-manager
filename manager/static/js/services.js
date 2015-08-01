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
        getAllAccounts: function(id) {
            var resp = $http.get('/api/v1/all-accounts').
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
        },
        getCronjobs: function(name) {
            var resp = $http.get('/api/v1/accounts/'+name+'/cronjobs').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getApps: function(name) {
            var resp = $http.get('/djapi/v1/accounts/'+name+'/apps').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getApp: function(name, id) {
            var resp = $http.get('/djapi/v1/accounts/'+name+'/apps/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        executeAppAction: function(name, id, action) {
            var resp = $http.post('/djapi/v1/accounts/'+name+'/apps/'+id+'/'+action+'/ajax').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getCronjob: function(name, id) {
            var resp = $http.get('/api/v1/accounts/'+name+'/cronjobs/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        addCronjob: function(name, params) {
            var resp = $http.post('/api/v1/accounts/'+name+'/cronjobs', params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        editCronjob: function(name, id, params) {
            var resp = $http.put('/api/v1/accounts/'+name+'/cronjobs/'+id, params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getUsers: function() {
            var resp = $http.get('/api/v1/users').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getUser: function(id) {
            var resp = $http.get('/api/v1/users/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getUserAccess: function(id) {
            var resp = $http.get('/api/v1/users/'+id+'/access').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        setUserAccess: function(id, account, params) {
            var resp = $http.post('/api/v1/users/'+id+'/access/'+account, params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        removeUserAccess: function(id, account) {
            var resp = $http.delete('/api/v1/users/'+id+'/access/'+account).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
    }
    return managerServices;
});
