package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nadig-google/cluster-director-mcp/pkg/genericCore"
	compute "google.golang.org/api/compute/v0.alpha"
)

var authToken string
var regions []*compute.Region
var regions2Zones = make(map[string][]string)

type regionsAndClustersStruct struct {
	region      string
	clusterName string
	clusterData Cluster
}

var regionAndClusterArr []regionsAndClustersStruct

// *********************************
// This command works
// *********************************
// gcloud compute ssh cluster0vk-login-001 --project=hpc-toolkit-dev --zone=us-central1-c --tunnel-through-iap --command 'sinfo'
var lastClusterInfo Cluster

// The root struct that holds the list of clusters.
type ClustersResponse struct {
	Clusters []Cluster `json:"clusters"`
}

var region2Clusters map[string][]Cluster

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
	var zonesList []string

	// The filter string tells the API to return only zones whose region name
	// matches the one we specified.
	filter := fmt.Sprintf("name=%s-*", region)

	// Call the Zones.List method with the project ID and the filter.
	req1 := computeService.Zones.List(projectID).Filter(filter)

	// The 'Do' method handles pagination for you. We process each page of results.
	if err := req1.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			genericCore.WriteToLog(zone.Name)
			zonesList = append(zonesList, zone.Name)
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
	/*
		$ curl     -H "Content-Type:application/json"     -H "Authorization: Bearer $(gcloud auth print-access-token)"     https://hypercomputecluster.googleapis.com/v1alpha/projects/hpc-toolkit-dev/locations/
		{
		"locations": [
		{
		"name": "projects/hpc-toolkit-dev/locations/asia-southeast1",
		"locationId": "asia-southeast1"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/europe-north1",
		"locationId": "europe-north1"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/europe-west1",
		"locationId": "europe-west1"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/europe-west4",
		"locationId": "europe-west4"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/us-central1",
		"locationId": "us-central1"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/us-east1",
		"locationId": "us-east1"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/us-east4",
		"locationId": "us-east4"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/us-east5",
		"locationId": "us-east5"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/us-east7",
		"locationId": "us-east7"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/us-west1",
		"locationId": "us-west1"
		},
		{
		"name": "projects/hpc-toolkit-dev/locations/us-west4",
		"locationId": "us-west4"
		}
		]
		}
	*/

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

// Cluster flow
// The moment we have project (either during install or first call to MCP), do the following
// - Get clusters in the project
// - Run in parallel: For each cluster
//   - Get zone, region and any other meta data and store it in a cache to be used as default later
//   -

func getClustersInAllRegions(projectID string) string {
	var listOfClusters string = "["
	for region, _ := range regions2Zones {
		getClustersInRegionIfExists(region, projectID)
		var clusterList []Cluster = region2Clusters[region]
		for _, clusterStruct := range clusterList {
			listOfClusters += string("\"" + clusterStruct.Name + "\"")
		}
	}
	listOfClusters += "]"
	return listOfClusters
}

func getClustersInRegionIfExists(region string, projectID string) {
	// Equivalent CURL command:
	// curl \
	// -H "Content-Type:application/json" \
	// -H "Authorization: Bearer $(gcloud auth print-access-token)" \
	// https://hypercomputecluster.googleapis.com/v1alpha/projects/cloud-hypercomp-dev/locations/us-central1/clusters
	url := "https://hypercomputecluster.googleapis.com/v1alpha/projects/" + projectID + "/locations/" + region + "/clusters"
	/*
	   	$ curl \
	       -H "Content-Type:application/json" \
	       -H "Authorization: Bearer $(gcloud auth print-access-token)" \
	       https://hypercomputecluster.googleapis.com/v1alpha/projects/hpc-toolkit-dev/locations/us-central1/clusters
	   {
	     "clusters": [
	       {
	         "name": "projects/hpc-toolkit-dev/locations/us-central1/clusters/quadrant",
	         "createTime": "2025-07-29T18:11:10.750875543Z",
	         "updateTime": "2025-07-29T18:26:11.691043831Z",
	         "networks": [
	           {
	             "network": "projects/hpc-toolkit-dev/global/networks/quadrant-net",
	             "initializeParams": {
	               "network": "projects/hpc-toolkit-dev/global/networks/quadrant-net"
	             },
	             "subnetwork": "projects/hpc-toolkit-dev/global/networks/quadrant-net"
	           }
	         ],
	         "storages": [
	           {
	             "storage": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/quadrant-fs",
	             "initializeParams": {
	               "filestore": {
	                 "fileShares": [
	                   {
	                     "capacityGb": "1024",
	                     "fileShare": "nfsshare"
	                   }
	                 ],
	                 "tier": "TIER_ZONAL",
	                 "filestore": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/quadrant-fs",
	                 "protocol": "PROTOCOL_NFSV3"
	               }
	             },
	             "id": "home"
	           },
	           {
	             "storage": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/quadrant-fs-1",
	             "initializeParams": {
	               "filestore": {
	                 "fileShares": [
	                   {
	                     "capacityGb": "1024",
	                     "fileShare": "nfsshare"
	                   }
	                 ],
	                 "tier": "TIER_ZONAL",
	                 "filestore": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/quadrant-fs-1",
	                 "protocol": "PROTOCOL_NFSV3"
	               }
	             },
	             "id": "shared0"
	           }
	         ],
	         "compute": {
	           "resourceRequests": [
	             {
	               "id": "quadrant-rr1",
	               "zone": "us-central1-c",
	               "machineType": "n2-standard-2",
	               "guestAccelerators": [
	                 {}
	               ],
	               "disks": [
	                 {
	                   "type": "pd-balanced",
	                   "sizeGb": "100",
	                   "boot": true,
	                   "sourceImage": "projects/hpc-toolkit-dev/global/images/family/common-slurm-image"
	                 }
	               ],
	               "provisioningModel": "PROVISIONING_MODEL_STANDARD"
	             }
	           ]
	         },
	         "orchestrator": {
	           "slurm": {
	             "nodeSets": [
	               {
	                 "id": "nodeset1",
	                 "resourceRequestId": "quadrant-rr1",
	                 "storageConfigs": [
	                   {
	                     "id": "home",
	                     "localMount": "/home"
	                   },
	                   {
	                     "id": "shared0",
	                     "localMount": "/shared0"
	                   }
	                 ],
	                 "staticNodeCount": "1",
	                 "allowAutomaticUpdate": true,
	                 "enableOsLogin": true
	               }
	             ],
	             "partitions": [
	               {
	                 "id": "part1",
	                 "nodeSetIds": [
	                   "nodeset1"
	                 ]
	               }
	             ],
	             "defaultPartition": "part1",
	             "loginNodes": {
	               "machineType": "n2-standard-2",
	               "zone": "us-central1-c",
	               "count": "1",
	               "disks": [
	                 {
	                   "type": "pd-balanced",
	                   "sizeGb": "100",
	                   "boot": true,
	                   "sourceImage": "projects/hpc-toolkit-dev/global/images/family/common-slurm-image"
	                 }
	               ],
	               "enableOsLogin": true,
	               "enablePublicIps": true,
	               "instances": [
	                 {
	                   "instance": "projects/hpc-toolkit-dev/zones/us-central1-c/instances/quadrant-login-001"
	                 }
	               ],
	               "storageConfigs": [
	                 {
	                   "id": "home",
	                   "localMount": "/home"
	                 },
	                 {
	                   "id": "shared0",
	                   "localMount": "/shared0"
	                 }
	               ]
	             }
	           }
	         },
	         "reconciling": false
	       },
	       {
	         "name": "projects/hpc-toolkit-dev/locations/us-central1/clusters/clusterum7",
	         "createTime": "2025-07-28T17:18:00.818807445Z",
	         "updateTime": "2025-07-28T17:33:08.231509624Z",
	         "networks": [
	           {
	             "network": "projects/hpc-toolkit-dev/global/networks/clusterum7-net",
	             "initializeParams": {
	               "network": "projects/hpc-toolkit-dev/global/networks/clusterum7-net"
	             },
	             "subnetwork": "projects/hpc-toolkit-dev/global/networks/clusterum7-net"
	           }
	         ],
	         "storages": [
	           {
	             "storage": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/clusterum7-fs",
	             "initializeParams": {
	               "filestore": {
	                 "fileShares": [
	                   {
	                     "capacityGb": "1024",
	                     "fileShare": "nfsshare"
	                   }
	                 ],
	                 "tier": "TIER_ZONAL",
	                 "filestore": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/clusterum7-fs",
	                 "protocol": "PROTOCOL_NFSV3"
	               }
	             },
	             "id": "home"
	           }
	         ],
	         "compute": {
	           "resourceRequests": [
	             {
	               "id": "cjdcluster1",
	               "zone": "us-central1-c",
	               "machineType": "n2-standard-2",
	               "guestAccelerators": [
	                 {}
	               ],
	               "disks": [
	                 {
	                   "type": "pd-balanced",
	                   "sizeGb": "100",
	                   "boot": true,
	                   "sourceImage": "projects/hpc-toolkit-dev/global/images/family/common-slurm-image"
	                 }
	               ],
	               "provisioningModel": "PROVISIONING_MODEL_STANDARD"
	             }
	           ]
	         },
	         "orchestrator": {
	           "slurm": {
	             "nodeSets": [
	               {
	                 "id": "nodeset1",
	                 "resourceRequestId": "cjdcluster1",
	                 "storageConfigs": [
	                   {
	                     "id": "home",
	                     "localMount": "/home"
	                   }
	                 ],
	                 "staticNodeCount": "6",
	                 "enableOsLogin": true
	               }
	             ],
	             "partitions": [
	               {
	                 "id": "part1",
	                 "nodeSetIds": [
	                   "nodeset1"
	                 ]
	               }
	             ],
	             "defaultPartition": "part1",
	             "loginNodes": {
	               "machineType": "n2-standard-2",
	               "zone": "us-central1-c",
	               "count": "1",
	               "disks": [
	                 {
	                   "type": "pd-balanced",
	                   "sizeGb": "100",
	                   "boot": true,
	                   "sourceImage": "projects/hpc-toolkit-dev/global/images/family/common-slurm-image"
	                 }
	               ],
	               "enableOsLogin": true,
	               "enablePublicIps": true,
	               "instances": [
	                 {
	                   "instance": "projects/hpc-toolkit-dev/zones/us-central1-c/instances/clusterum7-login-001"
	                 }
	               ],
	               "storageConfigs": [
	                 {
	                   "id": "home",
	                   "localMount": "/home"
	                 }
	               ]
	             }
	           }
	         },
	         "reconciling": false
	       },
	       {
	         "name": "projects/hpc-toolkit-dev/locations/us-central1/clusters/cluster0vk",
	         "createTime": "2025-07-30T08:26:31.572323371Z",
	         "updateTime": "2025-07-30T08:40:05.631214829Z",
	         "networks": [
	           {
	             "network": "projects/hpc-toolkit-dev/global/networks/cluster0vk-net",
	             "initializeParams": {
	               "network": "projects/hpc-toolkit-dev/global/networks/cluster0vk-net"
	             },
	             "subnetwork": "projects/hpc-toolkit-dev/global/networks/cluster0vk-net"
	           }
	         ],
	         "storages": [
	           {
	             "storage": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/cluster0vk-fs",
	             "initializeParams": {
	               "filestore": {
	                 "fileShares": [
	                   {
	                     "capacityGb": "1024",
	                     "fileShare": "nfsshare"
	                   }
	                 ],
	                 "tier": "TIER_ZONAL",
	                 "filestore": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/cluster0vk-fs",
	                 "protocol": "PROTOCOL_NFSV3"
	               }
	             },
	             "id": "home"
	           }
	         ],
	         "compute": {
	           "resourceRequests": [
	             {
	               "id": "cluster0vk-rr1",
	               "zone": "us-central1-c",
	               "machineType": "n2-standard-2",
	               "guestAccelerators": [
	                 {}
	               ],
	               "disks": [
	                 {
	                   "type": "pd-balanced",
	                   "sizeGb": "100",
	                   "boot": true,
	                   "sourceImage": "projects/hpc-toolkit-dev/global/images/family/common-slurm-image"
	                 }
	               ],
	               "provisioningModel": "PROVISIONING_MODEL_STANDARD"
	             }
	           ]
	         },
	         "orchestrator": {
	           "slurm": {
	             "nodeSets": [
	               {
	                 "id": "nodeset1",
	                 "resourceRequestId": "cluster0vk-rr1",
	                 "storageConfigs": [
	                   {
	                     "id": "home",
	                     "localMount": "/home"
	                   }
	                 ],
	                 "staticNodeCount": "2",
	                 "enableOsLogin": true
	               }
	             ],
	             "partitions": [
	               {
	                 "id": "part1",
	                 "nodeSetIds": [
	                   "nodeset1"
	                 ]
	               }
	             ],
	             "defaultPartition": "part1",
	             "loginNodes": {
	               "machineType": "n2-standard-2",
	               "zone": "us-central1-c",
	               "count": "1",
	               "disks": [
	                 {
	                   "type": "pd-balanced",
	                   "sizeGb": "100",
	                   "boot": true,
	                   "sourceImage": "projects/hpc-toolkit-dev/global/images/family/common-slurm-image"
	                 }
	               ],
	               "enableOsLogin": true,
	               "enablePublicIps": true,
	               "instances": [
	                 {
	                   "instance": "projects/hpc-toolkit-dev/zones/us-central1-c/instances/cluster0vk-login-001"
	                 }
	               ],
	               "storageConfigs": [
	                 {
	                   "id": "home",
	                   "localMount": "/home"
	                 }
	               ]
	             }
	           }
	         },
	         "reconciling": false
	       },
	       {
	         "name": "projects/hpc-toolkit-dev/locations/us-central1/clusters/harsclus",
	         "createTime": "2025-06-18T20:13:05.824451278Z",
	         "updateTime": "2025-06-18T20:25:48.817877731Z",
	         "networks": [
	           {
	             "network": "projects/hpc-toolkit-dev/global/networks/harsclus-net",
	             "initializeParams": {
	               "network": "projects/hpc-toolkit-dev/global/networks/harsclus-net"
	             },
	             "subnetwork": "projects/hpc-toolkit-dev/global/networks/harsclus-net"
	           }
	         ],
	         "storages": [
	           {
	             "storage": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/harsclus-fs",
	             "initializeParams": {
	               "filestore": {
	                 "fileShares": [
	                   {
	                     "capacityGb": "1024",
	                     "fileShare": "nfsshare"
	                   }
	                 ],
	                 "tier": "TIER_ZONAL",
	                 "filestore": "projects/hpc-toolkit-dev/locations/us-central1-c/instances/harsclus-fs",
	                 "protocol": "PROTOCOL_NFSV3"
	               }
	             },
	             "id": "home"
	           }
	         ],
	         "compute": {
	           "resourceRequests": [
	             {
	               "id": "harsclus-rr1",
	               "zone": "us-central1-c",
	               "machineType": "n2-standard-2",
	               "guestAccelerators": [
	                 {}
	               ],
	               "provisioningModel": "PROVISIONING_MODEL_STANDARD"
	             }
	           ]
	         },
	         "orchestrator": {
	           "slurm": {
	             "nodeSets": [
	               {
	                 "id": "nodeset1",
	                 "resourceRequestId": "harsclus-rr1",
	                 "storageConfigs": [
	                   {
	                     "id": "home",
	                     "localMount": "/home"
	                   }
	                 ],
	                 "staticNodeCount": "2",
	                 "allowAutomaticUpdate": true,
	                 "enableOsLogin": true
	               }
	             ],
	             "partitions": [
	               {
	                 "id": "part1",
	                 "nodeSetIds": [
	                   "nodeset1"
	                 ]
	               }
	             ],
	             "defaultPartition": "part1",
	             "loginNodes": {
	               "machineType": "n2-standard-2",
	               "zone": "us-central1-c",
	               "count": "1",
	               "disks": [
	                 {
	                   "type": "pd-balanced",
	                   "sizeGb": "100",
	                   "boot": true
	                 }
	               ],
	               "enableOsLogin": true,
	               "enablePublicIps": true,
	               "instances": [
	                 {
	                   "instance": "projects/hpc-toolkit-dev/zones/us-central1-c/instances/harsclus-login-001"
	                 }
	               ]
	             }
	           }
	         },
	         "reconciling": false
	       }
	     ]
	   }
	*/
	//print("URL: " + url)
	genericCore.WriteToLog(fmt.Sprintf("Getting clusters in region %s URL : %s", region, url))

	// Remove all previous data about clusters in this region
	delete(region2Clusters, region)

	bodyString, success := genericCore.QueryURLAndGetResult(authToken, url)
	genericCore.WriteToLog("Body : " + string(bodyString))
	if success && strings.Contains(bodyString, "storages") {
		// If the body has "storages" than that means there is a cluster
		var parsedClusterData ClustersResponse
		err := json.Unmarshal([]byte(bodyString), &parsedClusterData)
		if err != nil {

			genericCore.WriteToLog(fmt.Sprintf("error unmarshalling JSON: %v", err))
			region2Clusters[region] = parsedClusterData.Clusters
		}
	} else {
		genericCore.WriteToLog("The response body does not contain the substring 'storages'.")
	}
}
