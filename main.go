package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"project-web/connection"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title":   "Personal web",
	"IsLogin": true,
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
}

var Projects = []structProject{}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnection()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET").Name("home")
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/project", projectPost).Methods("POST")
	route.HandleFunc("/project/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/project/d/{id}", projectDelete).Methods("GET")
	route.HandleFunc("/project/e/{id}", projectEdit).Methods("GET")
	route.HandleFunc("/project/e/{id}", projectEditPost).Methods("POST")
	route.HandleFunc("/contact", contactMe).Methods("GET")

	// port := 5000
	fmt.Println("Server running at localhost:5000")
	http.ListenAndServe("localhost:5000", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// var tmpl, err = template.ParseFiles("views/index.html")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	var result []structProject
	// rows, err := connection.Conn.Query(context.Background(), "SELECT id,name,start_date, end_date,description,technologies,image FROM tb_projects")
	rows, err := connection.Conn.Query(context.Background(), "SELECT * FROM tb_projects")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for rows.Next() {
		var each = structProject{}

		var err = rows.Scan(&each.ID, &each.Name, &each.Start, &each.End, &each.Description, &each.Tech, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Duration = getDuration(each.Start, each.End)
		result = append(result, each)
	}

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
		// "Projects": Projects,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project-form.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	dataPage := map[string]interface{}{
		"Title": "ADD MY PROJECT",
		"url":   "/project/",
	}

	DataDetail := map[string]interface{}{
		"Data": Data,
		"Page": dataPage,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, DataDetail)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	// Penampung
	ProjectDetail := structProject{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_projects WHERE id=$1", id).Scan(&ProjectDetail.ID, &ProjectDetail.Name, &ProjectDetail.Start, &ProjectDetail.End, &ProjectDetail.Description, &ProjectDetail.Tech, &ProjectDetail.Image)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// Project := structProject
	var Project structProject
	for _, data := range Projects {
		if data.ID == id {
			Project = data
			break
		}
	}

	DataDetail := map[string]interface{}{
		"Data":    Data,
		"Project": Project,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, DataDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func projectPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	ProjName := r.PostForm.Get("project-name")
	ProjStart := r.PostForm.Get("project-start")
	ProjEnd := r.PostForm.Get("project-end")
	ProjDescription := r.PostForm.Get("project-description")
	ProjTech := r.Form["project-tech"]

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description, technologies, image) VALUES ($1, $2, $3, $4,$5, 'project.webp')", ProjName, ProjStart, ProjEnd, ProjDescription, ProjTech)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// Parsing data kedalam tipe time.Time
	timeStart, _ := time.Parse("2006-01-02", ProjStart)
	timeEnd, _ := time.Parse("2006-01-02", ProjEnd)

	// Hasilkan data Durasi berdasarkan variable yang telah diparsing
	ProjDuration := getDuration(timeStart, timeEnd)

	// Tampilkan Hasil inputnya
	fmt.Println("Project Name : ", ProjName)
	fmt.Println("Start Date   : ", ProjStart)
	fmt.Println("End Date     : ", ProjEnd)
	fmt.Println("Duration     : ", ProjDuration)
	fmt.Println("Description  : ", ProjDescription)
	fmt.Println("Technologies : ", ProjTech)
	fmt.Println("================================")

	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func projectEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project-form.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Project := structProject
	Project := structProject{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects WHERE id=$1", id).Scan(&Project.ID, &Project.Name, &Project.Start, &Project.End, &Project.Description, &Project.Tech)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Parsing data kedalam tipe time.Time
	// timeStart, _ := time.Parse("2006-01-02", Project.Start)
	// timeEnd, _ := time.Parse("2006-01-02", Project.End)
	// formatStart := Project.Start.Format("2006-01-02")
	// formatEnd := Project.End.Format("2006-01-02")

	dataPage := map[string]interface{}{
		"Title": "EDIT MY PROJECT",
		"url":   "/project/e/{{.Project.ID}}",
	}
	DataDetail := map[string]interface{}{
		"Data":    Data,
		"Page":    dataPage,
		"Project": Project,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, DataDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func projectDelete(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// fmt.Println(id)
	// Projects = append(Projects[:id], Projects[id+1:]...)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func projectEditPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	ProjName := r.PostForm.Get("project-name")
	ProjStart := r.PostForm.Get("project-start")
	ProjEnd := r.PostForm.Get("project-end")
	ProjDescription := r.PostForm.Get("project-description")
	ProjTech := r.Form["project-tech"]

	// Parsing data kedalam tipe time.Time
	timeStart, _ := time.Parse("2006-01-02", ProjStart)
	timeEnd, _ := time.Parse("2006-01-02", ProjEnd)

	// Hasilkan data Durasi berdasarkan variable yang telah diparsing
	// ProjDuration := getDuration(timeStart, timeEnd)

	fmt.Println("Project Name : ", ProjName)
	fmt.Println("Start Date   : ", ProjStart)
	fmt.Println("End Date     : ", ProjEnd)
	// fmt.Println("Duration     : ", ProjDuration)
	fmt.Println("Description  : ", ProjDescription)
	fmt.Println("Technologies : ", ProjTech)
	fmt.Println("================================")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name = $1, start_date = $2, end_date = $3, description = $4, technologies = $5 WHERE id=$6", ProjName, timeStart, timeEnd, ProjDescription, ProjTech, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func getDuration(start, end time.Time) string {

	// Store Date with the Format
	// DataStart, _ := time.Parse("2006-01-02", start)
	// DataEnd, _ := time.Parse("2006-01-02", end)

	// Get data Range
	DataRange := end.Sub(start)

	// Calc duration
	yearRange := int(DataRange.Hours() / (12 * 30 * 24))
	monthRange := int(DataRange.Hours() / (30 * 24))
	weekRange := int(DataRange.Hours() / (7 * 24))
	dayRange := int(DataRange.Hours() / 24)

	if yearRange != 0 {
		return "Duration - " + strconv.Itoa(yearRange) + " Year"
	}
	if monthRange != 0 {
		return "Duration - " + strconv.Itoa(monthRange) + " Month"
	}
	if weekRange != 0 {
		return "Duration - " + strconv.Itoa(weekRange) + " Week Left"
	}
	if dayRange != 0 {
		return "Duration - " + strconv.Itoa(dayRange) + " Day Left"
	}
	return "Duration - Today"
}
