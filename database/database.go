package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	DATABASE_URL = os.Getenv("DATABASE_URL")
)

// statement
var (
	stmtInsertUser    *sql.Stmt
	stmtInsertPost    *sql.Stmt
	stmtInsertComment *sql.Stmt
	stmtInsertReport  *sql.Stmt

	stmtGetAllPosts *sql.Stmt

	stmtGetWisdomPointByID *sql.Stmt
	stmtCheckWisdomPoint   *sql.Stmt
	stmtInsertWisdomPoint  *sql.Stmt

	stmtGetPsikologByID *sql.Stmt
	stmtInsertPsikolog  *sql.Stmt
)

type User struct {
	Id         int    `json:"user_id"`
	Email      string `json:"user_email"`
	Gender     string `json:"user_gender"`
	Age        int    `json:"user_age"`
	Profession string `json:"user_profession"`
}

type Psikolog struct {
	Id       int    `json:"psikolog_id"`
	Email    string `json:"psikolog_email"`
	Name     string `json:"psikolog_name"`
	ImageURL string `json:"psikolog_image_url"`
	Wisdom   int    `json:"psikolog_wisdom,string"`
	Bio      string `json:"psikolog_bio"`
}

type Post struct {
	Id          int        `json:"post_id"`
	UserId      int        `json:"post_user_id"`
	PsikologId  int        `json:"post_psikolog_id"`
	Date        *time.Time `json:"post_date"`
	Title       string     `json:"post_title"`
	Category    string     `json:"post_category"`
	Content     string     `json:"post_content"`
	ImageURL    string     `json:"post_image_url"`
	ReportCount int        `json:"post_report_count"`
}

type Comment struct {
	Id         int        `json:"comment_id"`
	UserId     int        `json:"comment_user_id"`
	PsikologId int        `json:"comment_psikolog_id"`
	PostId     int        `json:"comment_post_id"`
	Text       string     `json:"comment_text"`
	Date       *time.Time `json:"comment_date"`
}

type Report struct {
	Id     int `json:"report_id"`
	UserId int `json:"report_user_id"`
	PostId int `json:"report_post_id"`
}

type Database struct {
	Conn *sql.DB
}

type WisdomPoint struct {
	UserID     int `json:"user_id"`
	PsikologID int `json:"psikolog_id"`
}

type PsikologPoint struct {
	PsikologID string `json:"psikolog_id"`
	Point      string `json:"psikolog_wisdom_point"`
}

type WisdomPointStatus struct {
	Status string `json:"wisdom_point_status"`
}

// response /v0/reliever?reliever_id=1
type Reliever struct {
	Name string `json:"reliever_name"`
	Bio  string `json:"reliever_bio"`
}

func New() (*Database, error) {
	db, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	// insert user statement
	stmtInsertUser, err = db.Prepare(`INSERT INTO users(user_email, user_gender, user_age, user_profession) VALUES ($1,$2,$3,$4)`)
	if err != nil {
		log.Printf("Error insert user statement: %v\n", err)
	}

	// Psikolog/reliever
	// insert psikolog statement
	stmtInsertPsikolog, err = db.Prepare(`INSERT INTO psikologs(psikolog_email, psikolog_name, psikolog_image_url, psikolog_wisdom, psikolog_bio) VALUES ($1,$2,$3,$4,$5)`)
	if err != nil {
		log.Printf("Error insert psikolog statement: %v\n", err)
	}
	// get psikolog by ID
	stmtGetPsikologByID, err = db.Prepare(`SELECT psikolog_name, psikolog_bio FROM psikologs WHERE psikolog_id=$1`)
	if err != nil {
		log.Printf("Error stmtGetPsikologByID: %v\n", err)
	}

	// insert post statement
	stmtInsertPost, err = db.Prepare(`INSERT INTO posts(post_user_id, post_psikolog_id, post_title, post_category, post_content, post_image_url) VALUES ($1,$2,$3,$4,$5,$6)`)
	if err != nil {
		log.Printf("Error insert post statement: %v\n", err)
	}

	// insert comment statement
	stmtInsertComment, err = db.Prepare(`INSERT INTO comments(comment_user_id, comment_psikolog_id, comment_post_id, comment_text) VALUES ($1,$2,$3,$4)`)
	if err != nil {
		log.Printf("Error insert comment statement: %v\n", err)
	}

	// insert report statement
	stmtInsertReport, err = db.Prepare(`INSERT INTO reports(report_user_id, report_post_id) VALUES ($1,$2)`)
	if err != nil {
		log.Printf("Error insert report statement: %v\n", err)
	}

	// get all posts
	stmtGetAllPosts, err = db.Prepare(`SELECT * FROM posts`)
	if err != nil {
		log.Printf("Error get all posts statement: %v\n", err)
	}

	// get the sum of psikolog wisdom points
	stmtGetWisdomPointByID, err = db.Prepare(`SELECT SUM(wisdom_point) FROM wisdom_points WHERE wisdom_psikolog_id=$1`)
	if err != nil {
		log.Printf("Error get wisdom point by ID statement: %v\n", err)
	}
	// check wisdom point if exists
	stmtCheckWisdomPoint, err = db.Prepare(`SELECT EXISTS(SELECT 1 FROM wisdom_points WHERE wisdom_user_id=$1 AND wisdom_psikolog_id=$2)`)
	if err != nil {
		log.Printf("Error check wisdom point statement: %v\n", err)
	}
	stmtInsertWisdomPoint, err = db.Prepare(`INSERT INTO wisdom_points(wisdom_user_id,wisdom_psikolog_id) VALUES ($1,$2)`)
	if err != nil {
		log.Printf("Error insert wisdom point statement: %v\n", err)
	}
	return &Database{db}, nil
}

func (db *Database) InsertUser(user *User) error {
	// insert data to database
	_, err := stmtInsertUser.Exec(user.Email, user.Gender, user.Age, user.Profession)
	if err != nil {
		log.Printf("Error while insert data to users table: %v\n", err)
		return err
	}
	return nil
}

func (db *Database) InsertPsikolog(p *Psikolog) error {
	// insert data to database
	_, err := stmtInsertPsikolog.Exec(p.Email, p.Name, p.ImageURL, p.Wisdom, p.Bio)
	if err != nil {
		log.Printf("Error while insert data to psikologs table: %v\n", err)
		return err
	}
	return nil
}

func (db *Database) InsertPost(p *Post) error {
	// insert data to database
	_, err := stmtInsertPost.Exec(p.UserId, p.PsikologId, p.Title, p.Category, p.Content, p.ImageURL)
	if err != nil {
		log.Printf("Error while insert data to posts table: %v\n", err)
		return err
	}
	return nil
}

func (db *Database) InsertComment(c *Comment) error {
	// insert data to database
	_, err := stmtInsertComment.Exec(c.UserId, c.PsikologId, c.PostId, c.Text)
	if err != nil {
		log.Printf("Error while insert data to comments table: %v\n", err)
		return err
	}
	return nil
}

func (db *Database) InsertReport(r *Report) error {
	// insert data to database
	_, err := stmtInsertReport.Exec(r.UserId, r.PostId)
	if err != nil {
		log.Printf("Error while insert data to reports table: %v\n", err)
		return err
	}
	return nil
}

func (db *Database) GetAllPosts() ([]Post, error) {
	var posts []Post
	rows, err := stmtGetAllPosts.Query()
	defer rows.Close()
	if err != nil {
		log.Printf("Error while get data all posts: %v\n", err)
		return nil, err
	}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Id, &post.UserId, &post.PsikologId, &post.Date, &post.Title, &post.Category, &post.Content, &post.ImageURL, &post.ReportCount)
		if err != nil {
			log.Printf("Error while iterating a rows on get all posts: %v\n", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetWisdomPointByID return the sum of psikolog wisdom point
func (db *Database) GetWisdomPointByID(id string) (PsikologPoint, error) {
	var p PsikologPoint
	p.PsikologID = id
	err := stmtGetWisdomPointByID.QueryRow(id).Scan(&p.Point)
	if err != nil {
		return p, err
	}
	return p, nil
}

// CheckWisdomPoint return a WisdomPointStatus if record exists.
func (db *Database) CheckWisdomPoint(user_id string, psikolog_id string) (WisdomPointStatus, error) {
	var ws WisdomPointStatus
	err := stmtCheckWisdomPoint.QueryRow(user_id, psikolog_id).Scan(&ws.Status)
	if err != nil {
		return ws, err
	}
	return ws, nil
}

// InsertWisdomPoint insert new records on wisdom_points table.
func (db *Database) InsertWisdomPoint(w *WisdomPoint) error {
	_, err := stmtInsertWisdomPoint.Exec(w.UserID, w.PsikologID)
	if err != nil {
		return err
	}
	return nil
}

// psikolog
// GetPsikologByID get psikolog data with specified psikolog_id.
// return Reliever if only if error is nil.
func (db *Database) GetPsikologByID(psikolog_id string) (Reliever, error) {
	var r Reliever
	err := stmtGetPsikologByID.QueryRow(psikolog_id).Scan(&r.Name, &r.Bio)
	if err != nil {
		return r, err
	}
	return r, nil
}
