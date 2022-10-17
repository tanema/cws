package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func UpdateBytes(path, version string) ([]byte, error) {
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

	manifestBytes, err = json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling updated manifest: %v", err)
	}
	return manifestBytes, nil
}

// Update will update the manifest version, and remove any existing dev key
func Update(path, version string) error {
	manifestBytes, err := UpdateBytes(path, version)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, manifestBytes, 0)
}
