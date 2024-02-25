package models

import "gorm.io/gorm"

type MerkleTree struct {
	gorm.Model `swaggerignore:"true"`
	Matrix     [][]string
	// TODO: Define structure
}
