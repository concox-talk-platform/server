/**
* @Author: yanKoo
* @Date: 2019/3/11 10:39
* @Description: main
 */
package main

import (
	cfgWs "configs/web_server"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"log"
	"net/http"
	"os"
	"server/web/backend/src/controllers"
	"strings"
	"time"
)

func main() {
	engine := Prepare()
	if err := engine.Run(":" + cfgWs.WebPort); err != nil {
		log.Println("listen is error", err)
	}
}

func Prepare() *gin.Engine {

	// 禁用控制台颜色
	gin.DisableConsoleColor()
	// 创建记录日志的文件
	//f, _ := os.Create("backend-web.log")
	//gin.DefaultWriter = io.MultiWriter(f)

	engine := gin.Default()

	// 日志， 解决跨域问题
	engine.Use(Logger(), Cors())

	// runTls error
	/*if err := engine.RunTLS(":8888", "fullchain.pem", "privkey.pem");err != nil {
		log.Println("")
	}*/

	// 注册路由
	// account
	engine.POST("/account/login.do/:account_name", controllers.SignIn)

	engine.POST("/account/logout.do/:account_name", controllers.SignOut)

	engine.POST("/account", controllers.CreateAccountBySuperior)

	engine.GET("/account/:account_name", controllers.GetAccountInfo)

	engine.POST("/account/info/update", controllers.UpdateAccountInfo)

	engine.POST("/account/pwd/update", controllers.UpdateAccountPwd)

	engine.GET("/account_class/:accountId/:searchId", controllers.GetAccountClass)

	engine.GET("/account_device/:accountId/:getAdviceId", controllers.GetAccountDevice)

	engine.POST("/account_device/:accountId", controllers.TransAccountDevice)

	// group
	engine.POST("/group", controllers.CreateGroup)

	engine.POST("/group/update", controllers.UpdateGroup)

	engine.POST("/group/delete", controllers.DeleteGroup)

	engine.POST("/group/devices/update", controllers.UpdateGroupDevice)

	// device
	engine.POST("/device/import/:account_name", controllers.ImportDeviceByRoot)

	// upload file
	engine.POST("/upload", controllers.UploadFile)

	// im server
	engine.GET("/im-server/:accountId", controllers.ImPush)

	return engine
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
func Logger() gin.HandlerFunc {
	logClient := logrus.New()

	//禁止logrus的输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	logClient.Out = src
	logClient.SetLevel(logrus.DebugLevel)
	apiLogPath := "web_server.log"
	logWriter, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(apiLogPath),       // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	logClient.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logClient.Infof("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
