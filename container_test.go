package ioc

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type Tp struct {
	Name    string
	FootNum int
}

func (t *Tp) IOCStart() error {
	t.Name = "monkey"
	t.FootNum = 2
	return nil
}

type Animal struct {
	Age   int
	Color string
	Tp    *Tp `inject:""`
}

func TestIOC(t *testing.T) {
	a := &Animal{}

	c := NewContainer()
	if err := c.Init(a); err != nil {
		t.Fatal(err)
	}

	fmt.Println(a)
	fmt.Println(a.Tp)
}

func TestPointer(t *testing.T) {
	a := &Animal{}
	fmt.Println(a)
	// 若不进行任何操作，a的值应该是&{0 <nil>}
	// 但是现在要对a的Tp进行赋值，且需要使用自动初始化的方式，先取a的type和value
	rt := reflect.TypeOf(a)
	// 当前已知Tp的index为1，取Tp的structField
	tpSf := rt.Elem().FieldByIndex([]int{1})
	tpTp := tpSf.Type
	// new一个同类型的val
	val := reflect.New(tpTp)
	// 计算该字段对应的地址，tpPtr为地址（uintptr：将指针类型转换为地址；unsafe.Pointer：将地址转换为指针类型）
	tpPtr := unsafe.Pointer(uintptr(unsafe.Pointer(a)) + tpSf.Offset)
	// 地址的指针指向val的地址，内存地址的指针指向=>val的地址
	*(*unsafe.Pointer)(tpPtr) = unsafe.Pointer(val.Pointer())

	fmt.Println()
	fmt.Println(a)
	fmt.Println(a.Tp)
}

func TestStringPoint(t *testing.T) {
	type Person struct {
		Name *string
	}

	n := "John"
	person := &Person{Name: &n}

	// 将person.Name的值修改为"Jane"
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(person)) + unsafe.Offsetof(person.Name))
	v := "Jane"

	rv := reflect.ValueOf(&v)
	namePtr := (*unsafe.Pointer)(ptr)
	*namePtr = unsafe.Pointer(rv.Pointer())

	fmt.Println("结果：", *person.Name) // 输出结果: Jane
}
