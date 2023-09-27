package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/google/uuid"
)

// 根据 open_id 和  sessionKey  去生成 token
func GenerateToken(openID, sessionKey string) (string, error) {
	// Concatenate openID and sessionKey
	data := openID + sessionKey

	// Create a new HMAC-SHA256 hasher
	hasher := hmac.New(sha256.New, []byte("your_secret_key_here")) // Replace with your secret key

	// Write the data to the hasher
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}

	// Calculate the HMAC and encode it in base64
	token := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return token, nil
}

// 用 OpenID 生成 UserID/UUID
func GenerateUUID(openID string) (string, error) {
	// Create a UUID based on the openID
	uuidBasedOnOpenID := uuid.NewSHA1(uuid.NameSpaceURL, []byte(openID))

	// Convert the UUID to string representation
	uuidString := uuidBasedOnOpenID.String()

	return uuidString, nil
}
