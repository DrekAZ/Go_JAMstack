package convert

import (
	"crypto/rand"
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

func Rand2str(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

// Byte2str Emergency!!!!! No Append This String
func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
