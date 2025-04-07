package handlers

import (
	"blockchain/blockchain"
	"fmt"
	"strconv"
)

type CLIHandler struct {
	Bc *blockchain.Blockchain
}

func NewCLIHandler(bc *blockchain.Blockchain) *CLIHandler {
	return &CLIHandler{Bc: bc}
}

func (cli *CLIHandler) Run(cmd string, args []string) {
	switch cmd {
	case "addBlock":
		if len(args) != 1 {
			println("Usage: addBlock <data>")
			return
		}
		data := args[0]
		cli.Bc.AddBlock(data)
	case "printChain":
		cli.PrintChain()
	default:
		println("Unknown command")
	}

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