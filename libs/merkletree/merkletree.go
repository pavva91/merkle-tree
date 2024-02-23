package merkletree

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func ComputeMerkleProof(file *os.File, merkleTree [][]string) []string {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	hashFile := fmt.Sprintf("hash file: %x\n", h.Sum(nil))
	fmt.Println(hashFile)

	// odd element
	merkleProof := createMerkleProof(hashFile, merkleTree)

	return merkleProof
}

func createMerkleProof(hashFile string, merkleTree [][]string) []string {
	merkleProof := []string{}
	treeDepth := len(merkleTree)
	hash := hashFile
	fmt.Println(hash)
	fmt.Println("d", treeDepth)

	// NOTE: works 1 (nested cycles)
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
	for _, f := range files {
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

		hashFile := fmt.Sprintf("hash file: %x\n", h.Sum(nil))
		fmt.Println(hashFile)
		hashLeaves = append(hashLeaves, hashFile)
	}

	merkleTree := createMerkleTreeAsMatrix(hashLeaves)

	return merkleTree[len(merkleTree)-1][0], nil
}

func ComputeMerkleTreeAsMatrix(files ...*os.File) ([][]string, error) {
	var hashLeaves []string
	for _, f := range files {
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

		hashFile := fmt.Sprintf("hash file: %x\n", h.Sum(nil))
		fmt.Println(hashFile)
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

			pair := merkleTree[i][j] + merkleTree[i][j+1]
			h := sha256.New()
			h.Write([]byte(pair))
			nextHash := fmt.Sprintf("%x", h.Sum(nil))
			upperLevelMerkleTree = append(upperLevelMerkleTree, nextHash)
		}
		merkleTree = append(merkleTree, upperLevelMerkleTree)
	}
	return merkleTree
}
