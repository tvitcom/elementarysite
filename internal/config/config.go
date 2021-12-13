package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	APP_MODE             string
	GIN_MODE             string
	MAIL_SMTPHOST        string
	MAIL_SMTPPORT        string
	MAIL_USERNAME        string
	MAIL_PASSWORD        string
	KEY32                string
	LOGS_PATH            string
	UI_THEME             string
	SSLKEYS_PATH         string
	WEBSERV_NAME         string
	BROWSER_LANG         string
	STORAGE_DRV          string
	STORAGE_DSN          string
	APP_ENTRYPOINT       string
	APP_BRANDNAME        string
	APP_FQDN             string
	APP_FQDN_ADDITIONAL  string
	APP_FQDN_ADDITIONAL2 string
	APP_SSL_ENTRYPOINT   string
	GOOGLE_CREDFILE      string
	GOOGLE_REDIRURL      string
	ADMIN_USER           string
	ADMIN_EMAIL          string
	ADMIN_PASSWORD       string
	BIZ_NAME             string
	BIZ_SHORTNAME        string
	BIZ_EMAIL            string
	BIZ_Tel            string
	BIZ_Tel2           string
	BIZ_LOGO             string
	UPLOAD_DIR           string
)

func init() {
	// Check is jpegoptim installed in system
	lsCmd := exec.Command("jpegoptim", "--version")
	lsOut, err := lsCmd.Output()
	if err != nil {
		log.Fatal("Env command dont respond")
	}
	cmdStatus := string(lsOut)
	if strings.Contains(cmdStatus, "linux") && strings.Contains(cmdStatus, "Kokkonen") {
		fmt.Println("OS ok and jpegoptim ok")
	} else {
		log.Fatal("Environment dont requirement prerequisits.txt file")
	}

	// loading .env settings
	err = godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	APP_MODE = os.Getenv("APP_MODE")
	GIN_MODE = os.Getenv("GIN_MODE")
	UI_THEME = os.Getenv("UI_THEME")           //Default "themeweb"
	APP_BRANDNAME = os.Getenv("app_brandname") // should be withoud finalize dot
	APP_FQDN = os.Getenv("app_fqdn")           // should be withoud finalize dot
	APP_FQDN_ADDITIONAL = os.Getenv("app_fqdn_additional")
	APP_FQDN_ADDITIONAL2 = os.Getenv("app_fqdn_additional2")
	LOGS_PATH = os.Getenv("logs_path")
	WEBSERV_NAME = os.Getenv("webserv_name")
	SSLKEYS_PATH = os.Getenv("sslkeys_path")
	MAIL_SMTPHOST = os.Getenv("mail_smtphost")
	MAIL_SMTPPORT = os.Getenv("mail_smtpport")
	MAIL_USERNAME = os.Getenv("mail_username")
	MAIL_PASSWORD = os.Getenv("mail_password")
	GOOGLE_CREDFILE = os.Getenv("google_credential_file")
	GOOGLE_REDIRURL = "https://" + APP_FQDN + os.Getenv("google_redirect_path")
	APP_ENTRYPOINT = os.Getenv("app_entrypoint")
	APP_SSL_ENTRYPOINT = os.Getenv("app_ssl_entrypoint")
	ADMIN_USER = os.Getenv("admin_user")
	ADMIN_EMAIL = os.Getenv("admin_email")
	ADMIN_PASSWORD = os.Getenv("admin_password")
	STORAGE_DRV = os.Getenv("db_type")
	STORAGE_DSN = os.Getenv("db_user") + ":" + os.Getenv("db_pass") + "@/" + os.Getenv("db_name") + "?parseTime=true"
	KEY32 = os.Getenv("app_secret_key")
	BIZ_NAME = os.Getenv("biz_name")
	BIZ_SHORTNAME = os.Getenv("biz_shortname")
	BIZ_EMAIL = os.Getenv("biz_email")
	BIZ_Tel = os.Getenv("biz_Tel")
	BIZ_Tel2 = os.Getenv("biz_Tel2")
	BIZ_LOGO = os.Getenv("biz_logo")
	UPLOAD_DIR = "filestore/uploaded/"

}

// print debug indormation only if debug mode available
func DD(args ...interface{}) {
	if APP_MODE == "debug" {
		debugee := []interface{}{0: "DEBUG:"}
		args = append(debugee, args)
	}
	fmt.Println(args...)
}
