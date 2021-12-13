
First sync
---
0. Go to the code:
cd projectnamedir;

1. You will edit configs files in the system dir:
./deploy.sh [confname in the system dir]

2. and go to the remoteserver and run:
cd projectnamedir;
./install_projectname.sh

3. then create and load the database;   

Deployment
---
0. Go to the code:
cd projectnamedir;

1. You will edit configs files in the system dir:
./deploy.sh [conf_for_hostname_file in the system dir]

2. and go to the remoteserver and run:
cd projectnamedir;
./system/reloadcode.sh

3. then empty database with ./system/cleardb.sh and load actual database;   