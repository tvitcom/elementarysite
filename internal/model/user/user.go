package user

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"

	conf "github.com/tvitcom/elementarysite/internal/config"
)

var (
	db *sql.DB
)

const (
	ADMIN  = "adm"
	MODER  = "moder"
	USER   = "guest"
)

//Define modules behaviors:
type (
	UserTemplate struct {
		Name, Tel, Email, Password_hash, Approve_token string
		Created int64
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
	User struct {
		UserId int64
		Name string
		Email string
		Tel  string
		Impp  string
		AuthKey      string
		PasswordHash string
		ApproveToken string
		Created       string
		Lastlogin     string
		Roles string
		Notes string
		Picture string
		Address string
		TosAgree string
		PrivacyAgree string
	}
	Profile struct {
		Profile_id             int64
		UserId                int64
		Avatar_mediacontent_id int64
		Meet_for               string
		Country                string
		Province               string
		City                   string
		Town_citizen           int
		Job                    string
		Auto                   string
		Hobby                  string
		Childs                 string
		Pets                   string
		Dob                    int64
		Gender                 int
		Created               int64
		Likes                  int64
		Moderator_id           int64
	}
)

func NewConnection(drv, dsn string) (*sql.DB, error) {
	db, err := sql.Open(drv, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	return db, nil
}

func IsUserByTelnumber(db *sql.DB, Tel string) int64 {
	sqlStatement := `SELECT user_id FROM user WHERE Tel = ?`
	var u User
	row := db.QueryRow(sqlStatement, Tel)
	err := row.Scan(
		&u.UserId,
	)
	switch err {
	case sql.ErrNoRows:
		return 0
	case nil:
		return u.UserId
	default:
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return 0
}

func GetUserById(db *sql.DB, id string) User {
	sqlStatement := `SELECT 
	user_id, name, email, tel, 
	impp, auth_key, password_hash, approve_token, 
	created, lastlogin, roles, notes, 
	picture, address, tos_agree, privacy_agree
	FROM user WHERE user_id =  ? limit 1`
	var u User
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(
		&u.UserId,
		&u.Name,
		&u.Email,
		&u.Tel,
		&u.Impp,
		&u.AuthKey,
		&u.PasswordHash,
		&u.ApproveToken,
		&u.Created,
		&u.Lastlogin,
		&u.Roles,
		&u.Notes,
		&u.Picture,
		&u.Address,
		&u.TosAgree,
		&u.PrivacyAgree,
	)
	switch err {
	case sql.ErrNoRows:
		if conf.GIN_MODE == "debug" {
			conf.DD("GetUserById:", "No rows were returned!")
		}
		return u
	case nil:
		return u
	default:
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return u
}

func GetUserByEmail(db *sql.DB, email string) User {
	sqlStatement := `SELECT 
	user_id, name, email, tel, 
	impp, auth_key, password_hash, approve_token, 
	created, lastlogin, roles, notes, 
	picture, address, tos_agree, privacy_agree
	FROM user WHERE email =  ? limit 1`
	var u User
	row := db.QueryRow(sqlStatement, email)
	err := row.Scan(
		&u.UserId,
		&u.Name,
		&u.Email,
		&u.Tel,
		&u.Impp,
		&u.AuthKey,
		&u.PasswordHash,
		&u.ApproveToken,
		&u.Created,
		&u.Lastlogin,
		&u.Roles,
		&u.Notes,
		&u.Picture,
		&u.Address,
		&u.TosAgree,
		&u.PrivacyAgree,
	)
	switch err {
	case sql.ErrNoRows:
		if conf.GIN_MODE == "debug" {
			conf.DD("GetUserByEmail", "No rows were returned!")
		}
		return u
	case nil:
		return u
	default:
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return u
}

func GetUserByEmailRoles(db *sql.DB, email, roles string) User {
	sqlStatement := `SELECT 
	user_id, name, email, tel, 
	impp, auth_key, password_hash, approve_token, 
	created, lastlogin, roles, notes, 
	picture, address, tos_agree, privacy_agree
	FROM user WHERE email = ? and roles = ?`
	var u User
	row := db.QueryRow(sqlStatement, email, roles)
	err := row.Scan(
		&u.UserId,
		&u.Name,
		&u.Email,
		&u.Tel,
		&u.Impp,
		&u.AuthKey,
		&u.PasswordHash,
		&u.ApproveToken,
		&u.Created,
		&u.Lastlogin,
		&u.Roles,
		&u.Notes,
		&u.Picture,
		&u.Address,
		&u.TosAgree,
		&u.PrivacyAgree,
	)
	switch err {
	case sql.ErrNoRows:
		if conf.GIN_MODE == "debug" {
			conf.DD("GetUserByEmailRoles:", "No rows were returned!")
		}
		return u
	case nil:
		return u
	default:
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return u
}

func GetAllUsers(db *sql.DB) ([]*User, error) {
	rows, err := db.Query(`SELECT 
			user_id, name, email, tel, 
			impp, auth_key, password_hash, approve_token, 
			created, lastlogin, roles, notes, 
			picture, address, tos_agree, privacy_agree 
		FROM user`)
	if err != nil {
		// return nil, err
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		u := new(User)
		err := rows.Scan(
			&u.UserId,
			&u.Name,
			&u.Email,
			&u.Tel,
			&u.Impp,
			&u.AuthKey,
			&u.PasswordHash,
			&u.ApproveToken,
			&u.Created,
			&u.Lastlogin,
			&u.Roles,
			&u.Notes,
			&u.Picture,
			&u.Address,
			&u.TosAgree,
			&u.PrivacyAgree,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return users, nil
}

// Paged data using:
// storage as provider imported!!!
// limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
// offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
// users, err := storage.GetUsersPaged(db, limit, offset)
func GetUsersPaged(db *sql.DB, limit int64, offset int64) ([]*User, error) {
	rows, err := db.Query(`SELECT 
			user_id, name, email, tel, 
			impp, auth_key, password_hash, approve_token, 
			created, lastlogin, roles, notes, 
			picture, address, tos_agree, privacy_agree 
		FROM user ORDER BY user_id DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		u := new(User)
		err := rows.Scan(
			&u.UserId,
			&u.Name,
			&u.Email,
			&u.Tel,
			&u.Impp,
			&u.AuthKey,
			&u.PasswordHash,
			&u.ApproveToken,
			&u.Created,
			&u.Lastlogin,
			&u.Roles,
			&u.Notes,
			&u.Picture,
			&u.Address,
			&u.TosAgree,
			&u.PrivacyAgree,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	return users, nil
}

// Return rows affected only!!!
func ApproveNewUser(db *sql.DB, email, token, setRole string) int64 {
	stmt, err := db.Prepare(`update user set approve_token = "", roles = ? where email = ? and approve_token = ?`)
	defer stmt.Close()
	if err != nil {
		log.Panicln(err)
	}
	res, err := stmt.Exec(setRole, email, token)
	if err != nil {
		log.Panicln(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Panicln(err)
	}
	return rowCnt
}

func NewUserWithApproveToken(db *sql.DB, u User) (newApproveToken string, lastId int64) {
	stmt, err := db.Prepare(`
		INSERT INTO user(
			user_id, name, email, tel, password_hash, approve_token, Created, roles) 
		VALUES
		(null, ?, ?, ?, ?, ?, ?, ?);
    `)
	defer stmt.Close()
	if err != nil {
		log.Panicln(err)
	}
	res, err := stmt.Exec(
		u.Name, 
		u.Email, 
		u.Tel, 
		u.PasswordHash, 
		u.ApproveToken, 
		u.Created, 
		u.Roles,
	)
	if err != nil {
		log.Panicln(err)
	}
	lastId, err = res.LastInsertId()
	if err != nil {
		log.Panicln(err)
	}
	// rowCnt, err := res.RowsAffected()
	// if err != nil {
	// 	log.Panicln(err)
	// }
	// log.Println("RowsAffected():",rowCnt)
	return u.ApproveToken, lastId
}

func UpdateUserInfo(db *sql.DB, u User) (rowCnt int64) {
	stmt, err := db.Prepare(`UPDATE user SET 
    	name=?, email= ?, Tel=?, password_hash=?, picture=?, birthday=?, Created=? where user_id=?;
    `)
	defer stmt.Close()
	if err != nil {
		log.Panicln(err)
	}
	res, err := stmt.Exec(u.Name, u.Email, u.Tel, u.PasswordHash, u.Picture, u.Created, u.UserId)
	if err != nil {
		log.Panicln(err)
	}
	rowCnt, err = res.RowsAffected()
	if err != nil {
		log.Panicln(err)
	}
	return rowCnt
}
