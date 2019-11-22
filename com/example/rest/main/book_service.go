package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json; charset=UTF-8"
)

func AllBooksHandler(writer http.ResponseWriter, request *http.Request) {
	books := SelectAllBooks()
	writer.Header().Set(contentType, applicationJson)
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(books)
	CheckError(err, "Could not encode books into json format.")
}

func SingleBookHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	CheckError(err, "Could not parse id param.")
	book := SelectBookById(id)
	writer.Header().Set(contentType, applicationJson)
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(book)
	CheckError(err, fmt.Sprintf("Could not encode book with id %d into json format.", id))
}

func InsertBookHandler(writer http.ResponseWriter, request *http.Request) {
	var book Book
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 10485768))
	CheckError(err, "Could not read body from request.")
	err = request.Body.Close()
	CheckError(err, "Could not close request body.")
	err = json.Unmarshal(body, &book)
	if err = json.Unmarshal(body, &book); err != nil {
		writer.Header().Set(contentType, applicationJson)
		writer.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(writer).Encode(err)
		CheckError(err, "Could not parse error to json.")
	}
	fmt.Println(book)
	id := InsertBook(book)
	writer.Header().Set(contentType, applicationJson)
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(id)
	CheckError(err, "Could not response.")
}

func UpdateBookHandler(writer http.ResponseWriter, request *http.Request) {
	var book Book
	body, err := ioutil.ReadAll(io.LimitReader(request.Body, 1048576))
	CheckError(err, "Could not read body from request.")
	err = request.Body.Close()
	CheckError(err, "Could not close request body.")
	err = json.Unmarshal(body, &book)
	if err = json.Unmarshal(body, &book); err != nil {
		writer.Header().Set(contentType, applicationJson)
		writer.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(writer).Encode(err)
		CheckError(err, "Could not parse error to json.")
	}
	id := UpdateBook(book)
	writer.Header().Set(contentType, applicationJson)
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(id)
	CheckError(err, "Could not response.")
}

func DeleteBookHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	CheckError(err, "Could not parse id param.")
	DeleteBook(id)
	writer.Header().Set(contentType, applicationJson)
	writer.WriteHeader(http.StatusOK)
}
