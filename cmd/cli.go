/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"blockchain/blockchain"
	"blockchain/handlers"
	"fmt"

	"github.com/spf13/cobra"
)




var cliHandler *handlers.CLIHandler

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Interface de linha de comando para a blockchain",
	Long: `Adicionar a desc longa aqui`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Pode usar esse hook para setups comuns
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			defer cliHandler.Bc.Db.Close()
		},
}

var addBlockCmd = &cobra.Command{
	Use:   "addBlock [data]",
	Short: "Adiciona um bloco à blockchain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data := args[0]
		cliHandler.Bc.AddBlock(data)
		fmt.Println("Bloco adicionado com sucesso!")
	},
}

var printChainCmd = &cobra.Command{
	Use:   "printChain",
	Short: "Imprime toda a blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		cliHandler.PrintChain()
	},
}

func init() {
	bc := blockchain.NewBlockchain()
	cliHandler = handlers.NewCLIHandler(bc)
	cliCmd.AddCommand(addBlockCmd)
	rootCmd.AddCommand(cliCmd)
}
