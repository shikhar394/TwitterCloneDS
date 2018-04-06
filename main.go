package main

import(
  "html/template"
  "net/http"
  "fmt"
  "container/list"
  "strings"
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
  http.HandleFunc("/home", homeHandler)
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

  //ckeck if user is in the gloabl Users List
  for u:=Users.Front(); u!=nil; u=u.Next(){
    if strings.EqualFold(u.Value.(*User).Email, user_email) && u.Value.(*User).Password == user_password{
      http.Redirect(w,r,"/home",302)
      return
    }
  }
  //If no user found
  tmpl.Execute(w, struct{ Authfail bool }{true})

  fmt.Print(user_email + " : " + user_password + "\n")

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

  Users.PushBack(newUser)
  fmt.Print(Users.Front().Value.(*User).FirstName)
  tmpl.Execute(w, struct{ Success bool }{true})

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("home.html"))
  if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
}
