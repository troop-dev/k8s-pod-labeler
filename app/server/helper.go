package server

import (
	"encoding/base64"
	"encoding/json"
)

// parseB64Map accepts a base64-encoded json object and
// unmarshals it to a map of strings
func parseB64Map(in string) (map[string]string, error) {
	// decode from base64
	decoded, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	// unmarshal json
	if len(decoded) > 0 {
		if err := json.Unmarshal(decoded, &out); err != nil {
			return nil, err
		}
	}
	return out, nil
}
