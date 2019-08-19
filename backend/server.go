package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type nodeData struct {
	LinkTo   string `json:"linkTo"`
	NodeName string `json:"nodeName"`
}

var father map[string][]string
var allNodes map[string]int

func addToDict(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Unmarshal
	var names nodeData
	err = json.Unmarshal(b, &names)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, ok := allNodes[names.NodeName]
	fmt.Println("Total Entries: ", len(father))
	if len(father) == 0 || ok {
		fmt.Printf("Added %v to %v\n", names.LinkTo, names.NodeName)
		father[names.NodeName] = append(father[names.NodeName], names.LinkTo)
		fmt.Println("Added to DB")
		allNodes[names.LinkTo] = 1
		addToAllNodes(names.NodeName)
	} else {
		retErr := map[string]string{"status": "Name is not present " + names.NodeName}
		json.NewEncoder(w).Encode(&retErr)
		return
	}
	ret := map[string]string{"status": "Added Successfully"}
	json.NewEncoder(w).Encode(&ret)
	updateStoredData()
}

func giveHimDict(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	fmt.Println("Gave the dict")
	json.NewEncoder(w).Encode(father)
}

func addToAllNodes(name string) {
	_, ok := allNodes[name]
	if ok {
		allNodes[name] = allNodes[name] + 1
	} else {
		allNodes[name] = 1
	}
}

func readStoredData() {
	var stream map[string][]string
	f, err := ioutil.ReadFile("nodes.json")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	err = json.Unmarshal(f, &stream)
	for node := range stream {
		father[node] = stream[node]
		addToAllNodes(node)
		for i := range stream[node] {
			addToAllNodes(stream[node][i])
		}
	}
}

func updateStoredData() {
	file, _ := json.MarshalIndent(father, "", " ")
	err := ioutil.WriteFile("nodes.json", file, 0644)
	if err != nil {
		fmt.Println("Updating file failed ", err)
	}
}

func main() {
	father = make(map[string][]string)
	allNodes = make(map[string]int)
	readStoredData()
	router := mux.NewRouter()
	router.HandleFunc("/addme", addToDict).Methods("POST")
	router.HandleFunc("/giveme", giveHimDict).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./../www"))))

	log.Fatal(http.ListenAndServe(":5000", router))
}
