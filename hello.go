package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var tpl = template.Must(template.ParseGlob("*.html"))

func homePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("app.html")
	t.Execute(w, nil)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("login.html")
	t.Execute(w, nil)

	if r.Method != "POST" {
		fmt.Println("hello")

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("create.html")
	t.Execute(w, nil)
}
func create1(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("fname")

	db, err := sql.Open("mysql", "root:rootroot@tcp(127.0.0.1:3306)/testdb")

	if err != nil {

		panic(err.Error())
	}
	defer db.Close()

	result, err := db.Query("SELECT email FROM users WHERE email=?;", email)

	var count int = 0
	var emailid string

	if err != nil {
		panic(err.Error())
	} else {
		for result.Next() {
			err = result.Scan(&emailid)
			if err != nil {
				panic(err.Error())
			}
			if emailid == email {
				count = count + 1
			}
		}
	}
	defer result.Close()

	fmt.Println(count)
	if count == 1 {
		tpl.ExecuteTemplate(w, "create.html", "User Already Exists")
	} else {

		//insert, err := db.Query("INSERT INTO users VALUES('" + email + "','" + password + "');")
		insert, err := db.Query("INSERT INTO users VALUES(?,?,?);", email, password, name)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()

		fmt.Println("Successfully inserted")

		tpl.ExecuteTemplate(w, "app.html", "Account Created")
	}

}
func user(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")
	db, err := sql.Open("mysql", "root:rootroot@tcp(127.0.0.1:3306)/testdb")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//result, err := db.Query("SELECT email FROM users WHERE email='" + email + "'and password='" + password + "';")
	result, err := db.Query("SELECT email,name FROM users WHERE email=? and password=?;", email, password)

	var count int = 0
	var emailid string
	var name string

	if err != nil {
		panic(err.Error())
	} else {
		for result.Next() {
			err = result.Scan(&emailid, &name)
			if err != nil {
				panic(err.Error())
			}
			if emailid == email {
				count = count + 1
			}
		}
	}
	defer result.Close()

	fmt.Println(count)
	if count == 1 {

		var str string
		str = "" + name + ""
		tpl.ExecuteTemplate(w, "user.html", str)
	} else {
		tpl.ExecuteTemplate(w, "login.html", "Username/Password Incorrect")

	}

}

func update1(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	oldpassword := r.FormValue("oldpassword")
	newpassword := r.FormValue("newpassword")
	db, err := sql.Open("mysql", "root:rootroot@tcp(127.0.0.1:3306)/testdb")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result, err := db.Query("SELECT email FROM users WHERE email=? and password=?;", email, oldpassword)

	var count int = 0
	var emailid string

	if err != nil {
		panic(err.Error())
	} else {
		for result.Next() {
			err = result.Scan(&emailid)
			if err != nil {
				panic(err.Error())
			}
			if emailid == email {
				count = count + 1
			}
		}
	}
	defer result.Close()
	fmt.Println(count)
	if count == 1 {
		//result, err := db.Query("update users set password='" + password + "' WHERE email='" + email + "';")
		result, err := db.Query("update users set password=? WHERE email=?;", newpassword, email)

		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		fmt.Println("Updated")

		tpl.ExecuteTemplate(w, "app.html", "Password Updated")
	} else {
		tpl.ExecuteTemplate(w, "update.html", "Username/Password Incorrect")

	}

}
func update(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("update.html")
	t.Execute(w, nil)

}
func delete1(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")
	db, err := sql.Open("mysql", "root:rootroot@tcp(127.0.0.1:3306)/testdb")
	defer db.Close()

	if err != nil {
		panic(err.Error())
	}
	result, err := db.Query("SELECT email FROM users WHERE email=? and password=?;", email, password)

	var count int = 0
	var emailid string

	if err != nil {
		panic(err.Error())
	} else {
		for result.Next() {
			err = result.Scan(&emailid)
			if err != nil {
				panic(err.Error())
			}
			if emailid == email {
				count = count + 1
			}
		}
	}
	defer result.Close()

	if count == 0 {
		tpl.ExecuteTemplate(w, "delete.html", "Username/Password Incorrect")
	} else {
		//result, err := db.Query("DELETE FROM users WHERE email='" + email + "'and password='" + password + "';")
		result, err := db.Query("DELETE FROM users WHERE email=? and password=?;", email, password)

		if err != nil {
			panic(err.Error())
		}
		defer result.Close()

		fmt.Println("Deleted")

		tpl.ExecuteTemplate(w, "app.html", "Account Deleted")
	}

}
func delete(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("delete.html")
	if err != nil {
		panic(err.Error())
	}
	t.Execute(w, nil)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/create", create)
	http.HandleFunc("/create1", create1)

	http.HandleFunc("/login", login)
	http.HandleFunc("/user", user)
	http.HandleFunc("/update", update)
	http.HandleFunc("/update1", update1)

	http.HandleFunc("/delete", delete)
	http.HandleFunc("/delete1", delete1)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
	db, err := sql.Open("mysql", "root:rootroot@tcp(127.0.0.1:3306)/testdb")

	if err != nil {
		fmt.Println("HERE")

		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Succesfully connected to mysql database")

}
