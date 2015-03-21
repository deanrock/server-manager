# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0013_cronjob'),
    ]

    operations = [
        migrations.AlterUniqueTogether(
            name='app',
            unique_together=set([('name', 'account')]),
        ),
    ]
