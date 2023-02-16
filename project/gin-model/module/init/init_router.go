package start

import (
	"github.com/gin-gonic/gin"
	"github.com/zqlpaopao/tool/project/gin-model/module/web/router"
)

func loadRouter(g *gin.Engine) {
	router.Ping{}.Load(g)
}
