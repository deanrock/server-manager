import json
import codecs
import shutil
from django.conf import settings
from django.core.validators import RegexValidator
import os
from django.db import models
from django.contrib.auth.models import User

account_validator = RegexValidator(r'^[a-zA-Z][0-9a-zA-Z-]*$', 'Only alphanumeric characters and \'-\' are allowed.')
db_validator = RegexValidator(r'^[a-zA-Z][0-9a-zA-Z_]*$', 'Only alphanumeric characters and underscore are allowed.')
appname_validator = RegexValidator(r'^[a-zA-Z][0-9a-zA-Z_-]*$', 'Only alphanumeric characters, underscore and \'-\' are allowed.')

class UserSSHKey(models.Model):
    user = models.ForeignKey(User, related_name='ssh_keys')

    name = models.CharField(max_length=255)
    ssh_key = models.TextField(max_length=1000)

    added_at = models.DateTimeField(auto_now_add=True, null=True, blank=True)
    added_by = models.ForeignKey(User)


class Account(models.Model):
    name = models.CharField(max_length=32, help_text='max 32 chars!', validators=[account_validator])
    description = models.CharField(max_length=1000)

    added_at = models.DateTimeField(auto_now_add=True, null=True, blank=True)
    added_by = models.ForeignKey(User)

    users = models.ManyToManyField(User, related_name='users')

    def variables(self):
        v = {
            "user": self.name
        }

        return v


    def __str__(self):
        return '%s (%s)' % (self.name, self.description)


class Domain(models.Model):
    account = models.ForeignKey(Account, related_name='domains')
    name = models.CharField(max_length=255)
    redirect_url = models.CharField(max_length=255, default=None, blank=True)
    nginx_config = models.TextField(default=None, blank=True)
    apache_config = models.TextField(blank=True)

    apache_enabled = models.BooleanField(default=False, blank=True)

    ssl_enabled = models.BooleanField(default=False, blank=True)

    added_at = models.DateTimeField(auto_now_add=True, null=True, blank=True)
    added_by = models.ForeignKey(User)

class ImageVariable(models.Model):
    name = models.CharField(max_length=255)
    description = models.TextField()
    default = models.TextField(null=True, blank=True)
    filename = models.CharField(null=True, blank=True, max_length=255)

    image = models.ForeignKey('Image', related_name='variables')

class ImagePort(models.Model):
    types = (
        ('fastcgi', 'FastCGI'),
        ('http', 'HTTP'),
        ('uwsgi', 'UWSGI'),
        ('tcp', 'TCP'),
    )
    type = models.CharField(max_length=255, choices=types)
    port = models.IntegerField()

    image = models.ForeignKey('Image', related_name='ports')

class Image(models.Model):
    name = models.TextField(max_length=255)
    description = models.TextField()
    types = (
        ('application', 'Application'),
        ('database', 'Database'),
    )
    type = models.CharField(max_length=255, choices=types)

    added_at = models.DateTimeField(auto_now_add=True, null=True, blank=True)

    def folder(self):
        return os.path.abspath(os.path.join(settings.BASE_DIR, '../', 'images/', self.name)) #'%s/../images/%s' % (settings.BASE_DIR, self.name)

    def __str__(self):
        return self.name

class AppImageVariable(models.Model):
    app = models.ForeignKey('App', related_name='variables')
    name = models.CharField(max_length=255, default='')
    value = models.TextField()

class App(models.Model):
    account = models.ForeignKey(Account, related_name='apps')
    name = models.CharField(max_length=50, validators=[appname_validator])
    image = models.ForeignKey(Image)
    container_id = models.TextField(max_length=255, null=True, blank=True)
    memory = models.IntegerField(default=256)

    added_at = models.DateTimeField(auto_now_add=True, null=True, blank=True)
    added_by = models.ForeignKey(User)

    class Meta:
        unique_together = ('name', 'account',)

    def container_name(self):
        return 'app-%s-%s' % (self.account.name, self.name)

    def image_name(self):
        return 'manager/'+self.container_name()


class Database(models.Model):
    account = models.ForeignKey(Account, related_name='databases')
    types = (
        ('mysql', 'MySQL'),
    )
    type = models.CharField(max_length=50, choices=types)
    name = models.CharField(max_length=50, help_text='Max 50 characters', unique=True, validators=[db_validator])
    user = models.CharField(max_length=16, help_text='Max 16 characters!', unique=True, validators=[db_validator])
    password = models.CharField(max_length=255)

    added_at = models.DateTimeField(auto_now_add=True, null=True, blank=True)
    added_by = models.ForeignKey(User)

    def __str__(self):
        return self.name
