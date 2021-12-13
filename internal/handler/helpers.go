package handler

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	//- "database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"time"
	//- "net/http"
	"bytes"
	"database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/smtp"
	"os"
	"strings"

	conf "github.com/tvitcom/elementarysite/internal/config"
	usermodel "github.com/tvitcom/elementarysite/internal/model/user"
)

type (
	PageData struct {
		Name             string
		Lang             string
		BaseUrl          string
		Seo              *Seo
		Biz              *BizInfo
		User             *UserInfo
		LeftMenuSelected int
		Data             interface{}
	}
	UserInfo struct {
		User_id  int64
		Name     string
		Picture  string
		Roles    []string
		Tel    string
		Email    string
		Birthday string
		Data     interface{}
	}
)

func GetBizInfo() *BizInfo {
	biz := new(BizInfo)
	biz.ShortName = conf.BIZ_SHORTNAME
	biz.Name = conf.BIZ_NAME
	biz.Email = conf.BIZ_EMAIL
	biz.Tel = conf.BIZ_Tel
	biz.Tel2 = conf.BIZ_Tel2
	return biz
}

//TODO: make normal jsonld builder
func GetSeoDefault(c *gin.Context, title, descr, typ, url, logo string) *Seo {
	if url == "" {
		url = c.FullPath()
	}
	if typ == "" {
		typ = "website"
	}
	if descr == "" {
		descr = "simple useful project application for people"
	}
	seo := new(Seo)
	seo.Jsonld = ""
	seo.Og = new(OG)
	seo.Og.Title = title
	seo.Og.Description = descr
	seo.Og.Type = typ
	seo.Og.Url = url
	seo.Og.Image = logo
	seo.Keywords = "simple useful project application for peaple"
	seo.Description = descr
	return seo
}

//TODO : make normal getting browser lang
func GetBrowserLang() string {
	return "en"
}

func SendSmtpMailByPlainAuth(smtpHost, smtpPort, fromMail, fromPassword, toMail, tplFile string, mailData interface{}) error {
	to := []string{toMail}
	t, err := template.ParseFiles(tplFile)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, mailData); err != nil {
		return err
	}
	auth := smtp.PlainAuth("", fromMail, fromPassword, smtpHost)
	addr := smtpHost + ":" + smtpPort
	return smtp.SendMail(addr, auth, fromMail, to, buf.Bytes())
}

func SessionUserId(c *gin.Context) (string) {
		session := sessions.Default(c)
		user_id := session.Get("user_id")
		uid, _ := user_id.(string)
		return uid
}

func GetUserInfo(c *gin.Context, db *sql.DB) *UserInfo {
	info := new(UserInfo)

	// get current user_id
	uid := SessionUserId(c)
	// userId, _ := uid.(int64)
	// get user info from db
	u := usermodel.GetUserById(db, uid)
	info.Name = u.Name
	info.Tel = u.Tel
	info.Email = u.Email
	info.Picture = GetUserpicURL(u.Picture)
	return info
}

func getSigninURL(state string) string {
	return oauthConf.AuthCodeURL(state)
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func renderOnAccept(c *gin.Context, statusCode int, templateName string, datas gin.H) {

	loggedIn, _ := c.Get("user_id")
	if loggedIn != nil {
		datas["is_logged"] = loggedIn.(bool)
	} else {
		datas["is_logged"] = false
	}
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(statusCode, datas["data"])
	case "application/xml":
		// Respond with XML
		c.XML(statusCode, datas["data"])
	default:
		// Respond with HTML
		c.HTML(statusCode, templateName, datas)
	}
}

// Google OAUTH clients functions
func RandToken32bit() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func GetMD5Hash(text string) string {
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func GetSha256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetSaltedSha256(salt, password string) string {
	return GetSha256(salt + password)
}

func IsEqualSaltedSha256(salt, password, hashDatabase string) bool {
	saltedPasswordHash := GetSaltedSha256(salt, password)
	if saltedPasswordHash == hashDatabase {
		return true
	}
	return false
}

func TelNormalisation(n string) string {
	//remove `()-\b+` in Tel number string
	r := strings.NewReplacer("-", "", "(", "", ")", "", "+", "", " ", "")
	return r.Replace(n)
}

func GetUserpicURL(filepath string) string {
	if filepath == "" {
		filepath = "0_60x60.jpeg"
	}
	return "/userpic/" + filepath
}

// import("strings")
func GetEmailsAlias(email string) string {
	splitted := strings.Split(email, "@")
	return splitted[0]
}

// Parse exif data in image-bytes and get datetime string if exist.
// If not - get file datetime info. Output photo format is:
// IMG_20200301_175147.jpg
func getDateTimeFromExif(imageBytes []byte) string {
	curr_dt := time.Now().Format("20060102_150405")
	exif.RegisterParsers(mknote.All...)
	imgReader := bytes.NewReader(imageBytes)
	x, err := exif.Decode(imgReader)
	if err != nil {
		log.Println(err)
		return curr_dt
		// dt, err := x.Get("DateTime")
		// DD("Exif:Alternate:Datetime:", dt)
	}
	dt, err := x.DateTime()
	if err != nil {
		log.Println(err)
		return curr_dt
	}
	return dt.Format("20060102_150405")
}

// Make new dir with name as the eid in current prevPath directory
// For example prevPath="upload/", uid="123"
// will make new directory "upload/123"
func MakeUploadDirByUserId(prevPath, uid string) string {
	newDir := prevPath + uid
	dir, err := os.Open(prevPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()
	_, err = os.Stat(newDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(newDir, 0766)
		if errDir != nil {
			log.Fatal(err)
		}
	}
	return newDir
}

func getImageSize(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
