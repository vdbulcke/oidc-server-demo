package oidcserver

import "testing"

func TestValidateConfig(t *testing.T) {

	filename := "../example/config.yaml"

	config, err := ParseConfig(filename)
	if err != nil {
		t.Fatal(err)
	}

	if !ValidateConfig(config) {
		t.Log("invalid config")
		t.Fail()
	}

	t.Log(config.VaultCryptoBackend)

}
