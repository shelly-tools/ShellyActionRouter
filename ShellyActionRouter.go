/**
 * @package   ShellyActionRouter
 * @copyright Thorsten Eurich
 * @license   GNU Affero General Public License (https://www.gnu.org/licenses/agpl-3.0.de.html)
 * @version  0.1
 *
 * @todo lots of documentation
 *
 * Simple Action Proxy written in Golang which makes it possible to execute multiple actions
 */

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/ini.v1"
)

type ActionRoute struct {
	Id     int
	Name   string
	Action []Action
}

type ActionUrl struct {
	Id   int
	Url  string
	Name string
}
type Action struct {
	Id         int
	ActionType string
	Content    string
	IdUrl      int
}

var netClient = &http.Client{
	Timeout: time.Second * 5,
}

type Unreachable struct {
	Status string `json:"status"`
}

var tmpl = template.Must(template.ParseGlob("views/*"))

func dbConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./ActionRouter.db")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Index(w http.ResponseWriter, r *http.Request) {

	IpAddress := ""

	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	port := cfg.Section("server").Key("port").String()

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			IpAddress = "http://" + ipv4.String() + ":" + port
		}
	}

	db := dbConn()
	selDBIndex, err := db.Query("SELECT idurl, urlname FROM urls")
	if err != nil {
		panic(err.Error())
	}
	emp := ActionUrl{}
	res := []ActionUrl{}
	for selDBIndex.Next() {

		var idurl int
		var urlname string

		err = selDBIndex.Scan(&idurl, &urlname)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = idurl
		emp.Url = IpAddress + "/api/action/" + urlname
		emp.Name = urlname

		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func AddUrl(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "AddUrl", nil)
}
func DeleteUrl(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("idurl")
	delForm, err := db.Prepare("DELETE FROM urls WHERE idurl=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("deleted URL id" + emp)
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func InsertUrl(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		urlname := r.FormValue("urlname")
		insForm, err := db.Prepare("INSERT INTO urls (urlname) VALUES(?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(urlname)
		log.Println("INSERT: URL: " + urlname + "")
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func EditUrl(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("idurl")
	selDB, err := db.Query("SELECT * FROM urls WHERE idurl=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := ActionUrl{}

	for selDB.Next() {
		var idurl int
		var urlname string
		err = selDB.Scan(&idurl, &urlname)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = idurl
		emp.Name = urlname
	}
	tmpl.ExecuteTemplate(w, "EditUrl", emp)
	defer db.Close()
}

func ShowUrl(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("idurl")
	selDB, err := db.Query("SELECT * FROM urls WHERE idurl=?", nId)
	if err != nil {
		panic(err.Error())
	}
	empindex := ActionRoute{}

	for selDB.Next() {
		var idurl int
		var urlname string
		err = selDB.Scan(&idurl, &urlname)
		if err != nil {
			panic(err.Error())
		}
		empindex.Id = idurl
		empindex.Name = urlname

		selDBAction, err := db.Query("SELECT * FROM actions WHERE idurl=?", nId)
		if err != nil {
			panic(err.Error())
		}
		emp := Action{}
		res := []Action{}

		for selDBAction.Next() {
			var id int
			var actiontype string
			var content string
			var idurl int

			err = selDBAction.Scan(&id, &actiontype, &content, &idurl)
			if err != nil {
				panic(err.Error())
			}
			emp.Id = id
			emp.ActionType = actiontype
			emp.IdUrl = idurl
			emp.Content = content
			res = append(res, emp)
		}
		empindex.Action = res

	}
	fmt.Println(empindex)
	tmpl.ExecuteTemplate(w, "ShowUrl", empindex)
	defer db.Close()
}

func UpdateUrl(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		urlname := r.FormValue("urlname")
		idurl := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE urls SET urlname=? WHERE idurl=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(urlname, idurl)
		log.Println("UPDATE URL: " + urlname + " ")
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func AddAction(w http.ResponseWriter, r *http.Request) {
	idurl := r.URL.Query().Get("idurl")
	tmpl.ExecuteTemplate(w, "AddAction", idurl)
}
func InsertAction(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		actiontype := r.FormValue("actiontype")
		content := r.FormValue("content")
		idurl := r.FormValue("idurl")
		insForm, err := db.Prepare("INSERT INTO actions (actiontype, content, idurl) VALUES(?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(actiontype, content, idurl)
		log.Println("INSERT Action: " + content + "")
	}
	defer db.Close()
	http.Redirect(w, r, "/url/show?idurl="+r.FormValue("idurl"), 301)
}

func DeleteAction(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	id := r.URL.Query().Get("id")
	idurl := r.URL.Query().Get("idurl")

	delForm, err := db.Prepare("DELETE FROM actions WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(id)
	log.Println("deleted Action id" + id)
	defer db.Close()
	http.Redirect(w, r, "/url/show?idurl="+idurl, 301)
}

func EditAction(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM actions WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Action{}

	for selDB.Next() {
		var id int
		var actiontype string
		var content string
		var idurl int

		err = selDB.Scan(&id, &actiontype, &content, &idurl)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.ActionType = actiontype
		emp.Content = content
		emp.IdUrl = idurl
	}
	tmpl.ExecuteTemplate(w, "EditAction", emp)
	defer db.Close()
}
func UpdateAction(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	url := r.URL.Query().Get("idurl")
	if r.Method == "POST" {
		id := r.FormValue("id")
		//idurl := r.FormValue("idurl")
		content := r.FormValue("content")
		actiontype := r.FormValue("actiontype")
		insForm, err := db.Prepare("UPDATE actions SET actiontype=?, content=?  WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(actiontype, content, id)
		log.Println("UPDATE Action " + id + " ")
	}
	defer db.Close()
	http.Redirect(w, r, "/url/show?idurl="+url, 301)
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	port := cfg.Section("server").Key("port").String()

	r := mux.NewRouter()
	staticDir := "/assets/"
	// Create the route
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	// URLs handler
	r.HandleFunc("/", Index)
	r.HandleFunc("/url/add", AddUrl)
	r.HandleFunc("/url/delete", DeleteUrl)
	r.HandleFunc("/url/edit", EditUrl)
	r.HandleFunc("/url/insert", InsertUrl)
	r.HandleFunc("/url/update", UpdateUrl)
	r.HandleFunc("/url/show", ShowUrl)

	// Actions handler
	r.HandleFunc("/action/add", AddAction)
	r.HandleFunc("/action/insert", InsertAction)
	r.HandleFunc("/action/delete", DeleteAction)
	r.HandleFunc("/action/edit", EditAction)
	r.HandleFunc("/action/update", UpdateAction)

	// API-Handler
	r.HandleFunc("/api/action/{title}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]

		db := dbConn()
		selDB, err := db.Query("SELECT a.id, a.actiontype, a.content, a.idurl FROM actions AS a LEFT JOIN urls AS u ON (a.idurl = u.idurl) WHERE urlname = ?", title)
		if err != nil {
			panic(err.Error())
		}
		emp := Action{}
		res := []Action{}

		for selDB.Next() {
			var id int
			var actiontype string
			var content string
			var idurl int

			err = selDB.Scan(&id, &actiontype, &content, &idurl)
			if err != nil {
				panic(err.Error())
			}
			emp.Id = id
			emp.ActionType = actiontype
			emp.Content = content
			emp.IdUrl = idurl
			res = append(res, emp)
		}
		defer db.Close()

		for _, v := range res {
			if v.ActionType == "sleep" {
				sl, _ := strconv.Atoi(v.Content)
				time.Sleep(time.Duration(sl) * time.Millisecond)
			} else {
				curl := fmt.Sprintf("%s", v.Content)

				res, err := netClient.Get(curl)
				if err != nil {
					emp := Unreachable{Status: "offline"}
					var jsonData []byte
					jsonData, _ = json.Marshal(emp)

					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonData)

					log.Println(curl + " unreachable")
				} else {
					defer res.Body.Close()
					body, _ := ioutil.ReadAll(res.Body)
					fmt.Fprintf(w, string(body))
				}
			}

		}
	})

	http.ListenAndServe(":"+port, r)
}
