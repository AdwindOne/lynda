package main

import (
	"io"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", authenticate)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":9000", nil)
}

func authenticate(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("logged-in")
	// no cookie
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "logged-in",
			Value: "0",
		}
	}

	// check log in
	if req.Method == "POST" {
		password := req.FormValue("password")
		if password == "secret" {
			cookie = &http.Cookie{
				Name:  "logged-in",
				Value: "1",
			}
		}
	}

	// if logout, then logout
	if req.URL.Path == "/logout" {
		cookie = &http.Cookie{
			Name:   "logged-in",
			Value:  "0",
			MaxAge: -1,
		}
	}

	http.SetCookie(res, cookie)
	var html string

	// not logged in
	if cookie.Value == strconv.Itoa(0) {
		html = `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title></title>
			</head>
			<body>
			<h1>LOG IN</h1>
			<form method="post" action="http://localhost:9000/">
				<h3>User name</h3>
				<input type="text" name="userName">
				<h3>Password</h3>
				<input type="text" name="password">
				<br>
				<input type="submit">
				<input type="submit" name="logout" value="logout">
			</form>
			</body>
			</html>`
	}

	// logged in
	if cookie.Value == strconv.Itoa(1) {
		html = `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title></title>
			</head>
			<body>
			<h1><a href="http://localhost:9000/logout">LOG OUT</a></h1>
			</body>
			</html>`
	}

	io.WriteString(res, html)
}

/*
All material is licensed under the Apache License Version 2.0, January 2004
http://www.apache.org/licenses/LICENSE-2.0
*/
