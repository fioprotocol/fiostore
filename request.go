package fiostore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fioprotocol/fio-go"
	"github.com/fioprotocol/fio-go/eos"
	"log"
)

// Request is the expected input from a client
type Request struct {
	FioAddress  fio.Address `json:"fio_address"`
	Payee       string      `json:"payee"`
	ChainCode   string      `json:"chain_code"`
	TokenCode   string      `json:"token_code"`
	Amount      float32     `json:"amount"`
	Memo        string      `json:"memo"`
	AccessToken string      `json:"access_token"`
}

// Response is the returned data, HTTP code is set to value in Code field
type Response struct {
	Sent    bool   `json:"sent"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Txid    string `json:"txid"`
}

// parseRequest validates that the incoming response is valid, if err is non-nil, a response is provided
// with an error message.
func parseRequest(body []byte) (r *Request, resp *Response, err error) {
	resp = &Response{
		Code: 400,
		Sent: false,
	}

	r = &Request{}
	err = json.Unmarshal(body, r)
	if err != nil {
		resp.Code = 400
		resp.Message = fmt.Sprintf("could not decode request body, please check the input. %q", err.Error())
		return
	}

	authorized := false
	for _, valid := range tokens {
		if r.AccessToken == valid {
			authorized = true
			break
		}
	}
	if !authorized {
		resp.Message = "unauthorized"
		resp.Code = 403
		err = errors.New("did not provide a valid authorization token")
		return nil, resp, err
	}

	switch true {
	case !r.FioAddress.Valid():
		resp.Message = "invalid FIO address"
		err = errors.New(fmt.Sprintf("could not validate FIO address: %q", r.FioAddress))
		return nil, resp, err
	case r.ChainCode == "", r.TokenCode == "", r.Amount == 0, r.Memo == "", r.Payee == "":
		resp.Message = "request fields cannot be blank"
		err = errors.New("one or more inputs was empty")
		return nil, resp, err
	}
	return r, nil, nil
}

// sendFioRequest handles encrypting the request and sending to the chain
func sendFioRequest(r *Request) (resp *Response, err error) {
	if ok := api.RefreshFees(); !ok {
		log.Println("Warning: could not refresh FIO fees.")
	}
	resp = &Response{}
	pubKey, found, err := api.PubAddressLookup(r.FioAddress, "FIO", "FIO")
	if err != nil {
		resp.Code = 500
		resp.Message = "server error while retrieving address"
		switch err.(type) {
		case eos.APIError:
			err = fmt.Errorf("%s: %+v", err.Error(), err.(eos.APIError).ErrorStruct)
			log.Println(err)
		}
		return resp, err
	}

	if !found {
		resp.Code = 404
		resp.Message = fmt.Sprintf("no FIO key found for %q", r.FioAddress)
		return resp, errors.New("fio address not found")
	}

	encrypted, err := fio.ObtRequestContent{
		PayeePublicAddress: r.Payee,
		Amount:             fmt.Sprintf("%f", r.Amount),
		ChainCode:          r.ChainCode,
		TokenCode:          r.TokenCode,
		Memo:               r.Memo,
	}.Encrypt(account, pubKey.PublicAddress)
	if err != nil {
		resp.Code = 500
		resp.Message = "could not create encrypted payload"
		return resp, err
	}

	result, err := api.SignPushActions(fio.NewFundsReq(account.Actor, string(r.FioAddress), sender, encrypted))
	if err != nil {
		resp.Code = 500
		resp.Message = fmt.Sprintf("could not send transaction: %q", err.Error())
		switch err.(type) {
		case eos.APIError:
			err = fmt.Errorf("%s: %+v", err.Error(), err.(eos.APIError).ErrorStruct)
			log.Println(err)
		}
		return resp, err
	}

	resp.Txid = result.TransactionID
	resp.Code = 200
	resp.Sent = true
	resp.Message = "success"

	return resp, nil
}
