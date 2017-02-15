var app = angular.module('managerApp', ['ngRoute', 'managerControllers', 'managerServices', 'managerDirectives', 'angular-blocks', 'ui.ace', 'ui.bootstrap', 'angularMoment', 'ngSanitize']);

app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {
    $routeProvider.
    when('/', {
        templateUrl: 'static/templates/accounts.html',
        controller: 'accounts',
        activetab: 'accounts'
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
        activetab: 'users'
        }).
    when('/users/:id', {
        templateUrl: 'static/templates/user/overview.html',
        controller: 'userOverview',
        activetab: 'users'
        }).
    when('/users/:id/access', {
        templateUrl: 'static/templates/user/access.html',
        controller: 'userAccess',
        activetab: 'users'
        }).
    when('/containers', {
        templateUrl: 'static/templates/containers.html',
        controller: 'containers',
        activetab: 'containers'
    }).
    when('/pull/images', {
        templateUrl: 'static/templates/sync/images.html',
        controller: 'pullImages',
    }).
    when('/sync/:action', {
        templateUrl: 'static/templates/sync.html',
        controller: 'sync',
    }).
    when('/a/:account', {
        templateUrl: 'static/templates/account/overview.html',
        controller: 'accountOverview',
        activetab: 'accounts'
        }).
    when('/a/:account/cronjobs', {
        templateUrl: 'static/templates/account/cronjobs.html',
        controller: 'accountCronjobs',
        activetab: 'accounts'
    }).
    when('/a/:account/cronjobs/add', {
        templateUrl: 'static/templates/account/cronjobs_edit.html',
        controller: 'accountCronjobAdd',
        activetab: 'accounts'
    }).
    when('/a/:account/cronjobs/:id', {
        templateUrl: 'static/templates/account/cronjobs_edit.html',
        controller: 'accountCronjobEdit',
        activetab: 'accounts'
    }).
    when('/a/:account/cronjobs/:id/logs', {
        templateUrl: 'static/templates/account/cronjobs_logs.html',
        controller: 'accountCronjobLogs',
        activetab: 'accounts'
    }).
    when('/a/:account/apps', {
        templateUrl: 'static/templates/account/apps.html',
        controller: 'accountApps',
        activetab: 'accounts'
    }).
    when('/a/:account/apps/add', {
        templateUrl: 'static/templates/account/apps_edit.html',
        controller: 'accountAppEdit',
        activetab: 'accounts'
    }).
    when('/a/:account/apps/:id', {
        templateUrl: 'static/templates/account/apps_edit.html',
        controller: 'accountAppEdit',
        activetab: 'accounts'
    }).
    when('/a/:account/apps/:action', {
        templateUrl: 'static/templates/account/frame.html',
        controller: 'account',
        activetab: 'accounts'
    }).
    when('/a/:account/apps/:id/logs', {
        templateUrl: 'static/templates/account/apps_logs.html',
        controller: 'accountAppLogs',
        activetab: 'accounts'
    }).
    when('/a/:account/domains', {
        templateUrl: 'static/templates/account/domains.html',
        controller: 'accountDomains',
        activetab: 'accounts'
    }).
    when('/a/:account/domains/add', {
        templateUrl: 'static/templates/account/domains_edit.html',
        controller: 'accountDomainEdit',
        activetab: 'accounts'
    }).
    when('/a/:account/domains/:id', {
        templateUrl: 'static/templates/account/domains_edit.html',
        controller: 'accountDomainEdit',
        activetab: 'accounts'
    }).
    when('/a/:account/databases', {
        templateUrl: 'static/templates/account/databases.html',
        controller: 'accountDatabases',
        activetab: 'accounts'
    }).
    when('/a/:account/databases/add', {
        templateUrl: 'static/templates/account/databases_edit.html',
        controller: 'accountDatabaseEdit',
        activetab: 'accounts'
    }).
    when('/a/:account/databases/:id', {
        templateUrl: 'static/templates/account/databases_edit.html',
        controller: 'accountDatabaseEdit',
        activetab: 'accounts'
    }).
    when('/a/:account/settings', {
        templateUrl: 'static/templates/account/settings.html',
        controller: 'accountSettings',
        activetab: 'accounts'
    }).
    when('/a/:account/settings/passwords/add', {
        templateUrl: 'static/templates/account/settings_passwords_edit.html',
        controller: 'accountSettingsPasswordAdd',
        activetab: 'accounts'
    }).
    when('/a/:account/:action', {
        templateUrl: 'static/templates/account/frame.html',
        controller: 'account',
        activetab: 'accounts'
    }).
    when('/profile/ssh-keys', {
        templateUrl: 'static/templates/profile/ssh_keys.html',
        controller: 'userSshKeys',
        activetab: 'profile'
    }).
    when('/profile/ssh-keys/add', {
        templateUrl: 'static/templates/profile/ssh_keys_edit.html',
        controller: 'userSshKeysEdit',
        activetab: 'profile'
    }).
    when('/profile/ssh-keys/:id', {
        templateUrl: 'static/templates/profile/ssh_keys_edit.html',
        controller: 'userSshKeysEdit',
        activetab: 'accounts'
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
