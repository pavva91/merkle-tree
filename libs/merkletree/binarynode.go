package merkletree

import (
	"crypto/sha256"
	"fmt"
)

// Leaf Node has only HashValue. LeftNode and RightNode will be nil
// RootHash will have HashValue and LefNode and RightNode will not be null (for n>1)
type BinaryNode struct {
	HashValue string
	LeftNode  *BinaryNode
	RightNode *BinaryNode
}

func createLeavesNodes(hashLeaves []string) []*BinaryNode {
	leavesNodes := []*BinaryNode{}
	n := len(hashLeaves)
	if n%2 != 0 {
		hashLeaves = append(hashLeaves, hashLeaves[n-1])
	}
	for i := 0; i < n; i = i + 1 {

		leafNode := *&BinaryNode{
			HashValue: hashLeaves[i],
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
		pair := calculateHashPair(node1.HashValue, node2.HashValue)

		h := sha256.New()
		h.Write([]byte(pair))
		nextHash := fmt.Sprintf("%x", h.Sum(nil))

		newNode.HashValue = nextHash
		newNode.LeftNode = node1
		newNode.RightNode = node2
		higherLevelNodes = append(higherLevelNodes, &newNode)
	}

	return calcMT(higherLevelNodes)
}

func CalculateMerkleTree(hashLeaves []string) (merkleTree *BinaryNode) {
	merkleArray := calcMT(createLeavesNodes(hashLeaves))
	merkleTree = merkleArray[0]
	return merkleTree
}
