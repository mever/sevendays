#!/bin/bash
# This script assumes a clean install of Ubuntu 14.04,
# it installs all stuff necessary to develop and deploy
# a Juju charm containing the 7 Days to Die server.
# This script assumes root.
set -ux

# Install updates, Git and Juju.
apt-get update
apt-get upgrade
add-apt-repository ppa:juju/stable
apt-get install -y git juju-core juju-local

# Download and install Go.
wget -q https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.5.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo -e '\n\nexport PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Setup workspace in home directory, choose the vagrant user
# as this is the default user Vagrant uses to login with.
homeDir=/home/vagrant
mkdir $homeDir/go $homeDir/charms
chown vagrant:vagrant $homeDir/go $homeDir/charms
cat >> $homeDir/.bashrc <<'EOF'

export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export JUJU_REPOSITORY=$HOME/charms
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

EOF

# let Juju 'think' we're running as normal user
sudo -u vagrant -i <<'JUJU_SH'

# configure Juju
juju generate-config
juju switch local
juju bootstrap

# setup environment
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export JUJU_REPOSITORY=$HOME/charms
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# build and deploy Juju Charm
go get github.com/juju/gocharm/cmd/gocharm
go get github.com/mever/sevendays/charms/sevendays
gocharm github.com/mever/sevendays/charms/sevendays
juju deploy local:trusty/sevendays

JUJU_SH