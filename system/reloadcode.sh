#!/bin/sh

# remove old binary and copy new binary from dist then reload own service
PROJ_NAME="elementarylearn"
DISTR_DIR="/build/i686/";
PROJ_DIR="/home/user/Go/src/github.com/tvitcom/";
rm -f ${PROJ_DIR}${PROJ_NAME}/${PROJ_NAME};
cp -f ${PROJ_DIR}${PROJ_NAME}${DISTR_DIR}${PROJ_NAME} ${PROJ_DIR}${PROJ_NAME}/${PROJ_NAME};
systemctl restart ${PROJ_NAME};
systemctl status ${PROJ_NAME};
