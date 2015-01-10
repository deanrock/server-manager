from subprocess import Popen, PIPE
import MySQLdb
from django.conf import settings
from django.template.loader import render_to_string

from manager import docker_api
from manager.models import Account, Domain, Database
from manager.utils import exec_command, Logs
import os
import re


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
            "logs"
        ]

        o, e = exec_command(logs, "sudo adduser nginx %s" % account.name)
        o, e = exec_command(logs, "sudo adduser apache %s" % account.name)

        for dir in dirs:
            o, e = exec_command(logs, "sudo mkdir -p /home/%s/%s" % (account.name, dir))

            o, e = exec_command(logs, "sudo chmod 750 /home/%s/%s" % (account.name, dir))
            o, e = exec_command(logs, "sudo chown %s:%s /home/%s/%s" % (account.name, account.name, account.name, dir))



    return logs


def rebuild_base_image():
    logs = Logs()

    images = [
        "debian7base",
        "debian7basehosting",
        "php56-base-hosting",
        "shell-base",
        "shell-php56",
        "python34-base-hosting",
    ]

    for image in images:
        folder = os.path.abspath(os.path.join(settings.BASE_DIR, '../', 'images/', image))

        #build image
        for line in docker_api.cli.build(
           path=folder, rm=True, tag='manager/%s' % image
        ):
            logs.add(line)

    return logs

def update_nginx_config():
    logs = Logs()

    #docker to ip mapping
    docker_mapping={}

    containers = docker_api.cli.containers()

    for c in containers:
        name = None
        try:
            name = c['Names'][0].strip('/')
        except:
            continue

        ip = None
        try:
            ip = docker_api.cli.inspect_container(name)['NetworkSettings']['IPAddress']
        except:
            continue

        docker_mapping[name] = ip



    #remove existing files
    o, e = exec_command(logs, "sudo rm /etc/nginx/manager/*")
    o, e = exec_command(logs, "sudo rm /etc/apache2/manager/*")

    for account in Account.objects.all():
        domains = account.domains.all()
        apps = account.apps.all()

        port_mapping={}

        for app in apps:
            for port in app.image.ports.all():
                ip = '255.255.255.255'

                if app.container_name() in docker_mapping:
                    ip = docker_mapping[app.container_name()]

                port_mapping['%s_%d_ip' % (app.name, port.port)] = ip

        logs.add(port_mapping)

        for domain in domains:
            logs.add("domain %s, account %s" % (domain.name, domain.account.name))

            nginx = ''
            apache = domain.apache_config

            if domain.redirect_url:
                #redirect to url

                nginx = render_to_string('system/nginx_redirect.conf', {
                    'url': domain.redirect_url
                })
            else:
                #parse nginx_config
                nginx = domain.nginx_config

                for port in port_mapping:
                    nginx = nginx.replace('#%s#' % port, port_mapping[port])

                    if domain.apache_enabled:
                        apache = apache.replace('#%s#' % port, port_mapping[port])

                nginx = nginx.replace('#apache#', '127.0.0.1:8080')

            logs.add(nginx)

            conf = render_to_string('system/nginx_vhost.conf', {
                'domain': domain.name,
                'config': nginx
            })

            fname = os.path.join('/etc/nginx/manager', '%s_%s.conf' % (account.name, domain.name))

            with open(fname, 'w') as f:
                f.write(conf)

            #test config
            o, e = exec_command(logs, "sudo nginx -t")

            if 'emerg' in e:
                logs.add('NGINX CONFIG TEST FAILED!')
                return logs


            #apache
            if domain.apache_enabled:
                conf = render_to_string('system/apache_vhost.conf', {
                    'domain': domain.name,
                    'config': apache
                })

                #find <php></php>
                m = re.search('<php>(.*)<\/php>', conf)
                for x in m.groups():
                    name = domain.name.replace('.', '-')
                    txt = """AddType application/x-httpd-fastphp5 .php
            Action application/x-httpd-fastphp5 /php5-fcgi
            Alias /php5-fcgi /usr/lib/cgi-bin/php5-fcgi-%s
            FastCgiExternalServer /usr/lib/cgi-bin/php5-fcgi-%s -host %s:9000 -pass-header Authorization""" % (name, name, x)
                    conf = conf.replace('<php>%s</php>' % x, txt)

                logs.add(conf)

                fname = os.path.join('/etc/apache2/manager', '%s_%s.conf' % (account.name, domain.name))

                with open(fname, 'w') as f:
                    f.write(conf)

                #test config
                o, e = exec_command(logs, "sudo apachectl configtest")

            if 'failed' in o:
                logs.add('APACHE CONFIG TEST FAILED!')
                return logs


    #reload nginx
    o, e = exec_command(logs, "sudo service nginx reload")

    #reload apache2
    o, e = exec_command(logs, "sudo service apache2 reload")

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
