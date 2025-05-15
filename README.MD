基于goFrame框架，快速开发api接口应用，集成常用模块、utils等

当前处于开发及测试阶段，请谨慎用于生产，我将在实际开发应用中不断进行优化

# 使用说明

> golang版本 >= 1.20
> 
> gf版本 = v2.9.0
> 

gf官方文档：https://goframe.org/

# 安装

go get -u -v github.com/zhangxa/gfcore

# 使用

internal/cmd/cmd.go

````
import (
    "github.com/zhangxa/gfcore/core"
    .....
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			core.Server().InitServer(s)
			core.Server().AddModule("api", func(group *ghttp.RouterGroup) {
				//TODO add this module logic
				group.Bind(api.NewV1())
			})
			// more module support
			......
			return core.Server().Start(ctx, s)
		},
	}
)
````

# 更多说明，待完善...
