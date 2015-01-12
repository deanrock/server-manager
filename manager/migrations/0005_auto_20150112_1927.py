# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations
import django.core.validators


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0004_auto_20150112_1926'),
    ]

    operations = [
        migrations.AlterField(
            model_name='account',
            name='description',
            field=models.CharField(max_length=1000),
        ),
        migrations.AlterField(
            model_name='app',
            name='name',
            field=models.CharField(max_length=50, validators=[django.core.validators.RegexValidator(b'^[a-zA-Z][0-9a-zA-Z_-]*$', b"Only alphanumeric characters, underscore and '-' are allowed.")]),
        ),
        migrations.AlterField(
            model_name='usersshkey',
            name='ssh_key',
            field=models.TextField(max_length=1000),
        ),
    ]
