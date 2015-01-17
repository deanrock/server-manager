import json
from optparse import make_option
import datetime
import os
from django.conf import settings
from os import walk
from django.core.management.base import BaseCommand
from manager.models import Image, ImageVariable, ImagePort


class Command(BaseCommand):
    option_list = BaseCommand.option_list + (
    )

    help = ''

    def handle(self, **options):
        for var in ImageVariable.objects.all():
            var.delete()

        for port in ImagePort.objects.all():
            port.delete()

        images = {}

        for (dirpath, dirnames, filenames) in walk(settings.IMAGES_PATH):
            for dir in dirnames:
                path = '%s/%s/config.json' % (settings.IMAGES_PATH, dir)
                if os.path.exists(path):
                    with open(path, 'r') as f:
                        j = json.loads(f.read())
                        print j

                        obj =  Image.objects.filter(id=j['id']).first()

                        if not obj:
                            obj = Image()
                            obj.id = j['id']
                            obj.added_at = datetime.datetime.now()

                        obj.description = j['description']
                        obj.name = j['name']
                        obj.type = j['type']

                        obj.save()

                        for port in j['ports']:
                            p = ImagePort()
                            p.type = port['type']
                            p.port = port['port']
                            p.image = obj
                            p.save()

                        if 'variables' in j:
                            for var in j['variables']:
                                v = ImageVariable()
                                v.name = var['name']
                                v.description = var['description']

                                if 'default' in var:
                                    v.default = var['default']
                                v.image = obj
                                v.save()
