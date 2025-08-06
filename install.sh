#!/bin/sh

sudo -s <<EOF


echo "--- Upgrading packages - takes 30mins the first time ---"
apt update -y
apt upgrade -y

echo "--- Installing gemin-cli ---"
export PATH=$PATH:/opt/gradle/bin:/opt/maven/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/usr/local/node_packages/node_modules/.bin:/usr/local/rvm/bin:/home/nadig/.gems/bin:/usr/local/rvm/bin:/home/nadig/gopath/bin:/google/gopath/bin:/google/flutter/bin:/usr/local/nvm/versions/node/v22.17.1/bin; 
/usr/local/nvm/versions/node/v22.17.1/bin/npm install -g @google/gemini-cli
/usr/local/nvm/versions/node/v22.17.1/bin/npm install -g @google/gemini-cli to update

echo "--- All tasks complete ---"
EOF

make clean
make
