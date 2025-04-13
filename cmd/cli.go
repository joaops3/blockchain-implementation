package cmd

import (
	"blockchain/handlers"

	"github.com/spf13/cobra"
)




var cliHandler *handlers.Handler

var address string
var from string
var to string
var amount int

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Interface de linha de comando para a blockchain",
	Long: ``,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			
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

var createWalletCmd = &cobra.Command{
	Use:   "createWallet",
	Short: "Cria uma nova carteira",
	Run: func(cmd *cobra.Command, args []string) {
		cliHandler.CreateWallet()
	},
}

var createBlockchainCmd = &cobra.Command{
	Use:   "createBlockchain",
	Short: "Cria a blockchain com um bloco gênesis",
	Run: func(cmd *cobra.Command, args []string) {
		cliHandler.CreateBlockChain(address)
	},
}

var CountTransactionsCmd = &cobra.Command{
	Use:   "countTransactions",
	Short: "Conta o número de transações na blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		cliHandler.ReindexUTXO()
	},
}


func init() {
	
	cliHandler = handlers.NewHandler()
	cliHandler.CreateBlockChain("12Q5pnzrQUzun1EtRdjKvbbMou7WUtfSRD6QieC9XTXoE8FQMQJ")
	cliCmd.AddCommand(addBlockCmd)
	cliCmd.AddCommand(printChainCmd)
	cliCmd.AddCommand(getBalanceCmd)
	cliCmd.AddCommand(createWalletCmd)
	cliCmd.AddCommand(sendCmd)
	cliCmd.AddCommand(createBlockchainCmd)
	cliCmd.AddCommand(CountTransactionsCmd)
	cliCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "Endereço da carteira")
	cliCmd.PersistentFlags().StringVarP(&from, "from", "f", "", "Endereço da carteira")
	cliCmd.PersistentFlags().StringVarP(&to, "to", "t", "", "Endereço da carteira")
	cliCmd.PersistentFlags().IntVarP(&amount, "amount", "v", 0, "Endereço da carteira")
	rootCmd.AddCommand(cliCmd)
}
