package middleware

import (
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/zhangxa/gfcore/core"
	"github.com/zhangxa/gfcore/fixed"
	"github.com/zhangxa/gfcore/store"
	"mime"
	"net/http"
	"time"
)

const (
	contentTypeEventStream  = "text/event-stream"
	contentTypeOctetStream  = "application/octet-stream"
	contentTypeMixedReplace = "multipart/x-mixed-replace"
)

var (
	streamContentType = []string{contentTypeEventStream, contentTypeOctetStream, contentTypeMixedReplace}
)

type sMiddleware struct {
}

func init() {
	core.RegisterMiddleware(NewMiddleware())
}

// NewMiddleware 中间件服务
func NewMiddleware() core.IMiddleware {
	return &sMiddleware{}
}

// Base 基础中间件
func (s *sMiddleware) Base(r *ghttp.Request) {
	core.Context().Init(r.Context())
	return
}

func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func (s *sMiddleware) VisitLimit(r *ghttp.Request) {
	module := core.Context().GetModule(r.Context())
	if core.IpLimit().IsLimited(r.Context(), module) {
		r.Response.WriteJson(store.DefaultHandlerResponse{
			Code:    gcode.CodeBusinessValidationFailed,
			Msg:     "request forbidden!",
			Data:    nil,
			Success: false,
		})
		r.ExitAll()
	} else {
		r.Middleware.Next()
	}
}

func (s *sMiddleware) I18n(r *ghttp.Request) {
	lang := r.GetHeader("Accept-Language")
	if lang == "" {
		module := core.Context().GetModule(r.Context())
		lang = core.Modules().GetConfig("defaultLang", fixed.LanguageDefault, module).String()
	}
	r.SetCtx(gi18n.WithLanguage(r.Context(), lang))
	r.Middleware.Next()
}

// HandlerResponse is the default middleware handling handler response object and its error.
func (s *sMiddleware) HandlerResponse(r *ghttp.Request) {
	r.Middleware.Next()
	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 {
		return
	}
	mediaType, _, _ := mime.ParseMediaType(r.Response.Header().Get("Content-Type"))
	for _, ct := range streamContentType {
		if mediaType == ct {
			return
		}
	}
	var (
		ctx      = r.Context()
		err      = r.GetError()
		localCtx = core.Context().Get(ctx)
		response interface{}
		code     = gerror.Code(err)
		resCode  interface{}
		resData  interface{}
		msg      = ""
		act      = ""
		url      = ""
		data     = r.GetHandlerResponse()
		module   = ""
	)
	if localCtx != nil {
		response = localCtx.HandlerResponse
		msg = localCtx.ResMsg
		act = localCtx.ResAct
		url = localCtx.ResUrl
		resCode = localCtx.ResCode
		resData = localCtx.ResData
		module = localCtx.Module
	}
	// There's custom buffer content, it then exits current handler.
	if response != nil {
		r.Response.WriteJson(response)
		return
	}
	if err != nil {
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		if msg == "" {
			msg = err.Error()
			if act != "" {
				msg = fmt.Sprintf("%s失败", act)
				if core.Modules().IsDebug(module) {
					msg = fmt.Sprintf("%s:%s", msg, err.Error())
				}
			}
		} else {
			if core.Modules().IsDebug(module) {
				msg = fmt.Sprintf("%s:%s", msg, err.Error())
			}
		}
	} else {
		if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
			msg = http.StatusText(r.Response.Status)
			switch r.Response.Status {
			case http.StatusNotFound:
				code = gcode.CodeNotFound
			case http.StatusForbidden:
				code = gcode.CodeNotAuthorized
			default:
				code = gcode.CodeUnknown
			}
			// It creates an error as it can be retrieved by other middlewares.
			err = gerror.NewCode(code, msg)
			r.SetError(err)
		} else {
			code = gcode.CodeOK
			if msg == "" {
				msg = code.Message()
				if act != "" {
					msg = fmt.Sprintf("%s成功", act)
				}
			}
		}
	}
	res := &store.DefaultHandlerResponse{
		Code:      code.Code(),
		Msg:       msg,
		Data:      data,
		Url:       url,
		Success:   code.Code() == 0,
		Timestamp: time.Now().Unix(),
	}
	if resCode != nil {
		res.Code = resCode
	}
	if resData != nil {
		res.Data = resData
	}
	//service.Context().SetResponse(ctx, res)
	// g.Dump("--response res:", res)
	// if r.IsAjaxRequest() || strings.EqualFold(r.Header.Get("Content-Type"), "application/json") {}
	r.Response.WriteJson(res)
}
