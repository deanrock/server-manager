from django.forms import ModelForm
from manager.models import App, Domain, Database, UserSSHKey


class AppForm(ModelForm):
    class Meta:
        model = App
        fields = ['name', 'image', 'memory']


class DomainForm(ModelForm):
    class Meta:
        model = Domain
        fields = ['name', 'redirect_url', 'nginx_config', 'apache_config', 'apache_enabled', 'ssl_enabled']

class DatabaseForm(ModelForm):
    class Meta:
        model = Database
        fields = ['type', 'name', 'user', 'password']

class UserSSHKeyForm(ModelForm):
    class Meta:
        model = UserSSHKey
        fields = ['name', 'ssh_key']
