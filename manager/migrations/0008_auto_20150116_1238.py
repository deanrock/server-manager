# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('manager', '0007_auto_20150116_0951'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='appimagevariable',
            name='image_variable',
        ),
        migrations.AddField(
            model_name='appimagevariable',
            name='name',
            field=models.CharField(default=b'', max_length=255),
            preserve_default=True,
        ),
        migrations.AlterField(
            model_name='imageport',
            name='type',
            field=models.CharField(max_length=255, choices=[(b'fastcgi', b'FastCGI'), (b'http', b'HTTP'), (b'uwsgi', b'UWSGI')]),
        ),
    ]
