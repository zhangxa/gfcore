package jwt

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhangxa/gfcore/core"
	"github.com/zhangxa/gfcore/store"
	"strings"
	"time"
)

type sJWT struct {
	configDict map[string]*jwtConfig
	whiteList  *gcache.Cache
}

// jwtConfig jwt配置信息
type jwtConfig struct {
	SecretKey      []byte
	TimeoutHour    time.Duration
	MaxRefreshHour time.Duration
	IdentityKey    string
	TokenHeadName  string
	TokenScheme    string
}

type userClaims struct {
	*store.JWTAuthInfo
	jwt.RegisteredClaims
}

func init() {
	core.RegisterJWT(NewJWT())
}

// NewJWT jwt服务
func NewJWT() core.IJWT {
	res := &sJWT{
		configDict: make(map[string]*jwtConfig),
	}
	res.whiteList = gcache.New()
	if g.Config().MustGet(context.Background(), "redis.status").Bool() {
		redis, err := gredis.New()
		if err == nil {
			res.whiteList.SetAdapter(gcache.NewAdapterRedis(redis))
		}
	}
	return res
}

// getConfig 获取主键
func (s *sJWT) getConfig(module string) *jwtConfig {
	if config, ok := s.configDict[module]; ok {
		return config
	}
	pattern := fmt.Sprintf("jwt.%s", module)
	conf := &jwtConfig{
		SecretKey:      []byte("secret key"),
		TimeoutHour:    24,
		MaxRefreshHour: 0,
		IdentityKey:    "id",
		TokenHeadName:  "Authorization",
		TokenScheme:    "Bearer",
	}
	_ = g.Config().MustGet(context.Background(), pattern).Struct(&conf)
	s.configDict[module] = conf
	return conf
}

// GenToken 生产token等信息
func (s *sJWT) GenToken(module string, authInfo *store.JWTAuthInfo) (res *store.JWTAuthResult, err error) {
	// 生成token
	conf := s.getConfig(module)
	now := gtime.Now()
	expire := now.Add(conf.TimeoutHour * time.Hour)
	res = &store.JWTAuthResult{}
	res.ExpiresIn = expire.Timestamp()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		JWTAuthInfo: authInfo,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    module,
			ExpiresAt: jwt.NewNumericDate(expire.Time),
			IssuedAt:  jwt.NewNumericDate(now.Time),
		},
	})
	res.AccessToken, err = accessToken.SignedString(conf.SecretKey)
	if err != nil {
		return
	}
	if conf.MaxRefreshHour > 0 {
		// 刷新token
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
			JWTAuthInfo: authInfo,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    module,
				ExpiresAt: jwt.NewNumericDate(now.Add(conf.MaxRefreshHour * time.Hour).Time),
				IssuedAt:  jwt.NewNumericDate(now.Time),
			},
		})
		res.RefreshToken, err = refreshToken.SignedString(conf.SecretKey)
		if err != nil {
			return
		}
	}
	err = s.setWhiteList(res.AccessToken, conf.TimeoutHour)
	return
}

func (s *sJWT) Parse(ctx context.Context, module string) (token *jwt.Token, err error) {
	conf := s.getConfig(module)
	headerToken := g.RequestFromCtx(ctx).Header.Get(conf.TokenHeadName)
	if conf.TokenScheme != "" && strings.HasPrefix(headerToken, conf.TokenScheme) {
		headerToken = strings.Replace(headerToken, conf.TokenScheme, "", 1)
		headerToken = strings.TrimSpace(headerToken)
	}
	token, err = jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	if err != nil {
		return
	}
	if !token.Valid {
		return nil, gerror.New("token is invalid")
	}
	return
}

func (s *sJWT) ParseWithClaims(ctx context.Context, module string) (authInfo *store.JWTAuthInfo, expiresIn int64, err error) {
	conf := s.getConfig(module)
	headerToken := g.RequestFromCtx(ctx).Header.Get(conf.TokenHeadName)
	if conf.TokenScheme != "" && strings.HasPrefix(headerToken, conf.TokenScheme) {
		headerToken = strings.Replace(headerToken, conf.TokenScheme, "", 1)
		headerToken = strings.TrimSpace(headerToken)
	}
	return s.getAuthInfoByToken(headerToken, conf)
}

// parseWithClaims 方法描述
func (s *sJWT) getAuthInfoByToken(headerToken string, conf *jwtConfig) (authInfo *store.JWTAuthInfo, expiresIn int64, err error) {
	var tokenClaims *jwt.Token
	tokenClaims, err = jwt.ParseWithClaims(headerToken, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return conf.SecretKey, nil
	})
	if err != nil {
		return
	}
	if claims, ok := tokenClaims.Claims.(*userClaims); ok && tokenClaims.Valid {
		authInfo = claims.JWTAuthInfo
		expiresIn = claims.ExpiresAt.Time.Unix()
		return
	}
	err = errors.New("token parse with claims fail")
	return
}

// RefreshToken check if token expire
func (s *sJWT) RefreshToken(module string, refreshToken string) (res *store.JWTAuthResult, err error) {
	conf := s.getConfig(module)
	authInfo, _, wr := s.getAuthInfoByToken(refreshToken, conf)
	if wr != nil {
		return nil, wr
	}
	return s.GenToken(module, authInfo)
}

// getCacheKey 获取缓存key
func (s *sJWT) getCacheKey(tokenStr string) (key string, err error) {
	// The goal of MD5 is to reduce the key length.
	var tokenRaw string
	tokenRaw, err = gmd5.EncryptString(tokenStr)
	if err != nil {
		return
	}
	key = "JWT:TOKEN:" + tokenRaw
	return
}

func (s *sJWT) setWhiteList(tokenStr string, expireHour time.Duration) error {
	key, err := s.getCacheKey(tokenStr)
	if err != nil {
		return err
	}
	return s.whiteList.Set(context.Background(), key, true, expireHour*time.Hour)
}

func (s *sJWT) inWhiteList(tokenStr string) (bool, error) {
	// The goal of MD5 is to reduce the key length.
	key, err := s.getCacheKey(tokenStr)
	if err != nil {
		return false, nil
	}
	return s.whiteList.Contains(context.Background(), key)
}

// Logout 退出登录
func (s *sJWT) Logout(ctx context.Context, module string) error {
	token, err := s.Parse(ctx, module)
	if err != nil {
		return err
	}
	key, wr := s.getCacheKey(token.Raw)
	if wr != nil {
		return wr
	}
	_, err = s.whiteList.Remove(ctx, key)
	return err
}
