package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title":   "Personal web",
	"IsLogin": false,
}

// Array of objects
// nama = []string{"Abel", "Dandi", "Ilham", "Jody"}

// This is interface
// type persegi interface {
// 	panjang() float64
// 	lebar() float64
// }

type dataProject struct {
	ID          int
	Name        string
	Start       string
	End         string
	Description string
	Tech        []string
	Duration    string
}

var Projects = []dataProject{
	// {
	// 	Name:        "Dumbways Mobile App 2022",
	// 	Start:       "1 Des 2022",
	// 	End:         "9 Des 2022",
	// 	Duration:    " 1 Minggu",
	// 	Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	// 	Tech:        []string{"nodejs", "nextjs", "reactjs", "typescript"},
	// },
}

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", index).Methods("GET")
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

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// var tmpl, err = template.ParseFiles("views/index.html")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": Projects,
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

	var tmpl, err = template.ParseFiles("views/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	// Project := dataProject
	var Project dataProject
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
	ProjDuration := getDuration(ProjStart, ProjEnd)

	var Project = dataProject{
		Name:        ProjName,
		Start:       ProjStart,
		End:         ProjEnd,
		Description: ProjDescription,
		Tech:        ProjTech,
		Duration:    ProjDuration,
	}

	fmt.Println("================================")

	fmt.Println("Project Name : ", Project.Name)
	fmt.Println("Start Date   : ", Project.Start)
	fmt.Println("End Date     : ", Project.End)
	fmt.Println("Duration     : ", Project.Duration)
	fmt.Println("Description  : ", Project.Description)
	fmt.Println("Technologies : ", Project.Tech)
	fmt.Println("================================")

	Projects = append(Projects, Project)
	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func projectEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project-form.html")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	// Project := dataProject
	var Project dataProject
	for _, data := range Projects {
		if data.ID == id {
			Project = data
			break
		}
	}

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

	Projects = append(Projects[:id], Projects[id+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
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
	ProjDuration := getDuration(ProjStart, ProjEnd)

	var Project = dataProject{
		Name:        ProjName,
		Start:       ProjStart,
		End:         ProjEnd,
		Description: ProjDescription,
		Tech:        ProjTech,
		Duration:    ProjDuration,
	}

	fmt.Println("================================")

	fmt.Println("Project Name : ", Project.Name)
	fmt.Println("Start Date   : ", Project.Start)
	fmt.Println("End Date     : ", Project.End)
	fmt.Println("Duration     : ", Project.Duration)
	fmt.Println("Description  : ", Project.Description)
	fmt.Println("Technologies : ", Project.Tech)
	fmt.Println("================================")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	Projects[id] = Project
	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
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

func getDuration(start, end string) string {

	// Store Date with the Format
	DataStart, err := time.Parse("2006-01-02", start)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	DataEnd, err := time.Parse("2006-01-02", end)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Get data Range
	DataRange := DataEnd.Sub(DataStart)

	// Calc duration
	monthRange := int(DataRange.Hours() / (30 * 24))
	weekRange := int(DataRange.Hours() / (7 * 24))
	dayRange := int(DataRange.Hours() / 24)

	if monthRange < 0 {
		return ""
	}
	if monthRange > 0 {
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
