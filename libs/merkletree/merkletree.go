package merkletree

import (
	"crypto/sha256"
	"fmt"
	"math"
	"slices"
)

// Leaf Node has only Value. LeftNode and RightNode will be nil
// Root will have Value and LefNode and RightNode will not be null (for n>1)
type BinaryNode struct {
	Value     string
	LeftNode  *BinaryNode
	RightNode *BinaryNode
}

type MerkleTree struct {
	RootHashNode *BinaryNode
	HashLeaves   []string
}

// merkleProofs is ordered from bottom to top (from leaves towards root-hash)
func createMerkleProof(hashLeaf string, merkleTree MerkleTree) (merkleProofs []string) {
	iProofs := []int{}
	indexLeaf := -1
	for k, v := range merkleTree.HashLeaves {
		if v == hashLeaf {
			indexLeaf = k
		}
	}

	// not found
	if indexLeaf == -1 {
		return merkleProofs
	}

	// 0 : Left
	// 1 : Right
	log2n := math.Log2(float64(len(merkleTree.HashLeaves)))
	lengthMerkleProofs := math.Ceil(log2n)

	for i := 0; i < int(lengthMerkleProofs); i++ {
		iProof := 1 - (indexLeaf % 2)
		iProofs = append(iProofs, iProof)
		indexLeaf = indexLeaf / 2
	}

	slices.Reverse(iProofs)

	nextNode := merkleTree.RootHashNode

	for _, v := range iProofs {
		if v == 0 {
			merkleProofs = append(merkleProofs, nextNode.LeftNode.Value)
			nextNode = nextNode.RightNode
		}
		if v == 1 {
			merkleProofs = append(merkleProofs, nextNode.RightNode.Value)
			nextNode = nextNode.LeftNode
		}
	}

	slices.Reverse(merkleProofs)

	return merkleProofs
}

func createLeavesNodes(hashLeaves []string) []*BinaryNode {
	leavesNodes := []*BinaryNode{}
	n := len(hashLeaves)
	if n%2 != 0 {
		hashLeaves = append(hashLeaves, hashLeaves[n-1])
	}
	for i := 0; i < n; i = i + 1 {

		leafNode := *&BinaryNode{
			Value:     hashLeaves[i],
			LeftNode:  nil,
			RightNode: nil,
		}
		leavesNodes = append(leavesNodes, &leafNode)
	}
	return leavesNodes
}

// recursive function
// user input: leavesNodes
// return: rootHashNode
func calcMT(hashNodes []*BinaryNode) []*BinaryNode {
	higherLevelNodes := []*BinaryNode{}
	n := len(hashNodes)
	if n == 1 {
		return hashNodes
	}
	if n%2 != 0 {
		hashNodes = append(hashNodes, hashNodes[n-1])
	}
	for i := 0; i < n; i = i + 2 {
		node1 := hashNodes[i]
		node2 := hashNodes[i+1]
		newNode := *&BinaryNode{}
		pair := calculateHashPair(node1.Value, node2.Value)

		h := sha256.New()
		h.Write([]byte(pair))
		nextHash := fmt.Sprintf("%x", h.Sum(nil))

		newNode.Value = nextHash
		newNode.LeftNode = node1
		newNode.RightNode = node2
		higherLevelNodes = append(higherLevelNodes, &newNode)
	}

	return calcMT(higherLevelNodes)
}

// returns a MerkleTree
// rootHash is merkleTree.RootHashNode.Value
func CalculateMerkleTree(hashLeaves []string) (merkleTree MerkleTree) {
	merkleTree.HashLeaves = hashLeaves
	merkleArray := calcMT(createLeavesNodes(hashLeaves))
	merkleTree.RootHashNode = merkleArray[0]
	return merkleTree
}
