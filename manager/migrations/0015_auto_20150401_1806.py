# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0014_auto_20150309_2009'),
    ]

    operations = [
        migrations.RenameField(
            model_name='cronjob',
            old_name='script_file',
            new_name='command',
        ),
        migrations.AlterField(
            model_name='image',
            name='type',
            field=models.CharField(max_length=255, choices=[(b'application', b'Application'), (b'database', b'Database')]),
        ),
        migrations.AlterField(
            model_name='imageport',
            name='type',
            field=models.CharField(max_length=255, choices=[(b'fastcgi', b'FastCGI'), (b'http', b'HTTP'), (b'uwsgi', b'UWSGI'), (b'tcp', b'TCP')]),
        ),
    ]
