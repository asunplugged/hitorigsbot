package schema

import (
	"encoding/json"
	"errors"
	"fmt"
)

// DBVersions -
type DBVersions struct {
	Versions []VersionInfo `json:"versions"`
	Test     []VersionInfo `json:"test"`
}

// VersionInfo -
type VersionInfo struct {
	File    string `json:"file"`
	Version int    `json:"version"`
}

// ReadVersions -
func ReadVersions() (dbver DBVersions, err error) {
	f, err := Asset("schema/dbVersion.json")
	if err != nil {
		return dbver, err
	}

	err = json.Unmarshal(f, &dbver)
	return
}

// ReadSchema -
func ReadSchema(name string) (q string, err error) {
	if len(name) == 0 {
		return "", errors.New("name is empty")
	}
	f, err := Asset(fmt.Sprintf("schema/%s", name))
	if err != nil {
		return
	}
	q = string(f)
	return
}
