# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0016_auto_20150406_1200'),
    ]

    operations = [
        migrations.AddField(
            model_name='cronjob',
            name='last_executed',
            field=models.DateTimeField(null=True),
            preserve_default=True,
        ),
        migrations.AddField(
            model_name='cronjob',
            name='last_execution_status',
            field=models.BooleanField(default=False),
            preserve_default=True,
        ),
    ]
