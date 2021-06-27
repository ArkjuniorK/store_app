// ======================
// This package is package to store model for cats
// One file could have multiple models defined using Struct, Slices
// ======================

package models

import (
	"time"

	"github.com/rs/xid"
)

// Cat type store an object for cat entity
type Cat struct {
	ID      xid.ID    `json:"id"`
	Name    string    `json:"name"`
	Variety string    `json:"variety"`
	Gender  string    `json:"gender"`
	Age     int16     `json:"age"`
	Address string    `json:"address"`
	ZipCode int16     `json:"zip_code"`
	Create  time.Time `json:"created_at"`
	Update  time.Time `json:"updated_at"`
	Image   *Picture  `json:"image"`
}

// Cat type to store data as map
// so it could be iterate using for loop
type CatMap map[string]interface{}

// Cats type store multiple Cat entities
type Cats []*Cat
