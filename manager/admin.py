from django import forms
from django.contrib import admin
from manager.models import Account, App, Database, Domain, Image, ImagePort, ImageVariable


class AccountAdmin(admin.ModelAdmin):
    change_list_template = 'admin/accounts.html'


admin.site.register(Account, AccountAdmin)

class AppAdmin(admin.ModelAdmin):
    change_list_template = 'admin/apps.html'
    list_display = ['name', 'account', 'image', 'container_id', 'memory', 'status', 'start', 'stop', 'redeploy']

    def status(self, form):
        return ''

    def start(self, form):
        return '<a href="/actions/container/%s/start">start</a>' % (form.id)

    start.allow_tags = True

    def stop(self, form):
        return '<a href="/actions/container/%s/stop">stop</a>' % (form.id)

    stop.allow_tags = True

    def redeploy(self, form):
        return '<a href="/actions/container/%s/redeploy">redeploy</a>' % (form.id)

    redeploy.allow_tags = True

admin.site.register(App, AppAdmin)


class DatabaseAdmin(admin.ModelAdmin):
    change_list_template = 'admin/databases.html'

admin.site.register(Database, DatabaseAdmin)

class DomainAdmin(admin.ModelAdmin):
    list_display = ['name', 'account', 'redirect_url', 'apache_enabled']

admin.site.register(Domain, DomainAdmin)

class ImagePortInline(admin.StackedInline):
    model = ImagePort

class ImageVariableInline(admin.StackedInline):
    model = ImageVariable


class ImageAdmin(admin.ModelAdmin):
    inlines = (ImagePortInline, ImageVariableInline, )

admin.site.register(Image, ImageAdmin)
