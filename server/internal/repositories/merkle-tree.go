package repositories

import (
	"github.com/pavva91/merkle-tree/server/internal/db"
	"github.com/pavva91/merkle-tree/server/internal/models"
)

var (
	MerkleTree MerkleTreeer = merkleTree{}
)

type MerkleTreeer interface {
	Create(merkleTree *models.MerkleTree) (*models.MerkleTree, error)
	UpdateMerkleTree(merkleTree *models.MerkleTree) (*models.MerkleTree, error)
	GetByID(id string) (*models.MerkleTree, error)
}

type merkleTree struct{}

func (r merkleTree) Create(merkleTree *models.MerkleTree) (*models.MerkleTree, error) {
	err := db.ORM.GetDB().Create(&merkleTree).Error
	if err != nil {
		return nil, err
	}
	return merkleTree, nil
}

func (r merkleTree) UpdateMerkleTree(merkleTree *models.MerkleTree) (*models.MerkleTree, error) {
	err := db.ORM.GetDB().Updates(&merkleTree).Error
	if err != nil {
		return nil, err
	}
	return merkleTree, nil
}

func (r merkleTree) GetByID(id string) (*models.MerkleTree, error) {
	var merkleTree *models.MerkleTree
	err := db.ORM.GetDB().Where("id = ?", id).First(&merkleTree).Error
	if err != nil {
		return nil, err
	}
	return merkleTree, nil
}
