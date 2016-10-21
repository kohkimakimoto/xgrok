package template

import (
	"os"
	"testing"
)

type TestStruct struct {
	Key1 string `yaml:"key1"`
	Key2 string `yaml:"key2"`
	Key3 string `yaml:"key3"`
	Key4 string `yaml:"key4"`
}

func TestUnmarshal(t *testing.T) {
	var out TestStruct
	if err := Unmarshal([]byte(`
key1: aaa
key2: bbb
key3: ccc
key4: ddd
`), nil, &out); err != nil {
		t.Error(err)
	}

	if out.Key1 != "aaa" {
		t.Errorf("invalid resutl: %v", out.Key1)
	}

	t.Logf("%v", out)
}

func TestUnmarshalWithEnv(t *testing.T) {
	var out TestStruct
	if err := UnmarshalWithEnv([]byte(`
key1: {{.HOME}}
key2: bbb
key3: ccc
key4: ddd
`), &out); err != nil {
		t.Error(err)
	}

	if out.Key1 != os.Getenv("HOME") {
		t.Errorf("invalid resutl: %v", out.Key1)
	}

	t.Logf("%v", out)
}
