package uapi

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/ucontext"
	"git.umu.work/AI/uglib/uhttp"
	"git.umu.work/AI/uglib/umetadata"
	"git.umu.work/be/goframework/auth"
	"git.umu.work/be/goframework/logger"
	"github.com/go-resty/resty/v2"
)

const KeyLoginToken = "umu-t"

type UserInfo struct {
	UmuId        int64  `json:"umu_id"`
	StudentId    int64  `json:"student_id"`
	TeacherId    int64  `json:"teacher_id"`
	EnterpriseId int64  `json:"enterprise_id"`
	UserName     string `json:"user_name"`
	Avatar       string `json:"avatar"`
}

type UserResponse struct {
	BaseResponse
	Data UserInfo `json:"data"`
}

type UserRequest struct {
}

type UserAPI struct {
	userBaseUrl string
	uAPPBaseUrl string
	uRPCUrl     string
	client      *resty.Client
}

type UserLoginInfo struct {
	Code     int `json:"code"`
	UserInfo struct {
		StudentId      string `json:"student_id"`
		UmuId          string `json:"umu_id"`
		UserName       string `json:"user_name"`
		Avatar         string `json:"avatar"`
		HomeUrl        string `json:"home_url"`
		Phone          string `json:"phone"`
		Email          string `json:"email"`
		UserMark       string `json:"user_mark"`
		EnterpriseInfo struct {
			EnterpriseId   string `json:"enterprise_id"`
			EnterpriseName string `json:"enterprise_name"`
			ShowName       string `json:"show_name"`
			LogoMobile     string `json:"logo_mobile"`
		} `json:"enterprise_info"`
		MedalInfo struct {
			ShowUserLevel    int `json:"show_user_level"`
			UserLevel        int `json:"user_level"`
			UserGrowthPoints int `json:"user_growth_points"`
			UserMedal        struct {
				Id        string `json:"id"`
				MedalType string `json:"medal_type"`
				MedalRank string `json:"medal_rank"`
			} `json:"user_medal"`
		} `json:"medal_info"`
	} `json:"user_info"`
	TeacherInfo struct {
		UmuId          string `json:"umu_id"`
		TeacherId      string `json:"teacherId"`
		TeacherId1     string `json:"teacher_id"`
		StudentId      string `json:"student_id"`
		EnterpriseId   string `json:"enterprise_id"`
		BindPhone      string `json:"bind_phone"`
		BindAreaCode   string `json:"bind_area_code"`
		RegisterStep   int    `json:"register_step"`
		RegisterFrom   string `json:"register_from"`
		LoginType      string `json:"login_type"`
		NeedBindMobile int    `json:"need_bind_mobile"`
	} `json:"teacher_info"`
}

type UserLoginResponse struct {
	Status    bool          `json:"status"`
	Errno     int           `json:"errno"`
	ErrorCode int           `json:"error_code"`
	Error     string        `json:"error"`
	Data      UserLoginInfo `json:"data"`
	Token     string        `json:"token"`
	PageToken string        `json:"page_token"`
	Version   string        `json:"version"`
	Config    struct {
		StepLength int    `json:"step_length"`
		Env        string `json:"env"`
		Lang       string `json:"lang"`
		SiteHost   string `json:"site_host"`
		System     string `json:"system"`
		Fixed      bool   `json:"fixed"`
	} `json:"config"`
}

type RefreshJwtRequest struct {
	Service string `json:"service"`
	Method  string `json:"method"`
	Request struct {
	} `json:"request"`
}

type RefreshJwtResponse struct {
	AutoGen bool   `json:"autoGen"`
	Cookie  string `json:"cookie"`
	Token   string `json:"token"`
}

func NewUserAPI(uAPIUrl, aPPUrl, uRPCUrl string, client *resty.Client) *UserAPI {
	userBaseUrl, err := uhttp.JoinPath(uAPIUrl, "v1", "user")
	if err != nil {
		panic(err)
	}
	return &UserAPI{
		uRPCUrl:     uRPCUrl,
		userBaseUrl: userBaseUrl,
		uAPPBaseUrl: aPPUrl,
		client:      client,
	}
}

func (api *UserAPI) GetUser(ctx context.Context, token string) (context.Context, *UserInfo, error) {
	logger.GetLogger(ctx).Info(fmt.Sprintf("user api get:token %s is %+v", KeyLoginToken, token))
	userGetUrl, err := uhttp.JoinPath(api.userBaseUrl, "get")
	if err != nil {
		return ctx, nil, err
	}
	result := &UserResponse{}
	resp, err := api.client.R().SetHeader(KeyLoginToken, token).ForceContentType("application/json").
		SetResult(result).Get(userGetUrl)
	if err != nil {
		return ctx, nil, err
	}
	if result.ErrorCode != 0 {
		return ctx, nil, fmt.Errorf("get user failed, code=%d, msg=%s", result.ErrorCode, result.ErrorMessage)
	}
	if result.Data.UmuId == 0 {
		return ctx, nil, fmt.Errorf("get user failed, umu_id=0")
	}
	ctx = umetadata.SetCookies(ctx, resp.Cookies())

	return ctx, &result.Data, nil
}

func (api *UserAPI) Login(ctx context.Context, username, password string) (context.Context, *UserLoginResponse, error) {
	loginUrl, err := uhttp.JoinPath(api.uAPPBaseUrl, "passport", "ajax", "account", "login")
	if err != nil {
		return ctx, nil, err
	}
	result := &UserLoginResponse{}
	resp, err := api.client.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"cache-control": "no-cache",
			"Accept":        "application/json",
		}).
		SetFormData(map[string]string{
			"username": username,
			"passwd":   password,
		}).
		SetResult(result).
		//SetError(Error{}).
		Post(loginUrl)
	if err != nil {
		return ctx, nil, err
	}
	ctx = umetadata.SetCookies(ctx, resp.Cookies())

	return ctx, result, nil
}

func (api *UserAPI) RefreshJwt(ctx context.Context, umuToken string) (string, error) {
	req := RefreshJwtRequest{
		Service: "umu.generic.auth",
		Method:  "AuthService.Cookie2Token",
		Request: struct{}{},
	}
	resp := &RefreshJwtResponse{}
	refreshJwtResp, err := api.client.R().SetHeader("Umu-T", umuToken).
		ForceContentType("application/json").SetBody(req).
		SetResult(resp).Post(api.uRPCUrl)
	if err != nil {
		return "", err
	}
	if refreshJwtResp.IsError() {
		err = fmt.Errorf("refresh jwt resp is %+v\n", refreshJwtResp)
		return "", err
	}

	return resp.Token, nil
}

func (api *UserAPI) InitAuthentication(ctx context.Context, username, password string) (context.Context, error) {
	ctx = ucontext.NewUValueContext(ctx)
	ctx, userInfo, err := api.Login(ctx, username, password)
	if err != nil {
		return nil, err
	}
	var token string
	token, err = api.RefreshJwt(ctx, userInfo.Token)
	if err != nil {
		return nil, err
	}
	fmt.Printf("token is %s\n", token)
	u, err := auth.DecodeAuthUser(token)
	if err != nil {
		return nil, err
	}
	ctx = auth.SetUser(ctx, u)
	ctx, err = umetadata.SetAuthorization(ctx)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (api *UserAPI) InitUserCookies(ctx context.Context, username, password string) (context.Context, error) {
	ctx = ucontext.NewUValueContext(ctx)
	ctx, _, err := api.Login(ctx, username, password)
	if err != nil {
		return nil, err
	}
	//ctx, _, err = api.GetUser(ctx, userInfo.Token)
	//if err != nil {
	//	return nil, err
	//}

	return ctx, nil
}
