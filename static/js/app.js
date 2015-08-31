var app = angular.module('managerApp', ['ngRoute', 'managerControllers', 'managerServices', 'managerDirectives', 'angular-blocks', 'ui.ace', 'ui.bootstrap']);

app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {
    $routeProvider.
    when('/', {
        templateUrl: 'static/templates/accounts.html',
        controller: 'accounts',
    }).
    when('/tasks/:id', {
        templateUrl: 'static/templates/tasks/task.html',
        controller: 'getTask',
    }).
    when('/tasks', {
        templateUrl: 'static/templates/tasks/tasks.html',
        controller: 'tasks',
    }).
    when('/users', {
        templateUrl: 'static/templates/users.html',
        controller: 'users',
    }).
    when('/users/:id', {
        templateUrl: 'static/templates/user/overview.html',
        controller: 'userOverview',
    }).
    when('/users/:id/access', {
        templateUrl: 'static/templates/user/access.html',
        controller: 'userAccess',
    }).
    when('/containers', {
        templateUrl: 'static/templates/containers.html',
        controller: 'containers',
    }).
    when('/sync/images', {
        templateUrl: 'static/templates/sync/images.html',
        controller: 'syncImages',
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
    when('/a/:account/cronjobs/add', {
        templateUrl: 'static/templates/account/cronjobs_edit.html',
        controller: 'accountCronjobAdd',
    }).
    when('/a/:account/cronjobs/:id', {
        templateUrl: 'static/templates/account/cronjobs_edit.html',
        controller: 'accountCronjobEdit',
    }).
    when('/a/:account/apps', {
        templateUrl: 'static/templates/account/apps.html',
        controller: 'accountApps',
    }).
    when('/a/:account/apps/add', {
        templateUrl: 'static/templates/account/apps_edit.html',
        controller: 'accountAppEdit',
    }).
    when('/a/:account/apps/:id', {
        templateUrl: 'static/templates/account/apps_edit.html',
        controller: 'accountAppEdit',
    }).
    when('/a/:account/apps/:action', {
        templateUrl: 'static/templates/account/frame.html',
        controller: 'account',
    }).
    when('/a/:account/apps/:id/logs', {
        templateUrl: 'static/templates/account/apps_logs.html',
        controller: 'accountAppLogs',
    }).
    when('/a/:account/apps/:id/:action', {
        templateUrl: 'static/templates/account/apps_action.html',
        controller: 'accountAppAction',
    }).
    when('/a/:account/domains', {
        templateUrl: 'static/templates/account/domains.html',
        controller: 'accountDomains',
    }).
    when('/a/:account/domains/add', {
        templateUrl: 'static/templates/account/domains_edit.html',
        controller: 'accountDomainEdit',
    }).
    when('/a/:account/domains/:id', {
        templateUrl: 'static/templates/account/domains_edit.html',
        controller: 'accountDomainEdit',
    }).
    when('/a/:account/:action', {
        templateUrl: 'static/templates/account/frame.html',
        controller: 'account',
    }).
    when('/profile/ssh-keys', {
        templateUrl: 'static/templates/profile/ssh_keys.html',
        controller: 'userSshKeys',
    }).
    when('/profile/ssh-keys/add', {
        templateUrl: 'static/templates/profile/ssh_keys_edit.html',
        controller: 'userSshKeysEdit',
    }).
    when('/profile/ssh-keys/:id', {
        templateUrl: 'static/templates/profile/ssh_keys_edit.html',
        controller: 'userSshKeysEdit',
    });
    $locationProvider.html5Mode(true);
}]).
run(function($rootScope, $location, $route, managerServices) {
    $rootScope.go = function(path) {
        $location.path(path);
    };

    $rootScope.$on("$locationChangeStart", function(event, next, current) {
        var path = $location.path().split('/');
        if (path[1] == 'a') {
            managerServices.getAccountByName(path[2]).then(function(data){
                $rootScope.account = data;
            })
        }
        console.log($location.path())
        $rootScope.path = $location.path();
    });

    $rootScope.$on("$routeChangeSuccess", function(event, currentRoute, previousRoute) {});

    managerServices.getProfile().then(function(data){
        $rootScope.profile = data;
    });
});
