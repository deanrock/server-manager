# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations
from django.conf import settings


class Migration(migrations.Migration):

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
    ]

    operations = [
        migrations.CreateModel(
            name='Account',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('name', models.CharField(max_length=8)),
                ('description', models.CharField(max_length=1000)),
                ('added_at', models.DateTimeField(auto_now_add=True, null=True)),
                ('added_by', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
                ('users', models.ManyToManyField(related_name=b'users', to=settings.AUTH_USER_MODEL)),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='App',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('name', models.TextField(max_length=50)),
                ('container_id', models.TextField(max_length=255)),
                ('memory', models.IntegerField(default=256)),
                ('added_at', models.DateTimeField(auto_now_add=True, null=True)),
                ('account', models.ForeignKey(to='manager.Account')),
                ('added_by', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='AppImageVariable',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('value', models.TextField()),
                ('app', models.ForeignKey(related_name=b'variables', to='manager.App')),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Database',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('type', models.TextField(verbose_name=50)),
                ('name', models.TextField(verbose_name=50)),
                ('user', models.TextField(verbose_name=50)),
                ('password', models.TextField(verbose_name=255)),
                ('added_at', models.DateTimeField(auto_now_add=True, null=True)),
                ('account', models.ForeignKey(to='manager.Account')),
                ('added_by', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Domain',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('name', models.CharField(max_length=255)),
                ('redirect_url', models.CharField(default=None, max_length=255)),
                ('nginx_config', models.TextField()),
                ('apache_config', models.TextField()),
                ('apache_enabled', models.BooleanField(default=False)),
                ('added_at', models.DateTimeField(auto_now_add=True, null=True)),
                ('account', models.ForeignKey(to='manager.Account')),
                ('added_by', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Image',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('name', models.TextField(max_length=255)),
                ('description', models.TextField()),
                ('type', models.CharField(max_length=255, choices=[(b'application', b'Application')])),
                ('added_at', models.DateTimeField(auto_now_add=True, null=True)),
                ('added_by', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='ImagePort',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('type', models.CharField(max_length=255, choices=[(b'fastcgi', b'FastCGI'), (b'http', b'HTTP')])),
                ('port', models.IntegerField()),
                ('image', models.ForeignKey(related_name=b'ports', to='manager.Image')),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='ImageVariable',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('name', models.CharField(max_length=255)),
                ('description', models.TextField()),
                ('default', models.TextField()),
                ('image', models.ForeignKey(related_name=b'variables', to='manager.Image')),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.AddField(
            model_name='appimagevariable',
            name='image_variable',
            field=models.ForeignKey(to='manager.ImageVariable'),
            preserve_default=True,
        ),
        migrations.AddField(
            model_name='app',
            name='image',
            field=models.ForeignKey(to='manager.Image'),
            preserve_default=True,
        ),
    ]
