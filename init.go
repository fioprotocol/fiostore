package fiostore

import (
	"github.com/fioprotocol/fio-go"
	"log"
	"os"
	"strings"
)

var (
	privkey string   // PRIV_KEY
	nodeos  string   // NODEOS
	sender  string   // FIO_ADDRESS
	tokens  []string // TOKENS (comma separated list)
	xff     bool     // TRUST_XFF

	account *fio.Account
	api     *fio.API
	opts    *fio.TxOptions
)

func init() {
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	var err error

	privkey = os.Getenv("PRIV_KEY")
	nodeos = os.Getenv("NODEOS")
	sender = os.Getenv("FIO_ADDRESS")
	t := strings.ReplaceAll(os.Getenv("TOKENS"), " ", "")
	switch "" {
	case privkey, nodeos, sender, t:
		log.Fatal("Please set the PRIV_KEY, NODEOS, TOKENS, and FIO_ADDRESS environment variables.")
	}
	if os.Getenv("TRUST_XFF") != "" {
		xff = true
	}

	tokens = strings.Split(t, ",")
	if len(tokens) == 0 || tokens[0] == "" {
		log.Fatal("No tokens specified")
	}

	account, api, opts, err = fio.NewWifConnect(privkey, nodeos)
	if err != nil {
		log.Fatal(err)
	}

}
