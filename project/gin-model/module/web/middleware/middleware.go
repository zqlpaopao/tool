package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/zqlpaopao/tool/gin-model/common"
	strTool "github.com/zqlpaopao/tool/string-byte/src"
	log "github.com/zqlpaopao/tool/zap-log/src"
	"mime/multipart"
	"net/http"
	"time"
)

//Cors 设置请求头信息
func Cors(c *gin.Context) {
	method := c.Request.Method
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//服务器支持的所有跨域请求的方法
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
	//允许跨域设置可以返回其他子段，可以自定义字段
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
	// 允许浏览器（客户端）可以解析的头部 （重要）
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
	//设置缓存时间
	c.Header("Access-Control-Max-Age", "172800")
	//允许客户端传递校验信息比如 cookie (重要)
	c.Header("Access-Control-Allow-Credentials", "true")

	//允许类型校验
	if method == "OPTIONS" {
		c.JSON(http.StatusOK, "ok!")
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic info is: %v", err)
		}
	}()
	c.Next()
}

func mWareInterCheck(g *gin.Context) {
	var (
		tag   bool
		err   error
		token string
	)

	if _, ok := common.EnvConf.Env.Web.NoAuthUrlMap[g.Request.URL.Path]; ok {
		tag = true
	}

	token = g.Request.Header.Get("authorization")
	l := len(common.Bearer)
	if len(token) > l+1 && token[:l] == common.Bearer {
		token = token[l+1:]
	}
	if !tag && token == "" {
		err = errors.New("token is empty")
		goto END
	}

	return
END:
	g.Abort()
	g.JSON(http.StatusOK, gin.H{
		"message": err.Error(),
		"code":    common.CodeErrBackIn,
	})
	//g.Redirect(http.StatusMovedPermanently, g.Request.Host+"/login")

}


type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

//MiddleLog 记录请求、响应日志
func MiddleLog(g *gin.Context){
	// 开始时间
	startTime := time.Now()
	blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: g.Writer}
	g.Writer = blw
	// 处理请求
	g.Next()
	// 结束时间
	endTime := time.Now()
	// 执行时间
	latencyTime := endTime.Sub(startTime)
	//请求参数
	var (
		req string
		err error
	)
	if req ,err = requestParams(g);nil != err{
		req = err.Error()
	}
	reqInfo := &common.ReqLogInfo{
		StartTime: startTime.Format(common.TimeFormat),
		EndTime:   endTime.Format(common.TimeFormat),
		RunTime:   latencyTime,
		ReqMethod: g.Request.Method,
		ReqUrl:    g.Request.RequestURI,
		ClientIP:  g.ClientIP(),
		ReqArgs:   req,
		RespCode:  g.Writer.Status(),
		RespInfo:  blw.body.String(),
	}
	log.InfoAsync(common.WebLogKey,reqInfo).MsgAsync(common.WebLogKey)
}

//requestParams 请求参数整理
func requestParams(g *gin.Context)(str string,err error){
	if g.Request.Method == "GET"{
		return g.Request.RequestURI,nil
	}
	var (
		par *multipart.Form
		b []byte
	)
	if par ,err = g.MultipartForm();nil != err{
		return
	}

	if b, err = json.Marshal(par.Value); nil != err{
		return
	}
	return strTool.Bytes2String(b),nil
}

//InitContext build context.Value
func InitContext()gin.HandlerFunc{
	return func(c *gin.Context) {
		c.Set(common.ContextKey,uuid.NewV4().String())
	}
}

//TimeoutMiddleware middleware wraps the request context with a timeou
func TimeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {

				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
