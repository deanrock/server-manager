var srv = angular.module('managerServices', []);
srv.factory('managerServices', function($http) {
    var managerServices = {
        getTasks: function() {
            var resp = $http.get('/api/v1/tasks').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getTask: function(id) {
            var resp = $http.get('/api/v1/tasks/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getTaskLog: function(id) {
            var resp = $http.get('/api/v1/tasks/'+id+'/log').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getShells: function(id) {
            var resp = $http.get('/api/v1/shells').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getImages: function(id) {
            var resp = $http.get('/api/v1/images').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getContainers: function(id) {
            var resp = $http.get('/api/v1/containers/').
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
        addAccount: function(params) {
            var resp = $http.post('/api/v1/accounts', params).
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
            var r = null;
            var resp = $http.get('/api/v1/accounts/'+name).
            then(function(response) {
                r = response.data;
                return $http.get('/api/v1/profile/access/'+r.id);
            }).then(function(response) {
                r.access = response.data;
                return r;
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
        getSSHKeys: function() {
            var resp = $http.get('/api/v1/profile/ssh-keys').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getSSHKey: function(id) {
            var resp = $http.get('/api/v1/profile/ssh-keys/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        deleteSSHKey: function(id) {
            var resp = $http.delete('/api/v1/profile/ssh-keys/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        editSSHKey: function(id, params) {
            var resp = $http.put('/api/v1/profile/ssh-keys/'+id, params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        addSSHKey: function(params) {
            var resp = $http.post('/api/v1/profile/ssh-keys', params).
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
            var resp = $http.get('/api/v1/accounts/'+name+'/apps').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getApp: function(name, id) {
            var resp = $http.get('/api/v1/accounts/'+name+'/apps/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        addApp: function(name, params) {
            var resp = $http.post('/api/v1/accounts/'+name+'/apps', params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        editApp: function(name, id, params) {
            var resp = $http.put('/api/v1/accounts/'+name+'/apps/'+id, params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        deleteApp: function(name, id) {
            var resp = $http.delete('/api/v1/accounts/'+name+'/apps/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        startApp: function(name, id) {
            var resp = $http.post('/api/v1/accounts/'+name+'/apps/'+id+'/start', {}).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        stopApp: function(name, id) {
            var resp = $http.post('/api/v1/accounts/'+name+'/apps/'+id+'/stop', {}).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        redeployApp: function(name, id) {
            var resp = $http.post('/api/v1/accounts/'+name+'/apps/'+id+'/redeploy', {}).
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
        getCronjobLog: function(name, id) {
            var resp = $http.get('/api/v1/accounts/'+name+'/cronjobs/'+id+'/log').
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
        getDomains: function(name) {
            var resp = $http.get('/api/v1/accounts/'+name+'/domains').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getDomain: function(name, id) {
            var resp = $http.get('/api/v1/accounts/'+name+'/domains/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        addDomain: function(name, params) {
            var resp = $http.post('/api/v1/accounts/'+name+'/domains', params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        editDomain: function(name, id, params) {
            var resp = $http.put('/api/v1/accounts/'+name+'/domains/'+id, params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        deleteDomain: function(name, id) {
            var resp = $http.delete('/api/v1/accounts/'+name+'/domains/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        syncDomains: function(name, id) {
            var resp = $http.post('/api/v1/accounts/'+name+'/domains/sync').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getDatabases: function(name) {
            var resp = $http.get('/api/v1/accounts/'+name+'/databases').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getDatabase: function(name, id) {
            var resp = $http.get('/api/v1/accounts/'+name+'/databases/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        addDatabase: function(name, params) {
            var resp = $http.post('/api/v1/accounts/'+name+'/databases', params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        editDatabase: function(name, id, params) {
            var resp = $http.put('/api/v1/accounts/'+name+'/databases/'+id, params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getSSHPasswords: function(name) {
            var resp = $http.get('/api/v1/accounts/'+name+'/ssh-passwords').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        addSSHPassword: function(name, params) {
            var resp = $http.post('/api/v1/accounts/'+name+'/ssh-passwords', params).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        deleteSSHPassword: function(name, id) {
            var resp = $http.delete('/api/v1/accounts/'+name+'/ssh-passwords/'+id).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        getSyncImages: function() {
            var resp = $http.get('/api/v1/sync/images').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        syncImage: function(name) {
            var resp = $http.post('/api/v1/sync/images/'+name).
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        forceSyncImage: function(name) {
            var resp = $http.post('/api/v1/sync/images/'+name+'?no-cache=true').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
        syncWebServers: function(name) {
            var resp = $http.post('/api/v1/sync/web-servers').
            then(function(response) {
                return response.data;
            });
            return resp;
        },
    }
    return managerServices;
});
