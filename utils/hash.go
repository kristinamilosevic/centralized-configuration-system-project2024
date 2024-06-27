package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"projekat/model"
)

func ComputeHash(config model.Config2) (string, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
