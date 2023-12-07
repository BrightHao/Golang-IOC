package ioc

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Container ioc容器
type Container struct {
	stopperList []IOCStopper             // 停止实例列表，按初始化顺序保存并执行
	values      map[string]reflect.Value // 全局实例缓存，暂未启用
}

// NewContainer 构造
func NewContainer() *Container {
	return &Container{
		values: make(map[string]reflect.Value),
	}
}

// Init 初始化对象
func (c *Container) Init(obj interface{}) error {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	// 判断空值
	if t == nil || t.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("error kind")
	}

	// new field
	sf := reflect.StructField{Type: t}
	f := newField(sf, nil)
	f.val = v
	// 遍历开始
	if err := c.doTraversal(f); err != nil {
		return err
	}
	// IOCStart
	iStarter, ok := obj.(IOCStarter)
	if ok {
		if err := iStarter.IOCStart(); err != nil {
			return err
		}
	}
	// IOCStop
	iStopper, ok := obj.(IOCStopper)
	if ok {
		c.stopperList = append(c.stopperList, iStopper)
	}
	return nil
}

// doTraversal 遍历字段
func (c *Container) doTraversal(parent *field) error {
	// 非struct不能遍历，直接返回
	if parent.rawType.Kind() != reflect.Struct {
		return fmt.Errorf("nonsupport traversal type %s, is not struct", parent.rawType.String())
	}
	// 开始遍历每个字段
	for i := 0; i < parent.rawType.NumField(); i++ {
		sf := parent.rawType.Field(i)
		parentPtr := unsafe.Pointer(parent.val.Pointer()) // 父节点的指针
		child := newField(sf, parentPtr)                  // 构建子字段结构体（还没创建值,标记是注入还是不注入）
		child.val = parent.val.Elem().Field(i)            // structField

		// 构建结构体中的字段
		if err := c.buildField(child); err != nil {
			return err
		}
	}
	return nil
}

// buildField 生成字段值
func (c *Container) buildField(f *field) error {
	switch f.injType {
	case objectInjectType:
		return c.buildObject(f)
	}
	return nil
}

// buildObject 生成依赖对象字段
func (c *Container) buildObject(f *field) error {
	switch f.Type.Kind() {
	case reflect.Ptr:
	case reflect.Interface:
		return fmt.Errorf("暂未支持ptr和interface")
	default:
		return fmt.Errorf("invalid inject object type %s, is not pointer", f.Type.String())
	}

	// 先不考虑共享实例
	val := reflect.New(f.rawType) // new一个rawType类型value
	f.val = val                   // 赋值到f

	// 递归遍历字段
	if err := c.doTraversal(f); err != nil {
		return err
	}

	// 执行钩子函数
	iVal := val.Interface()
	starter, ok := iVal.(IOCStarter)
	if ok {
		if err := starter.IOCStart(); err != nil {
			return fmt.Errorf("%s.%s start failed , error: %s", f.rawType.PkgPath(), f.rawType.Name(), err.Error())
		}
	}
	stopper, ok := iVal.(IOCStopper)
	if ok {
		c.stopperList = append(c.stopperList, stopper)
	}

	// 核心代码：将val放到父指针的地址，完成初始化
	ptr := unsafe.Pointer(uintptr(f.parentPtr) + f.Offset)
	*(*unsafe.Pointer)(ptr) = unsafe.Pointer(val.Pointer())
	return nil
}
