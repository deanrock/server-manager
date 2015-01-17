# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0008_auto_20150116_1238'),
    ]

    operations = [
        migrations.AddField(
            model_name='domain',
            name='ssl_enabled',
            field=models.BooleanField(default=False),
            preserve_default=True,
        ),
    ]
