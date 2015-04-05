var ctrls = angular.module('managerControllers', []);


ctrls.controller('mainCtrl', ['$scope', '$rootScope', 'managerServices', function($scope, $rootScope, managerServices) {
    $scope.lol = "lal";
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
controller('accountCronjobs', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'cronjobs';

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
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
controller('userSshKeys', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.suburl = '/frame/profile/ssh-keys';
}]);
