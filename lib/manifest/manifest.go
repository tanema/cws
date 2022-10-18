package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func parseJSONChangeset(changeset string) (map[string]string, error) {
	set := map[string]string{}
	for _, change := range strings.Split(changeset, ",") {
		parts := strings.Split(change, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed json change %v, expected key:value", change)
		}
		set[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return set, nil
}

func UpdateBytes(path, version, jsonChangeset string) ([]byte, error) {
	changeset, err := parseJSONChangeset(jsonChangeset)
	if err != nil {
		return nil, err
	}

	manifestBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %v", err)
	}

	manifest := map[string]interface{}{}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("error unmarshalling manifest: %v", err)
	}

	delete(manifest, "key")
	manifest["version"] = version
	for key, val := range changeset {
		manifest[key] = fmt.Sprintf(val, manifest[key])
	}

	manifestBytes, err = json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling updated manifest: %v", err)
	}
	return manifestBytes, nil
}

// Update will update the manifest version, and remove any existing dev key
func Update(path, version, jsonChangeset string) error {
	manifestBytes, err := UpdateBytes(path, version, jsonChangeset)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, manifestBytes, 0)
}
