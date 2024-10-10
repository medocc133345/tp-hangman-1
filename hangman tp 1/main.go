package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

type Student struct {
    FirstName string
    LastName  string
    Age       int
    Gender    string
}

type Class struct {
    ClassName    string
    Field        string
    Level        string
    StudentCount int
    Students     []Student
}

type UserData struct {
    FirstName string
    LastName  string
    Birthdate string
    Gender    string
}

var (
    promoTemplate   = template.Must(template.ParseFiles("templates/promo.html"))
    changeTemplate  = template.Must(template.ParseFiles("templates/change.html"))
    formTemplate    = template.Must(template.ParseFiles("templates/form.html"))
    errorTemplate   = template.Must(template.ParseFiles("templates/error.html"))
    displayTemplate = template.Must(template.ParseFiles("templates/display.html"))
    viewCounter     = 0
    counterMutex    = sync.Mutex{}
    userData        UserData
)

func promoHandler(w http.ResponseWriter, r *http.Request) {
    class := Class{
        ClassName:   "B1 Informatique",
        Field:       "Informatique",
        Level:       "Bachelor 1",
        StudentCount: 3,
        Students: []Student{
            {FirstName: "Alice", LastName: "Dupont", Age: 20, Gender: "feminin"},
            {FirstName: "Bob", LastName: "Martin", Age: 21, Gender: "masculin"},
            {FirstName: "Charlie", LastName: "Leroy", Age: 19, Gender: "masculin"},
        },
    }
    promoTemplate.Execute(w, class)
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
    counterMutex.Lock()
    viewCounter++
    count := viewCounter
    counterMutex.Unlock()

    data := struct {
        Counter int
        IsEven  bool
    }{
        Counter: count,
        IsEven:  count%2 == 0,
    }
    changeTemplate.Execute(w, data)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        formTemplate.Execute(w, nil)
    } else if r.Method == http.MethodPost {
        r.ParseForm()
        firstName := r.FormValue("firstname")
        lastName := r.FormValue("lastname")
        birthdate := r.FormValue("birthdate")
        gender := r.FormValue("gender")

        if len(firstName) == 0 || len(lastName) == 0 || (gender != "masculin" && gender != "feminin" && gender != "autre") {
            http.Redirect(w, r, "/user/error", http.StatusSeeOther)
            return
        }

        userData = UserData{
            FirstName: firstName,
            LastName:  lastName,
            Birthdate: birthdate,
            Gender:    gender,
        }
        http.Redirect(w, r, "/user/display", http.StatusSeeOther)
    }
}

func displayHandler(w http.ResponseWriter, r *http.Request) {
    if userData.FirstName == "" || userData.LastName == "" {
        http.Redirect(w, r, "/user/error", http.StatusSeeOther)
        return
    }
    displayTemplate.Execute(w, userData)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
    errorTemplate.Execute(w, nil)
}

func main() {

    http.HandleFunc("/promo", promoHandler)
    http.HandleFunc("/change", changeHandler)
    http.HandleFunc("/user/form", formHandler)
    http.HandleFunc("/user/display", displayHandler)
    http.HandleFunc("/user/error", errorHandler)

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    log.Println("Server running on http://localhost:8080/")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
