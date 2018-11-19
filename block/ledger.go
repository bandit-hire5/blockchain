package block

import (
	"bufio"
	"encoding/json"
	. "naeltok/go-blockchain/config"
	"os"
)

type Ledger struct {
	Filepath string
}

var config = Config{}

func (l *Ledger) Init() {
	if l.createLedgerFile() {
		block, err := GenerateGenesisBlock()

		if err != nil {
			panic(err)
		}

		err = l.WriteLedgerFile(block)
		if err != nil {
			panic(err)
		}
	}
}

func (l *Ledger) createLedgerFile() bool {
	var _, err = os.Stat(l.Filepath)

	if os.IsNotExist(err) {
		var file, err = os.Create(l.Filepath)

		if err != nil {
			panic(err)
		}

		defer file.Close()

		return true
	}

	return false
}

func (l *Ledger) WriteLedgerFile(data Block) error {
	var file, err = os.OpenFile(l.Filepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	dataString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(dataString) + "\n")
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func NewLedger() Ledger {
	config.Read()

	return Ledger{
		Filepath: config.LedgerPath,
	}
}

func (l *Ledger) GetLatestBlock() (Block, error) {
	file, err := os.Open(l.Filepath)
	if err != nil {
		return Block{}, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var blockRaw string
	var block Block

	for scanner.Scan() {
		blockRaw = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return Block{}, err
	}

	json.Unmarshal([]byte(blockRaw), &block)

	return block, nil
}
