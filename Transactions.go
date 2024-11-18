package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

type TXInput struct {
	Txid      []byte
	Vout      int
	Scriptsig string
}

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

const subsidy = 50 // Reward for mining a new block

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())

	tx.ID = hash[:]

}

func NewCoinBaseTransaction(to, data string) *Transaction {
if data == "" {
	data = fmt.Sprintf("Reward to '%s'", to)
}

txin := TXInput{[]byte{} , -1 , data}
txout := TXOutput{subsidy , to}


tx:= Transaction{nil , []TXInput{txin} , []TXOutput{txout} }

tx.SetID()

return &tx

}
