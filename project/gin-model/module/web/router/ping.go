package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zqlpaopao/tool/gin-model/common"
	"net/http"
)

type Ping struct {

}

func (_ Ping) Load(r *gin.Engine){
	l := r.Group("/")
	l.POST("ping", func(g *gin.Context) {
		fmt.Println(g.GetString(common.ContextKey))
		g.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"code":    200,
			"data":    map[string]string{
				"aaaa":"ssss",
			},
		})
		return
	})

}