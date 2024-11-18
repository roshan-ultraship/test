package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

func (block *Block) Serialize() []byte {
	var results bytes.Buffer
	encoder := gob.NewEncoder(&results)

	// Encode the block and handle any potential errors
	if err := encoder.Encode(block); err != nil {
		log.Fatalf("Failed to serialize block: %v", err) // You can handle it differently, such as logging or returning an empty byte slice
	}

	return results.Bytes()
}



func DeserailizeBlock(d []byte) *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))

	if err := decoder.Decode(&block); err != nil {
		log.Fatalf("Failed to serialize block: %v", err) // You can handle it differently, such as logging or returning an empty byte slice
	}

	return &block

}
