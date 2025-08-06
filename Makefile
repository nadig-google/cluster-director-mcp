.phony: all clean install

all:	cluster-director-mcp install
	

clean:
	rm -f cluster-director-mcp
	rm -f .gemini/extensions/cluster-director-mcp/gemini-extension.json
	rm -f .gemini/extensions/cluster-director-mcp/gemini-extension.json.orig

install:
	./cluster-director-mcp install gemini-cli
	mv .gemini/extensions/cluster-director-mcp/gemini-extension.json .gemini/extensions/cluster-director-mcp/gemini-extension.json.orig
	./updateMcpServers.pl .gemini/extensions/cluster-director-mcp/gemini-extension.json.orig .gemini/extensions/cluster-director-mcp/gemini-extension.json

cluster-director-mcp:
	go build -o cluster-director-mcp .




