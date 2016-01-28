# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.network "public_network"
  config.vm.provision "shell", path: "vagrant-inst.sh"
  config.vm.provider "virtualbox" do |v|
    v.memory = 2048
  end
  config.vm.provider "vmware_fusion" do |v|
    v.vmx["memsize"] = "2048"
  end
end