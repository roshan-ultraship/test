package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"
	"github.com/boltdb/bolt"
)

type Block struct {
	Timestamp     int64
	Transaction   []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

const targetBits = 16 // Adjusted difficulty for faster mining

func NewBlock(transaction []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transaction, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (bc *Blockchain) AddBlock(transaction []*Transaction) {
	var lastHash []byte

	// View transaction to get the last block hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockBucket"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		lastHash = b.Get([]byte("l"))
		return nil
	})

	// Handle error from the View transaction
	if err != nil {
		log.Fatal(err)
	}

	// Create a new block
	newBlock := NewBlock(transaction, lastHash)

	// Update transaction to store the new block and update the last hash
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockBucket"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// Put the new block's serialized data into the bucket
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}

		// Update the 'l' key to point to the new block (latest block)
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}

		// Update the blockchain tip
		bc.tip = newBlock.Hash

		return nil
	})

	// Handle error from the Update transaction
	if err != nil {
		log.Fatal(err)
	}
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.HashTranstion(),
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}

func (block *Block) HashTranstion() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transaction {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]

}

func (pow *ProofOfWork) Validate() bool {
	var hashint big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)

	hashint.SetBytes(hash[:])

	isValid := hashint.Cmp(pow.target) == -1

	return isValid

}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	var maxNonce = math.MaxInt64


	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Transaction)
	
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func main() {


	cli := CLI{}
	cli.Run()
}
