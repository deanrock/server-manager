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
controller('accountOverview', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.action = 'overview';

    managerServices.getShells().then(function(data){
		$scope.shells = data;
	})

    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    })

    var shell = null;

    $scope.submit = function() {
        if (shell == null) {
            var selected = $scope.shell;
            shell = new AccountShell($scope.account.name, selected, document.getElementById('shell'));
            $('#interactive-shell-form select:first').prop('disabled', true);
            $('#interactive-shell-form button.submit').prop('disabled', true);
            $('#interactive-shell-form button.stop').show();
        }
    }

    $('#interactive-shell-form button.stop').click(function() {
        if (shell != null) {
            shell.stop();

            setTimeout(function() {
                $('#interactive-shell-form select:first').prop('disabled', false);
                $('#interactive-shell-form button.submit').prop('disabled', false);
                $('#interactive-shell-form button.stop').hide();
                shell = null;
            }, 0);
        }
    });
}]).
controller('account', ['$scope', 'managerServices', '$location', '$routeParams', function($scope, managerServices, $location, $routeParams) {
    $scope.suburl = '/frame' + $location.path();
    
    managerServices.getAccountByName($routeParams.account).then(function(data){
        $scope.account = data;
    })

    $scope.action = $routeParams.action;
}]);
