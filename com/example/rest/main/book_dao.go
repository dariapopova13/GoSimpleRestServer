package main

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	_ "github.com/lib/pq"
	"time"
)

type Book struct {
	Id     uint64
	Title  string
	Year   int
	Author *Author
}

type Author struct {
	Id       uint64
	Name     string
	Surname  string
	Birthday time.Time
}

// database properties
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "books"
	driver   = "postgres"
)

// sql queries
const (
	selectAllBooksQuery = `SELECT b.id, b.title, b.year
     , a.id as a_id, a.name as a_name, a.surname as a_surname, a.birthday as a_birthday
	FROM books b INNER JOIN author a ON a.id = author_id
	ORDER BY b.id;`
	selectBookByIdQuery = `SELECT b.id, b.title, b.year
     , a.id as a_id, a.name as a_name, a.surname as a_surname, a.birthday as a_birthday
	FROM books b INNER JOIN author a ON a.id = author_id
    WHERE b.id = $1;`
	insertBookQuery = `INSERT INTO books
		 (title, year, author_id)
		 VALUES ($1, $2, $3)
		 RETURNING id;`
	updateBookQuery = `UPDATE books SET
		 title = $1
		 , year = $2
		 , author_id = $3
		 WHERE id = $4
		 RETURNING id;;`
	deleteBookQuery = `DELETE FROM books WHERE id = $1`
)

func DeleteBook(id uint64) {
	db := getConnection()
	defer closeConnection(db)
	_, err := db.Exec(deleteBookQuery, id)
	CheckError(err, fmt.Sprintf("Could not delete book with id %d", id))
}

func UpdateBook(book Book) int {
	db := getConnection()
	defer closeConnection(db)
	row := db.QueryRow(updateBookQuery, book.Title, book.Year, book.Author.Id, book.Id)
	var id int
	err := row.Scan(&id)
	CheckError(err, fmt.Sprintf("Could not update a book with id %d", book.Id))
	return id
}

func InsertBook(book Book) uint64 {
	db := getConnection()
	defer closeConnection(db)
	row := db.QueryRow(insertBookQuery, book.Title, book.Year, book.Author.Id)
	var id uint64
	err := row.Scan(&id)
	CheckError(err, "Could not insert a new book into database.")
	return id
}

func SelectBookById(id uint64) (book Book) {
	db := getConnection()
	defer closeConnection(db)
	row := db.QueryRow(selectBookByIdQuery, id)
	book = parseBook(nil, row)
	return
}

func SelectAllBooks() (books []Book) {
	db := getConnection()
	defer closeConnection(db)
	rows, err := db.Query(selectAllBooksQuery)
	CheckError(err, "Could not get all books from database.")
	books = make([]Book, 0, 0)
	for rows.Next() {
		book := parseBook(rows, nil)
		books = append(books, book)
	}
	return
}

func parseBook(rows *sql.Rows, row *sql.Row) (book Book) {
	var id uint64
	var title string
	var year int
	var aId uint64
	var aName string
	var aSurname string
	var aBirthday time.Time
	var err error
	if rows != nil {
		err = rows.Scan(&id, &title, &year, &aId, &aName, &aSurname, &aBirthday)
	} else if row != nil {
		err = row.Scan(&id, &title, &year, &aId, &aName, &aSurname, &aBirthday)
	} else {
		panic("No rows were passed.")
	}
	CheckError(err, "Could not get values from book table.")
	book = Book{Id: id, Title: title, Year: year, Author: &Author{Id: aId, Name: aName, Surname: aSurname, Birthday: aBirthday}}
	return
}

func closeConnection(db *sql.DB) {
	err := db.Close()
	CheckError(err, "Could not close a connection with the database.")
	logger.Info("The connection was closed.")
}

func getConnection() (db *sql.DB) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open(driver, psqlInfo)
	CheckError(err, "Could not get a connection with the database.")
	err = db.Ping()
	CheckError(err, "Could not get a connection with the database.")
	logger.Info("Successfully connected to the database.")
	return
}
