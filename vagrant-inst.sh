#!/bin/bash
# This script assumes a clean install of Ubuntu 14.04,
# it installs all stuff necessary to develop and deploy
# a Juju charm containing the 7 Days to Die server.
# This script assumes root.
set -u

echo "Install updates, Git and Juju..."
apt-get update > /dev/null
apt-get upgrade  > /dev/null
add-apt-repository ppa:juju/stable
apt-get install -y git juju-core juju-local  > /dev/null

echo "Download and install Go."
wget -q https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.5.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo -e '\n\nexport PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Setup workspace in home directory, choose the vagrant user
# as this is the default user Vagrant uses to login with.
echo "Setup workspace..."
homeDir=/home/vagrant
mkdir $homeDir/go $homeDir/charms
chown vagrant:vagrant $homeDir/go $homeDir/charms
cat >> $homeDir/.bashrc <<'EOF'

export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export JUJU_REPOSITORY=$HOME/charms
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

EOF

# Juju refuses to run as root (which we are now). Therefore
# I continue as the vagrant user. It also allows us to go on
# directly after logging in with 'vagrant ssh'
sudo -u vagrant -i <<'JUJU_SH'

# configure Juju
juju generate-config
juju switch local
juju bootstrap

echo -n "Download and compile charm..."

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

# Wait until we get the ip from the sevendays unit.
echo -n "Deploying the sevendays service... (this takes a while)"
while [ -z "$(juju status sevendays/0 | grep public-address | grep -Po '[0-9]{1,3}')" ]; do sleep 10; done

JUJU_SH

# Vagrant uses eth0 as it's primary interface and
# eth1 as bridged interface. This may be changed
# if that assumption is false.
BRIDGE_IF=eth1

# Let's configure the Linux firewall to allow all
# traffic on the bridge to flow right into our Juju unit
echo -e "#!/bin/bash\n\nBRIDGE_IF=${BRIDGE_IF}\n\n" > /etc/network/if-up.d/iptables
cat >> /etc/network/if-up.d/iptables <<'IPCONF'
PATH=/sbin:/bin:/usr/sbin:/usr/bin
if [ "$IFACE" == "$BRIDGE_IF" ]; then
    export HOME=/home/vagrant
    IP=$(ifconfig "$IFACE" | /usr/bin/awk '/inet addr/{print substr($2,6)}')
    UNIT_IP=$(juju status sevendays/0 | grep public-address | grep -Po '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}')
    if [ ! -z "$UNIT_IP" ]; then
        iptables -t nat -C PREROUTING -d "$IP"/32 -j DNAT --to-destination "$UNIT_IP"
        if [ "$?" == 1 ]; then
            iptables -t nat -A PREROUTING -d "$IP"/32 -j DNAT --to-destination "$UNIT_IP"
        fi
    fi
fi

IPCONF

chmod +x /etc/network/if-up.d/iptables

echo -n "Restarting the network with new firewall configuration..."
ifdown -a 2> /dev/null; ifup -a 2> /dev/null

BRIDGE_IP=$(ifconfig "$BRIDGE_IF" | /usr/bin/awk '/inet addr/{print substr($2,6)}')
echo "|"
echo "|  All done. Enter '${BRIDGE_IP}' in your web browser!"
echo "|"