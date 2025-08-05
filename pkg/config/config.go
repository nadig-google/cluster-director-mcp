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

package config

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/nadig-google/cluster-director-mcp/pkg/genericCore"
)

type Config struct {
	userAgent        string
	defaultProjectID string
	defaultZone      string
	defaultRegion    string
}

func (c *Config) UserAgent() string {
	return c.userAgent
}

func (c *Config) GetDefaultProjectID() string {
	return c.defaultProjectID
}

func (c *Config) SetDefaultProjectID(p string) {
	c.defaultProjectID = p
}

func (c *Config) GetDefaultZone() string {
	return c.defaultZone
}

func (c *Config) SetDefaultZone(p string) {
	c.defaultZone = p
}

func (c *Config) GetDefaultRegion() string {
	return c.defaultRegion
}

func (c *Config) SetDefaultRegion(p string) {
	c.defaultRegion = p
}

func New(version string) *Config {
	return &Config{
		userAgent:        "cluster-director-mcp/" + version,
		defaultProjectID: getDefaultProjectID(),
	}
}

func getDefaultProjectID() string {
	/*
			ctx := context.Background()
			creds, err := google.FindDefaultCredentials:qcreds(ctx)
			if err != nil {
				genericCore.WriteToLog(fmt.Sprintf("Failed to find default credentials: %v", err))
				log.Fatalf("Failed to find default credentials: %v", err)
			}


		projectID, err := gcp.DefaultProjectID(nil) // nil for default credentials

		// The ProjectID field will be populated if ADC can find it.
		if err == nil {
			genericCore.WriteToLog("Could not find a default project ID. ")
			genericCore.WriteToLog("Set the GOOGLE_CLOUD_PROJECT environment variable or run 'gcloud config set project YOUR_PROJECT_ID'.")
			fmt.Println("Could not find a default project ID.")
			fmt.Println("Set the GOOGLE_CLOUD_PROJECT environment variable or run 'gcloud config set project YOUR_PROJECT_ID'.")
		} else {
			genericCore.WriteToLog("Default Project ID: " + string(projectID))
			fmt.Printf("Default Project ID: %s\n", string(projectID))
		}

		return string(projectID)
	*/

	out, err := exec.Command("gcloud", "config", "get", "core/project").Output()
	if err != nil {
		genericCore.WriteToLog(fmt.Sprintf("Failed to get default project: %v", err))
		return ""
	}
	projectID := strings.TrimSpace(string(out))
	log.Printf("Using default project ID: %s", projectID)
	genericCore.WriteToLog(fmt.Sprintf("Using default project ID: %s", projectID))
	return projectID
}
