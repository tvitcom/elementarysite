#!/bin/bash

## Install golang app as systemd instanse
## Preconditions: root, pwd: [projname]
## Alertpoints: STDERR
## RUN FORM: approot NOT from system dir
## RUN: ./system/install_unixlean.sh
apt-get update
apt-get -y install jpegoptim

APPNAME="elementarylearn"
GO_USER="user"

# Deploy nginx configs
if [ -f /etc/nginx/nginx.conf-original ]; then
    cp -f ./system/nginx.conf /etc/nginx/nginx.conf
    echo "The special file 'nginx.conf' placed successfully in /etc/nginx"
else
    mv /etc/nginx/nginx.conf /etc/nginx/nginx.conf-original
    cp -f nginx.conf /etc/nginx/nginx.conf
    echo "The default file 'nginx.conf' backuped in /etc/nginx with -original suffix"
fi

# Backup default nginx and default-site configs
if [ -f /etc/nginx/sites-available/default-original ] || [ -f /etc/nginx/sites-available/default.conf-original ]; then
	echo "file default config backup file is available"
else
    unlink /etc/nginx/sites-enabled/default*;
    mv /etc/nginx/sites-available/default* /etc/nginx/sites-available/default-original;
    echo "The default 'default' file buckuped to -original suffix"
fi

cp -f system/$APPNAME".conf" /etc/nginx/sites-available/$APPNAME;
chmod 644 /etc/nginx/sites-available/$APPNAME;
ln -s /etc/nginx/sites-available/$APPNAME /etc/nginx/sites-enabled/$APPNAME
echo "Sites nginx config is installed"

# start nginx with new configs
systemctl restart nginx

# Env file setting:
if [ -f ./.env ]; then
    echo "Env file is available. Sites .env params did installed."
else
    cp env.* ./.env;
    chmod 400 ./.env;
    chown $GO_USER:$GO_USER ./.env;
    echo "Env file installed now: Ok!"
fi

cp -f ./build/$(arch)/$APPNAME ./$APPNAME;
if [ -f ./$APPNAME ]; then
    echo "New build placed: Ok!"
else
    echo "New build is not placed: FAIL."
    exit 1
fi

if [ -f /lib/systemd/system/$APPNAME.service ]; then
    echo "Systemd has previous $APPNAME.service"
    file /lib/systemd/system/$APPNAME.service
else
    cp -f ./system/$APPNAME".service" /lib/systemd/system
    chmod 644 /lib/systemd/system/$APPNAME".service"
    systemctl enable $APPNAME
    systemctl start $APPNAME
    systemctl status $APPNAME
    echo "Sytemd Unit of $SERVICEUNIT now installed: Ok!"
fi

echo "Sytemd Unit of $SERVICEUNIT now installed: Ok!"
echo "You will:"
echo "1. load appropriate database"
exit 0

