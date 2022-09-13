package modules

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(path string) ([]Module, error) {
	var modules []Module
	if len(path) > 0 {
		f, err := os.Open(path)
		if err != nil {
			return modules, err
		}
		if err := json.NewDecoder(f).Decode(&modules); err != nil {
			return modules, fmt.Errorf("modulesConfPath(%s) is invalid JSON", path)
		}
		for i, _ := range modules {
			if modules[i].Form == "" {
				modules[i].Form = "git"
			}
		}
	}
	return modules, nil
}
