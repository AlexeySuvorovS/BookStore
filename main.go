package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"text/template"
)

type User struct {
	Id       int
	Name     string `json:"name"`
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string
}

func (u User) checkFields() bool {
	if u.Name > "" &&
		u.Email > "" &&
		u.Login > "" &&
		u.Password > "" {
		return true
	} else {
		return false
	}

}

type Users struct {
	users []User
}

func (u *Users) maxId() int {
	var maxId int
	for _, v := range u.users {
		if v.Id > maxId {
			maxId = v.Id
		}
	}
	return maxId
}

func (u *Users) add(user User) bool {
	if user.checkFields() {
		u.users = append(u.users, user)
		return true
	}

	return false
}

func (u *Users) update(user User) bool { // todo error
	if !user.checkFields() {
		return false
	}

	userUpd, _ := u.findById(user.Id)

	userUpd.Name = user.Name
	userUpd.Email = user.Email
	userUpd.Login = user.Login
	userUpd.Password = user.Password

	return true
}

func (u *Users) find(id int) (User, bool) {
	var user User
	for _, v := range u.users {
		if v.Id == id {
			return v, false
		}
	}

	return user, true
}

func (u *Users) findById(id int) (*User, bool) {
	var user User
	for i := range u.users {
		if u.users[i].Id == id {
			return &u.users[i], false
		}
	}

	return &user, true
}

func (u *Users) NewUsers() Users {
	var newUsers Users
	newUsers.users = make([]User, 0)
	return newUsers
}

var tpl *template.Template

const (
	tmplFileIndexPage    = "index.html"
	tmplFileAllUsers     = "allUsers.html"
	tmplFileAddUser      = "addUser.html"
	tmplFileAddUserPOST  = "addUserPOST.html"
	tmplFileEditUser     = "editUser.html"
	tmplFileEditUserPOST = "editUserPOST.html"
	tmplFooter           = "footer.tmpl"
)

func prepareTemplates() {

	tmplList := [...]string{
		tmplFileIndexPage,
		tmplFileAllUsers,
		tmplFileAddUser,
		tmplFileAddUserPOST,
		tmplFileEditUser,
		tmplFileEditUserPOST,
		tmplFooter,
	}

	var tmplS []string = make([]string, len(tmplList))
	for i, v := range tmplList {
		fp := path.Join("templates", v)
		tmplS[i] = fp
	}

	var err error
	tpl, err = template.ParseFiles(tmplS...)
	if err != nil {
		panic(err.Error())
	}
}

var mux *http.ServeMux

func makeRoutes() {
	mux = http.NewServeMux()
	mux.HandleFunc("/", getDefaultPage)
	mux.HandleFunc("/hello", getHelloPage)

	mux.HandleFunc("/allUsers", getAllUsers)

	mux.HandleFunc("/addUser", addUser)
	mux.HandleFunc("/addUserPOST", addUserPOST)

	mux.HandleFunc("/editUser", editUser)
	mux.HandleFunc("/editUserPOST", editUserPOST)

	http.ListenAndServe(":8888", mux)
}

func getDefaultPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := tpl.ExecuteTemplate(w, tmplFileIndexPage, allUsers.users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getHelloPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "Hello")
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := tpl.ExecuteTemplate(w, tmplFileAllUsers, allUsers.users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := tpl.ExecuteTemplate(w, tmplFileAddUser, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addUserPOST(w http.ResponseWriter, r *http.Request) {
	user := User{
		Id:       allUsers.maxId() + 1,
		Name:     r.FormValue("email"),
		Email:    r.FormValue("email"),
		Login:    r.FormValue("login"),
		Password: r.FormValue("psw"),
	}

	w.Header().Set("Content-Type", "text/html")
	if allUsers.add(user) {
		if err := tpl.ExecuteTemplate(w, tmplFileAddUserPOST, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "All mandatory fields should be filled in", http.StatusUnprocessableEntity)
	}

}

func editUser(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")

	if ids == "" {
		http.Error(w, "User ID not specified", http.StatusUnprocessableEntity)
		return
	}

	var id int
	id, _ = strconv.Atoi(ids)

	user, err := allUsers.find(id)
	if err {
		http.Error(w, "User "+ids+" not found", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tpl.ExecuteTemplate(w, tmplFileEditUser, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func editUserPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	ids := r.FormValue("id")

	var id int
	id, _ = strconv.Atoi(ids)

	if id <= 0 {
		http.Error(w, "User ID not specified", http.StatusUnprocessableEntity)
		return
	}

	user := User{
		Id:       id,
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Login:    r.FormValue("login"),
		Password: r.FormValue("psw"),
	}

	if user.checkFields() {
		allUsers.update(user)

		if err := tpl.ExecuteTemplate(w, tmplFileEditUserPOST, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "All mandatory fields should be filled in", http.StatusUnprocessableEntity)
	}
}

var allUsers Users

func initSomeUsers() {
	user := User{
		Id:       1,
		Name:     "name1",
		Email:    "email1@mail.mail",
		Login:    "login1",
		Password: "psw1",
	}

	allUsers.add(user)

	user = User{
		Id:       2,
		Name:     "name2",
		Email:    "email2@mail.mail",
		Login:    "login2",
		Password: "psw2",
	}

	allUsers.add(user)
}

func main() {

	allUsers = allUsers.NewUsers()
	initSomeUsers()
	prepareTemplates()
	makeRoutes()
}
