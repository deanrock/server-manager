import json
from optparse import make_option
import datetime
import os
from django.conf import settings
from os import walk
from django.core.management.base import BaseCommand
from manager.models import Image, ImageVariable, ImagePort
from manager import actions

class Command(BaseCommand):
    option_list = BaseCommand.option_list + (
    )

    help = ''

    def handle(self, **options):
        logs = actions.update_nginx_config()

        print(json.dumps(logs.logs))
