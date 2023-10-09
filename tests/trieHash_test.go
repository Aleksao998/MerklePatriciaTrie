package tests

import (
	"encoding/json"
	"github.com/Aleksao998/Merkle-Patricia-Trie/storage/mpt"
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ActionType string

const (
	Add    ActionType = "add"
	Delete ActionType = "delete"
)

type Action struct {
	Type  ActionType `json:"type"`
	Key   string     `json:"key"`
	Value string     `json:"value,omitempty"`
}

type Example struct {
	Actions []Action `json:"actions"`
	Root    string   `json:"root"`
}

func updateString(tr *trie.Trie, key, value string) {
	tr.Put([]byte(key), []byte(value))
}

func deleteString(tr *trie.Trie, key string) {
	tr.Del([]byte(key))
}

func loadExamplesFromJSON(path string) ([]Example, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var examples []Example
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&examples)
	if err != nil {
		return nil, err
	}
	return examples, nil
}

func TestFromJSON(t *testing.T) {
	examples, err := loadExamplesFromJSON("./examples.json")
	if err != nil {
		t.Fatal(err)
	}

	for idx, example := range examples {
		db := mpt.NewMPTMemoryStorage()
		tr := trie.NewTrie(db)

		for _, action := range example.Actions {
			if action.Type == Add {
				updateString(tr, action.Key, action.Value)
			} else if action.Type == Delete {
				deleteString(tr, action.Key)
			}
		}

		exp := common.HexToHash(example.Root)

		assert.Equal(t, exp.Bytes(), tr.Hash(), "Test #%d: Unexpected root hash", idx+1)
	}
}
