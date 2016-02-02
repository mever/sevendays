# About

The main purpose of this repository is to make it easier for players of the
game 7 Days to Die to host there own server. It's secondary objective is to
create a framework in which server admins and developers can configure and
build their ultimate 7 Days to Die experience.

This project uses the shiniest of technologies; [Go](https://golang.org/)
and [Juju](https://jujucharms.com/), a products of Google, Canonical and
their communities, respectively.

## tl;dr

1. Install [Vagrant](https://www.vagrantup.com/) on your machine.
2. Download this repository and place it somewhere.
3. Run 'vagrant up' in that place on your machine.
4. Skip to the 'After Vagrant Up' section.


## Why and How

When I was browsing around looking for resources to build my own server
I saw a lot of binaries and closed code. This project is an attempt to
unite available resources and setup great server software that anyone can
use and contribute to.

This repository contains the source files to compile a Juju charm. When compiled
the charm can be deployed and then it becomes a service. This seems like a hassle
but for future development it is not unthinkable to have it available
in the Juju Charm Store, which makes it very easy to deploy.

I added a Vagrantfile to make it easier for people that just want a server running.
If you are a developer, developing on Ubuntu then you may skip Vagrant and use
the instructions in the vagrant-inst.sh to customise your own environment.
Otherwise to create a development environment or server all you need is Vagrant.


## After Vagrant Up

After running 'vagrant up' the sevendays service is running (which is basically
a web server). It allows you to install and configure the 7 Days to Die server.
Just enter the IP address provided at the end of 'vagrant up' in the address bar
of your browser. Go to the Steam page (click 'Use steam') and click 'install'.
The installation process asks your Steam credentials (and your Steam Guard
code if it's your first install with this service) and it resumes installation.
During installation the page automatically refreshes, if it is done the status
changes and you are ready to connect with your 7 Days to Die client.


### SSL Warning

Do not enter your Steam credentials when you run this service on a public
address. The connection between your web-browser and this service is
not secured with SSL.


## On The Roadmap

- Nicer web UI, with fancy stuff like CSS ;-)
- Home / Status page
- A modding page for XML editing and scripting (Oxide?)
- Server admin page, for interaction with the running game
- SSL