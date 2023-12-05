package umetadata

import (
	"context"
	"git.umu.work/be/goframework/auth"
	"git.umu.work/be/goframework/metadata"
	"github.com/google/uuid"
	gmetadata "github.com/micro/go-micro/v2/metadata"
	"net/http"
)

type cookiesKey struct{}

const (
	KeyUmuSID   = "X-Umu-Sid"
	KeyDeviceID = "Umu-Di"
	KeyDebug    = "X-Umu-Debug"
)

// GetSID set和get valuectx key类型不一致，一个是内部`封装的一个是go-micro官方的，都是struct不是interface
func GetSID(ctx context.Context) (string, bool) {
	if sID, ok := metadata.Get(ctx, KeyUmuSID); ok {
		return sID, true
	} else {
		return uuid.NewString(), false
	}
}

func GetDeviceID(ctx context.Context) (string, bool) {
	if deviceID, ok := metadata.Get(ctx, KeyDeviceID); ok {
		return deviceID, true
	} else {
		return "", false
	}
}

func GetDebug(ctx context.Context) bool {
	if _, ok := metadata.Get(ctx, KeyDebug); ok {
		return true
	} else {
		return false
	}
}

// SetSID set和get valuectx key类型不一致，一个是内部封装的一个是go-micro官方的，都是struct不是interface
func SetSID(ctx context.Context, sID *string) (ctxWithSID context.Context) {
	if sID != nil {
		ctxWithSID = gmetadata.Set(ctx, KeyUmuSID, *sID)
	} else {
		ctxWithSID = gmetadata.Set(ctx, KeyUmuSID, uuid.NewString())
	}

	return ctxWithSID
}

func SetDeviceID(ctx context.Context, deviceID *string) (ctxWithDeviceID context.Context) {
	if deviceID != nil {
		ctxWithDeviceID = gmetadata.Set(ctx, KeyDeviceID, *deviceID)
	} else {
		ctxWithDeviceID = gmetadata.Set(ctx, KeyDeviceID, "")
	}

	return ctxWithDeviceID
}

func SetAuthorization(ctx context.Context) (ctxWithAuthorization context.Context, err error) {
	authenticateUser, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	jwtToken, err := auth.EncodeAuthUser(authenticateUser)
	if err != nil {
		return nil, err
	}
	ctxWithAuthorization = gmetadata.Set(ctx, auth.KeyUmuJwtUser, jwtToken)

	return ctxWithAuthorization, nil
}

func SetDebug(ctx context.Context) (ctxWithDeviceID context.Context) {
	ctxWithDeviceID = gmetadata.Set(ctx, KeyDebug, "debug")

	return ctxWithDeviceID
}

func SetCookies(ctx context.Context, rs []*http.Cookie) (ctxWithCookies context.Context) {
	ctxWithCookies = context.WithValue(ctx, cookiesKey{}, rs)

	return ctxWithCookies
}

func GetCookies(ctx context.Context) ([]*http.Cookie, bool) {
	cookies, ok := ctx.Value(cookiesKey{}).([]*http.Cookie)
	if !ok {
		return nil, false
	}

	return cookies, true
}
