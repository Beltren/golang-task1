package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
Напишите сервис для сокращения длинных URL. Требования:

    Сервер должен быть доступен по адресу: http://localhost:8080.
    Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
    Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
    Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
    Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
*/

/*


 */
var UrlMap = make(map[int]string)
var IdCounter = 0
var hostname = "localhost"
var port = "8080"

func getHandler(w http.ResponseWriter, r *http.Request) {
	//Возвращаем 307

	//Вычитываем URL строку
	q := r.URL.String()
	q = strings.Replace(q, "/", "", -1)
	idAsInteger, err := strconv.Atoi(q)
	if err != nil {
		// handle error
		log.Println(err)
	}
	w.Header().Set("Location", UrlMap[idAsInteger])
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	//возвращаем статус 201
	w.WriteHeader(http.StatusCreated)
	//Читаем тушку bodyString
	bodyString, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	UrlMap[IdCounter] = string(bodyString)
	w.Write([]byte("http://" + hostname + ":" + port + "/" + strconv.Itoa(IdCounter)))
	//fmt.Printf("Added:" + string(bodyString) + " to " + strconv.Itoa(IdCounter) + "\n")
	IdCounter++
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

	//Парсинг в зависимости от метода
	switch r.Method {

	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {

	fmt.Printf("Start\n")
	http.HandleFunc("/", mainHandler)

	server := &http.Server{
		Addr: hostname + ":" + port,
	}
	server.ListenAndServe()
}
