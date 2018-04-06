package main

import(
  "html/template"
  "net/http"
  "fmt"
  "container/list"
)

type User struct{
  Email string
  FirstName string
  LastName string
  Password string
}

var Users = list.New()

func main() {
  http.HandleFunc("/", loginHandler)
  http.HandleFunc("/signup", signupHandler)
  http.ListenAndServe(":8000", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("index.html"))
  if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

  user_email := r.FormValue("email")
  user_password := r.FormValue("password")
  fmt.Print(user_email + " : " + user_password + "\n")

}

func signupHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("signUp.html"))
  if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

  newUser := User{
    Email: r.FormValue("new_email"),
    FirstName: r.FormValue("new_fname"),
    LastName: r.FormValue("new_lname"),
    Password: r.FormValue("new_password"),
  }

  Users.PushBack(newUser)
  fmt.Print(Users.Front())
  tmpl.Execute(w, struct{ Success bool }{true})

}
