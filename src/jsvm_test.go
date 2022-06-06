// RedNaga / Tim Strazzere (c) 2018-*

package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestJSVMRun(t *testing.T) {
	jsFile := "pvvhnzyazwpzgkhv.js.original"
	file := filepath.Join("test", jsFile)
	// expected := DistilConfig{
	// 	Path:              "/pvvhnzyazwpzgkhv.js?PID=14CDB9B4-DE01-3FAA-AFF5-65BC2F771745",
	// 	XDistilAjax:       "twzvbatvrxzavsfzbzeyurav",
	// 	HeartbeatInterval: 270000,
	// }

	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("failed reading test file: %s", err)
	}

	err = Run(string(data))
	if err != nil {
		t.Fatalf("Failed performing a Run: %s", err)
	}
}
