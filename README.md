# server-manager

## Server requirements

* Debian Wheezy or Jessie
* Python 2.7 installed
* SSH access with user with sudo privileges

## Installation

* Install Debian Jessie
* Install python2.7
* Set-up SSH access with user (e.g. john) with sudo privileges
* Run ansible deploy.yml playbook
* Install libssl-dev (for now, until we fix pull.sh script)
```bash
	sudo apt-get install libssl-dev -t jessie-backports
```
* set mysql root password via mysql_secure_installation command
* log-in as `manager` and clone repo to /home/manager/server-manager/
* clone git submodules
```bash
	cd server-manager/
	git submodule update --init --recursive
```

* create `/home/manager/server-manager/manager/settings/production.py` with the contents:
```python
from manager.settings.base import *

# SECURITY WARNING: keep the secret key used in production secret!
SECRET_KEY = 'fewoigh5940659j0--0i-0i34y90u6(H$*Y0i0u%uu'

# SECURITY WARNING: don't run with debug turned on in production!
DEBUG = True

TEMPLATE_DEBUG = True

MYSQL_ROOT_PASSWORD = 'root_password'
```

* cd to `server-manager/`` folder and run `./pull.sh production`

## Create first user

Execute as `manager` user:
```bash
cd ~/server-manager/
source ./env/bin/activate
python manage.py createsuperuser --settings=manager.settings.production
```

## Development env setup

1. install ansible, vagrant, clone git and its submodules
2. ssh to VM
```bash
vagrant ssh
```
3. copy your public key to ~/.ssh/authorized_keys on vagrant
4. go to deployment/ subdirectory and run ansible playbook
```bash
cd deployment/
source ~/virtualenv/ansible/bin/activate #or wherever you have ansible env
ansible-playbook -i dev.hosts deploy.yml
```

5. ansible probably won't detect when VM reboots after kernel install and will stall at "docker | kernel - wait for reboot"; after waiting a minute or so it's safe to retry the playbook
6. after ansible finishes restarting the VM (you NEED to do this via vagrant halt/vagrant up, otherwise VBox extension won't be reinstalled)
```bash
vagrant halt
vagrant up
vagrant ssh
```

7. add vagrant user to docker, nginx and apache group
```bash
(vagrant)$ sudo adduser vagrant docker
(vagrant)$ sudo adduser vagrant apache
(vagrant)$ sudo adduser vagrant nginx
(vagrant)$ sudo adduser vagrant manager
(vagrant)$ exit
vagrant ssh
```

8. compile Go programs, migrate database, install Python requirements ...
```bash
(vagrant)$ cd files/
(vagrant)$ ./pull.sh dev
```
9. workaround because we are not using "manager" user:
```bash
sudo chmod -R 777 /var/log/manager/
```

10. change mysql password (set to 'password' for dev) and remove test data:
```bash
mysql_secure_installation
```

11. create symbolic link to images/ folder
```bash
sudo mkdir /home/manager/server-manager/
sudo ln -s /home/vagrant/files/images /home/manager/server-manager/images
```

12. create first user
```bash
source env/bin/activate
python manage.py createsuperuser --settings=manager.settings.dev --noreload
```

13. install `screen` via `sudo apt-get install screen`, and run each app in different screen; you need to start the following apps:
```bash
./dev.sh manager
./dev.sh ssh
./dev.sh cron
./dev.sh proxy
```

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
