package stubs

import "github.com/pavva91/merkle-tree/server/internal/models"

type MerkleTreeService struct {
	CreateFn            func() (*models.MerkleTree, error)
	IsValidFn           func() bool
	CreateMerkleProofFn func(string) ([]string, error)
}

func (stub MerkleTreeService) Create() (*models.MerkleTree, error) {
	return stub.CreateFn()
}

func (stub MerkleTreeService) IsValid() bool {
	return stub.IsValidFn()
}

func (stub MerkleTreeService) CreateMerkleProof(filePath string) ([]string, error) {
	return stub.CreateMerkleProofFn(filePath)
}
