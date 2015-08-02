var ctrls = angular.module('managerControllers', []);


ctrls.controller('mainCtrl', ['$scope', '$rootScope', 'managerServices', '$window', function($scope, $rootScope, managerServices, $window) {
    $scope.djangoAdmin = function() {
        $window.location = '/admin'
    }
}]).
controller('accounts', ['$scope', 'managerServices', '$location', function($scope, managerServices) {
    $scope.accounts = [];

    managerServices.getAccounts().then(function(data){
        console.log(data)
        $scope.accounts = data;
    })
}]).
controller('containers', ['$scope', 'managerServices', '$location', function($scope, managerServices) {
    $scope.apps = [];

    managerServices.getContainers().then(function(data){
        console.log(data)
        $scope.apps = data;
    })
}]).
controller('accountOverview', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'overview';
    $scope.started = false;

    managerServices.getShells().then(function(data){
		$scope.shells = data;

        $scope.shell = data[0];
	})

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    })

    var shell = null;

    $scope.submit = function() {
        if (shell == null) {
            var selected = $scope.shell;
            shell = new AccountShell($scope.account.name, selected, document.getElementById('shell'));
            $scope.started = true;
        }
    }

    $scope.stop = function() {
        if (shell != null) {
            shell.stop();
            
            $scope.started = false;
            shell = null;
        }
    }
}]).
controller('accountApps', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.apps = [];
    $scope.action = 'apps';
    console.log($routeParams.account)

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getApps($routeParams.account).then(function(data){
        console.log(data)
        $scope.apps = data;
    })
}]).
controller('accountAppLogs', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.apps = [];
    $scope.action = 'apps';
    $scope.logs = [];

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;

        managerServices.getApp($routeParams.account, $routeParams.id).then(function(data) {
            console.log(data)
            $scope.app = data;

            managerServices.getApps($routeParams.account).then(function(data){
                $scope.apps = data;

                var logs = new WebSocket(getWebsocketHost()+'/api/v1/accounts/'+$scope.account.name+'/apps/'+$scope.app.id+'/logs');

                logs.onopen = function() {
                    console.log('open container logs ws')
                }

                logs.onmessage = function(msg) {
                    $scope.$apply(function() {
                        $scope.logs.push(msg.data);
                    })
                    
                    console.log(msg)
                }
            });
        });
    });
}]).
controller('accountAppAction', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.apps = [];
    $scope.action = 'apps';

    $scope.appAction = $routeParams.action;
    $scope.appActionText = null;

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;

        managerServices.getApp($routeParams.account, $routeParams.id).then(function(data) {
            console.log(data)
            $scope.app = data;

            switch($scope.appAction) {
                case "redeploy":
                    $scope.appActionText = 'Redeploying '+$scope.app.name+' for '+$scope.account.name+' ...';
                    break;
                case "stop":
                    $scope.appActionText = 'Stopping '+$scope.app.name+' for '+$scope.account.name+' ...';
                    break;
                case "start":
                    $scope.appActionText = 'Starting '+$scope.app.name+' for '+$scope.account.name+' ...';
                    break;
                default:
                    $scope.appActionText = 'WRONG ACTION!';
                    $scope.appAction = null;
                    break;
            }

            if ($scope.appAction != null) {
                $scope.loader = true;

                $scope.fails = [];
                $scope.appActionResponse = '';
                managerServices.executeAppAction($scope.account.name, $scope.app.id,
                    $scope.appAction).then(function(data) {
                        console.log(data);

                        $scope.loader = false;
                        $scope.appActionResponse += '<b>Response:</b><br />';

                        for(var x in data) {
                            var line = data[x];

                            if (typeof line === 'string') {
                                line = line.replace(new RegExp('\n', 'g'), '<br />');

                                $scope.appActionResponse +='<br /><span style=\'color:red\'>-&gt; </span> '+ line;
                            }else{
                                var background = line.flag == 'fail' ? 'red;color:white;padding:10px 0 10px 0' : 'white';
                                $scope.appActionResponse +='<br /><div style="background:'+background+'"><span style=\'color:red\'>-&gt; </span>'+ line.message+'</div>';

                                if (line.flag=='fail') {
                                    $scope.fails.push(line);
                                }
                            }
                        }


                        if ($scope.fails.length > 0) {
                            $scope.error = '<div class="alert alert-danger" role="alert"><span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span><span class="sr-only">Error:</span>  '+fails.length+' failure(s)!</div>';
                        }
                    });
            }
        });
    });
}]).
controller('accountCronjobs', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'cronjobs';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getCronjobs($routeParams.account).then(function(data) {
        $scope.cronjobs = data;
    });
}]).
controller('accountCronjobAdd', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'cronjobs';
    $scope.form = {
        'timeout': 300
    };

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
        $scope.form.directory = '/home/'+data.name+'/';
    });

    managerServices.getShells().then(function(data){
        $scope.shells = data;
        $scope.form.image = data[0];
    });

    $scope.submit = function() {
        managerServices.addCronjob($routeParams.account, $scope.form).then(function(data) {
            $location.path('/a/'+$scope.account.name+'/cronjobs');
        },
        function(err) {
            $scope.errors = err.data.errors;
            console.log(err)
        });
    }
}]).
controller('accountCronjobEdit', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'cronjobs';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getShells().then(function(data){
        $scope.shells = data;
    });

    managerServices.getCronjob($routeParams.account, $routeParams.id).then(function(data){
        $scope.form = data;

        $scope.submit = function() {
            managerServices.editCronjob($routeParams.account, $scope.form.id, $scope.form).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/cronjobs');
            },
            function(err) {
                $scope.errors = err.data.errors;
                console.log(err)
            });
        }
    });
}]).
controller('sync', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    if ($routeParams.action == "images") {
        $scope.action_url = "/action/rebuild-base-image";
        $scope.action = "Rebuilding base image";
    }else if ($routeParams.action == "users") {
        $scope.action_url = "/action/sync-users";
        $scope.action = "Syncing users";
    }else if ($routeParams.action == "databases") {
        $scope.action_url = "/action/sync-databases";
        $scope.action = "Syncing databases";
    }else if ($routeParams.action == "nginx-apache") {
        $scope.action_url = "/action/update-nginx-config";
        $scope.action = "Updating nginx and apache config ...";
    }

    var fails = [];
    $(document).ready(function() {
        $.get($scope.action_url, function( data ) {
            $('#action-loader').hide();
            console.log(data);
            $('#action-response').html('<b>Response:</b><br />');

            for(var x in data) {
                var line = data[x];

                if (typeof line === 'string') {
                    line = line.replace(new RegExp('\n', 'g'), '<br />');

                    $('#action-response').append('<br /><span style=\'color:red\'>-&gt; </span> '+ line);
                }else{
                    var background = line.flag == 'fail' ? 'red;color:white;padding:10px 0 10px 0' : 'white';
                    $('#action-response').append('<br /><div style="background:'+background+'"><span style=\'color:red\'>-&gt; </span>'+ line.message+'</div>');

                    if (line.flag=='fail') {
                        fails.push(line);
                    }
                }
            }


            if (fails.length > 0) {
                $('#error').html('<div class="alert alert-danger" role="alert"><span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span><span class="sr-only">Error:</span>  '+fails.length+' failure(s)!</div>');
            }
        });
    });
}]).
controller('account', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.suburl = '/frame' + $location.path();
    
    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    })

    $scope.action = $routeParams.action;
}]).
controller('users', ['$scope', 'managerServices', '$location', function($scope, managerServices) {
    $scope.users = [];

    managerServices.getUsers().then(function(data){
        console.log(data)
        $scope.users = data;
    })
}]).
controller('userOverview', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.tab = 'overview';
    
    managerServices.getUser($routeParams.id).then(function(data){
        $scope.user = data;
    })

    $scope.action = $routeParams.action;
}]).
controller('userAccess', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.tab = 'access';

    managerServices.getUser($routeParams.id).then(function(data){
        $scope.user = data;
    })

    loadData();

    $scope.action = $routeParams.action;

    $scope.addAccount = function() {
        managerServices.setUserAccess($routeParams.id, $scope.addAccountName, {}).then(function(data) {
            loadData();
        })
    }

    $scope.delete = function(a) {
        managerServices.removeUserAccess($routeParams.id, a.id).then(function(data) {
            loadData();
        })
    }

    $scope.checkAll = function(a) {
        var access = a.access;

        if (access.ssh_access && access.shell_access && access.app_access &&  access.database_access && access.cronjob_access && access.domain_access) {
            access.ssh_access = false;
            access.shell_access = false;
            access.app_access = false;
            access.database_access = false;
            access.cronjob_access = false;
            access.domain_access = false;
        }else{
            access.ssh_access = true;
            access.shell_access = true;
            access.app_access = true;
            access.database_access = true;
            access.cronjob_access = true;
            access.domain_access = true;
        }

        $scope.update(a);
    }

    $scope.update = function(a) {
        managerServices.setUserAccess($routeParams.id, a.id, a.access).then(function(data) {
            loadData();
        })
    }

    function loadData() {
        managerServices.getUserAccess($routeParams.id).then(function(data) {
            $scope.access = data;

            managerServices.getAllAccounts().then(function(data){
                $scope.accounts = data;

                $scope.availableAccounts = [];
                $scope.accountAccess = [];

                angular.forEach($scope.accounts, function(account, key) {
                    var found = false;
                    angular.forEach($scope.access, function (access, k) {
                        if (access.account_id == account.id) {
                            found = true;
                            account.access = access;

                            if (access.ssh_access && access.shell_access && access.app_access &&  access.database_access && access.cronjob_access && access.domain_access) {
                                access.all = true;
                            }

                            $scope.accountAccess.push(account);
                        }
                    });

                    if (!found) {
                        $scope.availableAccounts.push(account);
                    }
                });

                if ($scope.availableAccounts.length > 0) {
                    $scope.addAccountName = $scope.availableAccounts[0].id;
                }
            })
        })
    }
}]).
controller('userSshKeys', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.suburl = '/frame/profile/ssh-keys';
}]).
controller('accountDomains', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'domains';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getDomains($routeParams.account).then(function(data) {
        $scope.domains = data;
    });
}]).
controller('accountDomainEdit', ['$scope', 'managerServices', '$location', '$routeParams', '$modal', function($scope, managerServices, $location, $routeParams, $modal) {
    $scope.action = 'domains';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    $scope.example = function(id, value) {
       var source   = $("#" + value).html();
       var template = Handlebars.compile(source);

       if (id == 'nginx') {
            $scope.form.nginx_config = template({});
       }else{
            $scope.form.apache_config = template({});
       }
    }

    $scope.form = {};
    if ($routeParams.id !== undefined) {
        managerServices.getDomain($routeParams.account, $routeParams.id).then(function(data) {
            $scope.form = data;
        });
    }

    managerServices.getApps($routeParams.account).then(function(data) {
        $scope.apps = data;

        managerServices.getImages().then(function(data) {
            $scope.images = data;
            $scope.variables = [];

            angular.forEach($scope.apps, function(app, key) {
                angular.forEach($scope.images, function(image, key2) {
                    if (image.name == app.image) {
                        angular.forEach(image.ports, function(port, key3) {
                            $scope.variables.push('#'+app.name+'_'+port.port+'_ip#')
                        });
                    }
                });
            });

            var s = '';
            for (v in $scope.variables) {
                s +=', ' + $scope.variables[v];
            }
            $scope.variablesString = s;
        });
    });

    $scope.nginxExamples = [
        {'name': 'PHP (wordpress, codeigniter... rewrite)', 'value': 'nginx1'},
        {'name': 'PHP (Drupal specific)', 'value': 'nginx2'},
        {'name': 'ProxyPass', 'value': 'nginx3'},
        {'name': 'ProxyPass + WebSockets', 'value': 'nginx5'},
        {'name': 'Proxy to Apache', 'value': 'nginx4'},
        {'name': 'Python (UWSGI)', 'value': 'nginx6'}
    ];

    $scope.apacheExamples = [
        {'name': 'PHP FPM', 'value': 'apache1'}
    ];

    $scope.deleteDialog = function() {
        var modalInstance = $modal.open({
            animation: true,
            templateUrl: 'domains_delete.html',
            controller: 'accountDomainEditDeleteDialog',
            size: '',
            resolve: {
                form: function () {
                    return $scope.form;
                }
            }
        });

        modalInstance.result.then(function (f) {
            managerServices.deleteDomain($routeParams.account, f.id).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/domains');
            }, function(err) {
                console.log(err);
            });
        }, function () {

        });
    };

    $scope.submit = function() {
        if ($scope.form.id === undefined) {
            managerServices.addDomain($routeParams.account, $scope.form).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/domains');
            },
            function(err) {
                $scope.errors = err.data.errors;
                console.log(err);
            });
        }else{
            managerServices.editDomain($routeParams.account, $scope.form.id, $scope.form).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/domains');
            },
            function(err) {
                $scope.errors = err.data.errors;
                console.log(err);
            });
        }
    };
}]).
controller('accountDomainEditDeleteDialog', ['$scope', '$modalInstance', 'form', function ($scope, $modalInstance, form) {
    $scope.form = form;
    $scope.confirm = false;

    $scope.delete = function () {
        if ($scope.confirm) {
            $modalInstance.close($scope.form);
        }
    };

    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
}]);