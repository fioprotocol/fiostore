package fiostore

import "testing"

func Test_parseRequest(t *testing.T) {
	tokens = []string{"abc", "123"}

	_, resp, _ := parseRequest([]byte(`{"fio_address":"","chain_code":"","token_code":"","amount":0.0,"memo":"","access_token":""}`))
	if resp == nil {
		t.Error("nil response")
		return
	}
	if resp.Code != 403 {
		t.Error("failed auth check")
	}

	_, resp, _ = parseRequest([]byte(`{"fio_address":"test@fiotestnet","chain_code":"","token_code":"","amount":0.0,"memo":"","access_token":"abc"}`))
	if resp == nil {
		t.Error("nil response")
		return
	}
	if resp.Code == 403 {
		t.Error("auth check did not allow valid token")
	}
	if resp.Message != `request fields cannot be blank` {
		t.Error("blank field check failed")
	}

	_, resp, _ = parseRequest([]byte(`{"fio_address":"b@d@fiotestnet","chain_code":"","token_code":"","amount":0.0,"memo":"","access_token":"abc"}`))
	if resp == nil {
		t.Error("nil response")
		return
	}
	if resp.Message != `invalid FIO address` {
		t.Error("allowed invalid fio address")
	}
}

// {"fio_address":"","chain_code":"","token_code":"","amount":0.0,"memo":"","access_token":""}