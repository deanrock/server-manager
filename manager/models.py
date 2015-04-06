import json
import codecs
import shutil
from django.conf import settings
from django.core.validators import RegexValidator
from manager import docker_api, utils
import os
from django.db import models
from django.contrib.auth.models import User
import utils

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

        logs = utils.Logs()
        userid, err = utils.exec_command(logs, 'id -u %s' % self.name)
        v["uid"] = userid.rstrip()

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

        #try:
        src_files = os.listdir(self.image.folder())
        for file_name in src_files:
            full_file_name = os.path.join(self.image.folder(), file_name)
            if (os.path.isfile(full_file_name)) and file_name != 'Dockerfile':
                shutil.copy(full_file_name, temp_folder)

        userid, err = utils.exec_command(logs, 'id -u %s' % self.account.name)

        def copy_file(name, append=None):
            print name
            if os.path.exists(os.path.join(self.image.folder(), name)):
                with codecs.open(os.path.join(self.image.folder(), name),'r', encoding='utf8') as read:
                    contents = read.read()

                    contents = contents.replace("#user#", self.account.name). \
                        replace("#uid#", userid.rstrip())

                    contents = contents.replace("#appname#", self.name)

                    variables = {}
                    names = []

                    for v in self.image.variables.all():
                        if v.default:
                            variables[v.name] = v.default

                        names.append(v.name)

                    for v in self.variables.all():
                        for n in names:
                            if v.name in n:
                                variables[v.name] = v.value

                    logs.add("variables:\n %s" % variables)

                    for v in variables:
                        contents = contents.replace('#variable_%s#' % v, variables[v])

                    if append:
                        contents = contents + "\n" + append + "\n"

                    logs.add("%s:\n%s" % (name, contents))

                    new_name = os.path.join(temp_folder, name)
                    with codecs.open(new_name, 'w', encoding='utf8') as write:
                        write.write(contents)


        copy_file('Dockerfile')
        copy_file('start.sh')

        file_variables = {}
        for v in self.image.variables.all():
            if v.filename:
                file_variables[v.name] = v

        found = []

        for v in self.variables.all():
            if v.name in file_variables:
                found.append(v.name)
                copy_file(file_variables[v.name].filename, append=v.value)

        for v in file_variables:
            if v not in found:
                copy_file(file_variables[v].filename)


        #except Exception as e:
        #    logs.add("error while preparing Dockerfile ... %s" % str(e))

        #    return logs

        #build image
        try:
            for line in docker_api.cli.build(
               path=temp_folder, rm=True, tag=self.image_name()
            ):
                logs.add(line)

            #utils.delete_temp_folder(temp_folder)
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

            hosts = {
                'mysql': '172.17.42.1'
            }

            docker_mapping={}

            containers = docker_api.cli.containers()

            for c in containers:
                name = None
                try:
                    name = c['Names'][0].strip('/')
                except:
                    continue

                ip = None
                try:
                    ip = docker_api.cli.inspect_container(name)['NetworkSettings']['IPAddress']
                except:
                    continue

                docker_mapping[name] = ip

            print docker_mapping

            for app in self.account.apps.filter(image__type='database'):
                    if app.container_name() in docker_mapping:
                        hosts[app.name] = docker_mapping[app.container_name()]

            print hosts

            try:
                response = docker_api.cli.start(container=container.get('Id'),
                                                restart_policy={
                                            "MaximumRetryCount": 0,
                                            "Name": "always"
                                        },
                                                extra_hosts=hosts,
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
