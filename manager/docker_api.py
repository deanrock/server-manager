from docker import Client
from io import BytesIO
import tarfile
import os
import re
from django.conf import settings

cli = Client(base_url='unix://var/run/docker.sock')

def is_integer(s):
    try:
        float(s)
        return True
    except ValueError:
        return False

def matches_filename(strg, search=re.compile(r'[^a-z0-9.]').search):
    return not bool(search(strg))

def containers():
    return cli.containers()
