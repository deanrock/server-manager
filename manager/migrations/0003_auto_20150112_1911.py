# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations
from django.conf import settings
import django.core.validators


class Migration(migrations.Migration):

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
        ('manager', '0002_auto_20150110_1718'),
    ]

    operations = [
        migrations.CreateModel(
            name='UserSSHKey',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('name', models.CharField(max_length=255)),
                ('ssh_key', models.CharField(max_length=1000)),
                ('added_at', models.DateTimeField(auto_now_add=True, null=True)),
                ('added_by', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
                ('user', models.ForeignKey(related_name=b'ssh_keys', to=settings.AUTH_USER_MODEL)),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.AlterField(
            model_name='account',
            name='name',
            field=models.CharField(help_text=b'max 32 chars!', max_length=32, validators=[django.core.validators.RegexValidator(b'^[a-zA-Z][0-9a-zA-Z-]*$', b"Only alphanumeric characters and '-' are allowed.")]),
        ),
        migrations.AlterField(
            model_name='app',
            name='name',
            field=models.TextField(max_length=50, validators=[django.core.validators.RegexValidator(b'^[a-zA-Z][0-9a-zA-Z_-]*$', b"Only alphanumeric characters, underscore and '-' are allowed.")]),
        ),
        migrations.AlterField(
            model_name='database',
            name='account',
            field=models.ForeignKey(related_name=b'databases', to='manager.Account'),
        ),
        migrations.AlterField(
            model_name='database',
            name='name',
            field=models.CharField(help_text=b'Max 50 characters', unique=True, max_length=50, validators=[django.core.validators.RegexValidator(b'^[a-zA-Z][0-9a-zA-Z_]*$', b'Only alphanumeric characters and underscore are allowed.')]),
        ),
        migrations.AlterField(
            model_name='database',
            name='user',
            field=models.CharField(help_text=b'Max 16 characters!', unique=True, max_length=16, validators=[django.core.validators.RegexValidator(b'^[a-zA-Z][0-9a-zA-Z_]*$', b'Only alphanumeric characters and underscore are allowed.')]),
        ),
    ]
