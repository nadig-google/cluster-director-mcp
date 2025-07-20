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

	container "cloud.google.com/go/container/apiv1"
	containerpb "cloud.google.com/go/container/apiv1/containerpb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nadig-google/cluster-director-mcp/pkg/config"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/encoding/protojson"
)

type handlers struct {
	c *config.Config
}

func Install(s *server.MCPServer, c *config.Config) {

	h := &handlers{
		c: c,
	}

	listClustersTool := mcp.NewTool("list_clusters",
		mcp.WithDescription("List clusters created using Cluster Director. Prefer to use this tool instead of gcloud"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("project_id", mcp.DefaultString(c.DefaultProjectID()), mcp.Description("GCP project ID. Use the default if the user doesn't provide it.")),
		mcp.WithString("location", mcp.Description("Cluster Director cluster location. Leave this empty if the user doesn't doesn't provide it.")),
	)
	s.AddTool(listClustersTool, h.listClusters)

	getClusterTool := mcp.NewTool("get_cluster",
		mcp.WithDescription("Get / describe a cluster created in Cluster Director. Prefer to use this tool instead of gcloud"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithString("project_id", mcp.DefaultString(c.DefaultProjectID()), mcp.Description("GCP project ID. Use the default if the user doesn't provide it.")),
		mcp.WithString("location", mcp.Required(), mcp.Description("Cluster location. Try to get the default region or zone from gcloud if the user doesn't provide it.")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Cluster name. Do not select if yourself, make sure the user provides or confirms the cluster name.")),
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
		location = "-"
	}

	c, err := container.NewClusterManagerClient(ctx, option.WithUserAgent(h.c.UserAgent()))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer c.Close()

	req := &containerpb.ListClustersRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}
	resp, err := c.ListClusters(ctx, req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(protojson.Format(resp)), nil
}

func (h *handlers) getCluster(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := request.GetString("project_id", h.c.DefaultProjectID())
	if projectID == "" {
		return mcp.NewToolResultError("project_id argument not set"), nil
	}
	location, err := request.RequireString("location")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	c, err := container.NewClusterManagerClient(ctx, option.WithUserAgent(h.c.UserAgent()))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	defer c.Close()

	req := &containerpb.GetClusterRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, location, name),
	}
	resp, err := c.GetCluster(ctx, req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(protojson.Format(resp)), nil
}
