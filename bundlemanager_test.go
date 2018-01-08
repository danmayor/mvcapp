/*
	Digivance MVC Application Framework - Unit Tests
	Bundle Manager Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.3.0 compatibility of bundlemanager.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in bundlemanager.go
*/

package mvcapp_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/digivance/mvcapp"
)

// TestBundleManager_CreateBundle ensures that the BundleManager.CreateBundle method operates as expected
func TestBundleManager_CreateBundle(t *testing.T) {
	bundleManager := mvcapp.NewBundleManager()
	if err := bundleManager.CreateBundle("TestBundle", "text/css", []string{}); err != nil {
		t.Errorf("Failed to create new bundle: %s", err)
	}

	if err := bundleManager.CreateBundle("TestBundle", "text/css", []string{}); err == nil {
		t.Errorf("Failed to prevent duplicate bundle creation: %s", err)
	}
}

// TestBundleManager_BuildBundle ensures that the BundleManager.BuildBundle and BundleManager.RebuildBundle method operates
// as expected
func TestBundleManager_BuildBundle(t *testing.T) {
	expectedResult := "html,body{font-size:14px}a{color:blue}div{border:0}"
	filename := []string{
		fmt.Sprintf("%s/a.css", mvcapp.GetApplicationPath()),
		fmt.Sprintf("%s/b.css", mvcapp.GetApplicationPath()),
		fmt.Sprintf("%s/c.css", mvcapp.GetApplicationPath()),
	}

	if err := ioutil.WriteFile(filename[0], []byte("html, body { font-size: 14px; }\n"), 0644); err != nil {
		t.Fatalf("Failed to write a.css: %s", err)
	}
	defer os.RemoveAll(filename[0])

	if err := ioutil.WriteFile(filename[1], []byte("\ta {color: blue; }\n"), 0644); err != nil {
		t.Fatalf("Failed to write a.css: %s", err)
	}
	defer os.RemoveAll(filename[1])

	if err := ioutil.WriteFile(filename[2], []byte("\n\ndiv   {    border:    0;\t\t}\n\n"), 0644); err != nil {
		t.Fatalf("Failed to write a.css: %s", err)
	}
	defer os.RemoveAll(filename[2])

	bundleManager := mvcapp.NewBundleManager()
	if err := bundleManager.BuildBundle("fail.js"); err == nil {
		t.Error("Failed to prevent build of unknown bundle")
	}

	if err := bundleManager.CreateBundle("fail.js", "text/javascript", []string{}); err != nil {
		t.Errorf("Failed to create empty bundle for testing (this is a good thing?): %s", err)
	}

	if err := bundleManager.BuildBundle("fail.js"); err == nil {
		t.Error("Failed to prevent build of empty bundle")
	}

	bundleManager.Bundles["fail.js"].Files = append([]string{}, "notfound.js")
	if err := bundleManager.BuildBundle("fail.js"); err == nil {
		t.Error("Failed to prevent build with missing filename")
	}

	if err := bundleManager.RemoveBundle("fail.js"); err != nil {
		t.Errorf("Failed to remove failing bundle: %s", err)
	}

	if err := bundleManager.RemoveBundle("fail.js"); err == nil {
		t.Error("Failed to prevent removing non existent bundle")
	}

	defer os.RemoveAll("bundle")
	if err := bundleManager.CreateBundle("styles.css", "text/css", filename); err != nil {
		t.Fatalf("Failed to create bundle: %s", err)
	}

	if err := bundleManager.BuildBundle("styles.css"); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("bundle/styles.css")

	if err := bundleManager.BuildBundle("«╗╣╟/failme"); err == nil {
		t.Error("Failed to fail when using an invalid bundle file name")
	}

	bundleFilename := fmt.Sprintf("%s/bundle/styles.css", mvcapp.GetApplicationPath())
	data, err := ioutil.ReadFile(bundleFilename)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if content != expectedResult {
		t.Fatalf("Failed to build test bundle:\n> Expected: %s\n Received: %s", expectedResult, content)
	}

	bundleName := "styles.css"
	if err := bundleManager.RebuildBundle(bundleName); err != nil {
		t.Fatalf("Failed to not rebuild test bundle: %s", err)
	}

	bundleManager.Bundles[bundleName].BuildDate = time.Time{}

	if err := bundleManager.RebuildBundle(bundleName); err != nil {
		t.Fatalf("Failed to not rebuild test bundle: %s", err)
	}

	bundleManager.Bundles[bundleName].BuildDate = time.Time{}.Add(48 * time.Hour)

	if err := bundleManager.RebuildBundle(bundleName); err != nil {
		t.Fatalf("Failed to not rebuild test bundle: %s", err)
	}
}

// TestBundleManager_BuildAllBundles ensures that the BundleManager.BuildAllBundles and BundleManager.RebuildAllBundles
// methods operate as expected
func TestBundleManagerBuildAllBundles(t *testing.T) {
	expectedResult := "html,body{font-size:14px}a{color:blue}div{border:0}"
	filename := []string{
		fmt.Sprintf("%s/a.css", mvcapp.GetApplicationPath()),
		fmt.Sprintf("%s/b.css", mvcapp.GetApplicationPath()),
		fmt.Sprintf("%s/c.css", mvcapp.GetApplicationPath()),
	}

	if err := ioutil.WriteFile(filename[0], []byte("html, body { font-size: 14px; }\n"), 0644); err != nil {
		t.Fatalf("Failed to write a.css: %s", err)
	}
	defer os.RemoveAll(filename[0])

	if err := ioutil.WriteFile(filename[1], []byte("\ta {color: blue; }\n"), 0644); err != nil {
		t.Fatalf("Failed to write a.css: %s", err)
	}
	defer os.RemoveAll(filename[1])

	if err := ioutil.WriteFile(filename[2], []byte("\n\ndiv   {    border:    0;\t\t}\n\n"), 0644); err != nil {
		t.Fatalf("Failed to write a.css: %s", err)
	}
	defer os.RemoveAll(filename[2])

	bundleManager := mvcapp.NewBundleManager()
	defer os.RemoveAll("bundle")
	if err := bundleManager.CreateBundle("styles.css", "text/css", filename); err != nil {
		t.Fatalf("Failed to create bundle: %s", err)
	}

	if err := bundleManager.BuildAllBundles(); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("bundle/styles.css")

	bundleFilename := fmt.Sprintf("%s/bundle/styles.css", mvcapp.GetApplicationPath())
	data, err := ioutil.ReadFile(bundleFilename)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if content != expectedResult {
		t.Fatalf("Failed to build test bundle:\n> Expected: %s\n Received: %s", expectedResult, content)
	}

	if err := bundleManager.RebuildAllBundles(); err != nil {
		t.Fatalf("Failed to not rebuild all bundles: %s", err)
	}

	bundleName := "styles.css"
	bundleManager.Bundles[bundleName].BuildDate = time.Time{}
	if err := bundleManager.RebuildAllBundles(); err != nil {
		t.Fatalf("Failed to rebuild all bundles: %s", err)
	}

	bundleManager.Bundles[bundleName].BuildDate = time.Time{}.Add(48 * time.Hour)
	if err := bundleManager.RebuildAllBundles(); err != nil {
		t.Fatalf("Failed to rebuild all bundles: %s", err)
	}
}
