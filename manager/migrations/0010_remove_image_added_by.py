# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0009_domain_ssl_enabled'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='image',
            name='added_by',
        ),
    ]
