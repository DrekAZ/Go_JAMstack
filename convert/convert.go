package convert

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"server_module/setting"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func Struct2map(data interface{}) (map[string]interface{}, error) {
	B, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal err", err)
		return nil, err
	}

	m, err := Unmarshal(B)
	if err != nil {
		return nil, err
	}

	return m, nil
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

func Unmarshal(b []byte) (map[string]interface{}, error) {
	var m map[string]interface{}

	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println("unmarshal err", err)
		return nil, err
	}
	return m, nil
}

func Str2bool(str string) bool {
	var b bool
	if str == "true" {
		b = true
	} else if str == "false" {
		b = false
	}

	return b
}

func BindJson2map(c *gin.Context, colName string) (map[string]interface{}, error) {
	var m map[string]interface{}

	if colName == "OnceTeam" {
		var j setting.OnceTeam
		err := c.BindJSON(&j)
		if err != nil {
			return nil, err
		}
		m, err = Struct2map(j)
		if err != nil {
			return nil, err
		}
	} else if colName == "Group" {
		var j setting.Group
		err := c.BindJSON(&j)
		if err != nil {
			return nil, err
		}
		m, err = Struct2map(j)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func UpdateBindJson2map(c *gin.Context, colName string) (map[string]interface{}, string, error) {
	var m map[string]interface{}

	if colName == "OnceTeam" {
		var j setting.UpdateOnceTeam
		err := c.BindJSON(&j)
		if err != nil {
			return nil, "", err
		}
		m, err = Struct2map(j)
		if err != nil {
			return nil, "", err
		}
	} else if colName == "Group" {
		var j setting.UpdateGroup
		err := c.BindJSON(&j)
		if err != nil {
			return nil, "", err
		}
		m, err = Struct2map(j)
		if err != nil {
			return nil, "", err
		}
	}
	return m["Data"].(map[string]interface{}), m["Page"].(string), nil
}
