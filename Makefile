.phony: all clean

all:	cluster-director-mcp
	./cluster-director-mcp install gemini-cli

clean:
	rm -f cluster-director-mcp

cluster-director-mcp:
	go build -o cluster-director-mcp .



