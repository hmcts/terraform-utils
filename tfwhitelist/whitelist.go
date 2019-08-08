package tfwhitelist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

type moduleCall struct {
	Source string `json:"source"`
}

type managedResource struct {
	Type     string `json:"type"`
	Provider string `json:"provider,omitempty"`
	Mode     string `json:"mode,omitempty"`
}

type whitelist struct {
	Resources   map[string]managedResource `json:"resources"`
	ModuleCalls []moduleCall               `json:"module_calls"`
}

func loadModule(dir string) *tfconfig.Module {
	module, _ := tfconfig.LoadModule(dir)

	if module.Diagnostics.HasErrors() {
		_, _ = fmt.Fprintf(os.Stderr, "error loading module [%s]: %v\n", dir, module.Diagnostics.Error())
		os.Exit(1)
	}

	return module
}

func loadWhitelist(path string) (*whitelist, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading file [%s]: %s\n", path, err)
		return nil, err
	}
	_, _ = fmt.Fprintf(os.Stdout, "Successfully opened %s\n", path)

	defer jsonFile.Close()

	var allowed whitelist
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &allowed)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error unmarshalling file [%s]: %s\n", path, err)
		return nil, err
	}

	return &allowed, nil
}

func matchModules(module *tfconfig.Module, whitelist *whitelist) error {
	var notAllowed []tfconfig.ModuleCall

	// TODO inefficient, use a map-based data structure instead
	for _, v := range module.ModuleCalls {
		allowed := false
		for i := 0; i < len(whitelist.ModuleCalls); i++ {
			if whitelist.ModuleCalls[i].Source == v.Source {
				allowed = true
				break
			}
		}
		if !allowed {
			notAllowed = append(notAllowed, *v)
		}
	}

	if len(notAllowed) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Error matching modules: %v\n", notAllowed)
		return fmt.Errorf("modules not allowed found: %v\n", notAllowed)
	}
	return nil
}

func matchResources(module *tfconfig.Module, whitelist *whitelist) error {
	var notAllowed []tfconfig.Resource

	for k, v := range module.ManagedResources {
		if _, ok := whitelist.Resources[k]; !ok {
			notAllowed = append(notAllowed, *v)
		}
	}

	if len(notAllowed) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Error matching resources: %v\n", notAllowed)
		return fmt.Errorf("resources not allowed found: %v\n", notAllowed)
	}
	return nil
}
