package ledger

import (
	"encoding/json"
	"naeltok/go-blockchain/utils"
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

func generateGenesisBlock() (Block, error) {
	data := Data{
		Message: "genesis block",
	}

	dataString, err := json.Marshal(&data)
	if err != nil {
		return Block{}, err
	}

	nextTime := time.Now().UTC()
	nextHash, err := utils.CalculateHash(0, "0", nextTime, dataString)

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
