# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure('2') do |config|
  config.vm.box = "debian8.1-opscode"
  config.vm.box_url = "http://opscode-vm-bento.s3.amazonaws.com/vagrant/virtualbox/opscode_debian-8.1_chef-provisionerless.box"
  config.vm.hostname = 'manager'
  
  config.vm.network "private_network", ip: "192.168.50.144"
  
  config.vm.synced_folder ".", "/home/vagrant/files"

  config.vm.provider :virtualbox do |vb|
      vb.customize ["modifyvm", :id, "--memory", "2048"]
      vb.customize ["modifyvm", :id, "--cpus", "1"]
      vb.customize ["modifyvm", :id, "--ioapic", "on"]
	  
      vb.customize ["setextradata", :id, "VBoxInternal2/SharedFoldersEnableSymlinksCreate/cross-compiler", "1"]
  end
end
