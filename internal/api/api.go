// Package api 提供系统REST api接口路由信息，接口包括
// TODO：用户登录
// TODO：用户转账
// 功能

package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/thunes/docs"
	"github.com/thunes/internal/controller"
	jwtauth "github.com/thunes/internal/jwt"
	"github.com/thunes/pkg/gindoc"
)

type Engine struct {
	*gin.Engine

	userHandler    *controller.User
	clientHandler  *controller.Client
	accountHandler *controller.Account
}

// NewEngine 初始化API路由信息，并返回执行引擎
func NewEngine() *Engine {
	engine := &Engine{
		Engine: gin.New(),
	}

	engine.init()

	baseRoute := engine.Group("v1")

	return engine.docs(baseRoute).
		logInOut(baseRoute).
		clientRoutes(baseRoute.Group("clients", jwtauth.JWTCheck))
}

func (E *Engine) init() {
	if E.userHandler == nil {
		E.userHandler = controller.NewUserHandler()
	}

	if E.accountHandler == nil {
		E.accountHandler = controller.NewAccountHandler()
	}

	if E.clientHandler == nil {
		E.clientHandler = &controller.Client{}
	}
}

func (E *Engine) docs(category *gin.RouterGroup) *Engine {
	gindoc.LoadDoc(category)
	return E
}

func (E *Engine) logInOut(category *gin.RouterGroup) *Engine {
	client := E.clientHandler

	category.POST("/login", client.Login)
	category.POST("/logout", jwtauth.JWTCheck, client.Logout)

	return E
}

func (E *Engine) clientRoutes(category *gin.RouterGroup) *Engine {
	client := E.clientHandler

	category.GET("/:cid", client.Get)
	category.GET("/:cid/accounts", client.FindAccounts)

	subGroup := category.Group("/:cid")
	// subGroup := category.Group("/:cid")
	return E.userRoutes(subGroup.Group("/user")).
		accountRoutes(subGroup.Group("/account"))
}

func (E *Engine) userRoutes(category *gin.RouterGroup) *Engine {
	// user := E.userHandler

	return E
}

func (E *Engine) accountRoutes(category *gin.RouterGroup) *Engine {
	account := E.accountHandler

	category.POST("/:aid/transfer", account.Transfer)
	category.GET("/:aid", account.Get)

	return E
}
