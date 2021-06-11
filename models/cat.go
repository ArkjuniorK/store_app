// ======================
// This package is package to store model for cats
// One file could have multiple models defined using Struct, Slices
// ======================

package models

import "github.com/google/uuid"

// Cat type store an object for cat entity
type Cat struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Type   string    `json:"type"`
	Gender string    `json:"gender"`
	Image  Picture   `json:"image"`
}

// Cats type store multiple Cat entities
type Cats []*Cat
