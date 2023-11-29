package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"iu9-networks/lab3/models"
	"iu9-networks/lab3/pkg/peer"
	"os"
)

func getDataForStart(nodeName string) (peer.Peer, error) {
	file, err := os.Open("a.json")
	if err != nil {
		fmt.Println(err)
		return peer.Peer{}, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return peer.Peer{}, err
	}

	var nodes map[string]models.Node
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		fmt.Println(err)
		return peer.Peer{}, err
	}

	neighbours := make(map[string]models.Node)

	for name, node := range nodes {
		if name != nodeName {
			neighbours[name] = node
		}
	}

	return peer.Peer{
		Name:      nodeName,
		Info:      nodes[nodeName],
		Neighbors: neighbours,
	}, nil
}

func main() {

	name := flag.String("name", "", "Node name")

	flag.Parse()

	peer, err := getDataForStart(*name)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(peer)

	go peer.StartWebServer()
	go peer.StartSocket()

	// Keep the main goroutine alive
	select {}
}
