# Cluster Director MCP Server

Use Cluster Director and deploy AI/ML clusters with GPUs using spoken english.

# Installation

1.  Check out code and other assets from github:

    ```sh
    git clone https://github.com/nadig-google/cluster-director-mcp.git
    cd cluster-director-mcp
    ```

2.  Install cluster-director-mcp as a `gemini-cli` extension:

    The dependencies for `cluster-director-mcp` including gemini-cli will be installed on your cloud shell.

    ```sh
    ./install.sh
    ```   

3. Authenticate yourself (run command and follow instructions - this step requires opening a new browser window)
  ```sh
   gcloud auth application-default login
  ```
  
4. Set the default project
  ```sh
  gcloud config set project <project-name>
  ```

5. Start gemini-cli
  ```sh
  gemini
  ```

## Tools

- `list_clusters`: List your clusters created using Cluster Director.
- `get_cluster`: Get detailed about a single Cluster.
- More to come soon....

## Context 

In addition to the above tools, this AI Assistant has additional fine-tuned and detailed information about Cluster Director.



