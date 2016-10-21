package template

import (
	"bytes"
	"gopkg.in/yaml.v1"
	"os"
	"strings"
	"text/template"
)

func Unmarshal(in []byte, data interface{}, out interface{}) error {
	tmpl, err := template.New("yaml-template").Parse(string(in))
	if err != nil {
		return err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return err
	}

	if err := yaml.Unmarshal(b.Bytes(), out); err != nil {
		return err
	}

	return nil
}

func UnmarshalWithEnv(in []byte, out interface{}) error {
	envMap := map[string]string{}

	env := os.Environ()
	for _, value := range env {
		v := strings.Split(value, "=")
		envMap[v[0]] = v[1]
	}

	return Unmarshal(in, envMap, out)
}
