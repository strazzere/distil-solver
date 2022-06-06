// RedNaga / Tim Strazzere (c) 2018-*

package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseCollectAndGo(t *testing.T) {
	jsFile := "ygaagxtnmspobnpr.js"
	file := filepath.Join("test", jsFile)
	expected := DistilConfig{
		Path:              "/ygaagxtnmspobnpr.js?PID=59A20418-3B6E-38B2-97BA-10E33819138C",
		XDistilAjax:       "ucarrzdv",
		HeartbeatInterval: 270000,
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("failed reading test file: %s", err)
	}

	distilConfig, err := parseFingerprints(data)
	if err != nil {
		t.Fatalf("Failed performing a parseFingerprints: %s", err)
	}

	if strings.Compare(expected.Path, distilConfig.Path) != 0 {
		t.Errorf("Expected %+v but got %+v", expected.Path, distilConfig.Path)
	}

	if strings.Compare(expected.XDistilAjax, distilConfig.XDistilAjax) != 0 {
		t.Errorf("Expected %+v but got %+v", expected.XDistilAjax, distilConfig.XDistilAjax)
	}

	if expected.HeartbeatInterval != distilConfig.HeartbeatInterval {
		t.Errorf("Expected %+v but got %+v", expected.HeartbeatInterval, distilConfig.HeartbeatInterval)
	}
}

func TestParseFingerprints(t *testing.T) {
	jsFile := "pvvhnzyazwpzgkhv.js.original"
	file := filepath.Join("test", jsFile)
	expected := DistilConfig{
		Path:              "/pvvhnzyazwpzgkhv.js?PID=14CDB9B4-DE01-3FAA-AFF5-65BC2F771745",
		XDistilAjax:       "twzvbatvrxzavsfzbzeyurav",
		HeartbeatInterval: 270000,
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("failed reading test file: %s", err)
	}

	distilConfig, err := parseFingerprints(data)
	if err != nil {
		t.Fatalf("Failed performing a parseFingerprints: %s", err)
	}

	if strings.Compare(expected.Path, distilConfig.Path) != 0 {
		t.Errorf("Expected %+v but got %+v", expected.Path, distilConfig.Path)
	}

	if strings.Compare(expected.XDistilAjax, distilConfig.XDistilAjax) != 0 {
		t.Errorf("Expected %+v but got %+v", expected.XDistilAjax, distilConfig.XDistilAjax)
	}

	if expected.HeartbeatInterval != distilConfig.HeartbeatInterval {
		t.Errorf("Expected %+v but got %+v", expected.HeartbeatInterval, distilConfig.HeartbeatInterval)
	}
}
