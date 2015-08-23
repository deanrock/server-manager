from django import forms
from django.contrib import admin
from django.core.urlresolvers import reverse
from manager.models import Account, App, Database, Domain, Image, ImagePort, ImageVariable


class MyModelAdmin(admin.ModelAdmin):
    pass


class AccountAdmin(MyModelAdmin):
    readonly_fields = ('added_by',)

    def save_model(self, request, obj, form, change):
        try:
            obj.added_by
        except:
            obj.added_by = request.user
        obj.save()


admin.site.register(Account, AccountAdmin)
