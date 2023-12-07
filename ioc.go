package ioc

// 每个结构体或者对象都可以具备IOCStart和IOCStop两个方法，可当做钩子使用

// IOCStarter 依赖注入模块启动接口
// 初始化工作，不该出现阻塞逻辑，否则会夯住
type IOCStarter interface {
	IOCStart() error
}

// IOCStopper 依赖注入模块停止接口
// 清理收尾工作，不该出现阻塞逻辑，否则会夯住
type IOCStopper interface {
	IOCStop() error
}
