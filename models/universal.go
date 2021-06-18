// ======================
// This package is package to store model that could be accessed by any file
// One file could have multiple models defined using Struct, Slices
// ======================
package models

import (
	"github.com/rs/xid"
)

// Object to hold information of image url
type Link struct {
	ID  xid.ID `json:"id"`
	URL string `json:"url"`
}

// Wrapper for Link object
type Picture []*Link
