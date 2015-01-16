# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0006_appimagevariable_create_file'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='appimagevariable',
            name='create_file',
        ),
        migrations.AddField(
            model_name='imagevariable',
            name='create_file',
            field=models.BooleanField(default=False),
            preserve_default=True,
        ),
    ]
