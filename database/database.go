package database

import (
	"database/sql"
	"log"
	"os"
	// "time"

	_ "github.com/lib/pq"
)

var (
	DATABASE_URL = os.Getenv("DATABASE_URL")
)

// statement
var (
	stmtInsertUser     *sql.Stmt
	stmtInsertPsikolog *sql.Stmt
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
	Wisdom   int    `json:"psikolog_wisdom"`
	Bio      string `json:"psikolog_bio"`
}

type Post struct {
	Id         int `json:"post_id"`
	UserId     int `json:"post_user_id"`
	PsikologId int `json:"post_psikolog_id"`
	// Date        *time.Time `json:"post_date"`
	Title       string `json:"post_title"`
	Category    string `json:"post_category"`
	Content     string `json:"post_content"`
	ImageURL    string `json:"post_image_url"`
	ReportCount int    `json:"post_report_count"`
}

type Comment struct {
	Id         int    `json:"comment_id"`
	UserId     int    `json:"comment_user_id"`
	PsikologId int    `json:"comment_psikolog_id"`
	PostId     int    `json:"comment_post_id"`
	Text       string `json:"comment_text"`
	// Date       *time.Time `json:"comment_date"`
}

type Report struct {
	Id     int `json:"report_id"`
	UserId int `json:"report_user_id"`
	PostId int `json:"report_post_id"`
}

type Database struct {
	Conn *sql.DB
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

	// insert psikolog statement
	stmtInsertPsikolog, err = db.Prepare(`INSERT INTO psikologs(psikolog_email, psikolog_name, psikolog_image_url, psikolog_wisdom, psikolog_bio) VALUES ($1,$2,$3,$4,$5)`)
	if err != nil {
		log.Printf("Error insert user statement: %v\n", err)
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
	log.Println("Insert user: succed")
	return nil
}

func (db *Database) InsertPsikolog(p *Psikolog) error {
	// insert data to database
	_, err := stmtInsertPsikolog.Exec(p.Email, p.Name, p.ImageURL, p.Wisdom, p.Bio)
	if err != nil {
		log.Printf("Error while insert data to psikologs table: %v\n", err)
		return err
	}
	log.Println("Insert psikolog: succed")
	return nil
}
