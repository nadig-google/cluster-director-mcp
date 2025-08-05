# cluster-director-mcp
Gemini-CLI compatible MCP Server for the Cluster Director product from Google Cloud

# Cluster Director MCP Server

Enable MCP-compatible AI agents to interact with Cluster Director.

# Installation

0. Install gemini-cli on your Cloud Shell Editor
    ```sh
    npm install -g @google/gemini-cli
    ```

    If the above command does not work, install it as root using the following four commands
    ```sh
    sudo -s
    ```

    ```sh
    export PATH=$PATH:/opt/gradle/bin:/opt/maven/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/usr/local/node_packages/node_modules/.bin:/usr/local/rvm/bin:/home/nadig/.gems/bin:/usr/local/rvm/bin:/home/nadig/gopath/bin:/google/gopath/bin:/google/flutter/bin:/usr/local/nvm/versions/node/v22.17.1/bin
    ```

    ```sh
    npm install -g @google/gemini-cli
    ```
    ```sh
    npm install -g @google/gemini-cli to update
    ```
    
    ```sh
    exit
    ```


1.  Install the tool:

    ```sh
    go install github.com/nadig-google/cluster-director-mcp@latest
    ```

    The `cluster-director-mcp` binary will be installed in the directory specified by the `GOBIN` environment variable. If `GOBIN` is not set, it defaults to `$GOPATH/bin` and, if `GOPATH` is also not set, it falls back to `$HOME/go/bin`.

    You can find the exact location by running `go env GOBIN`. If the command returns an empty value, run `go env GOPATH` to find the installation directory.

2.  Install it as a `gemini-cli` extension:

    ```sh
    cd cluster-director-mcp
    go build -o cluster-director-mcp .
    ./cluster-director-mcp install gemini-cli
    ```

    This will create a manifest file in `./.gemini/extensions/cluster-director-mcp` that points to the installed `cluster-director-mcp` binary.

3. Add an MCP extension to get detailed information on cluster director

   if ~/.gemini/settings.json already exists 

   ```sh
   echo '  ,
   "mcpServers": {
     "context7": {
     "httpUrl": "https://mcp.context7.com/mcp"
    }
   }' >> ~/.gemini/settings.json
   ```

   If ~/.gemini/settings.json does NOT exist, then run the following commands
```sh
   mkdir ~/.gemini
   echo '
   {
    "selectedAuthType": "cloud-shell",
    "theme": "Default",
    "mcpServers": {
       "context7": {
         "httpUrl": "https://mcp.context7.com/mcp"
        }
    }
   } ' >> ~/.gemini/settings.json
```


4. Authenticate yourself (run command and follow instructions - this step requires opening a new browser window)
  ```sh
   gcloud auth application-default login
  ```
  
5. Set the default project
  ```sh
  gcloud config set project hpc-toolkit-dev
  ```

6. Start gemini-cli
  ```sh
  gemini
  ```

## Tools

- `create_cluster`: Creates AI optimized Clusters.
- `list_clusters`: List your clusters created using Cluster Director.
- `get_cluster`: Get detailed about a single Cluster.

## Context 

In addition to the tools above, a lot of value is provided through the bundled context instructions.



