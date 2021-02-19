package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Node struct{
	ProjectID string	`json:"project_id"`
	Host string
	ComputeID string	`json:"compute_id"`
	Console int			`json:"console"`
	ConsoleAutoStart bool `json:"console_auto_start"`
	ConsoleHost string	`json:"console_host"`
	ConsoleType string	`json:"console_type"`
	// Custom adapters
	FirstPortName string`json:"first_port_name"`
	Height int 			`json:"height"`
	Width int			`json:"width"`
	// Label
	Locked bool			`json:"locked"`
	Name string			`json:"name"`
	NodeDirectory string`json:"node_directory"`
	ID string 			`json:"node_id"`
	NodeType string		`json:"node_type"`
	PortNameFormat string `json:"port_name_format"`
	PortSegmentSize int	`json:"port_segment_size"`
	// Ports 
	Properties map[string]interface{} `json:"properties"`
	Status string		`json:"status"`
	Symbol string 		`json:"symbol"`
	TemplateID string	`json:"template_id"`
	X int				`json:"x"`
	Y int				`json:"y"`
	Z int				`json:"z"`

	BaseProject *Project
}

type Nodes []Node

// This gets a file from the server, named filename, which is located in the folder of the node.
// It saves it to the disk in outputPath
func (n *Node) getFile(filename string, outputPath string) error {
	resp, err := http.Get(n.Host + "/v2/projects/" + n.ProjectID + "/nodes/" + n.ID + "/files/" + filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()
	
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) getNodeImage(outputPath string) error{
	image := n.Properties["hda_disk_image"].(string)
	url := n.Host + "/v2/compute/qemu/images/" + image
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil

}

func getNodes(dest string, project *Project) (Nodes, error){
	resp, err := http.Get(dest + "/v2/projects/" + project.ID + "/nodes")
	if err != nil {
		return nil, err
		//panic(err)
	}
	defer resp.Body.Close()
	var tmp Nodes
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(tmp); i++ {
		tmp[i].BaseProject = project
		tmp[i].Host = dest
	}
	return tmp, nil
}