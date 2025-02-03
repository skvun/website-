package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Users struct {
	Id       int
	Name     string
	Password string
}

type Articles struct {
	Name, LNFNP, Awards, MilitaryBackground string
	//Фамилия имя отчество(Last name first name patronymic), Награды, Боевое прошлое, даты жизни
}

type USS struct {
	U Users
	A []Articles
}

var US USS
var UserMain Users
var ArticlesArray []Articles

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT `Name`, `LNFNP`, `Awards`, `MilitaryBackground` FROM `articles`")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	ArticlesArray = []Articles{}
	for res.Next() {
		var article Articles
		err = res.Scan(&article.Name, &article.LNFNP, &article.Awards, &article.MilitaryBackground)
		if err != nil {
			panic(err)
		}
		ArticlesArray = append(ArticlesArray, article)
	}
	t.ExecuteTemplate(w, "index", ArticlesArray)
}

func registr(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/registr.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "registr", nil)
}

func entrance(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/entrance.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "entrance", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func saveUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `users` (`name`, `password`) VALUES('%s', '%s')", name, password))
	if err != nil {
		panic(err)
	}
	defer insert.Close()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func savePost(w http.ResponseWriter, r *http.Request) {
	LNFNP := r.FormValue("LNFNP")
	Awards := r.FormValue("Awards")
	MilitaryBackground := r.FormValue("MilitaryBackground")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`Name`, `LNFNP`, `Awards`, `MilitaryBackground`) VALUES('%s', '%s', '%s', '%s')", UserMain.Name, LNFNP, Awards, MilitaryBackground))
	if err != nil {
		panic(err)
	}
	defer insert.Close()
	http.Redirect(w, r, "/indexUser", http.StatusSeeOther)
}

func check(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/indexUser.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}

	nameform := strings.TrimSpace(r.FormValue("name"))
	passwordform := strings.TrimSpace(r.FormValue("password"))

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT `Name`, `LNFNP`, `Awards`, `MilitaryBackground` FROM `articles`")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	ArticlesArray = []Articles{}
	for res.Next() {
		var article Articles
		err = res.Scan(&article.Name, &article.LNFNP, &article.Awards, &article.MilitaryBackground)
		if err != nil {
			panic(err)
		}
		ArticlesArray = append(ArticlesArray, article)
	}

	resu, err := db.Query("SELECT `name`, `password` FROM `users`")
	if err != nil {
		panic(err)
	}
	defer resu.Close()

	for resu.Next() {
		var user Users
		err = resu.Scan(&user.Name, &user.Password)
		if err != nil {
			panic(err)
		}

		user.Name = strings.TrimSpace(user.Name)
		user.Password = strings.TrimSpace(user.Password)

		if nameform == user.Name && passwordform == user.Password {
			UserMain = user
			US.A = ArticlesArray
			US.U = UserMain
			err = t.ExecuteTemplate(w, "indexUser", US)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func indexUser(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/indexUser.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT `Name`, `LNFNP`, `Awards`, `MilitaryBackground` FROM `articles`")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	ArticlesArray = []Articles{}
	for res.Next() {
		var article Articles
		err = res.Scan(&article.Name, &article.LNFNP, &article.Awards, &article.MilitaryBackground)
		if err != nil {
			panic(err)
		}
		ArticlesArray = append(ArticlesArray, article)
	}
	US.A = ArticlesArray
	US.U = UserMain

	t.ExecuteTemplate(w, "indexUser", US)
}

func handeFunc() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", index)
	rtr.HandleFunc("/registr", registr)
	rtr.HandleFunc("/entrance", entrance)
	rtr.HandleFunc("/saveUser", saveUser)
	rtr.HandleFunc("/create", create)
	rtr.HandleFunc("/savePost", savePost)
	rtr.HandleFunc("/indexUser", indexUser)
	rtr.HandleFunc("/check", check)

	http.Handle("/", rtr)

	http.ListenAndServe(":1988", nil)
}
func main() {
	handeFunc()
}
