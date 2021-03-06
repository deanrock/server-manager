var ctrls = angular.module('managerControllers', []);


ctrls.controller('mainCtrl', ['$scope', '$rootScope', 'managerServices', '$window', '$route', function($scope, $rootScope, managerServices, $window, $route) {
    $scope.$route = $route;

    $scope.djangoAdmin = function() {
        $window.location = '/admin'
    }

    $scope.syncWebServers = function() {
        managerServices.syncWebServers().then(function(data) {
            console.log(data);
        }, function(err) {
            console.log(err);
        });
    }

    $scope.syncAccounts = function() {
        managerServices.syncAccounts().then(function(data) {
            console.log(data);
        }, function(err) {
            console.log(err);
        });
    }

    $scope.purgeOldLogs = function() {
        managerServices.purgeOldLogs().then(function(data) {
            console.log(data);
        }, function(err) {
            console.log(err);
        });
    }

    $scope.logout = function() {
        document.cookie = 'session=; expires=' + +new Date() + ';';
        location.href = '/';
    }
}]).
controller('tasksCtrl', ['$scope', '$rootScope', 'managerServices', '$window', '$interval', '$location', function($scope, $rootScope, managerServices, $window, $interval, $location) {
    $scope.tasks = [];

    function setTaskInfo(t) {
        t.info = "";
        var vars = JSON.parse(t.variables);

        if (t.name == "start-app" || t.name == "stop-app" || t.name == "redeploy-app") {
            t.info = "("+vars.app.name+")";
        }else if (t.name == "sync-image") {
            t.info = "("+vars.image_name+")";
        }
    }

    function updateTask(t) {
        var found = false;
        angular.forEach($scope.tasks, function(v,k) {
            if (v.id == t.id) {
                //update
                t.added_at_timestamp = Date.parse(t.added_at);
                $scope.tasks[$scope.tasks.indexOf(v)] = t;
                showTimeElapsed(t);
                setTaskInfo(t);
                found = true;
                return;
            }
        });

        if (found) {
            return;
        }

        t.added_at_timestamp = Date.parse(t.added_at);
        $scope.tasks.push(t);
        showTimeElapsed(t);
        setTaskInfo(t);
    }

    function showTimeElapsed(v) {
        var format = function(x) {
            x = Math.floor(x);
            if (("" + x).length == 1) {
                return "0" + x;
            }

            return x;
        }

        var e = 0;
        if (v.finished) {
            e = v.duration;
        } else {
            var now = new Date().getTime();
            e = (now - v.added_at_timestamp) / 1000;
        }

        var m = Math.floor(e / 60);
        var s = e - m * 60;

        v.elapsed_time = format(m) + ":" + format(s);
    }

    $interval(function() {
        angular.forEach($scope.tasks, function(v,k) {
            showTimeElapsed(v);
        });
    }, 1000);

    $scope.goTo = function(task) {
        $location.path('/tasks/'+task.id);
    }

    $scope.showTasks = function() {
        $location.path('/tasks')
    }

    var ws = new ReconnectingWebSocket(getWebsocketHost() + '/ws/', null, {debug: true, reconnectInterval: 3000});

    ws.onopen = function() {
        console.log('open')
    }

    ws.onmessage = function(msg) {
        console.log('messssage '+msg.data)

        try {
            var m = JSON.parse(msg.data);
        }catch(err){
            return;
        }

        if (m.type == 'my-running-tasks') {
            $scope.$apply(function() {
                angular.forEach(m.tasks, function(v, k) {
                    updateTask(v);
                });
            });
        }else if (m.type == 'update-task') {
            $scope.$apply(function() {
                updateTask(m.task);
            });
        }
    }
}]).
controller('tasks', ['$scope', 'managerServices', '$location', function($scope, managerServices) {
    $scope.tasks = [];

    managerServices.getTasks().then(function(data){
        data.reverse();
        $scope.tasks = data;

        angular.forEach($scope.tasks, function(v, k) {
            v.expandedVariables = JSON.parse(v.variables);
            v.roundedDuration = Math.floor(v.duration);
        });
    })
}]).
controller('getTask', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.task = [];
    $scope.logs = [];

    managerServices.getTask($routeParams.id).then(function(data){
        $scope.task = data;

        var toShow = JSON.parse(JSON.stringify(data));
        delete toShow.variables;
        toShow.variables = JSON.parse($scope.task.variables);

        function syntaxHighlight(json) {
            if (typeof json != 'string') {
                 json = JSON.stringify(json, undefined, 2);
            }
            json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
            return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
                var cls = 'number';
                if (/^"/.test(match)) {
                    if (/:$/.test(match)) {
                        cls = 'key';
                    } else {
                        cls = 'string';
                    }
                } else if (/true|false/.test(match)) {
                    cls = 'boolean';
                } else if (/null/.test(match)) {
                    cls = 'null';
                }
                return '<span class="' + cls + '">' + match + '</span>';
            });
        }

        $scope.vars = syntaxHighlight(toShow);
    });

    managerServices.getTaskLog($routeParams.id).then(function(data){
        $scope.logs = data;
    });
}]).
controller('accounts', ['$scope', 'managerServices', '$location', function($scope, managerServices) {
    $scope.accounts = [];
    $scope.newAccount = {};

    $scope.reloadAccounts = function() {
        managerServices.getAccounts().then(function(data) {
            $scope.accounts = data;
        });
    };
    $scope.reloadAccounts();

    $scope.toggleAddAccount = function() {
        $scope.showAddAccount = !$scope.showAddAccount;
        $scope.newAccountMessage = null;
    }

    $scope.addAccount = function() {
        managerServices.addAccount($scope.newAccount).then(function(data) {
            $scope.reloadAccounts();
            $scope.newAccount = {};
            $scope.errors = {};
            $scope.showAddAccount = false;
            $scope.newAccountMessage = 'Account added successfully.';
        }, function(err) {
            $scope.errors = err.data.errors;
        });
    }
}]).
controller('containers', ['$scope', 'managerServices', '$location', function($scope, managerServices) {
    $scope.apps = [];

    managerServices.getContainers().then(function(data){
        $scope.apps = data;
    });

    $scope.start = function(app) {
        managerServices.startApp(app.account_name, app.app_id).then(function(data) {
            console.log(data);
        }, function(err) {
            console.log(err);
        });
    }

    $scope.stop = function(app) {
        managerServices.stopApp(app.account_name, app.app_id).then(function(data) {
            console.log(data);
        }, function(err) {
            console.log(err);
        });
    }

    $scope.redeploy = function(app) {
        managerServices.redeployApp(app.account_name, app.app_id).then(function(data) {
            console.log(data);
        }, function(err) {
            console.log(err);
        });
    }
}]).
controller('accountOverview', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'overview';
    $scope.started = false;

    managerServices.getShells().then(function(data){
		$scope.shells = data;

        $scope.shell = {'name': data[0]};
	})

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    })

    var shell = null;

    $scope.submit = function() {
        if (shell == null) {
            var selected = $scope.shell.name;
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

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getImages().then(function(data) {
        $scope.images = data;

        managerServices.getApps($routeParams.account).then(function(data) {
            $scope.apps = data;

            managerServices.getContainers().then(function(data){
                angular.forEach($scope.apps, function (vApp, kApp) {
                    angular.forEach(data, function(v,k) {
                        if (vApp.id == v.app_id) {
                            vApp.status = v.status;
                            vApp.up = v.up;
                        }
                    });

                    angular.forEach($scope.images, function(vImg, kImg) {
                        if (vImg.id == vApp.image_id) {
                            vApp.image_name = vImg.name;
                        }
                    })
                });
            });
        });
    });

    $scope.start = function(app) {
        managerServices.startApp($routeParams.account, app.id).then(function(data) {
            console.log(data);
        }, function(err) {
            console.log(err);
        });
    }

    $scope.stop = function(app) {
        managerServices.stopApp($routeParams.account, app.id).then(function(data) {
            console.log(data);
        }, function(err) {
            console.log(err);
        });
    }

    $scope.redeploy = function(app) {
        managerServices.redeployApp($routeParams.account, app.id).then(function(data) {
            console.log(data);
        }, function(err) {
            console.log(err);
        });
    }
}]).
controller('accountAppEdit', ['$scope', 'managerServices', '$location', '$routeParams', '$modal', function($scope, managerServices, $location, $routeParams, $modal) {
    $scope.images = [];
    $scope.action = 'apps';
    $scope.app = {};
    $scope.variables = [];
    $scope.image = null;
    $scope.errors = {};

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    $scope.showVariables = function() {
        if ($scope.images.length > 0) {
            angular.forEach($scope.images, function(v,k) {
                if (v.id == $scope.app.image_id) {
                    $scope.image = v;
                    $scope.variables = $scope.image.variables;

                    angular.forEach($scope.variables, function(v,k) {
                        angular.forEach($scope.app.variables, function (myvar, mykey) {
                            if (myvar.name == v.name) {
                                v.value = myvar.value;
                            }
                        })
                    })
                    return;
                }
            });
        }
    }

    managerServices.getImages().then(function(data) {
        $scope.images = data;

        if ($routeParams.id !== undefined) {
            managerServices.getApp($routeParams.account, $routeParams.id).then(function(data){
                $scope.app = data;

                $scope.showVariables();
            });
        }else{
            $scope.app.memory=256;

            if ($scope.images.length > 0) {
                $scope.app.image_id = $scope.images[0].id;
            }
        }
    });

    $scope.deleteDialog = function() {
        var modalInstance = $modal.open({
            animation: true,
            templateUrl: 'apps_delete.html',
            controller: 'accountAppEditDeleteDialog',
            size: '',
            resolve: {
                form: function () {
                    return $scope.app;
                }
            }
        });

        modalInstance.result.then(function (f) {
            managerServices.deleteApp($routeParams.account, f.id).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/apps');
            }, function(err) {
                console.log(err);
            });
        }, function () {

        });
    };

    $scope.submit = function() {
        console.log($scope.variables);
        angular.forEach($scope.variables, function(v,k) {
            var found = false;
            angular.forEach($scope.app.variables, function (myvar, mykey) {
                if (myvar.name == v.name) {
                    found = true;
                    myvar.value = v.value;
                }
            });

            if (!found) {
                if ($scope.app.variables === undefined) {
                    $scope.app.variables = [];
                }

                $scope.app.variables.push({
                    'app_id': $scope.app.id,
                    'name': v.name,
                    'value': v.value
                });
            }
        });

        if ($scope.app.id === undefined) {
            managerServices.addApp($routeParams.account, $scope.app).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/apps');
            },
            function(err) {
                $scope.errors = err.data.errors;
                console.log(err);
            });
        }else{
            managerServices.editApp($routeParams.account, $scope.app.id, $scope.app).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/apps');
            },
            function(err) {
                $scope.errors = err.data.errors;
                console.log(err);
            });
        }
    }
}]).
controller('accountAppEditDeleteDialog', ['$scope', '$modalInstance', 'form', function ($scope, $modalInstance, form) {
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
controller('accountCronjobs', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'cronjobs';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getCronjobs($routeParams.account).then(function(data) {
        $scope.cronjobs = data;
    });
}]).
controller('accountCronjobLogs', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'cronjobs';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getCronjob($routeParams.account, $routeParams.id).then(function(data) {
        $scope.cronjob = data;

        managerServices.getCronjobLog($routeParams.account, $routeParams.id).then(function(data) {
            angular.forEach(data, function(v, k) {
                v.roundedDuration = Math.floor(v.elapsed_time);
            })
            $scope.logs = data;
        });
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
    if ($routeParams.action == "users") {
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
controller('pullImages', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    managerServices.getSyncImages().then(function(data) {
        $scope.images = data;
    });

    $scope.pull = function(image) {
        managerServices.pullImage(image).then(function(data) {
            console.log('starting to pull ... i guess')
        }, function(err) {
            console.log('an error occured')
        })
    }
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
    });

    $scope.save = function() {
        managerServices.editUser($routeParams.id, $scope.user).then(function(data){
            $scope.errors = null;

            managerServices.getUser($routeParams.id).then(function(data){
                $scope.user = data;
            });
        }, function(err) {
            $scope.errors = err.data.errors;
        });
    }

    $scope.action = $routeParams.action;
}]).
controller('userAdd', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.save = function() {
        managerServices.addUser($scope.user).then(function(data){
            $scope.errors = null;
            $scope.user = data.data;
        }, function(err) {
            $scope.errors = err.data.errors;
        });
    }

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
controller('userSshKeys', ['$scope', 'managerServices', '$location', '$routeParams', '$modal', '$route', function($scope, managerServices, $location, $routeParams, $modal, $route) {
    $scope.tab = 'ssh-keys';
    $scope.keys = [];

    managerServices.getSSHKeys().then(function(data) {
        $scope.keys = data;
    });

    $scope.deleteDialog = function(key) {
        var modalInstance = $modal.open({
            animation: true,
            templateUrl: 'ssh_keys_delete.html',
            controller: 'userSshKeysDeleteDialog',
            size: '',
            resolve: {
                key: function () {
                    return key;
                }
            }
        });

        modalInstance.result.then(function (f) {
            managerServices.deleteSSHKey(f.id).then(function(data) {
                $route.reload();
            }, function(err) {
                console.log(err);
            });
        }, function () {

        });
    };
}]).
controller('userSshKeysEdit', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.tab = 'ssh-keys';
    $scope.key = {};

    if ($routeParams.id !== undefined) {
        managerServices.getSSHKey($routeParams.id).then(function(data) {
            $scope.key = data;
        });
    }

    $scope.submit = function() {
        if ($scope.key.id === undefined) {
            managerServices.addSSHKey($scope.key).then(function(data) {
                $location.path('/profile/ssh-keys');
            },
            function(err) {
                $scope.errors = err.data.errors;
                console.log(err);
            });
        }else{
            managerServices.editSSHKey($scope.key.id, $scope.key).then(function(data) {
                $location.path('/profile/ssh-keys');
            },
            function(err) {
                $scope.errors = err.data.errors;
                console.log(err);
            });
        }
    };
}]).
controller('userSshKeysDeleteDialog', ['$scope', '$modalInstance', 'key', function ($scope, $modalInstance, key) {
    $scope.key = key;
    $scope.confirm = false;

    $scope.delete = function () {
        if ($scope.confirm) {
            $modalInstance.close($scope.key);
        }
    };

    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
}]).
controller('accountDomains', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'domains';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getDomains($routeParams.account).then(function(data) {
        $scope.domains = data;
    });

    $scope.syncDomains = function() {
        managerServices.syncDomains($routeParams.account).then(function(data) {

        });
    }
}]).
controller('accountDomainEdit', ['$scope', 'managerServices', '$location', '$routeParams', '$modal', function($scope, managerServices, $location, $routeParams, $modal) {
    $scope.action = 'domains';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    $scope.example = function(id, value) {
       var source   = $("#" + value).html().replace('\\/\\/', '//');
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
}]).
controller('accountDatabases', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'databases';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getDatabases($routeParams.account).then(function(data) {
        $scope.databases = data;
    });

}]).
controller('accountDatabaseEdit', ['$scope', 'managerServices', '$location', '$routeParams', '$modal', function($scope, managerServices, $location, $routeParams, $modal) {
    $scope.action = 'databases';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    $scope.database_options = [
      "mysql",
      "postgres"
    ];

    $scope.form = {
        type: 'mysql'
    };
    if ($routeParams.id !== undefined) {
        managerServices.getDatabase($routeParams.account, $routeParams.id).then(function(data) {
            $scope.form = data;
        });
    }

    $scope.deleteDialog = function() {
        var modalInstance = $modal.open({
            animation: true,
            templateUrl: 'apps_delete.html',
            controller: 'accountAppEditDeleteDialog',
            size: '',
            resolve: {
                form: function () {
                    return $scope.app;
                }
            }
        });

        modalInstance.result.then(function (f) {
            managerServices.deleteApp($routeParams.account, f.id).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/apps');
            }, function(err) {
                console.log(err);
            });
        }, function () {

        });
    };


    $scope.submit = function() {
        if ($scope.form.id === undefined) {
            managerServices.addDatabase($routeParams.account, $scope.form).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/databases');
            },
            function(err) {
                $scope.errors = err.data.errors;
                console.log(err);
            });
        }else{
            managerServices.editDatabase($routeParams.account, $scope.form.id, $scope.form).then(function(data) {
                $location.path('/a/'+$scope.account.name+'/databases');
            },
            function(err) {
                $scope.errors = err.data.errors;
                console.log(err);
            });
        }
    };
}]).
controller('accountSettings', ['$scope', 'managerServices', '$location', '$routeParams', '$modal', '$route', function($scope, managerServices, $location, $routeParams, $modal, $route) {
    $scope.action = 'settings';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    managerServices.getSSHPasswords($routeParams.account).then(function(data) {
        $scope.passwords = data;
    });

    $scope.deletePasswordDialog = function(password) {
        var modalInstance = $modal.open({
            animation: true,
            templateUrl: 'settings_passwords_delete.html',
            controller: 'accountSettingsPasswordDeleteDialog',
            size: '',
            resolve: {
                password: function () {
                    return password;
                }
            }
        });

        modalInstance.result.then(function (f) {
            managerServices.deleteSSHPassword($routeParams.account, f.id).then(function(data) {
                $route.reload();
            }, function(err) {
                console.log(err);
            });
        }, function () {

        });
    };
}]).
controller('accountSettingsPasswordAdd', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'settings';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    });

    $scope.form = {};

    $scope.submit = function() {
        managerServices.addSSHPassword($routeParams.account, $scope.form).then(function(data) {
            $location.path('/a/'+$scope.account.name+'/settings');
        },
        function(err) {
            $scope.errors = err.data.errors;
            console.log(err);
        });
    }
}]).
controller('accountSettingsPasswordDeleteDialog', ['$scope', '$modalInstance', 'password', function ($scope, $modalInstance, password) {
    $scope.password = password;
    $scope.confirm = false;

    $scope.delete = function () {
        if ($scope.confirm) {
            $modalInstance.close($scope.password);
        }
    };

    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
}]);
