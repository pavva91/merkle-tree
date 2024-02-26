package services

import (
	"log"
	"os"

	"github.com/pavva91/merkle-tree/libs/merkletree"
	"github.com/pavva91/merkle-tree/server/config"
	"github.com/pavva91/merkle-tree/server/internal/models"
	"github.com/pavva91/merkle-tree/server/internal/repositories"
)

var (
	MerkleTree MerkleTreer = merkleT{}
)

type MerkleTreer interface {
	Create() (*models.MerkleTree, error)
	IsValid() bool
	CreateMerkleProof(filePath string) ([]string, error)
}

type merkleT struct{}

func (s merkleT) Create() (*models.MerkleTree, error) {
	var fFiles []*os.File
	files, err := os.ReadDir(config.Values.Server.UploadFolder)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, f := range files {
		filePath := config.Values.Server.UploadFolder + "/" + f.Name()
		file, err := os.Open(filePath)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fFiles = append(fFiles, file)
		defer file.Close()
	}

	merkleTree := &models.MerkleTree{}
	merkleTree.Matrix, err = merkletree.ComputeMerkleTree(fFiles...)
	if err != nil {
		return nil, err
	}

	merkleTree, err = repositories.MerkleTree.Save(merkleTree)
	if err != nil {
		return nil, err
	}

	return merkleTree, nil
}

func (s merkleT) IsValid() bool {
	merkleTree, err := repositories.MerkleTree.Get()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return len(merkleTree.Matrix) != 0
}

func (s merkleT) CreateMerkleProof(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}
	defer file.Close()
	merkleTree, err := repositories.MerkleTree.Get()
	return merkletree.ComputeMerkleProof(file, merkleTree.Matrix), err
}
