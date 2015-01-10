from django import forms
from django.contrib import admin
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
        return '<a href="/actions/container/%s/start">start</a>' % (form.id)

    start.allow_tags = True

    def stop(self, form):
        return '<a href="/actions/container/%s/stop">stop</a>' % (form.id)

    stop.allow_tags = True

    def redeploy(self, form):
        return '<a href="/actions/container/%s/redeploy">redeploy</a>' % (form.id)

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

class DomainAdmin(MyModelAdmin):
    list_display = ['name', 'account', 'redirect_url', 'apache_enabled']
    readonly_fields = ('added_by',)
    change_list_template = 'admin/domains.html'

    def save_model(self, request, obj, form, change):
        try:
            obj.added_by
        except:
            obj.added_by = request.user
        obj.save()

admin.site.register(Domain, DomainAdmin)

class ImagePortInline(admin.StackedInline):
    model = ImagePort

class ImageVariableInline(admin.StackedInline):
    model = ImageVariable


class ImageAdmin(MyModelAdmin):
    inlines = (ImagePortInline, ImageVariableInline, )
    readonly_fields = ('added_by',)

    def save_model(self, request, obj, form, change):
        try:
            obj.added_by
        except:
            obj.added_by = request.user
        obj.save()

admin.site.register(Image, ImageAdmin)

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
