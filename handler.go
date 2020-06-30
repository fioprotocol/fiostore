package fiostore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ReqHandler(resp http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("%s: unable to read body, sending blank response. %v\n", req.RemoteAddr, err)
		resp.WriteHeader(400)
		resp.Write(nil)
		return
	}
	defer req.Body.Close()

	fioRequest, status, err := parseRequest(body)
	if err != nil {
		log.Printf("%s: %v\n", req.RemoteAddr, err)
		if status == nil {
			log.Printf("%s: parse request returned nil status! sending blank response\n", req.RemoteAddr)
			resp.WriteHeader(500)
			resp.Write(nil)
			return
		}
		j, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			log.Printf("%s: parse request returned invalid status struct! sending blank response %v\n", req.RemoteAddr, err)
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
		log.Printf("%s: %v\n", req.RemoteAddr, err)
	}
	if result == nil {
		log.Printf("%s: request to %s, send request provided nil result! sending empty response\n", req.RemoteAddr, fioRequest.FioAddress)
		resp.WriteHeader(500)
		resp.Write(nil)
		return
	}

	j, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("%s: %v\n", req.RemoteAddr, err)
		resp.WriteHeader(500)
		resp.Write(nil)
		return
	}
	if result.Code == 200 {
		log.Printf("%s: SUCCESS: sent request to %s with txid %s\n", req.RemoteAddr, fioRequest.FioAddress, result.Txid)
	} else {
		log.Printf("%s: request to %s failed %s\n", req.RemoteAddr, fioRequest.FioAddress, result.Message)
	}
	resp.WriteHeader(result.Code)
	resp.Write(j)

}
