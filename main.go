package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

type Config struct {
	Deployments []DeploymentPayload `json:"deployments"`
	Session     string              `json:"session"`
	Host        string              `json:"host"`
}

// DeploymentPayload represents the structure of the JSON payload
type DeploymentPayload struct {
	EnvironmentName string    `json:"envName"`
	NodeGroup       NodeGroup `json:"nodeGroup"`
}

// NodeGroup represents the structure of a node group in the JSON payload
type NodeGroup struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

func main() {

	configPath := flag.String("config", "config.json", "Location of the config file.")

	flag.Parse()

	fc, err := os.ReadFile(*configPath) // just pass the file name

	if err != nil {
		fmt.Println("Error reading config.json file:", err)
		return
	}

	var c Config
	err = json.Unmarshal(fc, &c)

	if err != nil {
		fmt.Println("Error unmarshaling config.json file:", err)
		return
	}

	apiEndpoint := fmt.Sprintf("%s/1.0/environment/control/rest/redeploycontainers", c.Host)
	apiSession := c.Session

	// deployments := []DeploymentPayload{
	// 	{
	// 		EnvironmentName: "stg-aprt-panel",
	// 		NodeGroup:       NodeGroup{Name: "mc", Tag: ""},
	// 	},
	// 	// {
	// 	// 	EnvironmentName: "prod-aprt-panel",
	// 	// 	NodeGroup:       NodeGroup{Name: "admin", Tag: "admin1.1.1"},
	// 	// },
	// 	// {
	// 	// 	EnvironmentName: "prod-aprt-panel",
	// 	// 	NodeGroup:       NodeGroup{Name: "mc", Tag: "mc1.1.1"},
	// 	// },
	// }

	for _, deployment := range c.Deployments {
		payload, err := json.Marshal(deployment)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		q := req.URL.Query()
		q.Add("session", apiSession)
		q.Add("useExistingVolumes", "true") // Adjust as needed
		q.Add("tag", deployment.NodeGroup.Tag)
		q.Add("envName", deployment.EnvironmentName)
		q.Add("nodeGroup", deployment.NodeGroup.Name)
		req.URL.RawQuery = q.Encode()

		req.Header.Set("Content-Type", "application/json")

		fmt.Printf("Trying to redeploy env: %s, nodegroup: %s, tag: %s \n", deployment.EnvironmentName, deployment.NodeGroup.Name, deployment.NodeGroup.Tag)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode == http.StatusOK {
			b, _ := httputil.DumpResponse(resp, true)
			fmt.Println("Redeploy result is: ")
			fmt.Println(string(b))
			// fmt.Printf("Container redeployed successfully for environment %s and node group %s.\n", deployment.EnvironmentName, deployment.NodeGroup.Name)
		} else {
			fmt.Printf("Error redeploying container. Status Code: %d\n", resp.StatusCode)
			// You might want to read the response body for more details on the error
			responseBody, _ := io.ReadAll(io.Reader(resp.Body))
			fmt.Println(string(responseBody))
		}
	}
}
