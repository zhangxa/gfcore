package utils

import (
	"fmt"
	"reflect"
)

type sFunction struct {
}

// Function 函数工具类
var Function = &sFunction{}

// CallClassFunc 执行指定类的指定名称函数
func (s *sFunction) CallClassFunc(myClass interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(funcName)
	if !m.IsValid() {
		return make([]reflect.Value, 0), fmt.Errorf("method not found \"%s\"", funcName)
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}
	out = m.Call(in)
	return
}
