from django.forms import ModelForm
from manager.models import App, Domain, Database


class AppForm(ModelForm):
    class Meta:
        model = App
        fields = ['name', 'image', 'memory']


class DomainForm(ModelForm):
    class Meta:
        model = Domain
        fields = ['name', 'redirect_url', 'nginx_config', 'apache_config', 'apache_enabled']

class DatabaseForm(ModelForm):
    class Meta:
        model = Database
        fields = ['type', 'name', 'user', 'password']
