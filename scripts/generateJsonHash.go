package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/trie"
	"math/rand"
	"os"
)

const (
	NumExamples = 100
	MaxActions  = 10
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
	Root    string   `json:"root"` // placeholder, you'll compute this during the test
}

func randomString() string {
	// generate a random length between 1 and 10
	n := rand.Intn(10) + 1

	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func updateString(tr *trie.Trie, key, value string) {
	tr.Update([]byte(key), []byte(value))
}

func deleteString(tr *trie.Trie, key string) {
	tr.Delete([]byte(key))
}

func generateExamples() []Example {
	examples := make([]Example, NumExamples)

	for i := 0; i < NumExamples; i++ {
		tr := trie.NewEmpty(trie.NewDatabase(rawdb.NewMemoryDatabase(), nil))
		numActions := rand.Intn(MaxActions) + 1
		actions := make([]Action, numActions)

		for j := 0; j < numActions; j++ {
			actionType := Add
			if j > 0 && rand.Float32() > 0.5 {
				actionType = Delete
			}

			action := Action{
				Type:  actionType,
				Key:   randomString(),
				Value: randomString(),
			}

			if actionType == Add {
				updateString(tr, action.Key, action.Value)
			} else if actionType == Delete {
				deleteString(tr, action.Key)
				action.Value = ""
			}
			actions[j] = action
		}

		rootHash := tr.Hash().String()
		examples[i] = Example{Actions: actions, Root: rootHash}
	}
	return examples
}

func main() {
	examples := generateExamples()
	file, err := os.Create("../tests/examples.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(examples)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
