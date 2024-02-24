package merkletree

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

// Verify if the file is correct and is not tampered.
// Returns true if the file is not tampered
// Returns false if the file is tampered
func Verify(hashFile string, merkleProofs []string, wantedRootHash string) bool {
	// TODO: Check if input string is hash (len)
	reconstructedRootHash := ReconstructRootHash(hashFile, merkleProofs)
	isCorrect := reconstructedRootHash == wantedRootHash
	return isCorrect
}

// Reconstruct the root hash given the hash of the file to check and the
// merkleProofs that are needed to reconstruct the merkle tree root hash:
// 1. hash of the file you want to check integrity
// 2. merkleProofs needed to reconstruct the root hash
// will return the root hash
func ReconstructRootHash(hashFile string, merkleProof []string) string {
	// TODO: Check if input string is hash (len)
	rootHash := ""
	hash := hashFile
	for k, mp := range merkleProof {
		mp = strings.TrimSpace(mp)
		fmt.Printf("hash %v: %s\n", k, hash)
		fmt.Printf("mp %v: %s\n", k, mp)
		pair := calculateHashPair(hash, mp)
		fmt.Printf("pair client %v: %s\n", k, pair)
		h := sha256.New()
		h.Write([]byte(pair))
		nextHash := fmt.Sprintf("%x", h.Sum(nil))
		hash = nextHash
	}
	fmt.Printf("root hash: %s\n", hash)
	rootHash = hash
	return rootHash
}

func calculateHashPair(h1 string, h2 string) string {
	pair := ""
	if h1 > h2 {
		pair = fmt.Sprintf("%s%s", h1, h2)

	} else {
		pair = fmt.Sprintf("%s%s", h2, h1)

	}
	return pair
}

func ComputeMerkleProof(file *os.File, merkleTree [][]string) []string {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	// odd element
	hashFile := fmt.Sprintf("%x", h.Sum(nil))

	fmt.Printf("hash file: %s\n", hashFile)
	merkleProof := createMerkleProof(hashFile, merkleTree)

	return merkleProof
}

func createMerkleProof(hashFile string, merkleTree [][]string) []string {
	merkleProof := []string{}
	hash := hashFile

	// NOTE: works 1 (nested cycles)
	// treeDepth := len(merkleTree)
	// fmt.Println("d", treeDepth)
	// for i := 0; i < len(merkleTree)-1; i++ {
	// 	for k, h := range merkleTree[i] {
	// 		fmt.Println("k", k)
	// 		if hash == h {
	//
	// 			if k%2 == 0 {
	// 				merkleProof = append(merkleProof, merkleTree[i][k+1])
	// 			} else {
	// 				merkleProof = append(merkleProof, merkleTree[i][k-1])
	// 			}
	// 			if i+1 < treeDepth {
	// 				hash = merkleTree[i+1][k/2]
	// 			}
	// 			break
	// 		}
	// 	}
	// }

	// NOTE: without nested cycles
	found := false
	for k, h := range merkleTree[0] {
		if hash == h {
			found = true
			for i := 0; i < len(merkleTree)-1; i++ {
				fmt.Println("k", k)
				if k%2 == 0 {
					merkleProof = append(merkleProof, merkleTree[i][k+1])
				} else {
					merkleProof = append(merkleProof, merkleTree[i][k-1])
				}
				k = k / 2
			}
		}
		if found {
			break
		}
	}

	return merkleProof
}

func ComputeRootHash(files ...*os.File) (string, error) {
	var hashLeaves []string
	for k, f := range files {
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

		hashFile := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Printf("hash file %v: %s\n", k+1, hashFile)
		hashLeaves = append(hashLeaves, hashFile)
	}

	merkleTree := createMerkleTreeAsMatrix(hashLeaves)

	return merkleTree[len(merkleTree)-1][0], nil
}

func ComputeMerkleTreeAsMatrix(files ...*os.File) ([][]string, error) {
	var hashLeaves []string
	for k, f := range files {
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

		hashFile := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Printf("hash file %v: %s\n", k+1, hashFile)
		hashLeaves = append(hashLeaves, hashFile)
	}

	merkleTree := createMerkleTreeAsMatrix(hashLeaves)

	return merkleTree, nil
}

func createMerkleTreeAsMatrix(hashLeaves []string) [][]string {
	merkleTree := [][]string{
		hashLeaves,
	}

	log2n := math.Log2(float64(len(hashLeaves)))
	treeDepth := math.Ceil(log2n) + 1

	for i := 0; i < int(treeDepth); i++ {
		fmt.Println(len(merkleTree[i]))

		if len(merkleTree[i]) == 1 {
			fmt.Println("root hash:", merkleTree[i][0])
			return merkleTree
		}

		if len(merkleTree[i])%2 != 0 {
			merkleTree[i] = append(merkleTree[i], merkleTree[i][len(merkleTree[i])-1])
		}
		var upperLevelMerkleTree []string
		for j := 0; j < len(merkleTree[i]); j = j + 2 {
			fmt.Printf("hash i: %v, j: %v = %s\n", i, j, merkleTree[i][j])
			fmt.Printf("hash i: %v, j: %v = %s\n", i, j+1, merkleTree[i][j+1])

			pair := calculateHashPair(merkleTree[i][j], merkleTree[i][j+1])
			if merkleTree[i][j] > merkleTree[i][j+1] {
				pair = fmt.Sprintf("%s%s", merkleTree[i][j], merkleTree[i][j+1])
				// pair = merkleTree[i][j] + merkleTree[i][j+1]
			} else {
				pair = fmt.Sprintf("%s%s", merkleTree[i][j+1], merkleTree[i][j])
				// pair = merkleTree[i][j+1] + merkleTree[i][j]
			}
			fmt.Printf("pair %v: %s\n", j/2, pair)
			h := sha256.New()
			h.Write([]byte(pair))
			nextHash := fmt.Sprintf("%x", h.Sum(nil))
			upperLevelMerkleTree = append(upperLevelMerkleTree, nextHash)
			// pair 0: f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe60dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d

		}
		merkleTree = append(merkleTree, upperLevelMerkleTree)
	}
	return merkleTree
}
