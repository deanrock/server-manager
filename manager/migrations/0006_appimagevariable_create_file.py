# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0005_auto_20150112_1927'),
    ]

    operations = [
        migrations.AddField(
            model_name='appimagevariable',
            name='create_file',
            field=models.BooleanField(default=False),
            preserve_default=True,
        ),
    ]
