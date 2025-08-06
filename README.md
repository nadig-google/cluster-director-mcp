# Cluster Director MCP Server

Use Cluster Director and deploy AI/ML clusters with GPUs using spoken english. More information about Cluster Director can be found here: https://cloud.google.com/ai-hypercomputer/docs/cluster-director

Link to gemini-cli github page: https://github.com/google-gemini/gemini-cli

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

6. Ask questions (QA Assistant) / prompt the AI agent do something (take action - Agentic Assistant):
  ```sh
  - Agentic prompt: Show me the clusters in my GCP project in Cluster Director
  - Agentic prompt: Show me information about my cluster
  - QA Assistant: What VM-types does Cluster Director support
  - QA Assistant: Does Cluster Director handle topology automatically during cluster creation
  - QA Assistant: What came first? The chicken or the egg
  ```

7. Some helpful gemini-cli commands
  ```sh
  /tools - shows the list of tools installed
  ```

  ```sh
  /extensions - shows the list of tools installed
  ```

  ```sh
  /quit - exit gemini-cli
  ```

8. To remove cluster-director-mcp and do a fresh install. Note this is will permanently delete all information from the previous installation
  ```sh
  cd ~
  ```

  ```sh
  # Use with caution
  rm -r cluster-director-mcp
  ```

  ```sh
  Go to step 1) in this document
  ```

## QA Assistant

This AI Assistant has a rich set of curated documents about Cluster Director to enable it to answer questions.

## Agentic Assistant (Current list of tools it can run)

- `list_clusters`: List your clusters created using Cluster Director.
- `get_cluster`: Get detailed about a single Cluster.
- More to come soon....

## Feedback
We'd love to hear from you. Please email nadig at-symbol google dot com 



