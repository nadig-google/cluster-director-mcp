# Cluster Director MCP Server

Use Cluster Director and deploy AI/ML clusters with GPUs using spoken english.

# Installation

1.  Check out code and other assets from github:
    ```sh
    git clone https://github.com/nadig-google/cluster-director-mcp.git
    cd cluster-director-mcp
    ```

2.  Install cluster-director-mcp and dependencies as `gemini-cli` extensions:
    ```sh
    ./install.sh
    ```   

3. Authenticate yourself (run command - follow instructions - requires opening browser):
  ```sh
   gcloud auth application-default login
  ```
  
4. Set the default GCP project in which your clusters exist or will be created:
  ```sh
  gcloud config set project <project-name>
  ```

5. Start gemini-cli
  ```sh
  gemini
  ```

6. Ask questions
  ```sh
  "Show me the clusters in my GCP project in Cluster Director"
  "Show me information about my cluster"
  ```

## Tools

- `list_clusters`: List your clusters created using Cluster Director.
- `get_cluster`: Get detailed about a single Cluster.
- More to come soon....

## Context 

In addition to the above tools, this AI Assistant has additional fine-tuned and detailed information about Cluster Director.



