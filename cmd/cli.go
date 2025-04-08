/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"blockchain/blockchain"
	"blockchain/handlers"

	"github.com/spf13/cobra"
)




var cliHandler *handlers.CLIHandler

var address string
var from string
var to string
var amount int

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
		cliHandler.AddBlock(data)
	},
}

var printChainCmd = &cobra.Command{
	Use:   "printChain",
	Short: "Imprime toda a blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		cliHandler.PrintChain()
	},
}

var getBalanceCmd = &cobra.Command{
	Use:   "getBalance",
	Short: "Imprime toda a blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		cliHandler.GetBalance(address)
	},
}


var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send an value",
	Run: func(cmd *cobra.Command, args []string) {
		cliHandler.Send(from, to, amount)
	},
}


func init() {
	bc := blockchain.NewBlockchain(address)
	cliHandler = handlers.NewCLIHandler(bc)
	cliCmd.AddCommand(addBlockCmd)
	cliCmd.AddCommand(printChainCmd)
	cliCmd.AddCommand(getBalanceCmd)
	cliCmd.AddCommand(sendCmd)
	cliCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "Endereço da carteira")
	cliCmd.PersistentFlags().StringVarP(&from, "from", "f", "", "Endereço da carteira")
	cliCmd.PersistentFlags().StringVarP(&to, "to", "t", "", "Endereço da carteira")
	cliCmd.PersistentFlags().IntVarP(&amount, "amount", "v", 0, "Endereço da carteira")
	rootCmd.AddCommand(cliCmd)
}
