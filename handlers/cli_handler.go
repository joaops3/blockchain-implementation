package handlers

import (
	"blockchain/blockchain"
	"fmt"
	"log"
	"strconv"
)

type CLIHandler struct {
	Bc *blockchain.Blockchain
}

func NewCLIHandler(bc *blockchain.Blockchain) *CLIHandler {
	return &CLIHandler{Bc: bc}
}

func (cli *CLIHandler) AddBlock(data string) {
	cli.Bc.AddBlock(data)
	fmt.Printf("Bloco adicionado com sucesso: %s\n", data)
}


func (cli *CLIHandler) PrintChain() {
	bci := cli.Bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}


func (cli *CLIHandler) GetBalance(address string){
	if address == "" {
		fmt.Println("Endereço da carteira não pode ser vazio")
		return
	}
	utxos := cli.Bc.FindUTXO(address)

	balance := 0


	for _, output := range utxos {
		balance += output.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLIHandler) Send(from, to string, amount int) {
	if from == "" || to == "" {
		fmt.Println("Endereço da carteira não pode ser vazio")
		return
	}
	if amount <= 0 {
		fmt.Println("Valor deve ser maior que zero")
		return
	}

	tx := blockchain.NewTransaction(from, to, amount, cli.Bc)
	if tx == nil {
		log.Fatalf("Erro ao criar transação")
		return
	}
	cli.Bc.MineBlock([]*blockchain.Transaction{tx})
	fmt.Printf("Transação enviada de %s para %s no valor de %d\n", from, to, amount)
	
}