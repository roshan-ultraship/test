package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)



type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func NewBlockchain (address string) *Blockchain {

	var tip []byte 

db , err := bolt.Open("blockchaindb", 0600 , nil)

if err != nil {
	fmt.Println("error" , err )
}

err = db.Update(func (tx *bolt.Tx) error {

	cbtx := NewCoinBaseTransaction(address, genesisCoinbaseData)


b := tx.Bucket([]byte("blockBucket"))

if b== nil {
	genesis := NewGenesisBlock(cbtx)

b, err := tx.CreateBucket([]byte("blockBucket"))
if err != nil {
	return err
}

err = b.Put(genesis.Hash, genesis.Serialize())
if err != nil {
	return err
}

err = b.Put([]byte("l") , genesis.Hash)
if err != nil {
	return err
}


tip = genesis.Hash

}else {
tip = b.Get([]byte("l"))
}

return nil
})

if err != nil {
	 log.Fatal(err)
}

bc := Blockchain{tip , db}

return &bc

}
