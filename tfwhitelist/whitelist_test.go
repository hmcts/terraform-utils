package tfwhitelist

import (
	"testing"
)

func TestLoadWhitelist(t *testing.T) {
	w, err := loadWhitelist("test-fixtures/module-not-allowed/whitelist.json")
	if err != nil {
		t.Errorf("Error %s", err)
	}
	if len(w.ModuleCalls) != 2 {
		t.Errorf("ModuleCalls = %v, want len(2)", w.ModuleCalls)
	}
	if len(w.Resources) != 0 {
		t.Errorf("Resources = %v, want len(0)", w.Resources)
	}
}

func TestModuleNotAllowed(t *testing.T) {
	w, err := loadWhitelist("test-fixtures/module-not-allowed/whitelist.json")
	if err != nil {
		t.Errorf("Error %s", err)
	}

	m := loadModule("test-fixtures/module-not-allowed")

	err = matchModules(m, w)
	if err == nil {
		t.Errorf("matchModules should return an error")
	}
}

func TestResourceNotAllowed(t *testing.T) {
	w, err := loadWhitelist("test-fixtures/resource-not-allowed/whitelist.json")
	if err != nil {
		t.Errorf("Error %s", err)
	}

	m := loadModule("test-fixtures/resource-not-allowed")

	err = matchResources(m, w)
	if err == nil {
		t.Errorf("matchResources should return an error")
	}
}

func TestAllAllowed(t *testing.T) {
	w, err := loadWhitelist("test-fixtures/all-allowed/whitelist.json")
	if err != nil {
		t.Errorf("Error %s", err)
	}

	m := loadModule("test-fixtures/all-allowed")

	errResources := matchResources(m, w)
	errModules := matchModules(m, w)
	if errResources != nil || errModules != nil {
		t.Errorf("matchResources and matchModules should not return an error")
	}
}
