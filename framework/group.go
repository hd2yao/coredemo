package framework

// IGroup 代表前缀分组
type IGroup interface {
	// 实现 HttpMethod 方法
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)

	// 实现嵌套 group
	Group(string) IGroup

	// 实现嵌套中间件
	Use(middlewares ...ControllerHandler)
}

// Group struct 实现了 IGroup
type Group struct {
	core        *Core               // 指向core结构
	parent      *Group              //指向上一个Group，如果有的话
	prefix      string              // 这个是group的通用前缀
	middlewares []ControllerHandler // 存放中间件
}

// NewGroup 初始化 Group
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:        core,
		parent:      nil,
		prefix:      prefix,
		middlewares: []ControllerHandler{},
	}
}

func (g *Group) Get(uri string, handler ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handler...)
	// 为当前 uri 注册 GET 方法的路由
	g.core.Get(uri, allHandlers...)
}

func (g *Group) Post(uri string, handler ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Post(uri, allHandlers...)
}

func (g *Group) Put(uri string, handler ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Put(uri, allHandlers...)
}

func (g *Group) Delete(uri string, handler ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Delete(uri, allHandlers...)
}

// 获取当前 group 的绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePrefix() + g.prefix
}

// 获取某个 group 的middleware
// 这里就是获取除了 Get/Post/Put/Delete 之外设置的 middleware
func (g *Group) getMiddlewares() []ControllerHandler {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}

// Group 实现 Group 方法
func (g *Group) Group(uri string) IGroup {
	cgroup := NewGroup(g.core, uri)
	cgroup.parent = g
	return cgroup
}

// 注册中间件
func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = middlewares
}
