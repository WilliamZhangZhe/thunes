package gindoc

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// LoadDoc 引入swagger doc文档支持，http 路径为 category/doc
func LoadDoc(category *gin.RouterGroup) {
	url := ginSwagger.URL("doc.json") // The url pointing to API definition
	category.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
