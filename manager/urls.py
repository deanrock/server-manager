from django.conf.urls import patterns, include, url
from django.contrib import admin

urlpatterns = patterns('',
    url(r'^accounts/login/$', 'django.contrib.auth.views.login',name="auth_login"),
    url(r'^accounts/logout/$', 'django.contrib.auth.views.logout',{'next_page': '/'}, name="auth_logout"),
    
    url(r'^action/(?P<action>.*)$', 'manager.views.action_ajax'),

    #admin
    url(r'^admin/', include(admin.site.urls)),
)
