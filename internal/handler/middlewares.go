package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	// "log"
	"fmt"
	"net/http"
	"time"

	conf "github.com/tvitcom/elementarysite/internal/config"
)

/*
//haha ok!
//Just use

var w http.ResponseWriter = c.Writer

var req *http.Request = c.Request

*/
func mwMonitoringUA() gin.HandlerFunc {
	return func(c *gin.Context) {
		if conf.APP_MODE == "monitoring" {
			// fmt.Sprintf("\n%s - [%s] %s \"%s\" %d \"%s\" %s\n",
			LoginPostForm, _ := c.GetPostForm("email")
			fmt.Println(
				time.Now().Format("2006-01-02 15:04:05"),
				c.Request.RemoteAddr,
				c.Request.Proto,
				c.Request.Method,
				c.Request.RequestURI,
				c.Request.UserAgent(),
				// c.Params,
				// c.Request.Form,
				`"`+LoginPostForm+`"`,
				// c.Request.MultipartForm,
				c.Request.ContentLength,
			)
		}
		c.Next()
	}
}

func mwIsUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user_id")
		if v != nil /*&& v.(int64) > 0*/ {
			fmt.Println("DEBUG:mwIsUser() Now user logged:v.(int64)= ", v)
			c.Redirect(http.StatusMovedPermanently, "/user/")
			c.Abort()
		}
		c.Next()
	}
}

func mwIsNotUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user_id")
		if v == nil /*&& v.(int64) > 0*/ {
			fmt.Println("DEBUG:mwIsNotUser() v=nil")
			c.Redirect(http.StatusMovedPermanently, "/auth/login.html")
			c.Abort()
		}
		c.Next()
	}
}

func confCORS(c *gin.Context) {
	// c.Header("server", WEBSERV_NAME)
	// Content-Security-Policy:

	c.Header("Cache-Control", "max-age=31536000") // Suggest by lighthouse for cache policy 
	c.Header("X-Powered-By", "PHP/8.1.12(joke)") // Its a joke for scriptkiddies
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("X-Frame-Options", "SAMEORIGIN") // Its a joke for scriptkiddies
	c.Header("X-Content-Type-Options", "nosniff")
	// c.Header("feature-policy", `
	// 	geolocation 'none'; 
	// 	midi 'none'; 
	// 	microphone 'none'; 
	// 	camera 'none'; 
	// 	magnetometer 'none'; 
	// 	gyroscope 'none'; 
	// 	speaker 'none'; 
	// 	fullscreen 'self'; 
	// 	payment 'none'
	// `)
	c.Header("Referrer-Policy", "no-referrer")
	c.Header("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS") // Suggest by lighthouse for cache policy 
	c.Header("Content-Security-Policy", `
base-uri 'self';
default-src 'none';
connect-src 'self' https://www.google-analytics.com/;
font-src 'self' https://fonts.gstatic.com/ https://fonts.googleapis.com/;
frame-src 'self' https://www.google.com/recaptcha/ https://www.google.com/maps/ youtube.com www.youtube.com youtu.be;
frame-ancestors 'self' youtu.be youtube.com www.youtube.com;
form-action 'self';
img-src 'self' https://lh3.googleusercontent.com/ https://images.unsplash.com data: blob: https://source.unsplash.com;
object-src 'none';
script-src 'self' 'unsafe-inline' www.googletagmanager.com;
style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;
worker-src 'none';
report-uri /cspcollector;
    `)

	if c.Request.Method == "OPTIONS" {
		if len(c.Request.Header["Access-Control-Request-Headers"]) > 0 {
			c.Header("Access-Control-Allow-Headers", c.Request.Header["Access-Control-Request-Headers"][0])
		}
		c.AbortWithStatus(http.StatusOK)
	}
}

func mwCaptcha() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var w http.ResponseWriter = c.Writer
		// var req *http.Request = c.Req
		// captcha.Server(captcha.StdWidth, captcha.StdHeight)
		// before request

		c.Next()

		// after request
	}
}
