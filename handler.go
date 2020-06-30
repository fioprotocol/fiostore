package fiostore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// ReqHandler is responsible for servicing the incoming request.
func ReqHandler(resp http.ResponseWriter, req *http.Request) {
	clientAddr := req.RemoteAddr
	if xff {
		xh := strings.Split(req.Header.Get("X-Forwarded-For"), ",")
		if len(xh) > 0 && xh[0] != "" {
			clientAddr = xh[0]
		}
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("%s: unable to read body, sending blank response. %v\n", clientAddr, err)
		resp.WriteHeader(400)
		resp.Write(nil)
		return
	}
	defer req.Body.Close()

	fioRequest, status, err := parseRequest(body)
	if err != nil {
		log.Printf("%s: %v\n", clientAddr, err)
		if status == nil {
			log.Printf("%s: parse request returned nil status! sending blank response\n", clientAddr)
			resp.WriteHeader(500)
			resp.Write(nil)
			return
		}
		j, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			log.Printf("%s: parse request returned invalid status struct! sending blank response %v\n", clientAddr, err)
			resp.WriteHeader(500)
			resp.Write(nil)
			return
		}
		resp.WriteHeader(status.Code)
		resp.Write(j)
		return
	}

	result, err := sendFioRequest(fioRequest)
	if err != nil {
		log.Printf("%s: %v\n", clientAddr, err)
	}
	if result == nil {
		log.Printf("%s: request to %s, send request provided nil result! sending empty response\n", clientAddr, fioRequest.FioAddress)
		resp.WriteHeader(500)
		resp.Write(nil)
		return
	}

	j, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("%s: %v\n", clientAddr, err)
		resp.WriteHeader(500)
		resp.Write(nil)
		return
	}
	if result.Code == 200 {
		log.Printf("%s: SUCCESS: sent request to %s with txid %s\n", clientAddr, fioRequest.FioAddress, result.Txid)
	} else {
		log.Printf("%s: request to %s failed %s\n", clientAddr, fioRequest.FioAddress, result.Message)
	}
	resp.WriteHeader(result.Code)
	resp.Write(j)

}
