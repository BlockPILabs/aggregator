package jsonrpc

import "fmt"

func Error(code int, msg string) []byte {
	return []byte(fmt.Sprintf(`{"error":{"code":%d,"message":"%s"}}`, code, msg))
}
