package safety

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GoPlusAddress(address string) string {
	address = strings.ToLower(address) + "gopluslabs"
	sha256b := sha256.Sum256([]byte(address))
	sha256hash := hex.EncodeToString(sha256b[0:])
	return sha256hash
}

func RpcHubAddress(address string) string {
	address = strings.ToLower(address) + "rpchub"
	md5b := md5.Sum([]byte(address))
	md5hash := hex.EncodeToString(md5b[0:])
	return md5hash
}

func SlowMistAddress(address string) string {
	address = strings.ToLower(address) + "SlowMist"
	md5b := md5.Sum([]byte(address))
	md5hash := hex.EncodeToString(md5b[0:])
	return md5hash
}
