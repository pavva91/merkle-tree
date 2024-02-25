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

// TODO: Check go docs best practices

// IsFileCorrect if the file is correct and is not tampered.
// Returns true if the file is not tampered
// Returns false if the file is tampered
func IsFileCorrect(file *os.File, merkleProofs []string, wantedRootHash string) (bool, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
		return false, err
	}

	hashFile := fmt.Sprintf("%x", h.Sum(nil))
	return isHashFileCorrect(hashFile, merkleProofs, wantedRootHash), nil
}

func isHashFileCorrect(hashFile string, merkleProofs []string, wantedRootHash string) bool {
	reconstructedRootHash := reconstructRootHash(hashFile, merkleProofs)
	return reconstructedRootHash == wantedRootHash
}

func ReconstructRootHash(file *os.File, merkleProofs []string) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
		return "", err
	}

	hashFile := fmt.Sprintf("%x", h.Sum(nil))
	return reconstructRootHash(hashFile, merkleProofs), nil
}

// Reconstruct the root hash given the sha256 hash of the file to check and the
// merkleProofs that are needed to reconstruct the merkle tree root hash:
// 1. hash of the file you want to check integrity (sha256)
// 2. merkleProofs needed to reconstruct the root hash
// will return the root hash
func reconstructRootHash(hashFile string, merkleProofs []string) string {
	rootHash := ""
	hash := hashFile
	for k, mp := range merkleProofs {
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
	merkleTreeLeaves := merkleTree[0]
	hash := hashFile

	found := false
	for i := 0; i < len(merkleTreeLeaves) && !found; i++ {
		if hash == merkleTree[0][i] {
			found = true
			for j := 0; j < len(merkleTree)-1; j++ {
				if i%2 == 0 {
					merkleProof = append(merkleProof, merkleTree[j][i+1])
				} else {
					merkleProof = append(merkleProof, merkleTree[j][i-1])
				}
				i = i / 2
			}
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
			return "", err
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
			return [][]string{}, err
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

	for i := 0; i < int(treeDepth) && len(merkleTree[i]) > 1; i++ {
		fmt.Println(len(merkleTree[i]))

		if len(merkleTree[i])%2 != 0 {
			merkleTree[i] = append(merkleTree[i], merkleTree[i][len(merkleTree[i])-1])
		}
		var upperLevelMerkleTree []string
		for j := 0; j < len(merkleTree[i]); j = j + 2 {
			fmt.Printf("hash i: %v, j: %v = %s\n", i, j, merkleTree[i][j])
			fmt.Printf("hash i: %v, j: %v = %s\n", i, j+1, merkleTree[i][j+1])

			pair := calculateHashPair(merkleTree[i][j], merkleTree[i][j+1])
			fmt.Printf("pair %v: %s\n", j/2, pair)
			h := sha256.New()
			h.Write([]byte(pair))
			nextHash := fmt.Sprintf("%x", h.Sum(nil))
			upperLevelMerkleTree = append(upperLevelMerkleTree, nextHash)
		}
		merkleTree = append(merkleTree, upperLevelMerkleTree)
	}
	fmt.Println("root hash:", merkleTree[len(merkleTree)-1][0])
	return merkleTree
}
