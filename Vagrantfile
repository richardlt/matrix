# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"

  config.vm.synced_folder "~/go", "/home/vagrant/go"

  config.vm.provider "virtualbox" do |vb|
      vb.memory = 2048
      vb.cpus = 2
      
      vb.customize ["modifyvm", :id, "--usb", "on"]
      vb.customize ['usbfilter', 'add', '0', '--target', :id, '--name',
        'Gamepad', '--vendorid', '0x0079', '--productid', '0x0011']
  end

  config.vm.provision "shell", inline: <<-SHELL
    # install tools
    apt-get update
    apt-get install -y git build-essential pkg-config libusb-1.0

    # install Golang
    export VERSION=1.9.4
    export OS=linux
    export ARCH=amd64
    wget -q https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz
    tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/vagrant/.bash_profile
    echo 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.bash_profile
    echo 'export GOROOT=/usr/local/go' >> /home/vagrant/.bash_profile
  SHELL
end