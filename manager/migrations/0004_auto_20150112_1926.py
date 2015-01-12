# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0003_auto_20150112_1911'),
    ]

    operations = [
        migrations.AlterField(
            model_name='account',
            name='description',
            field=models.TextField(max_length=1000),
        ),
    ]
