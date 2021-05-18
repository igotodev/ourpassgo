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

type FinalNews struct {
	TitleF   []string
	SnippetF []string
	LinkF    []string
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func balalaikaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Loadind radio-t...")
	tmpl, err := template.ParseFiles("balalaika.html")
	checkErr(err)

	resp, err := http.Get("https://news.radio-t.com/api/v1/news/last/5")
	checkErr(err)

	data, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	defer resp.Body.Close()

	news := []Article{}
	FNews := FinalNews{}
	err = json.Unmarshal(data, &news)
	checkErr(err)
	for _, v := range news {
		FNews.TitleF = append(FNews.TitleF, v.Title)
		FNews.SnippetF = append(FNews.SnippetF, v.Snippet)
		FNews.LinkF = append(FNews.LinkF, v.Link)
	}
	err = tmpl.Execute(w, FNews)
	checkErr(err)

}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	temp, err := template.ParseFiles("index.html")
	checkErr(err)
	t := time.Now().UTC().Local()
	passList := make([]string, 5)
	passListStrong := make([]string, 5)

	//fmt.Printf("%v success\n", passList)
	for i := 0; i < 5; i++ {
		passList[i], err = gspass.GetPassDL(35)
		checkErr(err)
	}

	for i := 0; i < 5; i++ {
		passListStrong[i], err = gspass.GetPass(35)
		checkErr(err)
	}

	resp, err := http.Get("https://api.myip.com/")
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
	mux.HandleFunc("/balalaika", balalaikaHandler)
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	server := http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("starting server on" + os.Getenv("PORT"))
	err := server.ListenAndServe()
	checkErr(err)

}
