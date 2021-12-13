#!/bin/sh

### input db|project name:
read  -p "Input the short alias of projectname (create name of db and username): " ALIAS
read  -p "!!! Input the Passfrase: " ALIAS_P
echo "Now enter mysq root password(if caching_sha2_password plugin enable):"
##DEBUG IT:
#cat <<_EOF_
mysql --user=root -p <<_EOF_
DROP USER IF EXISTS '${ALIAS}'@'localhost';
DROP USER IF EXISTS '${ALIAS}'@'%';
CREATE USER '${ALIAS}'@'%' IDENTIFIED WITH caching_sha2_password BY '${ALIAS_P}';
GRANT USAGE ON *.* TO '${ALIAS}'@'%';
ALTER USER '${ALIAS}'@'%' REQUIRE NONE WITH MAX_QUERIES_PER_HOUR 0 MAX_CONNECTIONS_PER_HOUR 0 MAX_UPDATES_PER_HOUR 0 MAX_USER_CONNECTIONS 0;
CREATE DATABASE IF NOT EXISTS ${ALIAS};
GRANT ALL PRIVILEGES ON ${ALIAS}.* TO '${ALIAS}'@'%';
FLUSH PRIVILEGES;
_EOF_

echo "Ok. Check it!";

