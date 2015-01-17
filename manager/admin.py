from django import forms
from django.contrib import admin
from django.core.urlresolvers import reverse
from manager.models import Account, App, Database, Domain, Image, ImagePort, ImageVariable


class MyModelAdmin(admin.ModelAdmin):
    class Media:
        js = (
            'js/admin.js',
        )

class AppAdmin(MyModelAdmin):
    change_list_template = 'admin/apps.html'
    list_display = ['name', 'account', 'image', 'container_id', 'memory', 'status', 'start', 'stop', 'redeploy']
    list_filter = ('account', 'memory', 'image', )
    readonly_fields = ('added_by',)

    def status(self, form):
        return ''

    def start(self, form):
        return '<a href="%s">start</a>' % reverse('manager.views.account_apps_action', kwargs={'name': form.account.name, 'app': form.id, 'action': 'start'})

    start.allow_tags = True

    def stop(self, form):
        return '<a href="%s">stop</a>' % reverse('manager.views.account_apps_action', kwargs={'name': form.account.name, 'app': form.id, 'action': 'stop'})

    stop.allow_tags = True

    def redeploy(self, form):
        return '<a href="%s">redeploy</a>' % reverse('manager.views.account_apps_action', kwargs={'name': form.account.name, 'app': form.id, 'action': 'redeploy'})

    redeploy.allow_tags = True

    def save_model(self, request, obj, form, change):
        try:
            obj.added_by
        except:
            obj.added_by = request.user
        obj.save()

admin.site.register(App, AppAdmin)


class DatabaseAdmin(MyModelAdmin):
    change_list_template = 'admin/databases.html'
    list_display = ['name', 'account']
    readonly_fields = ('added_by',)

    def save_model(self, request, obj, form, change):
        try:
            obj.added_by
        except:
            obj.added_by = request.user
        obj.save()

admin.site.register(Database, DatabaseAdmin)


class AccountAdmin(MyModelAdmin):
    change_list_template = 'admin/accounts.html'
    readonly_fields = ('added_by',)

    def save_model(self, request, obj, form, change):
        try:
            obj.added_by
        except:
            obj.added_by = request.user
        obj.save()


admin.site.register(Account, AccountAdmin)
