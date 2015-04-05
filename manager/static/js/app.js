var app = angular.module('managerApp', ['ngRoute', 'managerControllers', 'managerServices', 'ngLayout', 'managerDirectives']);

app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {
    $routeProvider.
    when('/', {
        templateUrl: 'static/templates/accounts.html',
        controller: 'accounts',
    }).
    when('/containers', {
        templateUrl: 'static/templates/containers.html',
        controller: 'containers',
    }).
    when('/sync/:action', {
        templateUrl: 'static/templates/sync.html',
        controller: 'sync',
    }).
    when('/a/:account', {
        templateUrl: 'static/templates/account/overview.html',
        controller: 'accountOverview',
    }).
    when('/a/:account/cronjobs', {
        templateUrl: 'static/templates/account/cronjobs.html',
        controller: 'accountCronjobs',
    }).
    when('/a/:account/:action', {
        templateUrl: 'static/templates/account/frame.html',
        controller: 'account',
    }).
    when('/profile/ssh-keys', {
        templateUrl: 'static/templates/ssh-keys.html',
        controller: 'userSshKeys',
    });
    $locationProvider.html5Mode(true);
}]).
run(function($rootScope, $location, $route, managerServices) {
    $rootScope.go = function(path) {
        $location.path(path);
    };

    $rootScope.$on("$locationChangeStart", function(event, next, current) {
        console.log($location.path())
        $rootScope.path = $location.path();
    });

    $rootScope.$on("$routeChangeSuccess", function(event, currentRoute, previousRoute) {});

    managerServices.getProfile().then(function(data){
        $rootScope.profile = data;
    });
});
