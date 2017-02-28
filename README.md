# server-manager

## Requirements

* Debian Jessie
* Python 2.7 installed (for ansible)
* SSH access with user with sudo privileges

## Installation

* Install Debian Jessie
* Install python2.7 (for ansible)
* Set-up SSH access with user (e.g. john) with sudo privileges
* Run ansible deploy.yml playbook
* set mysql root password via mysql_secure_installation command

## First run

Execute as `manager` user:
```bash
./server-manager first-run
```

## Development

### Requirements

* Vagrant
* Ansible

### Installation

```bash
git clone https://github.com/deanrock/server-manager.git
git submodule update --init --recursive
vagrant up
ansible-playbook \
  --private-key=.vagrant/machines/default/virtualbox/private_key \
  -u vagrant -i deployment/dev.hosts deployment/development.yml

# ansible will hang-out on 'update grub and reboot' task;
# you need to stop ansible (ctrl+c), and do:
vagrant halt
vagrant up
ansible-playbook \
  --private-key=.vagrant/machines/default/virtualbox/private_key \
  -u vagrant -i deployment/dev.hosts deployment/development.yml

vagrant ssh

# inside vagrant run:
sudo adduser vagrant docker
sudo adduser vagrant apache
sudo adduser vagrant nginx
sudo adduser vagrant manager
sudo chmod -R 777 /var/log/manager/

# services can then be started via:
./dev.sh ssh
# or
./dev.sh proxy
# or
./dev.sh cron
```

Services
========

a) *proxy* (API + webapp)
- API server
- websocket server
- web SSH access
- serving static files

b) *ssh* (SSH server)
- used for SSH and SFTP access to docker environments

c) *cron* (cronjob daemon)
- used for executing scheduled cron jobs
