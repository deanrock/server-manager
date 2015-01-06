from django.conf.urls import patterns, include, url
from django.contrib import admin

urlpatterns = patterns('',
    # Examples:
    # url(r'^$', 'manager.views.home', name='home'),
    # url(r'^blog/', include('blog.urls')),

    url(r'^$', 'manager.views.index'),
    url(r'^sync-users/$', 'manager.views.sync_users'),
    url(r'^actions/container/(?P<id>[0-9]+)/(?P<action>[a-z]+)', 'manager.views.actions_container'),
    url(r'^rebuild-base-image/$', 'manager.views.rebuild_base_image'),

    url(r'^admin/', include(admin.site.urls)),
)
