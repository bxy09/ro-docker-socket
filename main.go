package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"net/http"
	"strconv"
)

func main() {
	var (
		ep   string
		port int
	)
	flag.StringVar(&ep, "ep", "unix:///var/run/docker.sock", "entrypoint for socket")
	flag.IntVar(&port, "port", 4244, "read only socket port")
	flag.Parse()
	client, err := docker.NewClient(ep)
	if err != nil {
		fmt.Errorf("Cannot connect to local docker socket")
		return
	}
	http.HandleFunc("/containers/json", func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"from": r.RemoteAddr,
		}).Info("Request conatiners/json")
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // parameter erro
			logrus.Info("Bad request, cannot parse")
			return
		}

		opts := docker.ListContainersOptions{}
		if all := r.FormValue("all"); all == "1" || all == "True" || all == "true" {
			opts.All = true
		}
		if size := r.FormValue("Size"); size == "1" || size == "True" || size == "true" {
			opts.Size = true
		}
		if limit, err := strconv.Atoi(r.FormValue("limit")); err != nil {
			opts.Limit = limit
		}
		opts.Since = r.FormValue("since")
		opts.Since = r.FormValue("before")

		if filter := r.FormValue("filter"); len(filter) > 0 {
			err := json.Unmarshal([]byte(filter), &opts.Filters)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				logrus.Info("Bad request, cannot parse the filter")
				return
			}
		}
		containers, err := client.ListContainers(docker.ListContainersOptions{})
		if err != nil {
			logrus.Error("Cannot get local container status, ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(containers)
		if err != nil {
			logrus.Error("Cannot marshal the containers, ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
		return
	})
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	fmt.Errorf("Serve quit: %s", err.Error())
	return
}
