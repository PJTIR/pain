package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Snapshot struct{
	CreatedAt uint64		`json:"created_at"`
	Name string				`json:"name"`
	ProjectID string 		`json:"project_id"`
	SnapshotID string 		`json:"snapshot_id"`
}

type Snapshots []Snapshot

func getSnapshots(dest string, projectid string) (Snapshots, error){
	url := dest + "/v2/projects/" + projectid + "/snapshots"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var tmp Snapshots
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}