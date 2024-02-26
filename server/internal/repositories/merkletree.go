package repositories

import "github.com/pavva91/merkle-tree/server/internal/models"

var (
	MerkleTree MerkleTreer = merkleT{}
	// TODO: will use DB (now is mocked in memory)
	merkleTree = &models.MerkleTree{
		Matrix: [][]string{},
	}
)

type MerkleTreer interface {
	Save(*models.MerkleTree) (*models.MerkleTree, error)
	Get() (*models.MerkleTree, error)
}

type merkleT struct{}

func (r merkleT) Save(merkleTreeIn *models.MerkleTree) (*models.MerkleTree, error) {
	merkleTree = merkleTreeIn
	return merkleTree, nil
}

func (r merkleT) Get() (*models.MerkleTree, error) {
	// TODO: Get from real DB
	return merkleTree, nil
}
