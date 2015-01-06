from subprocess import Popen, PIPE
from django.conf import settings

from manager import docker_api
from manager.models import Account
from manager.utils import exec_command, Logs
import os


def sync_accounts():
    logs = Logs()

    #get all accounts
    accounts = Account.objects.all()

    for account in accounts:
        logs.add("-----------------------")
        logs.add("Account %s" % account)

        o, e = exec_command(logs, "sudo adduser --disabled-password --gecos \"\" %s" % account.name)

        o, e = exec_command(logs, "sudo chmod 750 /home/%s" % account.name)

        dirs = [
            "apps"
        ]

        for dir in dirs:
            o, e = exec_command(logs, "sudo mkdir -p /home/%s/%s" % (account.name, dir))

            o, e = exec_command(logs, "sudo chmod 750 /home/%s/%s" % (account.name, dir))



    return logs


def rebuild_base_image():
    logs = Logs()

    folder = os.path.abspath(os.path.join(settings.BASE_DIR, '../', 'images/', 'debian7basehosting'))

    #build image
    for line in docker_api.cli.build(
       path=folder, rm=True, tag='manager/debian7basehosting'
    ):
        logs.add(line)

    return logs