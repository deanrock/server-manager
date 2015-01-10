# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0001_initial'),
    ]

    operations = [
        migrations.AlterField(
            model_name='account',
            name='name',
            field=models.CharField(help_text=b'max 32 chars!', max_length=32),
        ),
        migrations.AlterField(
            model_name='app',
            name='account',
            field=models.ForeignKey(related_name=b'apps', to='manager.Account'),
        ),
        migrations.AlterField(
            model_name='app',
            name='container_id',
            field=models.TextField(max_length=255, null=True, blank=True),
        ),
        migrations.AlterField(
            model_name='database',
            name='name',
            field=models.CharField(help_text=b'Max 50 characters', unique=True, max_length=50),
        ),
        migrations.AlterField(
            model_name='database',
            name='password',
            field=models.CharField(max_length=255),
        ),
        migrations.AlterField(
            model_name='database',
            name='type',
            field=models.CharField(max_length=50, choices=[(b'mysql', b'MySQL')]),
        ),
        migrations.AlterField(
            model_name='database',
            name='user',
            field=models.CharField(help_text=b'Max 16 characters!', unique=True, max_length=16),
        ),
        migrations.AlterField(
            model_name='domain',
            name='account',
            field=models.ForeignKey(related_name=b'domains', to='manager.Account'),
        ),
        migrations.AlterField(
            model_name='domain',
            name='apache_config',
            field=models.TextField(blank=True),
        ),
        migrations.AlterField(
            model_name='domain',
            name='nginx_config',
            field=models.TextField(default=None, blank=True),
        ),
        migrations.AlterField(
            model_name='domain',
            name='redirect_url',
            field=models.CharField(default=None, max_length=255, blank=True),
        ),
        migrations.AlterField(
            model_name='imagevariable',
            name='default',
            field=models.TextField(null=True, blank=True),
        ),
    ]
