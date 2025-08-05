package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	"net/http"

	"github.com/nadig-google/cluster-director-mcp/pkg/genericCore"
	compute "google.golang.org/api/compute/v0.alpha"
)

var authToken string
var regions []*compute.Region
var regions2Zones = make(map[string][]string)

type regionsAndClustersStruct struct {
	region  string
	cluster string
}

var regionAndClusterArr []regionsAndClustersStruct

// *********************************
// This command works
// *********************************
// gcloud compute ssh cluster0vk-login-001 --project=hpc-toolkit-dev --zone=us-central1-c --tunnel-through-iap --command 'sinfo'
var lastClusterInfo Cluster

// Cluster defines the top-level structure of the JSON object.
type Cluster struct {
	Name         string       `json:"name"`
	CreateTime   string       `json:"createTime"`
	UpdateTime   string       `json:"updateTime"`
	Networks     []Network    `json:"networks"`
	Storages     []Storage    `json:"storages"`
	Compute      Compute      `json:"compute"`
	Orchestrator Orchestrator `json:"orchestrator"`
	Reconciling  bool         `json:"reconciling"`
}

// Network corresponds to an object in the "networks" array.
type Network struct {
	Network          string `json:"network"`
	InitializeParams struct {
		Network string `json:"network"`
	} `json:"initializeParams"`
	Subnetwork string `json:"subnetwork"`
}

// Storage corresponds to an object in the "storages" array.
type Storage struct {
	Storage          string `json:"storage"`
	InitializeParams struct {
		Filestore struct {
			FileShares []struct {
				CapacityGb string `json:"capacityGb"`
				FileShare  string `json:"fileShare"`
			} `json:"fileShares"`
			Tier      string `json:"tier"`
			Filestore string `json:"filestore"`
			Protocol  string `json:"protocol"`
		} `json:"filestore"`
	} `json:"initializeParams"`
	ID string `json:"id"`
}

// Compute corresponds to the "compute" object.
type Compute struct {
	ResourceRequests []ResourceRequest `json:"resourceRequests"`
}

// ResourceRequest corresponds to an object in the "resourceRequests" array.
type ResourceRequest struct {
	ID                string                   `json:"id"`
	Zone              string                   `json:"zone"`
	MachineType       string                   `json:"machineType"`
	GuestAccelerators []map[string]interface{} `json:"guestAccelerators"`
	Disks             []Disk                   `json:"disks"`
	ProvisioningModel string                   `json:"provisioningModel"`
}

// Disk corresponds to a disk object.
type Disk struct {
	Type        string `json:"type"`
	SizeGb      string `json:"sizeGb"`
	Boot        bool   `json:"boot"`
	SourceImage string `json:"sourceImage"`
}

// Orchestrator corresponds to the "orchestrator" object.
type Orchestrator struct {
	Slurm Slurm `json:"slurm"`
}

// Slurm corresponds to the "slurm" object.
type Slurm struct {
	NodeSets         []NodeSet   `json:"nodeSets"`
	Partitions       []Partition `json:"partitions"`
	DefaultPartition string      `json:"defaultPartition"`
	LoginNodes       LoginNodes  `json:"loginNodes"`
}

// NodeSet corresponds to an object in the "nodeSets" array.
type NodeSet struct {
	ID                string          `json:"id"`
	ResourceRequestID string          `json:"resourceRequestId"`
	StorageConfigs    []StorageConfig `json:"storageConfigs"`
	StaticNodeCount   string          `json:"staticNodeCount"`
	EnableOsLogin     bool            `json:"enableOsLogin"`
}

// Partition corresponds to an object in the "partitions" array.
type Partition struct {
	ID         string   `json:"id"`
	NodeSetIDs []string `json:"nodeSetIds"`
}

// LoginNodes corresponds to the "loginNodes" object.
type LoginNodes struct {
	MachineType     string `json:"machineType"`
	Zone            string `json:"zone"`
	Count           string `json:"count"`
	Disks           []Disk `json:"disks"`
	EnableOsLogin   bool   `json:"enableOsLogin"`
	EnablePublicIps bool   `json:"enablePublicIps"`
	Instances       []struct {
		Instance string `json:"instance"`
	} `json:"instances"`
	StorageConfigs []StorageConfig `json:"storageConfigs"`
}

// StorageConfig corresponds to a storage configuration object.
type StorageConfig struct {
	ID         string `json:"id"`
	LocalMount string `json:"localMount"`
}

// getGCloudToken executes the 'gcloud auth print-access-token' command
// and returns the access token as a string.
func getGCloudToken() bool {

	if authToken != "" {
		return true
	}

	genericCore.WriteToLog("Executing 'gcloud auth print-access-token' to get bearer token...")

	// Prepare the command
	cmd := exec.Command("gcloud", "auth", "print-access-token")

	// Run the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		// If 'gcloud' is not installed or not in the PATH, this will fail.
		// It can also fail if the user is not authenticated.
		genericCore.WriteToLog(fmt.Sprintf("Error running gcloud command: %v", err))
		return false
	}

	// The output is a byte slice, so we convert it to a string and
	// trim any trailing newline or whitespace.
	authToken = strings.TrimSpace(string(output))
	genericCore.WriteToLog("Successfully retrieved access token.")
	return true
}

func getAllZonesInRegion(region string, projectID string, ctx context.Context, computeService *compute.Service) []string {
	zonesList := []string{}

	// The filter string tells the API to return only zones whose region name
	// matches the one we specified.
	filter := fmt.Sprintf("name=%s-*", region)

	// Call the Zones.List method with the project ID and the filter.
	req1 := computeService.Zones.List(projectID).Filter(filter)

	// The 'Do' method handles pagination for you. We process each page of results.
	if err := req1.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			genericCore.WriteToLog(zone.Name)
			append(zonesList, zone.Name)
		}
		return nil
	}); err != nil {
		genericCore.WriteToLog(fmt.Sprintf("Error getting zones for project %s in region %s : %v",
			projectID, region, err))
	}
	return zonesList
}

func getAllRegionsAndZonesSupportedByHCS(projectID string) bool {
	// // Location represents a single location object inside the array.
	type Location struct {
		Name       string `json:"name"`
		LocationID string `json:"locationId"`
	}

	// LocationList represents the top-level JSON object.
	type LocationList struct {
		Locations []Location `json:"locations"`
	}

	// Declare a variable of your top-level struct type
	var locationData LocationList

	// Equivalent CURL command:
	// curl \
	// -H "Content-Type:application/json" \
	// -H "Authorization: Bearer $(gcloud auth print-access-token)" \
	// https://hypercomputecluster.googleapis.com/v1alpha/projects/cloud-hypercomp-dev/locations/us-central1/clusters
	url := "https://hypercomputecluster.googleapis.com/v1alpha/projects/" + projectID + "/locations/"

	bodyJson, success := genericCore.QueryURLAndGetResult(authToken, url)
	if !success {
		genericCore.WriteToLog(fmt.Sprintf("Error "))
		return false
	}

	// Unmarshal the JSON data into the locationData variable
	// We pass the JSON string as a byte slice and a pointer to our variable.
	err := json.Unmarshal([]byte(bodyJson), &locationData)
	if err != nil {
		genericCore.WriteToLog(fmt.Sprintf("Error unmarshaling JSON: %v", err))
		return false
	}

	ctx := context.Background()
	computeService, err := compute.NewService(ctx)

	// Now you can access the data through the struct
	fmt.Println("Successfully parsed JSON!")
	fmt.Println("Found locations:")
	for _, loc := range locationData.Locations {
		fmt.Printf("- %s\n", loc.LocationID)
		getAllZonesInRegion(loc.LocationID, projectID, ctx, computeService)
	}

	return true
}

// Older obsolete code - calls GCE to get every region and spawns a thread per region
// to get all the zones in parallel
// HCS has an API that returns the list of regions it supports so this function is not
// needed, but nevertheless it could be useful later
func getAllRegionsAndZones(projectID string) bool {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		genericCore.WriteToLog(fmt.Sprintf("Error getting compute service for project %s : %v", projectID, err))
		return false
	}

	// Call the Regions.List method to get the list of regions.
	req := computeService.Regions.List(projectID)
	// The API returns a paginated list. The 'Do' method handles pagination automatically.
	if err := req.Pages(ctx, func(page *compute.RegionList) error {
		regions = append(regions, page.Items...)
		return nil
	}); err != nil {
		genericCore.WriteToLog(fmt.Sprintf("Error getting regions for project %s : %v", projectID, err))
		return false
	}

	genericCore.WriteToLog(fmt.Sprintf("Regions available in project %s:\n", projectID))
	for _, region := range regions {
		// Call the Zones.List method with the project ID and the filter.
		//req = computeService.Zones.List(projectID).Filter(filter)
		genericCore.WriteToLog(fmt.Sprintf("Zones available in project %s in %s are:", projectID, region.Name))
		getAllZonesInRegion(region.Name, projectID, ctx, computeService)

	}

	return true
}

// Cluster flow
// The moment we have project (either during install or first call to MCP), do the following
// - Get clusters in the project
// - Run in parallel: For each cluster
//   - Get zone, region and any other meta data and store it in a cache to be used as default later
//   -

func getClustersInAllRegions(projectID string) {
	for region := range regions2Zones {
		cluster, success := getClusterInRegionIfExists(region, projectID)
		if success {
			append(regionAndClusterArr(region2Cluster[region], cluster)
		}
	}
}

func getClusterInRegionIfExists(region string, projectID string) (string, bool) {
	defer wg.Done()

	// Equivalent CURL command:
	// curl \
	// -H "Content-Type:application/json" \
	// -H "Authorization: Bearer $(gcloud auth print-access-token)" \
	// https://hypercomputecluster.googleapis.com/v1alpha/projects/cloud-hypercomp-dev/locations/us-central1/clusters
	url := "https://hypercomputecluster.googleapis.com/v1alpha/projects/" + projectID + "/locations/" + region + "/clusters"

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
		genericCore.WriteToLog("Error making HTTP request")
		return "", false
	}
	// Defer the closing of the response body.
	// This is important to free up network resources.
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		genericCore.WriteToLog("http.Get() did NOT return StatusOK")
		return "", false
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		genericCore.WriteToLog("io.ReadAll(body) returned error. Returning ERROR")
		return "", false
	}

	bodyString := string(body)
	genericCore.WriteToLog("Body : " + string(body))
	if strings.Contains(bodyString, "storages") {
		// If the body has "storages" than that means there is a cluster
		return bodyString, true
	} else {
		genericCore.WriteToLog("The response body does not contain the substring 'storages'.")
		return string(body), false
	}
}
