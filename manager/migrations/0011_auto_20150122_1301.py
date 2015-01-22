# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0010_remove_image_added_by'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='imagevariable',
            name='create_file',
        ),
        migrations.AddField(
            model_name='imagevariable',
            name='filename',
            field=models.CharField(default=False, max_length=255),
            preserve_default=True,
        ),
    ]
