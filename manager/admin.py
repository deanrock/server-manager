from django import forms
from django.contrib import admin
from manager.models import Account, App, Database, Domain, Image, ImagePort, ImageVariable


class MyModelAdmin(admin.ModelAdmin):
    class Media:
        js = (
            'js/admin.js',
        )

class AccountAdmin(MyModelAdmin):
    change_list_template = 'admin/accounts.html'


admin.site.register(Account, AccountAdmin)

class AppAdmin(MyModelAdmin):
    change_list_template = 'admin/apps.html'
    list_display = ['name', 'account', 'image', 'container_id', 'memory', 'status', 'start', 'stop', 'redeploy']
    list_filter = ('account', 'memory', 'image', )

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


class DatabaseAdmin(MyModelAdmin):
    change_list_template = 'admin/databases.html'
    list_display = ['name', 'account']

admin.site.register(Database, DatabaseAdmin)

class DomainAdmin(MyModelAdmin):
    list_display = ['name', 'account', 'redirect_url', 'apache_enabled']
    change_list_template = 'admin/domains.html'

admin.site.register(Domain, DomainAdmin)

class ImagePortInline(admin.StackedInline):
    model = ImagePort

class ImageVariableInline(admin.StackedInline):
    model = ImageVariable


class ImageAdmin(MyModelAdmin):
    inlines = (ImagePortInline, ImageVariableInline, )

admin.site.register(Image, ImageAdmin)
