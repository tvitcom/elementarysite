#!/bin/sh

## Install golang app as systemd instanse
## Preconditions: root, remote-root, Makefile (build32 process)
## Alertpoints: STDERR

## USING: sh system/deploy.sh instance-1 OR: sh deploy.sh micro
HOSTNAME=$1
APPROOT=$(pwd)

test -n $HOSTNAME || exit 100

set -a
test -f $PWD/system/$HOSTNAME.conf && . $PWD/system/$HOSTNAME.conf
set +a

## Prepare build directory
#env GOOS=linux GOARCH=386 go build -o $PWD/build/i686/$APPNAME
#env GOOS=linux GOARCH=amd64 go build -o $PWD/build/x86_64/$APPNAME
make build32

# Sync process:
rsync -e "ssh -p $REMOTE_PORT" \
	--exclude=".gitignore" \
	--exclude="filestore/" \
	--exclude="*.swp"	\
	--exclude=".env*"	\
	--exclude="info/*"	\
	--exclude="run.sh"	\
	--exclude="logs/access.log" \
	--exclude="logs/server.log" \
	--exclude="logs/errors.log" \
	--exclude="123*" \
	--exclude=".git" \
	-PLSluvr --del --no-perms --no-t \
	$APPROOT"/" $REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR"/"

#rsync -e "ssh -p $REMOTE_PORT" \
#	--exclude="123*" \
#	-PLSluvr --no-perms --no-t \
#	$APPROOT"/filestore/" $REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR"/filestorage/"
## clean

echo "Transfered to webserver code-files $REMOTE_HOST:$REMOTE_PORT: Ok!"

# ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "rm $REMOTE_DIR/$APPNAME"
# echo "Deleted old build: Ok!"

#ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "cp -f $REMOTE_DIR/build/$(arch)/$APPNAME $REMOTE_DIR"
#echo "New buil created: Ok!"

#ssh -p $REMOTE_PORT root@$REMOTE_HOST "sudo systemctl try-restart $SERVICEUNIT"

echo "!!!ALSO: You will run projdir/system/install_elementarysite.sh on the $REMOTE_HOST!!!"
echo "!!!  OR: projdir/system/reloadcode.sh on the $REMOTE_HOST!!!"
exit 100
