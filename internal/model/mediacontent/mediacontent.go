package mediacontent

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

//Define modules behaviors:
const (
	MEDIACONTENT_USERLIMIT = 8
)

type (
	Mediacontent struct {
		Mediacontent_id    string
		Event_id           string
		UserId            string
		Filename           string
		Description        string
		Mediatype          string
		Storage_host       string
		Created           string
		Moderator_id       string
		Uploader_useragent string
		Uploader_ip        string
	}
	MediacontentAuthored struct {
		Mediacontent_id    int64
		Event_id           int64
		Author             string
		Filename           string
		Description        string
		Mediatype          string
		Storage_host       string
		Created           int64
		Moderator_id       int64
		Uploader_useragent string
		Uploader_ip        string
	}
)

func GetUsersMediacontents(db *sql.DB, uid int64) []*Mediacontent {
	rows, err := db.Query(`
		SELECT m.mediacontent_id, m.event_id, m.user_id, m.filename, 
		m.description, m.mediatype, 
		m.storage_host, m.created, m.uploader_useragent, m.uploader_ip, m.moderator_id
		FROM mediacontent m
		WHERE m.user_id = ?
		ORDER BY m.created DESC
		LIMIT ?
		`, uid, MEDIACONTENT_USERLIMIT)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	mediacontents := make([]*Mediacontent, 0)
	for rows.Next() {
		m := new(Mediacontent)
		err := rows.Scan(
			&m.Mediacontent_id,
			&m.Event_id,
			&m.UserId,
			&m.Filename,
			&m.Description,
			&m.Mediatype,
			&m.Storage_host,
			&m.Created,
			&m.Uploader_useragent,
			&m.Uploader_ip,
			&m.Moderator_id,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		mediacontents = append(mediacontents, m)
	}
	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error on rows.Err()%s\n", err)
		os.Exit(1)
	}
	return mediacontents
}

func NewMediacontent(db *sql.DB, m *Mediacontent) (lastId int64) {
	stmt, err := db.Prepare(`
		INSERT INTO mediacontent(
		event_id, user_id, filename, description, 
		mediatype, storage_host, 
		created, moderator_id, uploader_useragent, uploader_ip) 
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `)
	defer stmt.Close()
	if err != nil {
		log.Panicln(err)
	}
	res, err := stmt.Exec(m.Event_id, m.UserId, m.Filename, m.Description, m.Mediatype, m.Storage_host, m.Created, m.Moderator_id, m.Uploader_useragent, m.Uploader_ip)
	if err != nil {
		log.Panicln(err)
	}
	lastId, err = res.LastInsertId()
	if err != nil {
		log.Panicln(err)
	}
	return lastId
}

func GetMediacontentsFourItemByEventId(db *sql.DB, uid, eid int64) []*Mediacontent {
	limit := int64(4)
	rows, err := db.Query(`
		SELECT mediacontent_id, event_id, UserId, filename, description, mediatype, storage_host, created, moderator_id
		FROM mediacontent
		WHERE UserId = ? AND event_id = ?
		ORDER BY mediacontent_id DESC 
		LIMIT ?
		`, uid, eid, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	mediacontents := make([]*Mediacontent, 0)
	for rows.Next() {
		m := new(Mediacontent)
		err := rows.Scan(
			&m.Mediacontent_id,
			&m.Event_id,
			&m.UserId,
			&m.Filename,
			&m.Description,
			&m.Mediatype,
			&m.Storage_host,
			&m.Created,
			&m.Uploader_useragent,
			&m.Uploader_ip,
			&m.Moderator_id,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		mediacontents = append(mediacontents, m)
	}
	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error on rows.Err()%s\n", err)
		os.Exit(1)
	}
	return mediacontents
}

// Paged data using:
// func GetMediacontentsPagedByEventId(db *sql.DB, uid, eid, limit, offset int64) []*Mediacontent {
// 	rows, err := db.Query(`
// 		SELECT mediacontent_id, UserId, name, description, dt_start, dt_finish, tags, address, coords, Created
// 		FROM Mediacontent
// 		WHERE UserId = ?
// 		ORDER BY mediacontent_id DESC
// 		LIMIT ? OFFSET ?
// 		`,uid, limit, offset)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		os.Exit(1)
// 	}
// 	defer rows.Close()

// 	mediacontents := make([]*Mediacontent, 0)
// 	for rows.Next() {
// 		e := new(Mediacontent)
// 		err := rows.Scan(
// 			&e.Mediacontent_id,
// 			&e.UserId,
// 			&e.Name,
// 			&e.Description,
// 			&e.Dt_start,
// 			&e.Dt_finish,
// 			&e.Tags,
// 			&e.Address,
// 			&e.Coords,
// 			&e.Created,
// 		)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "%s\n", err)
// 			os.Exit(1)
// 		}
// 		Mediacontents = append(Mediacontents, e)
// 	}
// 	if err = rows.Err(); err != nil {
// 		fmt.Fprintf(os.Stderr, "Error on rows.Err()%s\n", err)
// 		os.Exit(1)
// 	}
// 	return mediacontents
// }

// func GetMediacontentById(db *sql.DB, uid, eid int64) Mediacontent {
// 	sqlStatement := `SELECT Mediacontent_id, UserId, name, description, dt_start, dt_finish, tags, address, coords, Created
// 	FROM Mediacontent WHERE UserId = ? AND Mediacontent_id =  ?`
// 	var e Mediacontent
// 	row := db.QueryRow(sqlStatement, uid, eid)
// 	err := row.Scan(
// 		&e.Mediacontent_id,
// 		&e.UserId,
// 		&e.Name,
// 		&e.Description,
// 		&e.Dt_start,
// 		&e.Dt_finish,
// 		&e.Tags,
// 		&e.Address,
// 		&e.Coords,
// 		&e.Created,
// 	)
// 	switch err {
// 		case sql.ErrNoRows:
// 			fmt.Println("mediacontent:225:No rows were returned!")
// 			return e
// 		case nil:
// 			return e
// 		default:
// 			fmt.Fprintf(os.Stderr, "Error on row.Scan:%s\n", err)
// 			os.Exit(1)
// 	}
// 	return e
// }

// func GetNewMediacontents(db *sql.DB, uid int64) int64 {
// 	sqlStatement := `SELECT Mediacontent_id FROM Mediacontent WHERE UserId = ? AND dt_start >= ?`
// 	var e Mediacontent
// 	row := db.QueryRow(sqlStatement, uid, dt_start)
// 	err := row.Scan(
// 		&e.UserId,
// 	)
// 	switch err {
// 		case sql.ErrNoRows:
// 			return 0
// 		case nil:
// 			return e.UserId
// 		default:
// 			fmt.Fprintf(os.Stderr, "%s\n", err)
// 			os.Exit(1)
// 	}
// 	return 0
// }

// func __GetAllUsers(db *sql.DB) ([]*Mediacontent, error) {
// 	rows, err := db.Query("SELECT Mediacontent_id, UserId, name, marketing_text, description, dt_start, dt_finish, tags, address, coords, Created FROM Mediacontent")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		os.Exit(1)
// 	}
// 	defer rows.Close()

// 	Mediacontents := make([]*Mediacontent, 0)
// 	for rows.Next() {
// 		e := new(Mediacontent)
// 		err := rows.Scan(
// 			&e.UserId,
// 			&e.Name,
// 			&e.Email,
// 			&e.Tel,
// 			&e.Impp,
// 			&e.AuthKey,
// 			&e.PasswordHash,
// 			&e.ApproveToken,
// 			&e.Picture,
// 			&e.Birthday,
// 			&e.Updated,
// 			&e.Lastlogin,
// 			&e.Roles,
// 		)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "Error on row.Scan:%s\n", err)
// 			os.Exit(1)
// 		}
// 		Mediacontents = append(Mediacontents, u)
// 	}
// 	if err = rows.Err(); err != nil {
// 		fmt.Fprintf(os.Stderr, "Error on rows.Err():%s\n", err)
// 		os.Exit(1)
// 	}
// 	return Mediacontents, nil
// }

// // Return rows affected only!!!
// func __ApproveNewUser(db *sql.DB, email, token, setRole string) int64 {

//     stmt, err := db.Prepare(`update user set approve_token = "", roles = ? where email = ? and approve_token = ?`)
//     defer stmt.Close()
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	res, err := stmt.Exec(setRole, email, token)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	rowCnt, err := res.RowsAffected()
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	return rowCnt
// }

// func ___UpdateUserInfo(db *sql.DB, u User) (rowCnt int64){
//     stmt, err := db.Prepare(`UPDATE user
//     	SET
//     	Mediacontent_id=?, UserId=?, name=?, marketing_text=?, description=?, dt_start=?, dt_finish=?, tags=?, address=?, coords=?, Created=? where UserId=?;
//     `)
//     defer stmt.Close()
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	res, err := stmt.Exec(e.Mediacontent_id, e.UserId, e.Name, e.Marketing_text, e.Description, e.Dt_start, e.Dt_finish, e.Tags, e.Address, e.Coords, e.Created)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	rowCnt, err = res.RowsAffected()
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	return rowCnt
// }
