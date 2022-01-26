package start

import (
	"github.com/gin-gonic/gin"
	"github.com/zqlpaopao/tool/gin-model/module/web/router"
)

func loadRouter(g *gin.Engine){
	router.Login{}.Load(g)
	//router.Project{}.Load(g)
	//router.User{}.Load(g)
	//router.Public{}.Load(g)
	//router.Task{}.Load(g)
	//router.Log{}.Load(g)
}