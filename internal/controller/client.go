package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thunes/internal/db"
	jwtauth "github.com/thunes/internal/jwt"
	"github.com/thunes/internal/model"
	"github.com/thunes/pkg/restapi"
)

type Client struct{}

//////////////////////////////////////////////////////////////////////////////

type LoginReq struct {
	Email string `json:"email" form:"email" binding:"required,email,min=1,max=64"`
	PWD   string `json:"pwd" form:"pwd" binding:"required,min=1,max=64"`
}

// Login 用户登录并分配JWT(json-web-token)身份信息
// @Summary 登录
// @Description 用户登录
// @Produce  json
// @Param request body controller.LoginReq true "账户及密码"
// @Success 200 {object} restapi.Response
// @Failure 406 {object} restapi.Response "参数错误"
// @Failure 404 {object} restapi.Response "账户不存在"
// @Failure 401 {object} restapi.Response "授权失败/未授权"
// @Failure 500 {object} restapi.Response "服务器错误"
// @Router /v1/login [post]
func (C *Client) Login(ctx *gin.Context) {
	var (
		req = LoginReq{}
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Client: Login param parse error, ", err)

		ginInvalidParam(ctx)
		return
	}

	user := UserModel{session: db.Cli().Table(&model.User{})}
	u, err := user.FindByEmail(req.Email)
	if err != nil {
		fmt.Println("Client:Login find user error, ", err)

		if err == ErrAccountNotFound {
			ginError(ctx, err)
			return
		}

		ginInternalErr(ctx)
		return
	}

	if !user.VerifyPWD(u, req.PWD) {
		ginResponse(ctx, restapi.Response{
			Code: http.StatusUnauthorized,
			Msg:  ErrWrongPassword.Error(),
		})
		return
	}

	ctx.Set(jwtauth.JWTValueKey, user)
	if !jwtauth.JWTSet(ctx) {
		ginError(ctx, ErrUnauthorized)
		return
	}

	ginSuccess(ctx, "login success")
}

type LogoutReq struct {
	CID model.UID `json:"cid" uri:"cid" binding:"required,numeric,min=1"`
}

// Logout 用户登出（TODO: 登出后将之前分配的JWT失效）
// @Summary 登出
// @Description 用户登出
// @Produce  json
// @Param user body controller.LogoutReq true "账户ID"
// @Success 200 {object} restapi.Response
// @Failure 406 {object} restapi.Response "参数错误"
// @Failure 404 {object} restapi.Response "账户不存在"
// @Failure 401 {object} restapi.Response "授权失败/未授权"
// @Failure 500 {object} restapi.Response "服务器错误"
// @Router /v1/logout [post]
func (C *Client) Logout(ctx *gin.Context) {
	var (
		req = LoginReq{}
	)

	if err := ctx.ShouldBindUri(&req); err != nil {
		fmt.Println("Client: Logout param parse error, ", err)

		ginInvalidParam(ctx)
		return
	}

	// TODO: 将指定的jwt invalid，防止登出后继续使用之前的token操作
	jwtauth.JWTInvalid(ctx)
	ginSuccess(ctx, "logout success")
}

//////////////////////////////////////////////////////////////

type ClientAccountsReq struct {
	CID model.UID `json:"cid" form:"cid" binding:"required,numeric,min=1,max=9999999999999999"`
}

type ClientAccountsResp []model.Account

// FindAccounts 按uid获取用户的所有账户
func (C *Client) FindAccounts(ctx *gin.Context) {
	var (
		req = ClientAccountsReq{}
	)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		fmt.Println("Use: FindAccounts error [invalid param], ", err)
		ginInvalidParam(ctx)
		return
	}

	m := &AccountModel{
		session: db.Cli().Table(&model.Account{}),
	}

	accounts, err := m.FindByUID(req.CID)
	if err != nil {
		if err == ErrAccountNotFound {
			ginError(ctx, err)
			return
		}

		ginInternalErr(ctx)
		return
	}

	ginSuccess(ctx, "success", accounts)
}

//////////////////////////////////////////////////////////////////////////////////////////

type GetClientResp struct {
	ID    model.UID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`

	Accounts []model.Account `json:"accounts"`
}

type GetClientReq struct {
	CID model.UID `json:"cid" uri:"cid" binding:"numeric,min=1"`
}

// Get
func (C Client) Get(ctx *gin.Context) {
	var (
		req = GetClientReq{}
	)

	if err := ctx.ShouldBindUri(&req); err != nil {
		fmt.Println(err)

		ginInvalidParam(ctx, "invalid id")
		return
	}

	userModel := &UserModel{session: db.Cli().Table(&model.User{})}
	user, err := userModel.Get(req.CID)
	if err != nil {
		fmt.Println("Client: Get user info error, ", err)

		if err == ErrAccountNotFound {
			ginError(ctx, err)
			return
		}

		ginInternalErr(ctx)
		return
	}

	accountModel := &AccountModel{session: db.Cli().Table(&model.Account{})}
	accounts, err := accountModel.FindByUID(req.CID)
	if err != nil && err != ErrAccountNotFound {
		fmt.Println("Client: Get user account error, ", err)

		ginInternalErr(ctx)
		return
	}

	resp := GetClientResp{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Accounts: accounts,
	}

	ginSuccess(ctx, "", resp)
}

////////////////////////////////////////////////////////////
