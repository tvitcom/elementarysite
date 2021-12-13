package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"html"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	// "github.com/dchest/captcha"

	"github.com/cnjack/throttle"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
	brotli "github.com/anargu/gin-brotli"

	conf "github.com/tvitcom/elementarysite/internal/config"
	mediamodel "github.com/tvitcom/elementarysite/internal/model/mediacontent"
	product "github.com/tvitcom/elementarysite/internal/model/product"
	usermodel "github.com/tvitcom/elementarysite/internal/model/user"
)

var (
	cliCredential GoogleOauth2ClientCredentials
	oauthConf     *oauth2.Config
	oauthState    string
	DbConn        *sql.DB
)

type (
	// Types for authentiacted user by google oauth2
	GoogleOauth2ClientCredentials struct {
		Cid     string `json:"cid"`
		Csecret string `json:"csecret"`
	}
	GoogleUser struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Locale        string `json:"gender"`
	}
	BizInfo struct {
		Name, ShortName, Email, Tel, Tel2 string
	}
	Seo struct {
		Jsonld      string
		Og          *OG
		Description string
		Keywords    string
	}
	OG struct { // Open Graph basic struct for SEO
		Title, Description, Type, Url, Image string
	}
)

func init() {
	//Google oauth2 init
	googlecredential, err := ioutil.ReadFile("./data/credentials.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(googlecredential, &cliCredential)

	oauthConf = &oauth2.Config{
		ClientID:     cliCredential.Cid,
		ClientSecret: cliCredential.Csecret,
		RedirectURL:  conf.GOOGLE_REDIRURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",   // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
			"https://www.googleapis.com/auth/userinfo.profile", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
	DbConn, err = usermodel.NewConnection(conf.STORAGE_DRV, conf.STORAGE_DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem with DbConn connection: %s\n", err)
		os.Exit(1)
	}
}

func SetupRoutes(gin_mode string) *gin.Engine {
	fmt.Println("get parsed GIN_MODE", conf.GIN_MODE)
	fmt.Println("Setup param gin_mode:", gin_mode)
	if gin_mode == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else if gin_mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	fmt.Println("GIN_MODE:", gin.Mode(), "APP_MODE:", conf.APP_MODE)

	router := gin.New()

	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Server settings
	router.Delims("%%", "%%")
	router.LoadHTMLGlob("web/templates/**/*")

	// Define common middlewares and routes
	router.Use(gin.Recovery())
	router.Use(confCORS)
	router.Use(brotli.Brotli(brotli.DefaultCompression))

	if gin.Mode() == gin.ReleaseMode {
		gin.DisableConsoleColor()
		// Sett log format:
		f, _ := os.Create(conf.LOGS_PATH)
		gin.DefaultWriter = io.MultiWriter(f)
		router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			//custom format for logging:
			return fmt.Sprintf("%s - [%s] %s \"%s\" %d \"%s\" %s\n",
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.ClientIP,
				param.Method,
				param.Path,
				param.StatusCode,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		}))
	}

	// Session store init
	sessionStore := cookie.NewStore([]byte("secretu"), []byte(conf.KEY32))
	router.Use(sessions.Sessions("bin", sessionStore)) //Название ключа в куках

	// Static assets
	router.Static("/assets", "web/assets")
	router.Static("/public", conf.UPLOAD_DIR)
	router.StaticFile("/robots.txt", "web/static/robots.txt")
	router.StaticFile("/sitemap.txt", "web/static/sitemap.txt")
	router.StaticFile("/humans.txt", "web/static/humans.txt")
	router.StaticFile("/favicon.ico", "web/assets/favicon.ico")
	router.StaticFile("/app.webmanifest", "web/static/app.webmanifest")
	router.StaticFile("/GDPR_POLICY_RU.txt", "GDPR_POLICY_RU.txt")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/index.html")
		return
	})
	router.GET("/index.html", func(c *gin.Context) {
		products := product.GetThreeProducts(DbConn, 14, 15, 16) // Predefined items ids
		proddesc := product.GetProductDetailGroupWithDetail(DbConn, 16)
		renderOnAccept(c, http.StatusOK, "devcourse/index.html", gin.H{
			"gridsproducts": products,
			"proddesc": proddesc,
		})
	})
	router.GET("/advert.html", func(c *gin.Context) {
		renderOnAccept(c, http.StatusOK, "devcourse/advert.html", gin.H{})
	})
	router.GET("/order.html", func(c *gin.Context) {
		// renderOnAccept(c, http.StatusOK, "devcourse/order.html", gin.H{})
		c.Redirect(http.StatusMovedPermanently, "/soon.html")
		return
	})
	router.POST("/cspcollector", func(c *gin.Context) {
		cspdata, err := c.GetRawData()
		if err != nil {
			fmt.Println("SCP report error")
		}
		conf.DD("ContentSecurePolicy:", string(cspdata))
	})

	router.POST("/order.php", func(c *gin.Context) {
		//Normalisation users parameters
		rawName := strings.TrimSpace(c.PostForm("given-name"))
		rawTel := strings.TrimSpace(c.PostForm("tel"))
		rawMsg := strings.TrimSpace(c.PostForm("msg"))

		normPhone := TelNormalisation(strings.TrimSpace(rawTel))
		var validTelStr = regexp.MustCompile(`[0-9]{9,12}`)
		if ok := validTelStr.MatchString(normPhone); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"warn": "invalid phone number"})
			return
		}
		var validNameStr = regexp.MustCompile(`[a-zA-Zа-яёА-ЯЁїЇіІґҐ\s\-]{2,45}`)
		if ok := validNameStr.MatchString(rawName); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"warn": "bad Name. Will min 2 max 150 symbols"})
			return
		}
		msg := html.EscapeString(rawMsg)
		fmt.Println("rawMsg:", rawMsg, len(rawMsg), "msg:", msg)
		// Send link to user email:
		tplFile := "web/templates/email/order_by_tel.emlt"
		emailData := struct {
			From  string
			Brand string
			Tel  string
			Name  string
			Note string
		}{
			From:       conf.MAIL_USERNAME,
			Brand:      conf.APP_BRANDNAME,
			Tel:        rawTel,
			Name:       rawName,
			Note:       rawMsg,
		}
		err := SendSmtpMailByPlainAuth(conf.MAIL_SMTPHOST, conf.MAIL_SMTPPORT, conf.MAIL_USERNAME, conf.MAIL_PASSWORD, conf.ADMIN_EMAIL, tplFile, emailData)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Emailed to admin email: ok")
		}

		renderOnAccept(c, http.StatusOK, "devcourse/thanks.html", gin.H{"msg": "в ближайшее время позвоню насчёт курса."})
	})
	router.GET("/reply.html", func(c *gin.Context) {
		renderOnAccept(c, http.StatusOK, "devcourse/reply.html", gin.H{})
	})
	router.GET("/support.html", func(c *gin.Context) {

		renderOnAccept(c, http.StatusOK, "devcourse/support.html", gin.H{})
	})
	router.POST("/support.php", func(c *gin.Context) {
		//Normalisation users parameters
		rawSubject := strings.TrimSpace(c.PostForm("subject"))
		rawName := strings.TrimSpace(c.PostForm("given-name"))
		rawEmail := strings.TrimSpace(c.PostForm("email"))
		img, err := imageupload.Process(c.Request, "support_picture")
		if err != nil {
			c.JSON(http.StatusBadRequest, "image not submited")
			panic(err)
		}
		rawMsg := strings.TrimSpace(c.PostForm("msg"))
		
		//validation
		var validSubjectStr = regexp.MustCompile(`[a-z]{3,6}`)
		if ok := validSubjectStr.MatchString(rawSubject); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"warn": "bad Subject. Will min 3 max 6 symbols"})
			return
		}
		var validNameStr = regexp.MustCompile(`[a-zA-Zа-яёА-ЯЁїЇіІґҐ\s\-]{2,45}`)
		if ok := validNameStr.MatchString(rawName); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"warn": "bad Name. Will min 2 max 150 symbols"})
			return
		}
		var validEmailStr = regexp.MustCompile(`[_a-zA-Z0-9\.]{2,60}@[-a-zA-Z0-9]{2,56}.[a-zA-Z]{2,6}`)
		if ok := validEmailStr.MatchString(rawEmail); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"warn": "bad Email address"})
			return
		}
		var validMsgStr = regexp.MustCompile(`[0-9a-zA-Zа-яёА-ЯЁїЇіІґҐ\s\-\.\,\:]{10,512}`)
		if ok := validMsgStr.MatchString(rawMsg); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"warn": "bad Message text. Will min 10 max 512 symbols"})
			return
		}

		// Follow by column table: 1-pic 2-vid 3-aud 4-vec
		var fileExt string
		var mtype int8
		switch img.ContentType {
		case "image/png":
			fileExt = ".png"
			mtype = 1
		case "image/jpeg":
			fileExt = ".jpg"
			mtype = 1
		case "video/mp4":
			mtype = 2
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"warn": "file has unsupport type"})
			return
		default:
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"warn": "file has unrecognized type"})
			return
		}

		dt_filename := getDateTimeFromExif(img.Data)

		// thumb, err := imageupload.ThumbnailJPEG(img, 1280, 1280, 50)
		// if err != nil {
		// 	panic(err)
		// }
		uploadsDir := "filestore/uploaded/"
		user_id := "0"
		newPath := MakeUploadDirByUserId(uploadsDir, user_id)
		filename := "IMG_" + dt_filename + fileExt
		err = img.Save(newPath + "/" + filename)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInsufficientStorage, gin.H{"warn": "server saving error"})
			return
		}
		// thumbFilename := newPath + "/" + "IMG_" + dt_filename + "_thumb" + fileExt
		// err = thumb.Save(thumbFilename)
		// if err != nil {
		// 	log.Panic(err)
		// 	c.JSON(http.StatusInsufficientStorage, gin.H{"warn": "server saving error"})
		// 	return
		// }

		lsCmd := exec.Command("bash", "-c", "file "+filename)
		lsOut, err := lsCmd.Output()
		if err != nil {
			panic(err)
		}
		fileStatus := string(lsOut)

		// test if progressive been formatted
		if !strings.Contains(fileStatus, "progressive") {
			ConvertingTimeout := 5 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), ConvertingTimeout)
			defer cancel()
			if err := exec.CommandContext(ctx, "jpegoptim", "--strip-all", "--all-progressive", "-ptm85", filename).Run(); err != nil {
				log.Println("SUCCESS : fail CAUSE: ConvertingTimeout occured")
			}
		}
		_ = mtype
		media := new(mediamodel.Mediacontent)
		media.Event_id = "0"
		media.UserId = user_id
		media.Filename = filename
		media.Description = ""
		media.Mediatype = "1"
		media.Storage_host = "localhost"
		media.Created = time.Now().Format("2006-01-02 15:04:05")
		media.Uploader_useragent = c.Request.Header.Get("User-Agent")
		media.Uploader_ip = c.ClientIP()
		media.Moderator_id = "0"//Status 0 for moderation
		lastId := mediamodel.NewMediacontent(DbConn, media)
		if lastId == 0 {
			c.JSON(http.StatusInsufficientStorage, gin.H{"warn": "media dont save"})
		}
		
		msg := html.EscapeString(rawMsg)
		//Order send to admin by mail:
		// Send link to user email:
		tplFile := "web/templates/email/support.emlt"
		emailData := struct {
			From  string
			Name  string
			Brand string
			Theme  string
			CustomerEmail  string
			Picture  string
			Message string
		}{
			From:       conf.MAIL_USERNAME,
			Name:       rawName,
			Brand:      conf.APP_BRANDNAME,
			Theme:       rawSubject,
			CustomerEmail: rawEmail,
			Picture:  "https://"+conf.APP_FQDN+"/public/"+user_id+"/"+filename,
			Message:  msg,
		}
		err = SendSmtpMailByPlainAuth(conf.MAIL_SMTPHOST, conf.MAIL_SMTPPORT, conf.MAIL_USERNAME, conf.MAIL_PASSWORD, conf.ADMIN_EMAIL, tplFile, emailData)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Sender:",rawEmail,"send support request to admin: ok")
		}

		renderOnAccept(c, http.StatusOK, "devcourse/thanks.html", gin.H{"msg": "обращение будет рассмотрено в ближайшее время"})
	})
	router.GET("/tos.html", func(c *gin.Context) {
		renderOnAccept(c, http.StatusOK, "devcourse/tos.html", gin.H{})
	})

	//auth router group
	authRouter := router.Group("/auth",
		mwMonitoringUA(),
		mwIsUser(),
		throttle.Policy(&throttle.Quota{
			Limit:  5, // 2 ones per 1 minute
			Within: time.Minute,
		}))
	// Login page display
	authRouter.GET("/login.html", func(c *gin.Context) {
		renderOnAccept(c, http.StatusOK, "devcourse/login.html", gin.H{})
	})

	// signup
	authRouter.GET("/signup.html", func(c *gin.Context) {
		referral := c.DefaultQuery("ref", "0")

		//Captcha init
		// captchaSign := struct {
		// 	CaptchaId string
		// }{
		// 	captcha.New(),
		// }

		var GoogleOAuthLink string = "/soon"
		oauthState = RandToken32bit()
		session := sessions.Default(c)
		session.Set("state", oauthState)
		session.Save()
		GoogleOAuthLink = oauthConf.AuthCodeURL(oauthState)
		renderOnAccept(c, http.StatusOK, "devcourse/signup.html", gin.H{
			"gOauthLink": GoogleOAuthLink,
			// "CaptchaId": captchaSign,
			"_csrf": RandToken32bit(),
			"ref":   referral,
		})
	}, mwCaptcha())

	authRouter.POST("/signup", func(c *gin.Context) {

		//Validation
		// if !captcha.VerifyString(f.CaptchaId, f.CaptchaSolution) {
		// 	c.String(403, "Wrong captcha solution! No robots allowed!\n")
		// 	return
		// }

		//Normalisation users parameters
		_ = c.PostForm("_csrf")
		_ = c.GetInt64("ref") //fRef
		// fCapcha := c.PostForm("captcha")
		fMeetFor := c.PostForm("meet_for")
		fGender := c.GetInt("gender")
		fName := c.PostForm("name")
		fEmail := c.PostForm("email")
		fPassword := c.PostForm("password")
		fPasswordRepeat := c.PostForm("password_repeat")
		fPictureFile, _ := c.FormFile("picture")
		if fPassword != fPasswordRepeat {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The passwords don't equal"})
			return
		}

		var validEmailStr = regexp.MustCompile(`[_a-zA-Z0-9]{2,60}@[-a-zA-Z0-9]{2,56}.[a-zA-Z]{2,6}`)
		if ok := validEmailStr.MatchString(fEmail); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"warn": "invalid email address"})
			return
		}

		//Check if values yet into DbConn:
		if isUser := usermodel.GetUserByEmail(DbConn, fEmail); isUser.UserId > 0 {
			c.JSON(http.StatusPreconditionFailed, gin.H{"warn": "bad Email for registration"})
			return
		}

		u := usermodel.User{}
		p := usermodel.Profile{}

		if fName == "" {
			u.Name = GetEmailsAlias(fEmail)
		} else {
			u.Name = fName
		}

		u.Email = fEmail
		u.Tel = ""
		u.Impp = ""
		u.AuthKey = ""
		u.PasswordHash = GetSaltedSha256(conf.KEY32, fPassword)
		u.Created = time.Now().Format("2006-01-02 15:04:05")
		u.ApproveToken = GetSaltedSha256(conf.KEY32, u.Email)
		u.Roles = usermodel.USER
		u.Notes = "from google login"
		link, uId := usermodel.NewUserWithApproveToken(DbConn, u)

		p.UserId = uId
		p.Meet_for = fMeetFor
		p.Gender = fGender

		// Upload the file to specific dst.
		fileExt := ".jpg"
		uIdStr := strconv.FormatInt(uId, 10)
		_ = c.SaveUploadedFile(fPictureFile, "./filestore/uploaded/"+uIdStr+fileExt)

		//Approve_token by mail:
		// Send link to user email:
		tplFile := "./templates/email/approvement.emlt"
		emailData := struct {
			From       string
			Brand      string
			Name       string
			ApproveURL string
		}{
			From:       conf.MAIL_USERNAME,
			Brand:      conf.APP_BRANDNAME,
			Name:       u.Name,
			ApproveURL: "https://" + conf.APP_FQDN + "/auth/approvement/" + link + "?signin=" + u.Email,
		}
		err := SendSmtpMailByPlainAuth(conf.MAIL_SMTPHOST, conf.MAIL_SMTPPORT, conf.MAIL_USERNAME, conf.MAIL_PASSWORD, u.Email, tplFile, emailData)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Emailed: ok")
		}
		// TODO: work user without approving own link...maybe captcha need

		//Redirect to user personal message:
		c.Redirect(http.StatusMovedPermanently, "/message?m=You should check your email and following by our link")
		return
	})

	// signin
	authRouter.GET("/signin.html", func(c *gin.Context) {

		//google oauth2 implementation
		var GoogleOAuthLink string = "#"
		oauthState = RandToken32bit()
		session := sessions.Default(c)
		session.Set("state", oauthState)
		session.Save()
		GoogleOAuthLink = getSigninURL(oauthState)
		renderOnAccept(c, http.StatusOK, "devcourse/login.html", gin.H{
			"gOauthLink": GoogleOAuthLink,
			"_csrf":      RandToken32bit(),
		})
	})
	authRouter.POST("/signin.php", func(c *gin.Context) {

		// _ = c.PostForm("_csrf")
		email := c.PostForm("username")
		password := c.PostForm("current-password")
		//validation
		var validStr = regexp.MustCompile(`[_a-zA-Z0-9]{2,60}@[a-zA-Z0-9]{2,56}.[a-zA-Z]{2,6}`)
		if ok := validStr.MatchString(email); !ok {
			c.AbortWithStatus(403)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"warn": "invalid login"})
			return
		}
		u := usermodel.GetUserByEmailRoles(DbConn, email, usermodel.ADMIN)
		
		conf.DD("DEBUG:u=", u, "dbuser.UserId > 0", bool(u.UserId > 0))
		
		if u.UserId == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"warn": "user not available"})
			return
		}
		passwordIsValid := IsEqualSaltedSha256(conf.KEY32, password, u.PasswordHash)
		
		conf.DD("DEBUG:email,password,passwordIsValid:", email, password, passwordIsValid)
		conf.DD("DEBUG:HashToDb=", GetSaltedSha256(conf.KEY32, password))
		
		if u.UserId > 0 && passwordIsValid {
			session := sessions.Default(c)
			session.Set("user_id", u.UserId)
			session.Save()
			c.Redirect(http.StatusMovedPermanently, "/user/index.html")
		} else {
			// c.Redirect(http.StatusMovedPermanently, "/auth/passwordrecover/")
			renderOnAccept(c, http.StatusOK, "error.html", gin.H{})
		}
		return
	})

	// approvement new signed user
	authRouter.GET("/approvement/:secretlink", func(c *gin.Context) {
		signin := c.Query("signin")
		lnk := c.Param("secretlink")
		affected := usermodel.ApproveNewUser(DbConn, signin, lnk, usermodel.USER)
		fmt.Println("DEBUG:ApproveNewUser:return affecte", affected)
		//Set session user_id
		if affected > 0 {
			u := usermodel.GetUserByEmailRoles(DbConn, signin, usermodel.USER)
			if u.UserId > 0 {
				session := sessions.Default(c)
				session.Set("user_id", u.UserId)
				session.Save()
				c.Redirect(http.StatusMovedPermanently, "/user/index.html")
				return
			}
		}
		c.JSON(http.StatusPreconditionFailed, gin.H{"warn": "bad link for approve email"})
		return
	})
	// restorepassword
	authRouter.GET("/passwordrecover", func(c *gin.Context) {
		renderOnAccept(c, http.StatusOK, "forgot-password.html", gin.H{})
	})

	authRouter.POST("/passwordrecover", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/soon")
		return
	})

	authRouter.GET("/passwordrecovered/:secretlink", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/soon")
		return
	})

	//oauth router group
	oauthRouter := router.Group("/oauth", throttle.Policy(&throttle.Quota{
		Limit:  4,
		Within: time.Minute,
	}))
	oauthRouter.GET("/googleuser", func(c *gin.Context) {
		// Request to Google about google account base user information
		// Handle the exchange code to initiate a transport.
		googleUser := GoogleUser{}
		session := sessions.Default(c)
		retrievedState := session.Get("state")
		if retrievedState != c.Query("state") {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
			return
		}

		tok, err := oauthConf.Exchange(oauth2.NoContext, c.Query("code"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		client := oauthConf.Client(oauth2.NoContext, tok)
		gAccount, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer gAccount.Body.Close()
		data, err := ioutil.ReadAll(gAccount.Body)
		if err != nil {
			c.AbortWithError(http.StatusPreconditionFailed, err)
			return
		}
		fmt.Println("gAccount (google says): ", string(data))

		err = json.Unmarshal(data, &googleUser)
		if err != nil {
			log.Fatalln("error:", err)
		}
		fmt.Println("DEBUG:googleUser.email: ", googleUser.Email)

		//Check if user yet not database
		dbuser := usermodel.GetUserByEmail(DbConn, googleUser.Email)
		fmt.Println("DEBUG:googleUser: Check ifUser from DbConn:", dbuser)
		if dbuser.UserId == 0 {
			u := usermodel.User{}
			u.Name = googleUser.Name
			u.Email = googleUser.Email
			u.Tel = ""
			u.Picture = googleUser.Picture
			u.PasswordHash = GetSaltedSha256(conf.KEY32, conf.KEY32) // Dummy password Only for secrecy
			u.Created = time.Now().Format("2006-01-02 15:04:05")
			u.ApproveToken = "" // Dummy token Only for secrecy
			u.Roles = usermodel.USER
			link, user_id := usermodel.NewUserWithApproveToken(DbConn, u)
			conf.DD("DEBUG:/signup POST: u.PasswordHash, link", user_id, link)
			session.Set("user_id", user_id)
		} else {
			session.Set("user_id", dbuser.UserId)
		}
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/user/")
	})

	// userRouter part
	userRouter := router.Group("/user", mwIsNotUser())

	userRouter.GET("/", func(c *gin.Context) {
		//Make Page Data
		c.Redirect(http.StatusMovedPermanently, "/user/index.html")
		return
	})

	userRouter.GET("/index.html", func(c *gin.Context) {
		//Make Page Data
		pd := new(PageData)
		pd.Name = "Dashboard"
		pd.BaseUrl = "/user/"
		pd.Lang = GetBrowserLang()
		pd.Biz = GetBizInfo()
		pd.Seo = GetSeoDefault(c, pd.Name, "", "", "", conf.BIZ_LOGO)
		pd.User = GetUserInfo(c, DbConn)
		renderOnAccept(c, http.StatusOK, "admin/index.html", gin.H{"Page": pd})
	})

	userRouter.GET("/payments.html", func(c *gin.Context) {
		//Make Page Data
		pd := new(PageData)
		pd.Name = "Payments"
		pd.BaseUrl = "/user/"
		pd.Lang = GetBrowserLang()
		pd.Biz = GetBizInfo()
		pd.Seo = GetSeoDefault(c, pd.Name, "", "", "", conf.BIZ_LOGO)
		pd.User = GetUserInfo(c, DbConn)
		renderOnAccept(c, http.StatusOK, "admin/payments.html", gin.H{"Page": pd})
	})
	
	userRouter.GET("/services.html", func(c *gin.Context) {
		//Make Page Data
		pd := new(PageData)
		pd.Name = "services"
		pd.BaseUrl = "/user/"
		pd.Lang = GetBrowserLang()
		pd.Biz = GetBizInfo()
		pd.Seo = GetSeoDefault(c, pd.Name, "", "", "", conf.BIZ_LOGO)
		pd.User = GetUserInfo(c, DbConn)
		renderOnAccept(c, http.StatusOK, "admin/services.html", gin.H{"Page": pd})
	})
	
	userRouter.GET("/topics.html", func(c *gin.Context) {
		//Make Page Data
		pd := new(PageData)
		pd.Name = "Topics"
		pd.BaseUrl = "/user/"
		pd.Lang = GetBrowserLang()
		pd.Biz = GetBizInfo()
		pd.Seo = GetSeoDefault(c, pd.Name, "", "", "", conf.BIZ_LOGO)
		pd.User = GetUserInfo(c, DbConn)
		renderOnAccept(c, http.StatusOK, "admin/topics.html", gin.H{"Page": pd})
	})
	
	userRouter.GET("/users.html", func(c *gin.Context) {
		//Make Page Data
		pd := new(PageData)
		pd.Name = "Users"
		pd.BaseUrl = "/user/"
		pd.Lang = GetBrowserLang()
		pd.Biz = GetBizInfo()
		pd.Seo = GetSeoDefault(c, pd.Name, "", "", "", conf.BIZ_LOGO)
		pd.User = GetUserInfo(c, DbConn)
		renderOnAccept(c, http.StatusOK, "admin/users.html", gin.H{"Page": pd})
	})

	userRouter.GET("/settings.html", func(c *gin.Context) {
		session := sessions.Default(c)
		user_id := session.Get("user_id")

		//Make Page Data
		pd := new(PageData)
		pd.Name = "settings"
		pd.BaseUrl = "/user/"
		pd.Lang = GetBrowserLang()
		pd.Biz = GetBizInfo()
		pd.Seo = GetSeoDefault(c, pd.Name, "", "", "", conf.BIZ_LOGO)
		pd.User = GetUserInfo(c, DbConn)
		photos := mediamodel.GetUsersMediacontents(DbConn, user_id.(int64))
		renderOnAccept(c, http.StatusOK, "admin/settings.html", gin.H{
			"Page":   pd,
			"photos": photos,
		})
	})

	userRouter.POST("/uploadmedia", func(c *gin.Context) {
		// // CHECK AVAILABLE EVENT_ID FOR CURRENT USER!!!
		// event_id := c.Param("event_id")
		// var validEventIdStr = regexp.MustCompile(`[0-9]{1,21}`)
		// if ok := validEventIdStr.MatchString(event_id); !ok {
		// 	c.JSON(http.StatusBadRequest, gin.H{"warn": "invalid event_id parameter"})
		// 	return
		// }

		// eid, err = strconv.ParseInt(event_id, 10, 64)
		// if err != nil {
		// 	panic(err)
		// }

		user_id := SessionUserId(c)
		conf.DD("user_id:", user_id,)
		// ok := storage.IsParticipationAvailabeByEventId(DbConn, uid, eid)
		// if ok == 0 {
		// 	c.JSON(http.StatusForbidden, gin.H{"warn": "forbidden event_id parameter"})
		// 	return
		// }

		img, err := imageupload.Process(c.Request, "file")
		if err != nil {
			c.String(http.StatusUnsupportedMediaType, "Only png, jpeg, jpg files!\n")
			return
		}
		// Follow by column table: 1-pic 2-vid 3-aud 4-vec
		var fileExt string
		var mtype int8
		switch img.ContentType {
		case "image/png":
			fileExt = ".png"
			mtype = 1
		case "image/jpeg":
			fileExt = ".jpg"
			mtype = 1
		case "video/mp4":
			c.String(http.StatusUnsupportedMediaType, "Only png, jpeg, jpg files!\n")
			return
		default:
			c.String(http.StatusUnsupportedMediaType, "Only png, jpeg, jpg files!\n")
			return
		}

		dt_filename := getDateTimeFromExif(img.Data)

		//TODO: lat, lng := getCoordsFromExif(img.Data)
		// DD("lat,lng =", lat, lng)

		newPath := MakeUploadDirByUserId(conf.UPLOAD_DIR, user_id)
		// Save original size
		filename := GetMD5Hash(dt_filename + user_id)
		err = img.Save(newPath + "/" + filename + fileExt)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInsufficientStorage, gin.H{"warn": "server saving error"})
			return
		}
		fw, fh := getImageSize(newPath + "/" + filename + fileExt)
		fmt.Println("DEBUG_IMAGE:", fw, fh)
		// Save thumbnail size
		thumb, err := imageupload.ThumbnailJPEG(img, 1280, 1280, 50)
		if err != nil {
			panic(err)
		}
		thumbFilename := newPath + "/" + filename + fileExt
		err = thumb.Save(thumbFilename)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInsufficientStorage, gin.H{"warn": "server saving error"})
			return
		}

		lsCmd := exec.Command("bash", "-c", "file "+thumbFilename)
		lsOut, err := lsCmd.Output()
		if err != nil {
			panic(err)
		}
		fileStatus := string(lsOut)

		// test if progressive been formatted
		if !strings.Contains(fileStatus, "progressive") {
			ConvertingTimeout := 5 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), ConvertingTimeout)
			defer cancel()
			if err := exec.CommandContext(ctx, "jpegoptim", "--strip-all", "--all-progressive", "-ptm85", thumbFilename).Run(); err != nil {
				log.Println("SUCCESS : fail CAUSE: ConvertingTimeout occured")
			}
		}
		_ = mtype
		media := new(mediamodel.Mediacontent)
		media.Event_id = "1"
		media.UserId = user_id
		media.Filename = filename + fileExt
		media.Description = ""
		media.Mediatype = "1"
		media.Storage_host = "localhost"
		media.Created = time.Now().Format("2006-01-02 15:04:05")
		media.Uploader_useragent = c.Request.Header.Get("User-Agent")
		media.Uploader_ip = c.ClientIP()
		media.Moderator_id = "0" //Status 0 for moderation
		lastId := mediamodel.NewMediacontent(DbConn, media)
		if lastId > 0 {
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", img.Filename))
		} else {
			c.JSON(http.StatusInsufficientStorage, gin.H{"warn": "media dont record into db"})
		}
		return
	})

	userRouter.POST("/profile.html&action=changepassword", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/soon.html")
		return
	})

	// signout:
	userRouter.GET("/signout", func(c *gin.Context) {
		session := sessions.Default(c)
		log.Println("running SignOut for user_id:", session.Get("user_id"))
		session.Set("user_id", nil)
		session.Save()
		session.Clear()
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	})
	router.GET("/development", func(c *gin.Context) {
		renderOnAccept(c, http.StatusOK, "development.html", gin.H{})
	})

	router.GET("/thanks.html", func(c *gin.Context) {
		renderOnAccept(c, http.StatusOK, "devcourse/thanks.html", gin.H{
			"msg": c.Query("m"),
		})
	})
	router.GET("/soon", func(c *gin.Context) {
		renderOnAccept(c, http.StatusServiceUnavailable, "devcourse/soon.html", gin.H{})
	})

	router.NoRoute(func(c *gin.Context) {
		renderOnAccept(c, http.StatusNotFound, "devcourse/error.html", gin.H{})
	})
	router.NoMethod(func(c *gin.Context) {
		//renderOnAccept(c, c *gin.Context, datas gin.H, templateName string)

		// 	cmd/webapp/main.go:571:17: not enough arguments in call to renderOnAccept
		// have (*gin.Context, number, gin.H)
		// want (*gin.Context, int, string, gin.H)

		renderOnAccept(c, http.StatusMethodNotAllowed, "devcourse/error.html", gin.H{
			"err": `405 MethodNotAllowed`,
		})
	})

	return router
}
