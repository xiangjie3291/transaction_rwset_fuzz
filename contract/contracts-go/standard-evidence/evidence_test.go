package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"chainmaker.org/chainmaker/contract-utils/standard"
)

func TestName(t *testing.T) {
	metadata := standard.Metadata{
		HashType:       "file",
		HashAlgorithm:  "sha256",
		Username:       "taifu",
		Timestamp:      "1672048892",
		ProveTimestamp: "",
	}
	data, _ := json.Marshal(&metadata)
	fmt.Println(string(data))

	evidence1 := standard.Evidence{
		Id:       "id1",
		Hash:     "hash1",
		Metadata: string(data),
	}

	metadata.HashType = "text"
	data, _ = json.Marshal(&metadata)
	evidence2 := standard.Evidence{
		Id:       "id2",
		Hash:     "hash2",
		Metadata: string(data),
	}
	evidences := make([]standard.Evidence, 0)
	evidences = append(evidences, evidence1)
	evidences = append(evidences, evidence2)

	data, _ = json.Marshal(&evidences)
	fmt.Println(string(data))

	// 入参
	inParam := struct {
		Evidences string `json:"evidences"`
	}{
		Evidences: string(data),
	}

	data, _ = json.Marshal(&inParam)
	fmt.Println(string(data))

}
