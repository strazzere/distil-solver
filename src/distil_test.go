// RedNaga / Tim Strazzere (c) 2018-*

package main

import (
	"strings"
	"testing"
)

func TestProofForward(t *testing.T) {
	input := "1554169628005:RzONDCgoMSnoo8zy3Pod"
	expected := "33:1554169628005:RzONDCgoMSnoo8zy3Pod"
	proof, err := workOnProof(input, 8)
	if err != nil {
		t.Errorf("Error trying to perform a workOnProof : %+v", err)
	}
	if strings.Compare(expected, proof) != 0 {
		t.Errorf("Expected %s but got %s", expected, proof)
	}
}
func TestProof(t *testing.T) {
	knownGood := "1554139566532:n3C64oACUKppKfAkD1n3"
	proof := getProof()

	// proof = workOnProof(proof, 8)

	if len(proof) != len(knownGood) {
		t.Errorf("Expected length of %d but got length of %d : %+v", len(knownGood), len(proof), proof)
	}
}

func TestGetProofQuery(t *testing.T) {
	expectedProof := "this_isnt_a_real_proof"
	expectedUA := "this_is_not_a_real_useragent"

	proofQuery, err := getProofQuery(expectedProof, expectedUA)
	if err != nil {
		t.Errorf("Unable to perform a getProofQuery : %+v", err)
	}

	proofFound := proofQuery.(map[string]interface{})["proof"].(string)
	if strings.Compare(expectedProof, proofFound) != 0 {
		t.Errorf("Expected %+v but got %+v", expectedProof, proofFound)
	}

	fp2 := proofQuery.(map[string]interface{})["fp2"]
	uaFound := fp2.(map[string]interface{})["userAgent"].(string)
	if strings.Compare(expectedUA, uaFound) != 0 {
		t.Errorf("Expected %+v but got %+v", expectedUA, uaFound)
	}
}
