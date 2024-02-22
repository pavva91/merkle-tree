package services

import (
	"strconv"

	"github.com/pavva91/merkle-tree/server/internal/models"
	"github.com/pavva91/merkle-tree/server/internal/repositories"
)

var (
	MerkleTree Merkletreer = merkleTree{}
)

type Merkletreer interface {
	Create() (*models.MerkleTree, error)
	GetProof(filename string) ([]string, error)
}

type merkleTree struct{}

func (s merkleTree) Create() (*models.MerkleTree, error) {
	var merkleTree *models.MerkleTree
	// TODO: Get files from local storage
	// TODO: Call library to calculate merkleTree
	// TODO: Save All Merkle Tree in DB
	return repositories.MerkleTree.Create(merkleTree)
}

func (s merkleTree) GetProof(filename string) ([]string, error) {
	// TODO: Get Merkle Tree from DB
	// TODO: Calculate proof slice array
	var merkleTree *models.MerkleTree
	strID := strconv.Itoa(int(filename))
	merkleTree, err := repositories.MerkleTree.GetByID(strID)
	if err != nil {
		return nil, err
	}
	return merkleTree, nil
}
