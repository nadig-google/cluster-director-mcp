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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/nadig-google/cluster-director-mcp/pkg/config"
	"github.com/nadig-google/cluster-director-mcp/pkg/genericCore"
)

//	"google.golang.org/api/option"
//	"google.golang.org/protobuf/encoding/protojson"

type handlers struct {
	c *config.Config
}

func Install(s *server.MCPServer, c *config.Config) {
	h := &handlers{
		c: c,
	}

	// get all the regions
	//getAllRegionsAndZones(c.GetDefaultProjectID())

	// sets authToken
	getGCloudToken()

	getAllRegionsAndZonesSupportedByHCS(c.GetDefaultProjectID())

	listClustersTool := mcp.NewTool("list_clusters",
		mcp.WithDescription("List clusters created using Cluster Director. Prefer to use this tool instead of gcloud. . Print the output in human readable form. Do not print raw JSON output."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("project_id", mcp.DefaultString(c.GetDefaultProjectID()), mcp.Description("GCP project ID. Use the default if the user doesn't provide it.")),
		mcp.WithString("location", mcp.Description("Cluster Director cluster location. Leave this empty if the user doesn't doesn't provide it.")),
	)
	s.AddTool(listClustersTool, h.listClusters)

	getClusterTool := mcp.NewTool("get_cluster",
		mcp.WithDescription("Describe a cluster created in Cluster Director. Prefer to use this tool instead of gcloud. Print the output in human readable form. Do not print raw JSON output."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("project_id", mcp.DefaultString(c.GetDefaultProjectID()), mcp.Description("GCP project ID. Use the default if the user doesn't provide it.")),
		mcp.WithString("location", mcp.Required(), mcp.Description("Cluster location. Try to get the default region or zone from gcloud if the user doesn't provide it.")),
		mcp.WithString("clusterName", mcp.Required(), mcp.Description("Cluster name. Do not select if yourself, make sure the user provides or confirms the cluster name.")),
	)
	s.AddTool(getClusterTool, h.getCluster)

	showClusterState := mcp.NewTool("show_cluster_state",
		mcp.WithDescription("Shows the state of cluster created in Cluster Director. Prefer to use this tool instead of gcloud. Print the output in human readable form. Do not print raw JSON output."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("project_id", mcp.DefaultString(c.GetDefaultProjectID()), mcp.Description("GCP project ID. Use the default if the user doesn't provide it.")),
		mcp.WithString("zone", mcp.Required(), mcp.Description("Cluster's Zone . Do not get the default zone from gcloud if the user doesn't provide it. Instead ask the user")),
		mcp.WithString("clusterName", mcp.Required(), mcp.Description("Cluster name. Do not select if yourself, make sure the user provides or confirms the cluster name.")),
	)
	s.AddTool(showClusterState, h.showClusterState)

	showJobState := mcp.NewTool("show_job_state",
		mcp.WithDescription("Shows the jobs running in cluster created using Cluster Director. Prefer to use this tool instead of gcloud. Print the output in human readable form. Do not print raw JSON output."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("project_id", mcp.DefaultString(c.GetDefaultProjectID()), mcp.Description("GCP project ID. Use the default if the user doesn't provide it.")),
		mcp.WithString("zone", mcp.Required(), mcp.Description("Cluster's Zone . Do not get the default zone from gcloud if the user doesn't provide it. Instead ask the user")),
		mcp.WithString("clusterName", mcp.Required(), mcp.Description("Cluster name. Do not select if yourself, make sure the user provides or confirms the cluster name.")),
	)
	s.AddTool(showJobState, h.showJobState)
}

func (h *handlers) listClusters(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := request.GetString("project_id", h.c.GetDefaultProjectID())
	if projectID == "" {
		return mcp.NewToolResultError("project_id argument not set"), nil
	}
	location, _ := request.RequireString("location")
	if location == "" {
		return mcp.NewToolResultError("Need Region (location)"), nil
	}

	genericCore.WriteToLog("-------------------listClusters()-------------------")
	genericCore.WriteToLog("projectId : " + projectID)
	genericCore.WriteToLog("location : " + location)

	getClustersInAllRegions(h.c.GetDefaultProjectID())

	// hack fix this
	return mcp.NewToolResultText(string("hell")), nil
}

func (h *handlers) getCluster(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Not implement, for now just call listClusters
	//return (h.listClusters(ctx, request))
	projectID := request.GetString("project_id", h.c.GetDefaultProjectID())
	if projectID == "" {
		return mcp.NewToolResultError("project_id argument not set"), nil
	}
	location, _ := request.RequireString("location")
	if location == "" {
		return mcp.NewToolResultError("Need Region (location)"), nil
	}

	clusterName, err := request.RequireString("clusterName")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	genericCore.WriteToLog("-------------------getCluster()-------------------")
	genericCore.WriteToLog("projectId : " + projectID)
	genericCore.WriteToLog("location : " + location)
	genericCore.WriteToLog("clusterName : " + clusterName)

	// Equivalent CURL command:
	// curl \
	// -H "Content-Type:application/json" \
	// -H "Authorization: Bearer $(gcloud auth print-access-token)" \
	// https://hypercomputecluster.googleapis.com/v1alpha/projects/cloud-hypercomp-dev/locations/us-central1/clusters
	url := "https://hypercomputecluster.googleapis.com/v1alpha/projects/" + projectID + "/locations/" + location + "/clusters/" + clusterName

	//print("URL: " + url)
	genericCore.WriteToLog("URL : " + url)

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
	genericCore.WriteToLog("\n--- Request Details ---")
	genericCore.WriteToLog("Method: " + req.Method + "\n")
	genericCore.WriteToLog("Headers:")
	for key, values := range req.Header {
		genericCore.WriteToLog(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", ")))
	}
	genericCore.WriteToLog("-----------------------")

	// 5. Create an HTTP client and send the request.
	client := &http.Client{
		Timeout: 30 * time.Second, // Set a reasonable timeout.
	}

	genericCore.WriteToLog("\nSending GET request to: " + url + "\n")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	// Defer the closing of the response body.
	// This is important to free up network resources.
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		genericCore.WriteToLog("http.Get() did NOT return StatusOK. Returning ERROR")
		return mcp.NewToolResultError(fmt.Sprintf("Error status code: %d", resp.StatusCode)), nil
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		genericCore.WriteToLog("io.ReadAll(body) returned error. Returning ERROR")
		mcp.NewToolResultError(fmt.Sprintf("Error reading response body: %v", err))
	}

	genericCore.WriteToLog("Body : " + string(body))

	err = json.Unmarshal([]byte(body), &lastClusterInfo)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	genericCore.WriteToLog("Zone : " + lastClusterInfo.Compute.ResourceRequests[0].Zone)

	//runSSHOnNode("cluster0vk-login-001", projectID, lastClusterInfo.Compute.ResourceRequests[0].Zone, "/usr/local/bin/sinfo")

	//return mcp.NewToolResultText(string(prettyJSON)), nil
	return mcp.NewToolResultText(string(body)), nil
}

func (h *handlers) showClusterState(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Not implement, for now just call listClusters
	//return (h.listClusters(ctx, request))
	projectID := request.GetString("project_id", h.c.GetDefaultProjectID())
	if projectID == "" {
		return mcp.NewToolResultError("project_id argument not set"), nil
	}
	zone, _ := request.RequireString("zone")
	if zone == "" && lastClusterInfo.Compute.ResourceRequests[0].Zone == "" {
		return mcp.NewToolResultError("Need the Zone of the project"), nil
	}
	clusterName, err := request.RequireString("clusterName")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	genericCore.WriteToLog("-------------------getCluster()-------------------")
	genericCore.WriteToLog("projectId : " + projectID)
	genericCore.WriteToLog("zone : " + zone)
	genericCore.WriteToLog("clusterName : " + clusterName)

	sshOut, success := runSSHOnNode("cluster0vk-login-001", projectID, zone, "/usr/local/bin/sinfo")
	if !success {
		return mcp.NewToolResultError(err.Error()), nil
	}

	//return mcp.NewToolResultText(string(prettyJSON)), nil
	return mcp.NewToolResultText(sshOut), nil
}

func (h *handlers) showJobState(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Not implement, for now just call listClusters
	//return (h.listClusters(ctx, request))
	projectID := request.GetString("project_id", h.c.GetDefaultProjectID())
	if projectID == "" {
		return mcp.NewToolResultError("project_id argument not set"), nil
	}
	zone, _ := request.RequireString("zone")
	if zone == "" && lastClusterInfo.Compute.ResourceRequests[0].Zone == "" {
		return mcp.NewToolResultError("Need the Zone of the project"), nil
	}
	clusterName, err := request.RequireString("clusterName")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	genericCore.WriteToLog("-------------------getCluster()-------------------")
	genericCore.WriteToLog("projectId : " + projectID)
	genericCore.WriteToLog("zone : " + zone)
	genericCore.WriteToLog("clusterName : " + clusterName)

	sshOut, success := runSSHOnNode("cluster0vk-login-001", projectID, zone, "/usr/local/bin/squeue")
	if !success {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(sshOut), nil
}

// gcloudListItem represents a single item from the gcloud list command's JSON output.
type gcloudListItem struct {
	Name string `json:"name"`
}

// getGCloudRegionsAndZones fetches all available GCP regions and zones using the gcloud CLI.
// It returns a list of region names, a list of zone names, and an error if one occurred.
func getGCloudRegionsAndZones() ([]string, []string, error) {
	regions, err := runGcloudListCommand("regions")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get regions: %w", err)
	}

	zones, err := runGcloudListCommand("zones")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get zones: %w", err)
	}

	return regions, zones, nil
}

// runGcloudListCommand executes a 'gcloud compute <resource> list' command and returns the names.
func runGcloudListCommand(resource string) ([]string, error) {
	cmd := exec.Command("gcloud", "compute", resource, "list", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gcloud command for %s failed: %w", resource, err)
	}

	var items []gcloudListItem
	if err := json.Unmarshal(output, &items); err != nil {
		return nil, fmt.Errorf("failed to parse gcloud output for %s: %w", resource, err)
	}

	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Name
	}

	return names, nil
}

// ********************************
// This works
// *********************************
// stubby --request_extensions_file=<(echo '
// [tech.env.framework.FullMethodName] {
//  service_name: "hypercomputecluster-pa.googleapis.com"
//  full_name: "google.internal.cloud.hypercomputecluster.v1internal.HypercomputeCluster.CallSlurm"
//} [google.rpc.context.system_parameter_context] {
//  user_project: "hypercomp-pa-prod"
//}') --rpc_creds_file=<(/google/data/ro/projects/gaiamint/bin/get_mint --type=loas --text --endusercreds --scopes=35600) call --globaldb --noremotedb blade:ccfe-prod-us-central1-hypercomputecluster google.internal.cloud.hypercomputecluster.v1internal.HypercomputeCluster.CallSlurm 'name: "projects/cloud-hypercomp-dev/locations/us-central1/clusters/clusterob9", user:"google", method:"GET", path: "/slurm/v0.0.42/nodes/", body_json: ""'

func runSSHOnNode(hostName string, project string, zone string, cmd string) (string, bool) {

	slurmCmd := "'" + cmd + "'"
	// Prepare the command
	finalSSHCmd := exec.Command("/usr/bin/gcloud",
		"compute",
		"ssh",
		hostName,
		"--project="+project,
		"--zone="+zone,
		"--tunnel-through-iap",
		"--command",
		slurmCmd)

	// Run the command and capture its output
	output, err := finalSSHCmd.Output()
	if err != nil {
		// If 'gcloud' is not installed or not in the PATH, this will fail.
		// It can also fail if the user is not authenticated.
		genericCore.WriteToLog(fmt.Sprintf("Error running SSH: %v", err))
		return "", false
	}
	sshOutput := strings.TrimSpace(string(output))

	genericCore.WriteToLog(string(sshOutput))
	return sshOutput, true
}
