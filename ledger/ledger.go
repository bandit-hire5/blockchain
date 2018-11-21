package ledger

import (
	"bufio"
	"encoding/json"
	"naeltok/go-blockchain/utils"
	"os"
	"time"
)

type Ledger struct {
	Filepath string
}

func NewLedger(filepath string) *Ledger {
	return &Ledger{
		Filepath: filepath,
	}
}

func (l *Ledger) Create() {
	_, err := os.Stat(l.Filepath)

	if os.IsNotExist(err) {
		file, err := os.Create(l.Filepath)
		if err != nil {
			panic(err)
		}

		defer file.Close()

		block, err := generateGenesisBlock()
		if err != nil {
			panic(err)
		}

		err = l.WriteLedgerFile(block)
		if err != nil {
			panic(err)
		}
	}
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

func (l *Ledger) GetFullLedger() ([]string, error) {
	var list []string

	file, err := os.Open(l.Filepath)
	if err != nil {
		return list, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return list, err
	}

	return list, nil
}

func (l *Ledger) GenerateNextBlock(data Data) (Block, error) {
	previousBlock, err := l.GetLatestBlock()
	if err != nil {
		return Block{}, err
	}

	nextIndex := previousBlock.Index + 1
	nextTime := time.Now().UTC()

	dataString, err := json.Marshal(&data)
	if err != nil {
		return Block{}, err
	}

	nextHash, err := utils.CalculateHash(nextIndex, previousBlock.Hash, nextTime, dataString)
	if err != nil {
		return Block{}, err
	}

	return Block{
		Index:        nextIndex,
		Hash:         nextHash,
		PreviousHash: previousBlock.Hash,
		Timestamp:    nextTime.String(),
		Data:         data,
	}, nil
}
