package share

import "github.com/yanyiwu/gojieba"

var Path string
var Address string
var TestMode bool
var Endpoint string

var JiebaPtr *gojieba.Jieba

func VariableWrapper[T any](anyValue T) T {
	return anyValue
}

func VariablePtrWrapper[T any](anyValue T) *T {
	return &anyValue
}
