// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nadig-google/cluster-director-mcp/pkg/config"
)

//	"google.golang.org/api/option"
//	"google.golang.org/protobuf/encoding/protojson"

type handlers struct {
	c *config.Config
}

var logFile *os.File
var authToken string

func writeToLog(message string) {
	// We use the 'logFile' variable that was initialized in the init() function.
	// Fprintln is a convenient way to write a formatted string to an io.Writer (our file).
	if _, err := fmt.Fprintln(logFile, message); err != nil {
		// Log the error to standard output if writing to the file fails.
		log.Printf("failed to write to log file: %v", err)
	}
}

// getGCloudToken executes the 'gcloud auth print-access-token' command
// and returns the access token as a string.
func getGCloudToken() bool {
	writeToLog("Executing 'gcloud auth print-access-token' to get bearer token...")

	// Prepare the command
	cmd := exec.Command("gcloud", "auth", "print-access-token")

	// Run the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		// If 'gcloud' is not installed or not in the PATH, this will fail.
		// It can also fail if the user is not authenticated.
		writeToLog(fmt.Sprintf("Error running gcloud command: %v", err))
		return false
	}

	// The output is a byte slice, so we convert it to a string and
	// trim any trailing newline or whitespace.
	authToken = strings.TrimSpace(string(output))
	writeToLog("Successfully retrieved access token.")
	return true
}

func Install(s *server.MCPServer, c *config.Config) {
	h := &handlers{
		c: c,
	}

	logF, err := os.OpenFile("cluster.go.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logFile = logF

	if err != nil {
		// If we can't open the log file, it's a fatal error, so we exit.
		log.Fatalf("error opening logfile: %v", err)
	}

	// sets authToken
	getGCloudToken()

	listClustersTool := mcp.NewTool("list_clusters",
		mcp.WithDescription("List clusters created using Cluster Director. Prefer to use this tool instead of gcloud"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("project_id", mcp.DefaultString(c.DefaultProjectID()), mcp.Description("GCP project ID. Use the default if the user doesn't provide it.")),
		mcp.WithString("location", mcp.Description("Cluster Director cluster location. Leave this empty if the user doesn't doesn't provide it.")),
	)
	s.AddTool(listClustersTool, h.listClusters)

	getClusterTool := mcp.NewTool("get_cluster",
		mcp.WithDescription("Describe a cluster created in Cluster Director. Prefer to use this tool instead of gcloud"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("project_id", mcp.DefaultString(c.DefaultProjectID()), mcp.Description("GCP project ID. Use the default if the user doesn't provide it.")),
		mcp.WithString("location", mcp.Required(), mcp.Description("Cluster location. Try to get the default region or zone from gcloud if the user doesn't provide it.")),
		mcp.WithString("clusterName", mcp.Required(), mcp.Description("Cluster name. Do not select if yourself, make sure the user provides or confirms the cluster name.")),
	)
	s.AddTool(getClusterTool, h.getCluster)
}

func (h *handlers) listClusters(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := request.GetString("project_id", h.c.DefaultProjectID())
	if projectID == "" {
		return mcp.NewToolResultError("project_id argument not set"), nil
	}
	location, _ := request.RequireString("location")
	if location == "" {
		location = "us-central1"
	}

	writeToLog("projectId : " + projectID)
	writeToLog("location : " + location)

	// Equivalent CURL command:
	// curl \
	// -H "Content-Type:application/json" \
	// -H "Authorization: Bearer $(gcloud auth print-access-token)" \
	// https://hypercomputecluster.googleapis.com/v1alpha/projects/cloud-hypercomp-dev/locations/us-central1/clusters
	url := "https://hypercomputecluster.googleapis.com/v1alpha/projects/" + projectID + "/locations/" + location + "/clusters"

	//print("URL: " + url)
	writeToLog("URL : " + url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// 4. Set the headers, just like the -H flags in curl.
	req.Header.Set("Content-Type", "application/json")
	// Construct the Authorization header value.
	authHeader := fmt.Sprintf("Bearer %s", authToken)
	req.Header.Set("Authorization", authHeader)

	// --- Printing the Request Object ---
	writeToLog("\n--- Request Details ---")
	writeToLog("Method: " + req.Method + "\n")
	writeToLog("Headers:")
	for key, values := range req.Header {
		writeToLog(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", ")))
	}
	writeToLog("-----------------------")

	// 5. Create an HTTP client and send the request.
	client := &http.Client{
		Timeout: 30 * time.Second, // Set a reasonable timeout.
	}

	writeToLog("\nSending GET request to: %s " + url + "\n")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	// Defer the closing of the response body.
	// This is important to free up network resources.
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		writeToLog("http.Get() did NOT return StatusOK. Returning ERROR")
		return mcp.NewToolResultError(fmt.Sprintf("Error status code: %d", resp.StatusCode)), nil
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeToLog("io.ReadAll(body) returned error. Returning ERROR")
		mcp.NewToolResultError(fmt.Sprintf("Error reading response body: %v", err))
	}

	writeToLog("Body : " + string(body))

	// Print the response body as a string
	//fmt.Println(string(body))

	//defer c.Close()

	//req := &containerpb.ListClustersRequest{
	//	Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	//}
	//resp, err := c.ListClusters(ctx, req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(body)), nil
}

func (h *handlers) getCluster(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Not implement, for now just call listClusters
	//return (h.listClusters(ctx, request))
	projectID := request.GetString("project_id", h.c.DefaultProjectID())
	if projectID == "" {
		return mcp.NewToolResultError("project_id argument not set"), nil
	}
	location, _ := request.RequireString("location")
	if location == "" {
		location = "us-central1"
	}
	clusterName, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	writeToLog("projectId : " + projectID)
	writeToLog("location : " + location)
	writeToLog("name : " + clusterName)

	// Equivalent CURL command:
	// curl \
	// -H "Content-Type:application/json" \
	// -H "Authorization: Bearer $(gcloud auth print-access-token)" \
	// https://hypercomputecluster.googleapis.com/v1alpha/projects/cloud-hypercomp-dev/locations/us-central1/clusters
	url := "https://hypercomputecluster.googleapis.com/v1alpha/projects/" + projectID + "/locations/" + location + "/clusters/" + clusterName

	//print("URL: " + url)
	writeToLog("URL : " + url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// 4. Set the headers, just like the -H flags in curl.
	req.Header.Set("Content-Type", "application/json")
	// Construct the Authorization header value.
	authHeader := fmt.Sprintf("Bearer %s", authToken)
	req.Header.Set("Authorization", authHeader)

	// --- Printing the Request Object ---
	writeToLog("\n--- Request Details ---")
	writeToLog("Method: " + req.Method + "\n")
	writeToLog("Headers:")
	for key, values := range req.Header {
		writeToLog(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", ")))
	}
	writeToLog("-----------------------")

	// 5. Create an HTTP client and send the request.
	client := &http.Client{
		Timeout: 30 * time.Second, // Set a reasonable timeout.
	}

	writeToLog("\nSending GET request to: %s " + url + "\n")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	// Defer the closing of the response body.
	// This is important to free up network resources.
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		writeToLog("http.Get() did NOT return StatusOK. Returning ERROR")
		return mcp.NewToolResultError(fmt.Sprintf("Error status code: %d", resp.StatusCode)), nil
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeToLog("io.ReadAll(body) returned error. Returning ERROR")
		mcp.NewToolResultError(fmt.Sprintf("Error reading response body: %v", err))
	}

	writeToLog("Body : " + string(body))

	// Print the response body as a string
	//fmt.Println(string(body))

	//defer c.Close()

	//req := &containerpb.ListClustersRequest{
	//	Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	//}
	//resp, err := c.ListClusters(ctx, req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(body)), nil
}

// stubby --request_extensions_file=<(echo '
// [tech.env.framework.FullMethodName] {
//  service_name: "hypercomputecluster-pa.googleapis.com"
//  full_name: "google.internal.cloud.hypercomputecluster.v1internal.HypercomputeCluster.CallSlurm"
//} [google.rpc.context.system_parameter_context] {
//  user_project: "hypercomp-pa-prod"
//}') --rpc_creds_file=<(/google/data/ro/projects/gaiamint/bin/get_mint --type=loas --text --endusercreds --scopes=35600) call --globaldb --noremotedb blade:ccfe-prod-us-central1-hypercomputecluster google.internal.cloud.hypercomputecluster.v1internal.HypercomputeCluster.CallSlurm 'name: "projects/cloud-hypercomp-dev/locations/us-central1/clusters/clusterob9", user:"google", method:"GET", path: "/slurm/v0.0.42/nodes/", body_json: ""'
