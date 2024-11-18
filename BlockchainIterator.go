package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) iterater() *BlockchainIterator{

bci := &BlockchainIterator{bc.tip , bc.db}

return bci

}


func (bci *BlockchainIterator) Next() *Block {
	var block *Block

	// Start a read-only transaction using the View method
	err := bci.db.View(func(tx *bolt.Tx) error {
		// Get the bucket named "blockBucket"
		b := tx.Bucket([]byte("blockBucket"))

		// Handle case where bucket does not exist
		if b == nil {
			fmt.Println("Error: bucket not found")
			return nil // Return nil so that the transaction ends, but no error is returned
		}

		// Retrieve the encoded block using the current hash
		encodedBlock := b.Get(bci.currentHash)

		// Handle case where the block is not found
		if encodedBlock == nil {
			fmt.Println("Error: block not found")
			return nil // No block found, but no error returned
		}

		// Deserialize the block from the encoded data
		block = DeserailizeBlock(encodedBlock)

		return nil // Success
	})

	// Log the error if needed, but don't return it
	if err != nil {
		fmt.Println("Error during database transaction:", err)
	}

	// Update the iterator's current hash to the previous block's hash if block is not nil
	if block != nil {
		bci.currentHash = block.PrevBlockHash
	}

	// Return the block (or nil if an error occurred)
	return block
}


