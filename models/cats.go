package models

// This package is package to store model for cats
// One file could have multiple models defined using Struct, Slices

// Cat type store an object for cat entity
type Cat struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Gender string `json:"gender"`
}

// Cats type store multiple Cat entities
type Cats []*Cat
