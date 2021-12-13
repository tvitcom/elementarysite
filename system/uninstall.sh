#!/bin/bash

## Install golang app as systemd instanse
## Preconditions: root, $APPNAME.service
## Alertpoints: STDERR

HOSTNAME=$(hostname)

set -a
test -f ./$HOSTNAME.conf && . ./$HOSTNAME.conf
set +a

# restore common nginx configs
if [ -f /etc/nginx/nginx.conf-original ]; then
    rm /etc/nginx/nginx.conf;
    mv /etc/nginx/nginx.conf-original /etc/nginx/nginx.conf;
    echo "The default file 'nginx.conf' restored successfully in /etc/nginx"
else
    echo "The backup of common file 'nginx.conf' is absent. Abort now"
    exit 100
fi

# restore nginx site configs
if [ -f /etc/nginx/sites-available/default-original ]; then
    unlink /etc/nginx/sites-enabled/$APPNAME;
    rm /etc/nginx/sites-available/$APPNAME;
    mv /etc/nginx/sites-available/default-original /etc/nginx/sites-available/default;
    ln -s /etc/nginx/sites-available/default /etc/nginx/sites-enabled/default
    echo "The default 'default' file restored"
else
    echo "The backup of 'default' nginx file is absent. Abort now"
    exit 100
fi

# start nginx with default configs
systemctl restart nginx

# Refresh application systemd services
systemctl stop $APPNAME
systemctl disable $APPNAME
rm -f /lib/systemd/system/$APPNAME*
systemctl daemon-reload


rm $REMOTE_DIR/$APPNAME
echo "Deleted binary application build: Ok!"
echo "Sytemd Unit of $SERVICEUNIT now uninstalled: Ok!"
exit 0