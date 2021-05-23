// a simple server to generate passwords
// and news parser just for a good time ( /matreshka )
package main

import (
	"bufio"
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

// logMiddleware implements a delay (not necessarily) and logging,
// but if you want it can implement more useful functionality
func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 19; i++ {
			time.Sleep(2 * time.Millisecond)
			fmt.Print(":")
			if i == 18 {
				fmt.Printf("\n%s opened %s\n", r.RemoteAddr, r.URL.Path)
			}
		}
		log.Printf("always good")
		next(w, r)
	}
}

// consolePrintASCII printed ASCII text from file to os.Stdout (not necessarily, it's for fun)
func consolePrintASCII(file string) {
	logo, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	info, err := logo.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := info.Size()
	reader := bufio.NewReader(logo)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(500 * time.Millisecond)

	for i := 0; i < int(size); i++ {
		time.Sleep(50 * time.Millisecond)
		fmt.Fprint(os.Stdout, string(b[i]))
	}
	time.Sleep(500 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "\n")
}

/*
// printLogo just printed logo (not necessarily, it's for fun)
func printLogo() {
	time.Sleep(400 * time.Millisecond)
	var logo string = ":::::OURPA55GO:::::"
	for i := 0; i < len(logo); i++ {
		time.Sleep(150 * time.Millisecond)
		fmt.Print(string(logo[i]))
		if i == len(logo)-1 {
			time.Sleep(800 * time.Millisecond)
		}
	}
	fmt.Printf("\n")
}
*/
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", logMiddleware(viewHandler))
	mux.HandleFunc("/matreshka", logMiddleware(matreshkaHandler))
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	server := http.Server{
		//Addr: ":" + os.Getenv("PORT"),
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	consolePrintASCII("warning.txt") // start logo with comments
	//printLogo()
	//log.Println("starting server on" + os.Getenv("PORT"))
	log.Println("starting server on port 8080")
	err := server.ListenAndServe()
	checkErr(err)

}
