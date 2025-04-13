package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

const walletFile = "wallet_%s.dat"

// Wallets stores a collection of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}


func NewWallets(nodeID string) (*Wallets, error) {
	ws := &Wallets{Wallets: make(map[string]*Wallet)}
	if err := ws.loadFromFile(nodeID); err != nil {
		return nil, err
	}
	return ws, nil
}


func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()
	ws.Wallets[string(address)] = wallet
	return string(address)
}


func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}


func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}


func (ws *Wallets) loadFromFile(nodeID string) error {
	filePath := getWalletFilePath(nodeID)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// arquivo ainda não existe — tudo certo
			ws.Wallets = make(map[string]*Wallet)
			return nil
		}
		return fmt.Errorf("failed to read wallet file: %w", err)
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(ws); err != nil {
		return fmt.Errorf("failed to decode wallet data: %w", err)
	}

	return nil
}


func (ws *Wallets) SaveToFile(nodeID string) {
	var content bytes.Buffer
	filePath := getWalletFilePath(nodeID)

	encoder := gob.NewEncoder(&content)
	if err := encoder.Encode(ws); err != nil {
		log.Panicf("failed to encode wallets: %v", err)
	}

	if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
		log.Panicf("failed to write wallet file: %v", err)
	}
}

func getWalletFilePath(nodeID string) string {
	return fmt.Sprintf(walletFile, nodeID)
}