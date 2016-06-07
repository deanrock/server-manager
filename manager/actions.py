from subprocess import Popen, PIPE
import MySQLdb
from django.conf import settings
from django.contrib.auth.models import User
from django.template.loader import render_to_string
from django.db import connection
from manager import docker_api
from manager.models import Account, Domain, Database
from manager.utils import exec_command, Logs
import os
import re
import utils


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
            "apps",
            "domains",
            ".ssh"
        ]

        o, e = exec_command(logs, "sudo adduser nginx %s" % account.name)
        o, e = exec_command(logs, "sudo adduser apache %s" % account.name)

        for dir in dirs:
            o, e = exec_command(logs, "sudo mkdir -p /home/%s/%s" % (account.name, dir))

            o, e = exec_command(logs, "sudo chmod 750 /home/%s/%s" % (account.name, dir))
            o, e = exec_command(logs, "sudo chown %s:%s /home/%s/%s" % (account.name, account.name, account.name, dir))

        # remove authorized_keys file & force nologin as shell
        o, e = exec_command(logs, "sudo rm /home/%s/.ssh/authorized_keys" % (account.name))
        o, e = exec_command(logs, "sudo chsh -s /usr/sbin/nologin %s" % (account.name))


    return logs


def sync_databases():
    logs = Logs()

    databases = Database.objects.all()

    db = MySQLdb.connect("localhost", "root", settings.MYSQL_ROOT_PASSWORD)

    cursor = db.cursor()




    for db in databases:
        logs.add("database %s" % db.name)

        try:
            logs.add(cursor.execute("""CREATE DATABASE `%s` CHARACTER SET utf8 COLLATE utf8_general_ci""" % (db.name)))
        except:
            logs.add('database already exists')

        logs.add(cursor.execute("""GRANT ALL ON `%s`.* TO `%s`@'%%' IDENTIFIED BY '%s'""" % (db.name, db.user, db.password)))

    logs.add(cursor.execute("""FLUSH PRIVILEGES;"""))

    return logs
