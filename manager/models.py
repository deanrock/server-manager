import json
import shutil
from django.conf import settings
from django.core.validators import RegexValidator
from manager import docker_api, utils
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

    def __str__(self):
        return '%s (%s)' % (self.name, self.description)


class Domain(models.Model):
    account = models.ForeignKey(Account, related_name='domains')
    name = models.CharField(max_length=255)
    redirect_url = models.CharField(max_length=255, default=None, blank=True)
    nginx_config = models.TextField(default=None, blank=True)
    apache_config = models.TextField(blank=True)

    apache_enabled = models.BooleanField(default=False, blank=True)

    added_at = models.DateTimeField(auto_now_add=True, null=True, blank=True)
    added_by = models.ForeignKey(User)

class ImageVariable(models.Model):
    name = models.CharField(max_length=255)
    description = models.TextField()
    default = models.TextField(null=True, blank=True)

    image = models.ForeignKey('Image', related_name='variables')

class ImagePort(models.Model):
    types = (
        ('fastcgi', 'FastCGI'),
        ('http', 'HTTP'),
    )
    type = models.CharField(max_length=255, choices=types)
    port = models.IntegerField()

    image = models.ForeignKey('Image', related_name='ports')

class Image(models.Model):
    name = models.TextField(max_length=255)
    description = models.TextField()
    types = (
        ('application', 'Application'),
    )
    type = models.CharField(max_length=255, choices=types)

    added_at = models.DateTimeField(auto_now_add=True, null=True, blank=True)
    added_by = models.ForeignKey(User)

    def folder(self):
        return os.path.abspath(os.path.join(settings.BASE_DIR, '../', 'images/', self.name)) #'%s/../images/%s' % (settings.BASE_DIR, self.name)

    def __str__(self):
        return self.name

class AppImageVariable(models.Model):
    app = models.ForeignKey('App', related_name='variables')
    image_variable = models.ForeignKey('ImageVariable')
    value = models.TextField()

class App(models.Model):
    account = models.ForeignKey(Account, related_name='apps')
    name = models.CharField(max_length=50, validators=[appname_validator])
    image = models.ForeignKey(Image)
    container_id = models.TextField(max_length=255, null=True, blank=True)
    memory = models.IntegerField(default=256)

    added_at = models.DateTimeField(auto_now_add=True, null=True, blank=True)
    added_by = models.ForeignKey(User)

    def container_name(self):
        return 'app-%s-%s' % (self.account.name, self.name)

    def image_name(self):
        return 'manager/'+self.container_name()

    def stop(self):
        return [docker_api.cli.stop(container=self.container_name())]

    def start(self):
        return [docker_api.cli.start(container=self.container_name())]

    def redeploy(self):
        logs = utils.Logs()
        failed = False

        try:
            logs.add(docker_api.cli.stop(self.container_name()))
        except:
            logs.add('Contaner not started; cannot stop it')

        try:
            logs.add(docker_api.cli.remove_container(self.container_name(), link=True))
        except Exception as e:
            logs.add('Container doesnt exists, cannot remove it %s' % e)

        try:
            out, err = utils.exec_command(logs, "sudo docker rm %s" % self.container_name())
        except Exception as e:
            logs.add('err %s', e)

        try:
            logs.add(docker_api.cli.remove_image(self.image_name()))
        except:
            logs.add('Docker image doesnt exists, cannot remove it')

        #create dockerfile folder
        temp_folder = utils.get_temp_folder()

        logs.add("image folder %s" % self.image.folder())
        logs.add("temp folder %s" % temp_folder)

        try:
            src_files = os.listdir(self.image.folder())
            for file_name in src_files:
                full_file_name = os.path.join(self.image.folder(), file_name)
                if (os.path.isfile(full_file_name)) and file_name != 'Dockerfile':
                    shutil.copy(full_file_name, temp_folder)

            dockerfile = os.path.join(temp_folder, "Dockerfile")

            userid, err = utils.exec_command(logs, 'id -u %s' % self.account.name)



            with open(os.path.join(self.image.folder(), "Dockerfile"),'r') as read:
                with open(dockerfile, 'w') as write:
                    for line in read:
                        write.write(line.replace("#user#", self.account.name). \
                            replace("#uid#", userid.rstrip()))
        except Exception as e:
            logs.add("error while preparing Dockerfile ... %s" % str(e))

            return logs

        #build image
        try:
            for line in docker_api.cli.build(
               path=temp_folder, rm=True, tag=self.image_name()
            ):
                logs.add(line)

            utils.delete_temp_folder(temp_folder)
        except Exception as e:
            logs.add("error while building image ... %s" % str(e))
            return logs

        #continue, create container
        m = '%dm' % self.memory
        try:
            container = docker_api.cli.create_container(image=self.image_name(),
                                                        user=self.account.name,
                                                        mem_limit=m,
                                                        name=self.container_name())

            logs.add(json.dumps(container))

            homefolder = '/home/%s' % self.account.name
            try:
                response = docker_api.cli.start(container=container.get('Id'),
                                                restart_policy={
                                            "MaximumRetryCount": 0,
                                            "Name": "always"
                                        },
                                                extra_hosts={
                                                    'mysql': '172.17.42.1'
                                                },
                                                binds={
                                homefolder: {
                                    'bind': homefolder
                                }})

                logs.add(response)
            except Exception as e:
                logs.add("error while starting container ... %s" % str(e))
                return logs

        except Exception as e:
            logs.add("error while creating container ... %s" % str(e))
            return logs

        return logs

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
