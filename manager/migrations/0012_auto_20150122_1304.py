# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0011_auto_20150122_1301'),
    ]

    operations = [
        migrations.AlterField(
            model_name='imagevariable',
            name='filename',
            field=models.CharField(max_length=255, null=True, blank=True),
        ),
    ]
