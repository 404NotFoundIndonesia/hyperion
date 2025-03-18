package obfuscator_test

import (
	"embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	obf "hyperion/backend/obfuscator"
	"testing"
)

//go:embed examples/*
var Examples embed.FS

func TestNewObfuscator(t *testing.T) {
	config := &obf.Config{}
	obfuscator := obf.NewObfuscator(config)

	assert.NotNil(t, obfuscator)
	assert.IsType(t, &obf.Obfuscator{}, obfuscator)
}

func Test_obfuscator_Obfuscate(t *testing.T) {
	files, err := Examples.ReadDir("examples")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}
