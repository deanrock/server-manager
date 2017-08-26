# server-manager

## Requirements

* Debian Jessie
* Python 2.7 installed (for ansible)
* SSH access with user with sudo privileges

## Installation

* Install Debian Jessie
* Install python2.7 (for ansible)
* Set-up SSH access with user (e.g. john) with passwordless sudo privileges
* Run ansible deploy.yml playbook
* SSH to the server and run:
```bash
sudo su
cd /home/manager/sm/

# Create initial database
./server-manager first-run

# Create admin user
./server-manager create-admin-user user password
```
* set mysql root password via mysql_secure_installation command

## Development

### Requirements

* Vagrant
* Ansible

### Installation

```bash
git clone https://github.com/deanrock/server-manager.git
cd server-manager/
git submodule update --init --recursive
vagrant up
ansible-playbook \
  --private-key=.vagrant/machines/default/virtualbox/private_key \
  -u vagrant -i deployment/dev.hosts deployment/development.yml

# VM should time-out after installing backported kernel; re-run ansible playbook
ansible-playbook \
  --private-key=.vagrant/machines/default/virtualbox/private_key \
  -u vagrant -i deployment/dev.hosts deployment/development.yml

vagrant ssh
cd ~/files/
./dev.sh first-run
./dev.sh create-admin-user user password

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
