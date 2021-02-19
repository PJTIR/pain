package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Project is a GNS3 project represented in a struct. It contains all (tm) variables returned on an API call, as well as additional things, such as Nodes
type Project struct{
	/* Returned values from API call */

	AutoClose bool				`json:"auto_close"`
	AutoOpen bool				`json:"auto_open"`
	AutoStart bool				`json:"auto_start"`
	Filename string				`json:"filename"`
	GridSize int 				`json:"grid_size"`
	Name string					`json:"name"`
	Path string 				`json:"path"`
	ID string					`json:"project_id"`
	SceneHeight int				`json:"scene_height"`
	SceneWidth int 				`json:"scene_width"`
	ShowGrid bool				`json:"show_grid"`
	ShowInterfaceLabels bool	`json:"show_interface_labels"`
	ShowLayers bool				`json:"show_layers"`
	SnapToGrid bool				`json:"snap_to_grid"`
	Status string 				`json:"status"`
	Zoom int					`json:"zoom"`
	/* There are two official values missing here, which are documented here http://api.gns3.net/en/2.1/api/v2/controller/project/projects.html#get-v2-projects 
	They are "supplier" and "variables". Currently I don't know what they mean and they might not be useful for now.
	*/
	
	/* Not actually returned, but necessary nonetheless for normal operation */
	Host string	
	Nodes Nodes
	Snapshots Snapshots
}

// Projects is a type defining a slice of project objects
type Projects []Project


/* 
	The following functions are defined in the API documentation here: http://api.gns3.net/en/2.1/api/v2/controller/project.html
	And therefor shall be implemented as methods
*/

func (p *Project) open() error{
	/* Honestly this function should do more error checking, will go in the todo */
	data := "{}"
	url := p.Host + "/v2/projects/" + p.ID + "/open"

	requestBody := bytes.NewBuffer([]byte(data))
	resp, err := http.Post(url, "application/json", requestBody)
	
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// createSnapshot creates a snapshot of an entire project.
// Please specify the name of the snapshot
func (p *Project) createSnapshot(name string) error{
	data, err := json.Marshal(
		struct{
			Name string		`json:"name"` // Snapshot name
		} {
			name,
		},
	)
	if err != nil {
		return err
	}
	url := p.Host + "/v2/projects/" + p.ID + "/snapshots"
	
	requestBody := bytes.NewBuffer([]byte(data))
	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var tmp Snapshot
	err = json.Unmarshal(response, &tmp)
	p.Snapshots = append(p.Snapshots, tmp)
	return nil
}

func getProjectByName(dest string, name string) (Project, error){
	resp, err := http.Get(dest + "/v2/projects")
	if err != nil {
		return Project{}, err
	}
	defer resp.Body.Close()
	var tmp Projects
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Project{}, err
	}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return Project{}, err
	}

	var pointer int
	for i := 0; i < len(tmp); i++ {
		if tmp[i].Name != name {
			continue
		}
		pointer = i
	}
	
	tmp[pointer].Nodes, err = getNodes(dest, &tmp[pointer])
	if err != nil {
		return Project{}, err
	}
	tmp[pointer].Snapshots, err = getSnapshots(dest, tmp[pointer].ID)
	if err != nil {
		return Project{}, err
	}
	tmp[pointer].Host = dest
	
	// Success!
	return tmp[pointer], nil
}

func (p *Project) getNodeByName(name string) (*Node, error){
	for i := 0; i < len(p.Nodes); i++ {
		if p.Nodes[i].Name == name {
			return &p.Nodes[i], nil
		}
	}
	return nil, errors.New("No node exists with such name as: " + name)
}


func listProjects(destination string) error{
	resp, err := http.Get(destination + "/v2/projects")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var tmp Projects
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("%s", body)
	json.Unmarshal(body, &tmp)
	

	for i := 0; i < len(tmp); i++ {
		tmp[i].Nodes, err = getNodes(destination, &tmp[i])
		if err != nil {
			return err
		}
		fmt.Printf("Project name: %s\n", tmp[i].Name)
		fmt.Printf(" Project status:%s\n", tmp[i].Status)
		for j := 0; j < len(tmp[i].Nodes); j++ {
			fmt.Printf("\tNode: %s\n\t  Console type:%s\n\t  Machine type:%s\n", tmp[i].Nodes[j].Name, tmp[i].Nodes[j].ConsoleType, tmp[i].Nodes[j].NodeType)
		}
	}
	return nil
}