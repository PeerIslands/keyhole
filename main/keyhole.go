// Copyright 2020 Kuei-chun Chen. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/simagix/keyhole/mdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var repo = "simagix/keyhole"
var version = "devel-xxxxxx"

// func main() {
// 	if version == "devel-xxxxxx" {
// 		version = "devel-" + time.Now().Format("20060102")
// 	}
// 	fullVersion := fmt.Sprintf(`%v %v`, repo, version)
// 	keyhole.Run(fullVersion)
// }

type Params struct {
	Uri       string
	Redaction bool
}

func AllInfoHandler(w http.ResponseWriter, r *http.Request) {
	var params Params

	if version == "devel-xxxxxx" {
		version = "devel-" + time.Now().Format("20060102")
	}
	fullVersion := fmt.Sprintf(`%v %v`, repo, version)

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uri := params.Uri

	var client *mongo.Client
	// connection string is required from here forward
	var connString connstring.ConnString
	if connString, err = mdb.ParseURI(uri); err != nil {
		log.Fatal(err)
	}

	uri = connString.String() // password can be injected if missing

	if client, err = mdb.NewMongoClient(uri); err != nil {
		log.Fatal(err)
	}

	stats := mdb.NewClusterStats(fullVersion)
	stats.SetRedaction(*&params.Redaction)
	stats.SetVerbose(true)
	if err = stats.GetClusterStats(client, connString); err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/allinfo", AllInfoHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
