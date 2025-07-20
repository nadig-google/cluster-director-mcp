# cluster-director-mcp
MCP Server for Cluster Director from Google Cloud
# Cluster Director MCP Server

Enable MCP-compatible AI agents to interact with Cluster Director.

# Installation

1.  Install the tool:

    ```sh
    go install github.com/nadig-google/cluster-director-mcp@latest
    ```

    The `cluster-director-mcp` binary will be installed in the directory specified by the `GOBIN` environment variable. If `GOBIN` is not set, it defaults to `$GOPATH/bin` and, if `GOPATH` is also not set, it falls back to `$HOME/go/bin`.

    You can find the exact location by running `go env GOBIN`. If the command returns an empty value, run `go env GOPATH` to find the installation directory.

2.  Install it as a `gemini-cli` extension:

    ```sh
    cluster-director-mcp install gemini-cli
    ```

    This will create a manifest file in `./.gemini/extensions/cluster-director-mcp` that points to the installed `cluster-director-mcp` binary.

## Tools

- `create_cluster`: Creates AI optimized Clusters.
- `list_clusters`: List your clusters created using Cluster Director.
- `get_cluster`: Get detailed about a single Cluster.
- `list_recommendations`: List recommendations for your clusters created using Cluster Director.

## Context 

In addition to the tools above, a lot of value is provided through the bundled context instructions.

## Development

To compile the binary and update the `gemini-cli` extension with your local changes, follow these steps:

1.  Build the binary from the root of the project:

    ```sh
    go build -o cluster-director-mcp .
    ```

2.  Run the installation command to update the extension manifest:

    ```sh
    ./cluster-director-mcp install gemini-cli
    ```

    This will make `gemini-cli` use your locally compiled binary.

