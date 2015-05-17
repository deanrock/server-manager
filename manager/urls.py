from django.conf.urls import patterns, include, url
from django.contrib import admin

urlpatterns = patterns('',
    # Examples:
    # url(r'^$', 'manager.views.home', name='home'),
    # url(r'^blog/', include('blog.urls')),

    url(r'^$', 'manager.views.index'),

    url(r'^frame/profile/ssh-keys/$', 'manager.views.profile_sshkeys'),
    url(r'^frame/profile/ssh-keys/new$', 'manager.views.profile_sshkeys_edit'),
    url(r'^frame/profile/ssh-keys/(?P<key>[0-9]+)$', 'manager.views.profile_sshkeys_edit'),
    url(r'^frame/profile/ssh-keys/(?P<key>[0-9]+)/delete$', 'manager.views.profile_sshkeys_delete'),

    url(r'^accounts/login/$', 'django.contrib.auth.views.login',name="auth_login"),
    url(r'^accounts/logout/$', 'django.contrib.auth.views.logout',{'next_page': '/'}, name="auth_logout"),

    #frame
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/$', 'manager.views.account'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/cronjobs/$', 'manager.views.account_cronjobs'),
    
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/apps/(?P<app>[0-9]+)$', 'manager.views.account_apps_edit'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/apps/add$', 'manager.views.account_apps_edit'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/domains/$', 'manager.views.account_domains'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/domains/(?P<domain>[0-9]+)$', 'manager.views.account_domains_edit'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/domains/add$', 'manager.views.account_domains_edit'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/domains/(?P<domain>[0-9]+)/delete$', 'manager.views.account_domains_delete'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/databases/$', 'manager.views.account_databases'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/databases/(?P<database>[0-9]+)$', 'manager.views.account_databases_edit'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/databases/add$', 'manager.views.account_databases_edit'),

    #ajax
    url(r'^api/v1.0/containers/', 'manager.views.containers'),
    
    url(r'^action/(?P<action>.*)$', 'manager.views.action_ajax'),

    url(r'^api/v1.0/images/(?P<id>[0-9]+)', 'manager.views.api_image'),

    #django apis
    url(r'^djapi/v1/accounts/(?P<name>[a-z0-9-]+)/apps/$', 'manager.views.account_apps'),
    url(r'^djapi/v1/accounts/(?P<name>[a-z0-9-]+)/apps/(?P<app>[0-9]+)$', 'manager.views.account_apps_info'),
    url(r'^djapi/v1/accounts/(?P<name>[a-z0-9-]+)/apps/(?P<app>[0-9]+)/(?P<action>[a-z]+)/ajax', 'manager.views.account_apps_action_ajax'),



    # old apis
    url(r'^api/v1.0/a/(?P<name>[a-z0-9-]+)/apps/(?P<id>[0-9]+)', 'manager.views.api_account_app'),
    url(r'^api/v1.0/a/(?P<name>[a-z0-9-]+)/variables', 'manager.views.api_account_variables'),
    
    #admin
    url(r'^admin/', include(admin.site.urls)),
)
