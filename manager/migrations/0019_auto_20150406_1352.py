# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0018_cronjob_enabled'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='cronjob',
            name='account',
        ),
        migrations.RemoveField(
            model_name='cronjob',
            name='added_by',
        ),
        migrations.DeleteModel(
            name='CronJob',
        ),
    ]
