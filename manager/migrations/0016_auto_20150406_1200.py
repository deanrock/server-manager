# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0015_auto_20150401_1806'),
    ]

    operations = [
        migrations.AlterField(
            model_name='cronjob',
            name='image',
            field=models.CharField(max_length=255),
        ),
    ]
