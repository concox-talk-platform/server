/**
* @Author: yanKoo
* @Date: 2019/3/11 10:39
* @Description: main
 */
package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
	cfgWs "server/web-api/configs/web_server"
	"server/web-api/controllers"
	"server/web-api/log"
	"strings"
)

func main() {
	engine := Prepare()
	if cfgWs.HttpWay == "https" {
		engine.Use(TlsHandler())
		if err := engine.RunTLS(":"+cfgWs.WebPort, cfgWs.CertFile, cfgWs.KeyFile); err != nil {
			log.Log.Printf("Read pem key file error: %+v", err)
		}
	} else if cfgWs.HttpWay == "http" {
		if err := engine.Run(":" + cfgWs.WebPort); err != nil {
			log.Log.Println("listen is error", err)
		}
	}

}

func Prepare() *gin.Engine {

	// 禁用控制台颜色
	gin.DisableConsoleColor()

	engine := gin.Default()

	// 日志， 解决跨域问题
	engine.Use(log.Logger(), Cors())

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

	engine.POST("/device/update", controllers.UpdateDeviceInfo)

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
			c.Header("Access-Control-Allow-Origin", "*")                                       // 允许访问所有域
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

func TlsHandler() gin.HandlerFunc {
	addr := flag.String("a", "localhost", "ssl 默认主机")
	flag.Parse()
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     *addr + ":" + cfgWs.WebPort,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}