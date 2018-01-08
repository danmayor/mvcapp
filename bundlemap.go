/*
	Digivance MVC Application Framework
	Bundle Map Object
	Dan Mayor (dmayor@digivance.com)

	This file defines a bundle map object used by the bundle manager system
*/

package mvcapp

import "time"

// BundleMap represents the slice of filenames, the mime type and last build
// date & time of a "Content Bundle"
type BundleMap struct {
	// Files is a slice of the full path and file name of the files to include
	// in this content bundle
	Files []string

	// MimeType is the MIME type of the files in this bundle and is used to execute
	// the appropriate minification methods in the bundle manager
	MimeType string

	// BuildDate is the time and date when this bundle was created, the bundle
	// manager uses this to determine if a rebuild is warranted
	BuildDate time.Time
}

// NewBundleMap returns a new BundleMap from the provided mime type and filename slice
func NewBundleMap(mimeType string, files []string) *BundleMap {
	return &BundleMap{
		Files:     files,
		MimeType:  mimeType,
		BuildDate: time.Time{},
	}
}
