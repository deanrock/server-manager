from django.conf.urls import patterns, include, url
from django.contrib import admin

urlpatterns = patterns('',
    # Examples:
    # url(r'^$', 'manager.views.home', name='home'),
    # url(r'^blog/', include('blog.urls')),

    url(r'^$', 'manager.views.index'),

    url(r'^a/(?P<name>[a-z0-9]+)/$', 'manager.views.account'),
    url(r'^a/(?P<name>[a-z0-9]+)/apps/$', 'manager.views.account_apps'),
    url(r'^a/(?P<name>[a-z0-9]+)/apps/(?P<app>[0-9]+)$', 'manager.views.account_apps_edit'),
    url(r'^a/(?P<name>[a-z0-9]+)/apps/add$', 'manager.views.account_apps_edit'),
    url(r'^a/(?P<name>[a-z0-9]+)/domains/$', 'manager.views.account_domains'),
    url(r'^a/(?P<name>[a-z0-9]+)/domains/(?P<domain>[0-9]+)$', 'manager.views.account_domains_edit'),
    url(r'^a/(?P<name>[a-z0-9]+)/domains/add$', 'manager.views.account_domains_edit'),
    url(r'^a/(?P<name>[a-z0-9]+)/domains/(?P<domain>[0-9]+)/delete$', 'manager.views.account_domains_delete'),
    url(r'^a/(?P<name>[a-z0-9]+)/databases/$', 'manager.views.account_databases'),
    url(r'^a/(?P<name>[a-z0-9]+)/databases/(?P<database>[0-9]+)$', 'manager.views.account_databases_edit'),
    url(r'^a/(?P<name>[a-z0-9]+)/databases/add$', 'manager.views.account_databases_edit'),

    url(r'containers/', 'manager.views.containers'),

    url(r'^sync-users/$', 'manager.views.sync_users'),
    url(r'^actions/container/(?P<id>[0-9]+)/(?P<action>[a-z]+)', 'manager.views.actions_container'),
    url(r'^rebuild-base-image/$', 'manager.views.rebuild_base_image'),
    url(r'^update-nginx-config/$', 'manager.views.update_nginx_config'),
    url(r'^sync-databases/$', 'manager.views.sync_databases'),

    url(r'^admin/', include(admin.site.urls)),
)
