# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations
from django.conf import settings


class Migration(migrations.Migration):

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
        ('manager', '0012_auto_20150122_1304'),
    ]

    operations = [
        migrations.CreateModel(
            name='CronJob',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('name', models.CharField(max_length=255)),
                ('directory', models.CharField(max_length=255)),
                ('script_file', models.CharField(max_length=255)),
                ('timeout', models.IntegerField(default=60)),
                ('cron_expression', models.CharField(max_length=255)),
                ('added_at', models.DateTimeField(auto_now_add=True, null=True)),
                ('account', models.ForeignKey(related_name=b'cronjobs', to='manager.Account')),
                ('added_by', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
                ('image', models.ForeignKey(to='manager.Image')),
            ],
            options={
            },
            bases=(models.Model,),
        ),
    ]
