server-manager
==============


Server requirements
===================

* Debian Wheezy or Jessie
* Python 2.7 installed
* SSH access with user with sudo privileges

Installation
============

* Install Debian Wheezy or Jessie
* Install python2.7
* Set-up SSH access with user (e.g. john) with sudo privileges
* Run ansible deploy.yml playbook
* Add user `manager` as sudo without password
	
	#add:
	manager ALL=(ALL) NOPASSWD:ALL
	#to /etc/sudoers

* Install libssl-dev (for now, until we fix pull.sh script)
* set mysql root password via mysql_secure_installation command
* log-in as `manager` and clone repo to /home/manager/server-manager/
* clone git submodules

	git submodule update --init --recursive

* create /home/manager/server-manager/manager/settings/production.py with the contents:
	
	from manager.settings.base import *

	# SECURITY WARNING: keep the secret key used in production secret!
	SECRET_KEY = 'fewoigh5940659j0--0i-0i34y90u6(H$*Y0i0u%uu'

	# SECURITY WARNING: don't run with debug turned on in production!
	DEBUG = True

	TEMPLATE_DEBUG = True

	MYSQL_ROOT_PASSWORD = 'root_password'

* cd to server-manager/ folder and run `./pull.sh production`

Create first user
=================
Execute as `manager` user:

	$ cd ~/server-manager/
	$ source ./env/bin/activate
	$ python manage.py createsuperuser --settings=manager.settings.production


Development env setup
=====================

1. clone git
2. ssh to VM

	$ vagrant ssh

3. copy your public key to ~/files/.ssh/authorized_keys on vagrant
4. go to deployment/ subdirectory and run ansible playbook
	
	$ cd deployment/
	$ source ~/virtualenv/ansible/bin/activate #or wherever you have ansible env
	$ ansible-playbook -i dev.hosts deploy.yml

5. ansible probably won't detect when VM reboots after kernel install and will stall at "docker | kernel - wait for reboot"; after waiting a minute or so it's safe to retry the playbook
6. after ansible finishes restarting the VM (you NEED to do this via vagrant halt/vagrant up, otherwise VBox extension won't be reinstalled)

	$ vagrant halt
	$ vagrant up
	$ vagrant ssh

7. add vagrant user to docker, nginx and apache group

	(vagrant)$ sudo adduser vagrant docker
	(vagrant)$ sudo adduser vagrant apache
	(vagrant)$ sudo adduser vagrant nginx
	(vagrant)$ sudo adduser vagrant manager
	(vagrant)$ exit
	$ vagrant ssh

8. compile Go programs, migrate database, install Python requirements ...
	
	(vagrant)$ cd files/
	(vagrant)$ ./pull.sh dev

9. workaround because we are not using "manager" user:

	(vagrant)$ sudo chmod -R 777 /var/log/manager/

10. change mysql password (set to 'password' for dev) and remove test data:
	
	(vagrant)$ mysql_secure_installation

11. create symbolic link to images/ folder
	
	(vagrant)$ sudo mkdir /home/manager/server-manager/
	(vagrant)$ sudo ln -s /home/vagrant/files/images /home/manager/server-manager/images


Your dev env should be kinda ready.


Services
========

a) *manager* (Django backend)
- exposes API for adding/modifying/deleting models (e.g. domains, apps, ...)
- manages database migrations


b) *proxy* (Golang backend)
- websocket server
- web SSH access
- serving static files
- proxies other request to manager

c) *ssh* (SSH server)
- used for SSH and SFTP access to docker environments


d) *cron* (cronjob daemon)
- used for executing scheduled cron jobs





How to generate SSL key?
========================

	$ cd . #to server-manager root (e.g. /home/manager/server-manager/)
	$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ssl.key -out ssl.crt

