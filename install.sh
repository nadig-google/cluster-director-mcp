#!/bin/sh


sudo -s <<EOF
echo "--- Updating package lists ---"
apt update -y

echo "--- Upgrading packages ---"
apt upgrade -y

echo "--- Installing npm ---"

echo "--- Installing gemin-cli ---"
export PATH=$PATH:/opt/gradle/bin:/opt/maven/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/usr/local/node_packages/node_modules/.bin:/usr/local/rvm/bin:/home/nadig/.gems/bin:/usr/local/rvm/bin:/home/nadig/gopath/bin:/google/gopath/bin:/google/flutter/bin:/usr/local/nvm/versions/node/v22.17.1/bin; 
/usr/local/nvm/versions/node/v22.17.1/bin/npm install -g @google/gemini-cli
/usr/local/nvm/versions/node/v22.17.1/bin/npm install -g @google/gemini-cli to update

echo "--- All tasks complete ---"
EOF

#git clone https://github.com/nadig-google/cluster-director-mcp.git

go build -o cluster-director-mcp .
./cluster-director-mcp install gemini-cli
echo '
   {
    "selectedAuthType": "cloud-shell",
    "theme": "Default",
    "mcpServers": {
       "context7": {
         "httpUrl": "https://mcp.context7.com/mcp"
        }
    }
   } ' >> .gemini/extensions/cluster-director-mcp/gemini-extension.json

