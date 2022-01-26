package router

import (
"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct {

}

func (_ Login) Load(r *gin.Engine){
	l := r.Group("/fire-timer")
	l.POST("/login", func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"code":    200,
			"data":    map[string]string{
				"aa":"bb",
			},
		})
		return
	})

}