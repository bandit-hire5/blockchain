package block

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"time"
)

type Block struct {
	Index        int64
	Hash         string
	PreviousHash string
	Timestamp    string
	Data         Data
	// ...
}

type Data struct {
	Message string `json:"message"`
}

var ledger Ledger

func Init() {
	ledger = NewLedger()
	ledger.Init()
}

func GenerateNextBlock(data Data) (Block, error) {
	previousBlock, err := ledger.GetLatestBlock()
	if err != nil {
		return Block{}, err
	}

	nextIndex := previousBlock.Index + 1
	nextTime := time.Now().UTC()

	nextHash, err := calculateHash(nextIndex, previousBlock.Hash, nextTime, data)
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

func WriteNewBlock(block Block) error {
	return ledger.WriteLedgerFile(block)
}

func calculateHash(index int64, prevBlockHash string, timestamp time.Time, data Data) (string, error) {
	dataString, err := json.Marshal(&data)
	if err != nil {
		return "", err
	}

	timeString := timestamp.String()
	str := []byte(string(index) + prevBlockHash + timeString + string(dataString))

	hasher := sha1.New()
	hasher.Write(str)

	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return sha, nil
}

func GenerateGenesisBlock() (Block, error) {
	data := Data{
		Message: "genesis block",
	}
	nextTime := time.Now().UTC()
	nextHash, err := calculateHash(0, "0", nextTime, data)

	if err != nil {
		return Block{}, err
	}

	return Block{
		Index:        0,
		Hash:         nextHash,
		PreviousHash: "",
		Timestamp:    nextTime.String(),
		Data:         data,
	}, nil
}
