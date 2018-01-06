/*
	Digivance MVC Application Framework
	Content Bundle Manager
	Dan Mayor (dmayor@digivance.com)

	This file defines the content bundle manager object for the Digivance mvcapp package
*/

package mvcapp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

// BundleManager is an object that we register our content bundles with and is used
// to compile said content bundles.
type BundleManager struct {
	// Bundles are a collection of filenames that are compiled into a single deliverable
	Bundles  map[string]*BundleMap
	Minifier *minify.M
}

// NewBundleManager returns a new instance of the bundle manager object
func NewBundleManager() *BundleManager {
	rtn := &BundleManager{
		Bundles:  make(map[string]*BundleMap, 0),
		Minifier: minify.New(),
	}

	rtn.Minifier.AddFunc("text/css", css.Minify)
	rtn.Minifier.AddFunc("text/javascript", js.Minify)

	return rtn
}

// CreateBundle is used to register a bundle by name, mime type and slice of filenames
// Once registered, bundles can be compiled using the BuildBundle method referencing
// the provided bundleName
func (bundleManager *BundleManager) CreateBundle(bundleName string, mimeType string, bundledFiles []string) error {
	if bundleManager.Bundles[bundleName] != nil {
		return errors.New("Failed to create new content bundle, there is already a bundle created using this name")
	}

	bundleManager.Bundles[bundleName] = &BundleMap{
		Files:    bundledFiles,
		MimeType: mimeType,
	}

	return nil
}

// RemoveBundle is used to remove a bundle from the manager by name.
func (bundleManager *BundleManager) RemoveBundle(bundleName string) error {
	if bundleManager.Bundles[bundleName] == nil {
		return fmt.Errorf("Failed to remove bundle, no bundles found for %s", bundleName)
	}

	bundleManager.Bundles[bundleName] = nil
	if bundleManager.Bundles[bundleName] != nil {
		return fmt.Errorf("Failed to remove bundle, %s seems to still exist", bundleName)
	}

	return nil
}

// doBuild is really just a micro-optimization, it can be used so that "ALL" methods
// of this object only have to iterate the bundles map once
func (bundleManager *BundleManager) doBuild(bundleMap *BundleMap, bundleName string) error {
	if bundleMap == nil {
		return errors.New("Failed to build bundle, none found for provided name")
	}

	if len(bundleMap.Files) <= 0 {
		return errors.New("Failed to build bundle, no files registered")
	}

	bundlePath := fmt.Sprintf("%s/bundle", GetApplicationPath())
	if err := os.Mkdir(bundlePath, 0644); err != nil && !strings.HasSuffix(err.Error(), "file already exists.") {
		return err
	}

	data := []byte{}

	for _, filename := range bundleMap.Files {
		if !strings.HasPrefix(filename, GetApplicationPath()) {
			// Is tested successfully, hard to demonstrate because of scoping when testing
			// (e.g. the GetApplicationPath is different between the unit test and the lib)
			filename = fmt.Sprintf("%s/%s", GetApplicationPath(), filename)
		}

		contentData, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("Failed to bundle %s : %s", filename, err)
		}

		data = append(data, contentData...)
	}

	bundleFilename := fmt.Sprintf("%s/bundle/%s", GetApplicationPath(), bundleName)
	if err := os.RemoveAll(bundleFilename); err != nil {
		// os.RemoveAll is normally pretty quiet, haven't tested as this is a very critical
		// failure that isn't very likely to ever execute.
		return fmt.Errorf("Failed to remove existing bundle file: %s", err)
	}

	bundleFile, err := os.Create(bundleFilename)
	if err != nil {
		return fmt.Errorf("Failed to create new bundle file: %s", err)
	}

	reader := bytes.NewReader(data)
	writer := bufio.NewWriter(bundleFile)
	defer func() {
		writer.Flush()
		bundleFile.Close()
	}()

	if err := bundleManager.Minifier.Minify(bundleMap.MimeType, writer, reader); err != nil {
		return fmt.Errorf("Failed to minify and write the bundle file: %s", err)
	}

	bundleMap.BuildDate = time.Now()

	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("Failed to launch build bundle: %s", err)
		}

		return err
	}

	return nil
}

// BuildBundle is used to compile the registered bundle defined by bundleName.
// This method will delete the existing bundle file if it exists and replacing
// with a newly built copy
func (bundleManager *BundleManager) BuildBundle(bundleName string) error {
	bundleMap := bundleManager.Bundles[bundleName]
	return bundleManager.doBuild(bundleMap, bundleName)
}

// BuildAllBundles is used to build all of the currently registered content bundles
func (bundleManager *BundleManager) BuildAllBundles() error {
	for buildName, buildMap := range bundleManager.Bundles {
		if err := bundleManager.doBuild(buildMap, buildName); err != nil {
			return err
		}
	}

	return nil
}

// RebuildBundle is used to compare the last modified times of the files in a content
// bundle to the creation time of this content bundle, if files have been modified
// since this bundle was built it will be built a new. Returns nil if no need to build
func (bundleManager *BundleManager) RebuildBundle(bundleName string) error {
	bundleMap := bundleManager.Bundles[bundleName]

	if bundleMap.BuildDate.IsZero() {
		return bundleManager.doBuild(bundleMap, bundleName)
	}

	for _, filename := range bundleMap.Files {
		si, err := os.Stat(filename)
		if err != nil {
			return err
		}

		if si.ModTime().After(bundleMap.BuildDate) {
			return bundleManager.doBuild(bundleMap, bundleName)
		}
	}

	return nil
}

// RebuildAllBundles is used to loop through all of the registered content bundles
// and compare file modification dates to the build time of the bundle. If files
// have been modified since this bundle was built, it will be built a new
func (bundleManager *BundleManager) RebuildAllBundles() error {
	for bundleName, bundleMap := range bundleManager.Bundles {
		if bundleMap.BuildDate.IsZero() {
			if err := bundleManager.doBuild(bundleMap, bundleName); err != nil {
				return err
			}
		}

		for _, filename := range bundleMap.Files {
			si, err := os.Stat(filename)
			if err != nil {
				return err
			}

			if si.ModTime().After(bundleMap.BuildDate) {
				if err := bundleManager.doBuild(bundleMap, bundleName); err != nil {
					return err
				}

				break
			}
		}
	}

	return nil
}
