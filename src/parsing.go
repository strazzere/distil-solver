// RedNaga / Tim Strazzere (c) 2018-*

package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/robertkrimen/otto"
)

const fingerPrintRegex = `FingerprintWrapper\([\{a-zA-Z\.\?\=\-\"\,\:0-9\/\_\}]+\)`

// parse for;
// FingerprintWrapper({path:"/pvvhnzyazwpzgkhv.js?PID=14CDB9B4-DE01-3FAA-AFF5-65BC2F771745",ajax_header:"twzvbatvrxzavsfzbzeyurav",interval:27e4})}
// Then build a struct of that shitty, non-revolable
func parseFingerprints(data []byte) (DistilConfig, error) {
	distilConfig := DistilConfig{}
	// get match
	rx := regexp.MustCompile(fingerPrintRegex)
	match := rx.FindStringSubmatch(string(data))
	if len(match) != 1 {
		log.Printf("Error attempting to find fingerprint, regex found : %+v", match)
		// TODO : This likely shouldn't throw an error to sentry?
		return distilConfig, errors.New("unexpected fingerprints resolved")
	}

	// strip down
	subData := match[0][len("FingerprintWrapper(") : len(match[0])-1]

	vm := otto.New()
	val, err := vm.Run(fmt.Sprintf("a = %s", subData))
	if err != nil {
		return distilConfig, err
	}

	temp, err := val.Object().Get("path")
	distilConfig.Path = temp.String()
	if err != nil {
		return distilConfig, err
	}

	temp, err = val.Object().Get("ajax_header")
	if err != nil {
		return distilConfig, err
	}
	distilConfig.XDistilAjax = temp.String()

	temp, err = val.Object().Get("interval")
	if err != nil {
		return distilConfig, err
	}
	distilConfig.HeartbeatInterval, err = temp.ToInteger()
	if err != nil {
		return distilConfig, err
	}

	return distilConfig, nil
}
