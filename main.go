package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Block struct {
	Index     int
	Timestamp time.Time
	Data      []float64
	PrevHash  string
	Hash      string
	Mean      float64 // Mittelwert
	Median    float64 // Median
	SDRange   float64 // 2-SD-Bereich
}

type Blockchain []Block

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp.String() + fmt.Sprintf("%v", block.Data) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(prevBlock Block, data []float64) Block {
	var index int
	if prevBlock.Index == -1 {
		index = 0
	} else {
		index = prevBlock.Index + 1
	}
	newBlock := Block{
		Index:     index,
		Timestamp: time.Now(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
	}
	newBlock.Hash = calculateHash(newBlock)
	newBlock.calculateStatistics()
	return newBlock
}

func isValidBlock(block Block) bool {
	for _, value := range block.Data {
		if value < 0.000 || value > 1.000 {
			return false
		}
	}
	calculatedHash := calculateHash(block)
	if calculatedHash != block.Hash {
		return false
	}
	return true
}

// Funktion zur Berechnung des Mittelwerts
func mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// Funktion zur Berechnung des Medians
func median(data []float64) float64 {
	sort.Float64s(data)
	if len(data)%2 == 0 {
		middle := len(data) / 2
		return (data[middle-1] + data[middle]) / 2
	}
	return data[len(data)/2]
}

// Funktion zur Berechnung des 2-SD-Bereichs
func sdRange(data []float64) float64 {
	mean := mean(data)
	sumSquaredDiff := 0.0
	for _, value := range data {
		sumSquaredDiff += (value - mean) * (value - mean)
	}
	variance := sumSquaredDiff / float64(len(data))
	standardDeviation := math.Sqrt(variance)
	return 2 * standardDeviation
}

// Funktion zur Berechnung der Statistiken und Speicherung in einem Block
func (block *Block) calculateStatistics() {
	block.Mean = mean(block.Data)
	block.Median = median(block.Data)
	block.SDRange = sdRange(block.Data)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var blockchain Blockchain

	genesisBlock := Block{
		Index:     -1,
		Timestamp: time.Now(),
		Data:      nil,
		PrevHash:  "",
	}
	genesisBlock.Hash = calculateHash(genesisBlock)
	blockchain = append(blockchain, genesisBlock)

	for i := 0; i < 10; i++ {
		var data []float64
		for j := 0; j < 10; j++ {
			data = append(data, rand.Float64())
		}

		newBlock := generateBlock(blockchain[len(blockchain)-1], data)

		if isValidBlock(newBlock) {
			blockchain = append(blockchain, newBlock)
		} else {
			fmt.Println("Invalid block detected. Skipping...")
		}
	}

	for _, block := range blockchain {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp.String())
		fmt.Printf("Data: %v\n", block.Data)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Mean: %f\n", block.Mean)
		fmt.Printf("Median: %f\n", block.Median)
		fmt.Printf("2-SD Range: %f\n", block.SDRange)
		fmt.Println("------------------")
	}
}
