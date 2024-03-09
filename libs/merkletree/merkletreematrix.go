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

// IsFileCorrect if the file is correct and is not tampered.
// Returns true if the file is not tampered
// Returns false if the file is tampered
func IsFileCorrect(file *os.File, merkleProofs []string, wantedRootHash string, fileOrder int) (bool, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
		return false, err
	}

	hashFile := fmt.Sprintf("%x", h.Sum(nil))
	return isHashFileCorrect(hashFile, merkleProofs, wantedRootHash, fileOrder), nil
}

func isHashFileCorrect(hashFile string, merkleProofs []string, wantedRootHash string, fileOrder int) bool {
	// nodePositions := getNodesPositions(fileOrder, numberOfLeaves)
	reconstructedRootHash := reconstructRootHash(hashFile, merkleProofs, fileOrder)
	return reconstructedRootHash == wantedRootHash
}

func ReconstructRootHash(file *os.File, merkleProofs []string, fileOrder int) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
		return "", err
	}

	hashFile := fmt.Sprintf("%x", h.Sum(nil))
	// TODO: call getNodesPositions
	// nodePositions := getNodesPositions(fileOrder, numberOfLeaves)
	return reconstructRootHash(hashFile, merkleProofs, fileOrder), nil
}

// Reconstruct the root hash given the sha256 hash of the file to check and the
// merkleProofs that are needed to reconstruct the merkle tree root hash:
// 1. hash of the file you want to check integrity (sha256)
// 2. merkleProofs needed to reconstruct the root hash
// will return the root hash
func reconstructRootHash(hashFile string, merkleProofs []string, fileOrder int) string {
	pair := ""
	rootHash := ""
	hash := hashFile
	filePosition := 0
	for k, mp := range merkleProofs {
		mp = strings.TrimSpace(mp)
		fmt.Printf("hash %v: %s\n", k, hash)
		fmt.Printf("mp %v: %s\n", k, mp)

		filePosition = fileOrder % 2

		// 0: file is left of proof
		if filePosition == 0 {
			pair = fmt.Sprintf("%s%s", hash, mp)
		}
		// 1: file is right of proof
		if filePosition == 1 {
			pair = fmt.Sprintf("%s%s", mp, hash)
		}

		fmt.Printf("pair client %v: %s\n", k, pair)
		h := sha256.New()
		h.Write([]byte(pair))
		nextHash := fmt.Sprintf("%x", h.Sum(nil))
		hash = nextHash

		fileOrder /= 2
	}
	fmt.Printf("root hash: %s\n", hash)
	rootHash = hash
	return rootHash
}

func ComputeMerkleProof(file *os.File, merkleTree [][]string) []string {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	// odd element
	hashFile := fmt.Sprintf("%x", h.Sum(nil))

	fmt.Printf("hash file: %s\n", hashFile)
	merkleProof := createMerkleProofMatrix(hashFile, merkleTree)

	return merkleProof
}

// merkleProofs is ordered from bottom to top (from leaves towards root-hash)
func createMerkleProofMatrix(hashFile string, merkleTree [][]string) (merkleProofs []string) {
	// merkleProof := []string{}
	merkleTreeLeaves := merkleTree[0]
	hash := hashFile

	found := false
	for i := 0; i < len(merkleTreeLeaves) && !found; i++ {
		if hash == merkleTree[0][i] {
			found = true
			for j := 0; j < len(merkleTree)-1; j++ {
				if i%2 == 0 {
					merkleProofs = append(merkleProofs, merkleTree[j][i+1])
				} else {
					merkleProofs = append(merkleProofs, merkleTree[j][i-1])
				}
				i /= 2
			}
		}
	}

	return merkleProofs
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

	merkleTree := createMerkleTree(hashLeaves)

	return merkleTree[len(merkleTree)-1][0], nil
}

func ComputeMerkleTree(files ...*os.File) ([][]string, error) {
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

	merkleTree := createMerkleTree(hashLeaves)

	return merkleTree, nil
}

func createMerkleTree(hashLeaves []string) [][]string {
	merkleTree := [][]string{
		hashLeaves,
	}

	log2n := math.Log2(float64(len(hashLeaves)))
	treeDepth := math.Ceil(log2n) + 1

	for i := 0; i < int(treeDepth) && len(merkleTree[i]) > 1; i++ {
		if len(merkleTree[i])%2 != 0 {
			merkleTree[i] = append(merkleTree[i], merkleTree[i][len(merkleTree[i])-1])
		}
		var upperLevelMerkleTree []string
		for j := 0; j < len(merkleTree[i]); j += 2 {
			fmt.Printf("hash i: %v, j: %v = %s\n", i, j, merkleTree[i][j])
			fmt.Printf("hash i: %v, j: %v = %s\n", i, j+1, merkleTree[i][j+1])

			pair := fmt.Sprintf("%s%s", merkleTree[i][j], merkleTree[i][j+1])
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
