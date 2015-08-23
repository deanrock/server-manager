from django.conf.urls import patterns, include, url
from django.contrib import admin

urlpatterns = patterns('',
    url(r'^accounts/login/$', 'django.contrib.auth.views.login',name="auth_login"),
    url(r'^accounts/logout/$', 'django.contrib.auth.views.logout',{'next_page': '/'}, name="auth_logout"),

    #frame
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/databases/$', 'manager.views.account_databases'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/databases/(?P<database>[0-9]+)$', 'manager.views.account_databases_edit'),
    url(r'^frame/a/(?P<name>[a-z0-9-]+)/databases/add$', 'manager.views.account_databases_edit'),
    
    url(r'^action/(?P<action>.*)$', 'manager.views.action_ajax'),

    #django apis
    url(r'^djapi/v1/accounts/(?P<name>[a-z0-9-]+)/apps/(?P<app>[0-9]+)/(?P<action>[a-z]+)/ajax', 'manager.views.account_apps_action_ajax'),

    #admin
    url(r'^admin/', include(admin.site.urls)),
)
