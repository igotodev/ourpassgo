// a simple server to generate passwords
// and news parser just for a good time ( /matreshka )
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/igotodev/gspass"
)

type GuestPage struct {
	MyTime      string
	Pass        []string
	PassStrong  []string
	YourIP      string
	YourCountry string
	YourCC      string
}

type ServerIP struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	CC      string `json:"cc"`
}
type Article struct {
	Title      string    `json:"title"`    // заголовок темы/новости
	Content    string    `json:"content"`  // полный текст новости
	Snippet    string    `json:"snippet"`  // короткое текстовое описание
	MainPic    string    `json:"pic"`      // ссылка на основную картинку
	Link       string    `json:"link"`     // ссылка на оригинал
	Author     string    `json:"author"`   // автор новости
	Ts         time.Time `json:"ts"`       // дата-время оригинла
	AddedTS    time.Time `json:"ats"`      // дата-время добавления на сайт
	Active     bool      `json:"active"`   // флаг текущей активности
	ActiveTS   time.Time `json:"activets"` // дата-время активации
	Geek       bool      `json:"geek"`     // флаг гиковской темы
	Votes      int       `json:"votes"`    // колличество голосов за тему
	Deleted    bool      `json:"del"`      // флаг удаления
	Archived   bool      `json:"archived"` // флаг архивации
	Slug       string    `json:"slug"`     // slug новости
	SourceFeed string    `json:"feed"`     // RSS фид источника
	Domain     string    `json:"domain"`   // домен новости
	Comments   int       `json:"comments"` // число комментариев
	Likes      int       `json:"likes"`    // число лайков
	ShowNumber int       `json:"show_num"` // номер выпуска
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// news parser just for a good time
func matreshkaHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/matreshka" {
		http.NotFound(w, r)
		return
	}
	fmt.Println("Loadind radio-t...")
	files := []string{
		"html/matreshka.html",
	}
	tmpl, err := template.ParseFiles(files...)
	checkErr(err)

	resp, err := http.Get("https://news.radio-t.com/api/v1/news/last/10") // it takes a json with the ten most recent news items
	checkErr(err)

	data, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	defer resp.Body.Close()

	news := []Article{}
	err = json.Unmarshal(data, &news)
	checkErr(err)
	err = tmpl.ExecuteTemplate(w, "news", news)
	checkErr(err)

}

// main page
func viewHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	temp, err := template.ParseFiles("html/index.html")
	checkErr(err)
	t := time.Now().UTC().Local()
	passList := make([]string, 5)
	passListStrong := make([]string, 5)

	// generates 5 passwords of medium complexity
	for i := 0; i < 5; i++ {
		passList[i], err = gspass.GetPassDL(35)
		checkErr(err)
	}
	// generates 5 passwords of high complexity
	for i := 0; i < 5; i++ {
		passListStrong[i], err = gspass.GetPass(35)
		checkErr(err)
	}

	resp, err := http.Get("https://api.myip.com/") // it takes a json with ip info about server
	checkErr(err)
	answ, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	checkErr(err)
	jsonAnsw := ServerIP{}
	err = json.Unmarshal(answ, &jsonAnsw)
	checkErr(err)

	myGuestPage := GuestPage{
		MyTime:      t.Format("2006-01-02"),
		Pass:        passList,
		PassStrong:  passListStrong,
		YourIP:      jsonAnsw.IP,
		YourCountry: jsonAnsw.Country,
		YourCC:      jsonAnsw.CC,
	}

	err = temp.Execute(writer, myGuestPage)
	checkErr(err)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/matreshka", matreshkaHandler)
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	server := http.Server{
		Addr: ":" + os.Getenv("PORT"),
		//Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("starting server on" + os.Getenv("PORT"))
	err := server.ListenAndServe()
	checkErr(err)

}
