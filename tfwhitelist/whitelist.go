package tfwhitelist

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
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
	Resources   []managedResource `json:"resources"`
	ModuleCalls []moduleCall      `json:"module_calls"`
}

func (w *whitelist) merge(part *whitelist) {
	if len(part.Resources) > 0 {
		w.Resources = append(w.Resources, part.Resources...)
	}
	if len(part.ModuleCalls) > 0 {
		w.ModuleCalls = append(w.ModuleCalls, part.ModuleCalls...)
	}
}

func LoadAndMatchAll(infraPath string, whitelistPaths []string) error {
	w, err := loadWhitelist(whitelistPaths)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading whitelist at %s\n", whitelistPaths)
		return fmt.Errorf("Error loading whitelist at %s\n", whitelistPaths)
	}

	m, err := loadModule(infraPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading infra definition at %s\n", infraPath)
		return fmt.Errorf("Error loading infra definition at %s\n", infraPath)
	}

	errResources := matchResources(m, w)
	errModules := matchModules(m, w)
	if errResources != nil || errModules != nil {
		return fmt.Errorf("matchResources and matchModules should not return an error")
	}
	return nil
}

func loadModule(dir string) (*tfconfig.Module, error) {
	module, _ := tfconfig.LoadModule(dir)

	if module.Diagnostics.HasErrors() {
		_, _ = fmt.Fprintf(os.Stderr, "error loading module [%s]: %v\n", dir, module.Diagnostics.Error())
		return nil, module.Diagnostics.Err()
	}

	return module, nil
}

func loadWhitelist(paths []string) (*whitelist, error) {
	var allowed whitelist

	// merge all whitelists
	for i := 0; i < len(paths); i++ {
		var allowedPart whitelist
		jsonFile, err := os.Open(paths[i])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error loading file [%s]: %s\n", paths[i], err)
			return nil, err
		}
		_, _ = fmt.Fprintf(os.Stdout, "Successfully opened %s\n", paths[i])
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(byteValue, &allowedPart)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error unmarshalling file [%s]: %s\n", paths, err)
			return nil, err
		}
		allowed.merge(&allowedPart)
	}

	return &allowed, nil
}

func matchModules(module *tfconfig.Module, whitelist *whitelist) error {
	var notAllowed []tfconfig.ModuleCall
	re, _ := regexp.Compile(`.*\?ref=.*`)

	for _, v := range module.ModuleCalls {
		allowed := false
		source := v.Source
		// use master branch as default
		if !re.MatchString(source) {
			source = strings.Join([]string{source, "?ref=master"}, "")
		}
		for i := 0; i < len(whitelist.ModuleCalls); i++ {
			if whitelist.ModuleCalls[i].Source == source {
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

	for _, v := range module.ManagedResources {
		allowed := false
		for i := 0; i < len(whitelist.Resources); i++ {
			if whitelist.Resources[i].Type == v.Type {
				allowed = true
				break
			}
		}
		if !allowed {
			notAllowed = append(notAllowed, *v)
		}
	}

	if len(notAllowed) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Error matching resources: %v\n", notAllowed)
		return fmt.Errorf("resources not allowed found: %v\n", notAllowed)
	}
	return nil
}
