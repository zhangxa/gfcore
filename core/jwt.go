// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package core

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zhangxa/gfcore/store"
)

type (
	IJWT interface {
		// GenToken 生产token等信息
		GenToken(module string, authInfo *store.AuthUserInfo) (res *store.JWTAuthResult, err error)
		Parse(ctx context.Context, module string) (token *jwt.Token, err error)
		ParseWithClaims(ctx context.Context, module string) (authInfo *store.AuthUserInfo, expiresIn int64, err error)
		// RefreshToken check if token expire
		RefreshToken(module string, refreshToken string) (res *store.JWTAuthResult, err error)
		// Logout 退出登录
		Logout(ctx context.Context, module string) error
	}
)

var (
	localJWT IJWT
)

func JWT() IJWT {
	if localJWT == nil {
		panic("implement not found for interface IJWT, forgot register?")
	}
	return localJWT
}

func RegisterJWT(i IJWT) {
	localJWT = i
}
