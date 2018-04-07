package main

import (
	"container/list"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type User struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

// TODO make maps pointing with email bool || pointer to node for constant lookups

var Email string

var UserMap = make(map[string]*User)

var UserList = list.New() // Only for distributing among servers while making distributed.
//Users := make(map[Emai])

func main() {
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/home", homeHandler)
	http.ListenAndServe(":8000", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	userEmail := r.FormValue("email")
	userPassword := r.FormValue("password")

	// Check if user is in the map keying using the email.
	User, ok := UserMap[userEmail]
	if ok {
		if User.Password == userPassword {
			http.Redirect(w, r, "/home", 302)
			return
		}
	}

	//If no user found
	tmpl.Execute(w, struct{ Authfail bool }{true})

	fmt.Print(userEmail + " : " + userPassword + "\n")

}

func sortedInsert(newUser *User) {
	//fmt.Printf("\n\nSorted insert on email %v\n\n", newUser.Email)
	if UserList.Front() != nil {
		for CurrentUser := UserList.Front(); CurrentUser != nil; CurrentUser = CurrentUser.Next() {
			if strings.Compare(CurrentUser.Value.(*User).Email, (*newUser).Email) == 1 {
				UserList.InsertBefore(newUser, CurrentUser)
				return
			}
		}
	} else {
		UserList.PushFront(newUser)
		return
	}
	UserList.PushBack(newUser)
	//fmt.Printf("\n\nDone inserting\n\n")
}

func readUsers() {
	for CurrentUser := UserList.Front(); CurrentUser != nil; CurrentUser = CurrentUser.Next() {
		fmt.Printf("User %v, Password %v\n", CurrentUser.Value.(*User).Email, CurrentUser.Value.(*User).Password)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("signUp.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	newUser := new(User)
	newUser.Email = r.FormValue("new_email")
	newUser.FirstName = r.FormValue("new_fname")
	newUser.LastName = r.FormValue("new_lname")
	newUser.Password = r.FormValue("new_password")

	//Users.PushBack(newUser) // TODO Make sorted instead and store in map for constant time lookup
	sortedInsert(newUser)
	UserMap[(newUser).Email] = newUser
	readUsers()
	fmt.Print(UserList.Front().Value.(*User).FirstName)
	tmpl.Execute(w, struct{ Success bool }{true})

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("home.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
}
