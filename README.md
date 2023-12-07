# Golang-IOC依赖注入

Golang-IOC提供了依赖注入功能，主要能力如下：
- 支持各种结构、接口的依赖注入
- 具备对象生命周期管理机制，可接管对象初始化和销毁

## 快速开始
以下示例将展示以下功能：
1. 注册结构体
2. 结构体初始化以及销毁方法注入
3. 结构体自动创建
```
package main

import (
	"fmt"

	ioc "github.com/BrightHao/golang-ioc"
)

type Tp struct {
	Name    string
	FootNum int
}

// 结构体初始化函数，注册后将在创建时自动调用
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

func main() {
	// 创建最上层结构体
	a := &Animal{}

	// 创建IOC容器
	c := ioc.NewContainer()
	// 将a初始化
	if err := c.Init(a); err != nil {
		fmt.Println(err)
		return
	}

	// &{0  0xc00000c090}，可以看到Tp并非nil，说明已经被注入
	fmt.Println(a)
	// &{monkey 2}，这里可以看到Tp的IOCStart函数被自动执行了
	fmt.Println(a.Tp)
}
```
