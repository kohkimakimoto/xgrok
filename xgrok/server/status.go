package server

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/kohkimakimoto/xgrok/xgrok/log"
)

func startStatusServer(addr string) {
	http.HandleFunc("/status", statusHandler)

	log.Info("Listening for status server on %v", addr)

	go http.ListenAndServe(addr, nil)
}


type StatusResp struct {
	Tunnels []*StatusRespTunnel
}

type StatusRespTunnel struct {
	Name string
	URL  string
}


func statusHandler(w http.ResponseWriter, r *http.Request) {
	resp := &StatusResp{
		Tunnels: []*StatusRespTunnel{},
	}

	for name, tunnel := range tunnelRegistry.tunnels {
		resp.Tunnels = append(resp.Tunnels, &StatusRespTunnel{
			Name: name,
			URL: tunnel.url,
		})
	}

	respBody, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(respBody))

}
