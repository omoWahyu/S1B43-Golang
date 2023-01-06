package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"project-web/connection"
	"project-web/middleware"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title     string
	IsLogin   bool
	IDUser    int
	NameUser  string
	FlashData string
}

var Data = MetaData{
	Title: "Personal Web",
}

type structProject struct {
	ID          int
	Name        string
	Start       time.Time
	End         time.Time
	Description string
	Tech        []string
	Duration    string
	Image       string
	ID_User     int
	IsLogin     bool
}

var Projects = []structProject{}

type structUser struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnection()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	route.HandleFunc("/", home).Methods("GET").Name("home")

	// Contact
	route.HandleFunc("/contact", contactMe).Methods("GET")

	// Project
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/project", middleware.UploadFile(projectPost)).Methods("POST")
	route.HandleFunc("/project/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/project/d/{id}", projectDelete).Methods("GET")
	route.HandleFunc("/project/e/{id}", projectEdit).Methods("GET")
	route.HandleFunc("/project/e/{id}", projectEditPost).Methods("POST")

	// Authentication
	route.HandleFunc("/auth/login", authLogin).Methods("GET")
	route.HandleFunc("/auth/login", authLoginPost).Methods("POST")
	route.HandleFunc("/auth/register", authRegister).Methods("GET")
	route.HandleFunc("/auth/register", authRegisterPost).Methods("POST")
	route.HandleFunc("/auth/logout", authLogout).Methods("GET")

	// port := 5000
	fmt.Println("Server running at localhost:5001")
	http.ListenAndServe("localhost:5001", route)
}

// Index Section
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// var tmpl, err = template.ParseFiles("views/index.html")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	var result []structProject

	var rows pgx.Rows
	// var err error
	if session.Values["IsLogin"] != true {
		// Session
		Data.IsLogin = false

		// Query
		rows, err = connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image  FROM tb_projects")
	} else {
		// Session
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.NameUser = session.Values["Name"].(string)
		Data.IDUser = session.Values["ID"].(int)

		// Query
		rows, err = connection.Conn.Query(context.Background(), "SELECT tb_projects.id, tb_projects.Name, start_date, end_date,description, technologies, image FROM tb_projects LEFT JOIN tb_users ON tb_projects.id_user = tb_users.id where tb_projects.id_user =$1", Data.IDUser)
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for rows.Next() {
		var db = structProject{}
		var err = rows.Scan(&db.ID, &db.Name, &db.Start, &db.End, &db.Description, &db.Tech, &db.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		db.Duration = getDuration(db.Start, db.End)
		result = append(result, db)
	}

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

// Contact Section
func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.NameUser = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	respData := map[string]interface{}{
		"Data": Data,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

// Project Management Section
func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project-add.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	DataDetail := map[string]interface{}{
		"Data": Data,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, DataDetail)
}

func projectPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.NameUser = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	Name := r.PostForm.Get("project-name")
	Start := r.PostForm.Get("project-start")
	End := r.PostForm.Get("project-end")
	Description := r.PostForm.Get("project-description")
	Tech := r.Form["project-tech"]
	dataContext := r.Context().Value("dataFile")
	Image := dataContext.(string)
	currentUser := session.Values["ID"].(int)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description, technologies, image, id_user) VALUES ($1, $2, $3, $4, $5, $6, $7)", Name, Start, End, Description, Tech, Image, currentUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// Parsing data kedalam tipe time.Time
	timeStart, _ := time.Parse("2006-01-02", Start)
	timeEnd, _ := time.Parse("2006-01-02", End)

	// Hasilkan data Durasi berdasarkan variable yang telah diparsing
	Duration := getDuration(timeStart, timeEnd)

	// Tampilkan Hasil inputnya
	fmt.Println("Project Name  : ", Name)
	fmt.Println("Start Date    : ", Start)
	fmt.Println("End Date      : ", End)
	fmt.Println("Duration      : ", Duration)
	fmt.Println("Description   : ", Description)
	fmt.Println("Technologies  : ", Tech)
	fmt.Println("Image         : ", Image)
	fmt.Println("User who post : ", currentUser)
	fmt.Println("================================")

	session.AddFlash("New Project Added!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.NameUser = session.Values["Name"].(string)
	}

	var tmpl, err = template.ParseFiles("views/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Penampung
	db := structProject{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image  FROM tb_projects WHERE id=$1", id).Scan(&db.ID, &db.Name, &db.Start, &db.End, &db.Description, &db.Tech, &db.Image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	db.Duration = getDuration(db.Start, db.End)

	// Tampilkan Hasil inputnya
	// fmt.Println("Project Name : ", db.Name)
	// fmt.Println("Start Date   : ", db.Start)
	// fmt.Println("End Date     : ", db.End)
	// fmt.Println("Duration     : ", db.Duration)
	// fmt.Println("Description  : ", db.Description)
	// fmt.Println("Technologies : ", db.Tech)
	// fmt.Println("================================")

	dataTime := map[string]interface{}{
		"fStart": db.Start.Format("2006-01-02"),
		"fEnd":   db.Start.Format("2006-01-02"),
	}
	DataDetail := map[string]interface{}{
		"Data":    Data,
		"Time":    dataTime,
		"Project": db,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, DataDetail)
}

func projectEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project-edit.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.NameUser = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Project := structProject
	db := structProject{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(&db.ID, &db.Name, &db.Start, &db.End, &db.Description, &db.Tech, &db.Image)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dataTime := map[string]interface{}{
		"fStart": db.Start.Format("2006-01-02"),
		"fEnd":   db.Start.Format("2006-01-02"),
	}
	DataDetail := map[string]interface{}{
		"Data":    Data,
		"Time":    dataTime,
		"Project": db,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, DataDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func projectEditPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	Name := r.PostForm.Get("project-name")
	Start := r.PostForm.Get("project-start")
	End := r.PostForm.Get("project-end")
	Description := r.PostForm.Get("project-description")
	Tech := r.Form["project-tech"]
	dataContext := r.Context().Value("dataFile")
	Image := dataContext.(string)

	// Parsing data kedalam tipe time.Time
	timeStart, _ := time.Parse("2006-01-02", Start)
	timeEnd, _ := time.Parse("2006-01-02", End)

	// Hasilkan data Durasi berdasarkan variable yang telah diparsing

	fmt.Println("Project Name : ", Name)
	fmt.Println("Start Date   : ", Start)
	fmt.Println("End Date     : ", End)
	// fmt.Println("Duration     : ", ProjDuration)
	fmt.Println("Description  : ", Description)
	fmt.Println("Technologies : ", Tech)
	fmt.Println("================================")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name = $1, start_date = $2, end_date = $3, description = $4, technologies = $5 image = $6 WHERE id=$7", Name, timeStart, timeEnd, Description, Tech, Image, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	session.AddFlash("Project are Updated", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func projectDelete(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func authRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// var tmpl, err = template.ParseFiles("views/index.html")

	var tmpl, err = template.ParseFiles("views/auth/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	respData := map[string]interface{}{
		"Data": Data,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func authRegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_users(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// if session.Values["IsLogin"] == true {
	// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	// }

	session.AddFlash("Successfully Register!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
}

// Auth Section
func authLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/auth/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// if session.Values["IsLogin"] == true {
	// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	// }

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	respData := map[string]interface{}{
		"Data": Data,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func authLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	db := structUser{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id,name,email,password FROM tb_users WHERE email=$1", email).Scan(
		&db.ID, &db.Name, &db.Email, &db.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(db.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = db.Name
	session.Values["ID"] = db.ID
	session.Options.MaxAge = 10800

	// println(db.ID)
	session.AddFlash("Successfully Login!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// Contact Section
func authLogout(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Logout")
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Global Func
func getDuration(start, end time.Time) string {

	// Get data Range
	DataRange := end.Sub(start)

	// Calc duration
	yearRange := int(DataRange.Hours() / (12 * 30 * 24))
	monthRange := int(DataRange.Hours() / (30 * 24))
	weekRange := int(DataRange.Hours() / (7 * 24))
	dayRange := int(DataRange.Hours() / 24)

	if yearRange != 0 {
		if yearRange >= 0 {
			return "Duration : " + strconv.Itoa(yearRange) + " Year"

		}
		str := strconv.Itoa(yearRange)
		return "Duration : " + strings.ReplaceAll(str, "-", "") + " Year"
	}
	if monthRange != 0 {
		if monthRange >= 0 {
			return "Duration : " + strconv.Itoa(monthRange) + " Month"

		}
		str := strconv.Itoa(monthRange)
		return "Duration : " + strings.ReplaceAll(str, "-", "") + " Month"
	}
	if weekRange != 0 {
		if weekRange >= 0 {
			return "Duration : " + strconv.Itoa(weekRange) + " Week"

		}
		str := strconv.Itoa(weekRange)
		return "Duration : " + strings.ReplaceAll(str, "-", "") + " Week"
	}
	if dayRange != 0 {
		if dayRange >= 0 {
			return "Duration : " + strconv.Itoa(dayRange) + " Day"

		}
		str := strconv.Itoa(dayRange)
		return "Duration : " + strings.ReplaceAll(str, "-", "") + " Day"

	}

	return "Duration : Today"
}
