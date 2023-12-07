package ioc

import (
	"reflect"
	"unsafe"
)

type injectType int

const (
	noInjectType     injectType = iota
	objectInjectType            // 对象注入
)

type field struct {
	reflect.StructField                // structField
	rawType             reflect.Type   // 原始类型
	val                 reflect.Value  // 原始值
	injType             injectType     // 注入类型
	parentPtr           unsafe.Pointer // 父节点地址
}

func newField(sf reflect.StructField, parentPtr unsafe.Pointer) *field {
	injType := noInjectType // 默认不注入

	// 注入对象
	_, ok := sf.Tag.Lookup("inject")
	if ok {
		injType = objectInjectType
	}

	var rawType reflect.Type           // 原始类型
	if sf.Type.Kind() == reflect.Ptr { // 类型为指针的，用原始类型覆盖
		rawType = sf.Type.Elem()
	}
	return &field{
		StructField: sf,
		injType:     injType,
		rawType:     rawType,
		parentPtr:   parentPtr,
	}
}
