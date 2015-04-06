# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0017_auto_20150406_1234'),
    ]

    operations = [
        migrations.AddField(
            model_name='cronjob',
            name='enabled',
            field=models.BooleanField(default=True),
            preserve_default=True,
        ),
    ]
