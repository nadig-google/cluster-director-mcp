# cluster-director-mcp
Gemini-CLI compatible MCP Server for the Cluster Director product from Google Cloud

# Cluster Director MCP Server

Enable MCP-compatible AI agents to interact with Cluster Director.

# Installation

1.  Check out from github:

    ```sh
    git clone https://github.com/nadig-google/cluster-director-mcp.git
    cd cluster-director-mcp
    ```

2.  Install cluster-director-mcp as a `gemini-cli` extension:

    The dependencies for `cluster-director-mcp` including gemini-cli will be installed on your cloud shell.

    Check if you see the cluster-director-mcp directory
    ```sh
    ./install.sh
    ```   

3. Authenticate yourself (run command and follow instructions - this step requires opening a new browser window)
  ```sh
   gcloud auth application-default login
  ```
  
4. Set the default project
  ```sh
  gcloud config set project hpc-toolkit-dev
  ```

5. Start gemini-cli
  ```sh
  gemini
  ```

## Tools

- `create_cluster`: Creates AI optimized Clusters.
- `list_clusters`: List your clusters created using Cluster Director.
- `get_cluster`: Get detailed about a single Cluster.

## Context 

In addition to the tools above, a lot of value is provided through the bundled context instructions.



