package main

import (
	"fmt"
	"testing"
	//"github.com/gin-gonic/gin"
)

func TestGetEnv(t *testing.T) {
	env, authEnv, err := GetEnv()
	fmt.Println(env)
	fmt.Println(authEnv)
	fmt.Println(err)

	//testStruct2Map(env)
	//testStruct2Map(authEnv)
}

func TestRandStr(t *testing.T) {
	bytes, err := RandStr(13)
	fmt.Println(bytes, err)
	bytes, err = RandStr(17)
	fmt.Println(bytes, err)
}

func testStruct2Map(data interface{}) {
	m := Struct2Map(data)
	fmt.Println(m)
}
