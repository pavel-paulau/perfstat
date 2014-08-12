package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Keeper struct {
	client *http.Client
	uri    string
}

func NewKeeper(host, snapshot, source string) *Keeper {
	return &Keeper{
		client: &http.Client{},
		uri:    fmt.Sprintf("http://%s/%s/%s", host, snapshot, source),
	}
}

func (k *Keeper) Store(header []string, values []float64) {
	sample := map[string]float64{}

	for i, metric := range header {
		sample[metric] = values[i]
	}
	b, err := json.Marshal(sample)
	if err != nil {
		log.Printf("%v", err)
	}
	j := bytes.NewReader(b)

	req, err := http.NewRequest("POST", k.uri, j)
	if err != nil {
		log.Printf("%v", err)
	}

	resp, err := k.client.Do(req)
	if err != nil {
		log.Printf("%v", err)
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)
}
