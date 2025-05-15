// Package swagger provides swagger UI resource files for swagger API service.
//
// Should be used with gf cli tool:
// gf pack public ./public-packed.go -p=swagger -y
package swagger

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/zhangxa/gfcore/core"
	"net/http"
	"time"
)

// Swagger is the struct for swagger feature management.
type Swagger struct {
	openapi *goai.OpenApiV3
	config  *swaggerConfig
	ctx     context.Context
}

// swaggerConfig is the information field for swagger.
type swaggerConfig struct {
	VisitPath      string    `json:"visitPath" dc:"访问路径，如果为空，则不启用"`
	BasicAuthUser  string    `json:"basicAuthUser" dc:"授权账号"`
	BasicAuthPass  string    `json:"basicAuthPass" dc:"授权密码"`
	ExcludePaths   []string  `json:"excludePaths" dc:"需要排除的接口前缀"`
	SecurityHeader string    `json:"securityHeader" dc:"安全校验header名称"`
	Info           goai.Info `json:"info" dc:"基础信息"`
}

const (
	Name               = "zhangxa"
	Author             = "zhangxa@yeah.net"
	Version            = "v1.0.0"
	Description        = "gf-swagger provides swagger API document feature for GoFrame project. https://github.com/gogf/gf-swagger"
	MaxAuthAttempts    = 10          // Max authentication count for failure try.
	AuthFailedInterval = time.Minute // Authentication retry interval after last failed.
)

// Name returns the name of the plugin.
func (sw *Swagger) Name() string {
	return Name
}

// Author returns the author of the plugin.
func (sw *Swagger) Author() string {
	return Author
}

// Version returns the version of the plugin.
func (sw *Swagger) Version() string {
	return Version
}

// Description returns the description of the plugin.
func (sw *Swagger) Description() string {
	return Description
}

// Install installs the swagger to server as a plugin.
func (sw *Swagger) Install(s *ghttp.Server) error {
	sw.ctx = context.Background()
	err := g.Cfg().MustGet(sw.ctx, "swagger").Struct(&sw.config)
	if err != nil {
		return err
	}
	if sw.config.VisitPath == "" {
		g.Log().Info(sw.ctx, "swagger config visitPath is empty, swagger will not start.")
		return nil
	}
	sw.openapi = goai.New()
	if sw.config.SecurityHeader != "" {
		sw.openapi.Components = goai.Components{
			SecuritySchemes: goai.SecuritySchemes{
				"Bearer": goai.SecuritySchemeRef{
					Ref: "", // 暂时还不知道该值是干什么用的
					Value: &goai.SecurityScheme{
						Type:        "apiKey",
						In:          "header",
						Name:        sw.config.SecurityHeader,
						Scheme:      "bearer",
						Description: "示例: Bearer --token--",
					},
				},
			},
		}
	}
	sw.openapi.Info = sw.config.Info
	// Initialize openapi.
	err = sw.initOpenApi(s)
	if err == nil {
		g.Log().Infof(sw.ctx, "swagger start success,visit path : %s", sw.config.VisitPath)
	}
	return err
}

// Remove uninstalls swagger feature from server.
func (sw *Swagger) Remove() error {
	return nil
}

// initOpenApi generates api specification using OpenApiV3 protocol.
func (sw *Swagger) initOpenApi(s *ghttp.Server) (err error) {
	var methods []string
	deployPath := core.Utils().GetDeployPath()
	for _, item := range s.GetRoutes() {
		switch item.Type {
		case ghttp.HandlerTypeMiddleware, ghttp.HandlerTypeHook:
			continue
		}
		if item.Handler.Info.IsStrictRoute {
			methods = []string{item.Method}
			if gstr.Equal(item.Method, "ALL") {
				methods = ghttp.SupportedMethods()
			}
			doExclude := false
			if len(sw.config.ExcludePaths) > 0 {
				for _, p := range sw.config.ExcludePaths {
					if gstr.HasPrefix(item.Route, p) {
						doExclude = true
						break
					}
				}
			}
			if doExclude {
				continue
			}
			if deployPath != "" {
				item.Route = fmt.Sprintf("/%s%s", deployPath, item.Route)
			}
			for _, method := range methods {
				api := goai.AddInput{
					Path:   item.Route,
					Method: method,
					Object: item.Handler.Info.Value.Interface(),
				}
				err = sw.openapi.Add(api)
				if err != nil {
					return err
				}
			}
		}
	}
	s.AddStaticPath(sw.config.VisitPath, "swagger")
	s.BindHookHandler(fmt.Sprintf("%s/*", sw.config.VisitPath), ghttp.HookBeforeServe, sw.swaggerUI)
	s.BindHandler(fmt.Sprintf("%s/swagger.json", sw.config.VisitPath), sw.openapiSpec)
	return nil
}

// openapiSpec is a build-in handler automatic producing for openapi specification json file.
func (sw *Swagger) openapiSpec(r *ghttp.Request) {
	r.Response.WriteJson(sw.openapi)
}

func (sw *Swagger) swaggerUI(r *ghttp.Request) {
	if sw.config.BasicAuthUser != "" {
		// Authentication security checks.
		var (
			authCacheKey = fmt.Sprintf(`swagger_auth_failed_%s`, r.GetClientIp())
			v, _         = gcache.Get(sw.ctx, authCacheKey)
			authCount    = v.Int()
		)
		if authCount > MaxAuthAttempts {
			r.Response.WriteStatus(
				http.StatusForbidden,
				"max authentication count exceeds, please try again in one minute!",
			)
			r.ExitAll()
		}
		// Basic authentication.
		if !r.BasicAuth(sw.config.BasicAuthUser, sw.config.BasicAuthPass) {
			_ = gcache.Set(sw.ctx, authCacheKey, authCount+1, AuthFailedInterval)
			r.ExitAll()
		}
	}
}
