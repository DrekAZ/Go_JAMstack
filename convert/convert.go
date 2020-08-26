package convert

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"unsafe"
)

func Struct2Map(data interface{}) map[string]interface{} {
	B, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal err", err)
		return nil
	}

	var m map[string]interface{}
	err = json.Unmarshal(B, &m)
	if err != nil {
		fmt.Println("unmarshal err", err)
		return nil
	}
	return m
}

func Rand2base64(digit uint32) (string, error) {
	// 乱数を生成
	b := make([]byte, digit)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	result := base64.StdEncoding.EncodeToString(b)
	return result, nil
}

// Byte2str Emergency!!!!! No Append This String
func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
