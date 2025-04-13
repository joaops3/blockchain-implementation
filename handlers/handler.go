package handlers

import (
	"blockchain/blockchain"
	"fmt"
	"log"
	"strconv"
)

type Handler struct {
	Bc *blockchain.Blockchain
}

func NewHandler() *Handler {
	return &Handler{Bc: nil}
}

func (cli *Handler) CreateBlockChain(address string) {
	if address == "" {
		fmt.Println("Endereço da carteira não pode ser vazio")
		return
	}
	bc := blockchain.NewBlockchain(address)
	cli.Bc = bc
	utxpSet := blockchain.NewUTXOSet(bc)
	utxpSet.Reindex()
	fmt.Printf("Blockchain criada com sucesso! Endereço: %s\n", address)
}

func (cli *Handler) AddBlock(data string) {
	if cli.Bc == nil {
		fmt.Println("Blockchain ainda não foi criada. Use 'createBlockchain' antes.")
		return
	}
	cli.Bc.AddBlock(data)
	fmt.Printf("Bloco adicionado com sucesso: %s\n", data)
}



func (cli *Handler) PrintChain() {
	if cli.Bc == nil {
		fmt.Println("Blockchain ainda não foi criada. Use 'createBlockchain' antes.")
		return
	}
	bci := cli.Bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		for _, tx := range block.Transactions {
			fmt.Printf("Transaction ID: %x\n", tx.ID)
			for _, vin := range tx.Vin {
				fmt.Printf("  Input - TXID: %x, Out: %d, Signature: %x\n", vin.Txid, vin.Vout, vin.Signature)
			}
			for _, vout := range tx.Vout {
				fmt.Printf("  Output - Value: %d, PubKeyHash: %x\n", vout.Value, vout.PubKeyHash)
			}
		}
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}


func (cli *Handler) GetBalance(address string){
	if cli.Bc == nil {
		fmt.Println("Blockchain ainda não foi criada. Use 'createBlockchain' antes.")
		return
	}
	if address == "" {
		fmt.Println("Endereço da carteira não pode ser vazio")
		return
	}
	utxoset := blockchain.NewUTXOSet(cli.Bc)
	pubkey := blockchain.GetPubKeyFromAddress(address)
	utxos := utxoset.FindUTXO(pubkey)

	balance := 0


	for _, output := range utxos {
		balance += output.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *Handler) Send(from, to string, amount int) {
	if cli.Bc == nil {
		fmt.Println("Blockchain ainda não foi criada. Use 'createBlockchain' antes.")
		return
	}
	if from == "" || to == "" {
		fmt.Println("Endereço da carteira não pode ser vazio")
		return
	}
	if amount <= 0 {
		fmt.Println("Valor deve ser maior que zero")
		return
	}
	
	utxoset := blockchain.NewUTXOSet(cli.Bc)
	
	tx := blockchain.NewUTXOTransaction(from, to, amount, utxoset)
	
	if tx == nil {
		log.Fatalf("Erro ao criar transação")
		return
	}
	cbtx := blockchain.NewCoinbaseTX(from, "award")
	newBlock := cli.Bc.MineBlock([]*blockchain.Transaction{tx, cbtx})
	
	utxoset.Update(newBlock)
	fmt.Printf("Transação enviada de %s para %s no valor de %d\n", from, to, amount)
	
}


func (cli *Handler) CreateWallet() {
	
	node := "node1"
	wallets, err := blockchain.NewWallets(node)
	if err != nil {
		log.Fatalf("Erro ao criar carteiras: %v", err)
		return
	}
	address := wallets.CreateWallet()
	wallets.SaveToFile(node)
	fmt.Printf("Endereço da nova carteira: %s\n", address)
}

func (cli *Handler) ReindexUTXO() {

	UTXOSet := blockchain.UTXOSet{cli.Bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}