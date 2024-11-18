package main

import (
	"encoding/hex"
)

func (In *TXInput) CanUnlockOutwith(unlockinhgdata string) bool {
	return In.Scriptsig == unlockinhgdata
}

func (out *TXOutput) CanBeUnlockedWith(unlockinhgdata string) bool {
	return out.ScriptPubKey == unlockinhgdata
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {

	var UnspentTxs []Transaction
	spentTx0s := make(map[string][]int)
	bci := bc.iterater()

	for {
		block := bci.Next()

		for _, tx := range block.Transaction {
			txId := hex.EncodeToString(tx.ID)

		Outputs:
			for outidx, out := range tx.Vout {
				if spentTx0s[txId] != nil {
					for _, spentout := range spentTx0s[txId] {
						if spentout == outidx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					UnspentTxs = append(UnspentTxs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutwith(address) {
						intxID := hex.EncodeToString(in.Txid)
						spentTx0s[intxID] = append(spentTx0s[intxID], in.Vout)
					}
				}
			}

		}

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}

	return UnspentTxs

}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput

	unspenttransaction := bc.FindUnspentTransactions(address)

	for _, tx := range unspenttransaction {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs

}



