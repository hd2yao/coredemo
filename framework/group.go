package framework

// IGroup 代表前缀分组
type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)
}

// Group struct 实现了 IGroup
type Group struct {
	core   *Core  // 指向core结构
	prefix string // 这个是group的通用前缀
}

// 初始化 Group
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

func (g *Group) Get(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	// 为当前 uri 注册 GET 方法的路由
	g.core.Get(uri, handler)
}

func (g *Group) Post(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	// 为当前 uri 注册 POST 方法的路由
	g.core.Post(uri, handler)
}

func (g *Group) Put(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	// 为当前 uri 注册 PUT 方法的路由
	g.core.Put(uri, handler)
}

func (g *Group) Delete(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	// 为当前 uri 注册 DELETE 方法的路由
	g.core.Delete(uri, handler)
}

// 实现 Group 方法
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}
