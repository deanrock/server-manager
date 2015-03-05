var app = angular.module('managerApp', ['ngRoute', 'managerControllers', 'managerServices', 'ngLayout', 'managerDirectives']);

app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {
    $routeProvider.
    when('/', {
        templateUrl: 'static/templates/accounts.html',
        controller: 'accounts',
    }).
     when('/a/:account', {
        templateUrl: 'static/templates/account/overview.html',
        controller: 'accountOverview',
    }).
    when('/a/:account/:action', {
        templateUrl: 'static/templates/account/frame.html',
        controller: 'account',
    });
    $locationProvider.html5Mode(true);
}]).
run(function($rootScope, $location, $route) {
    $rootScope.go = function(path) {
        $location.path(path);
    };

    $rootScope.$on("$locationChangeStart", function(event, next, current) {
        console.log($location.path())
        $rootScope.path = $location.path();
    });

    $rootScope.$on("$routeChangeSuccess", function(event, currentRoute, previousRoute) {});
});
