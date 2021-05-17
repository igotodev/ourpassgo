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

type GuestBook struct {
	MyTime      string
	Pass        []string
	PassStrong  []string
	YourIP      string
	YourCountry string
	YourCC      string
}

type YourIP struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	CC      string `json:"cc"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
	jsonAnsw := YourIP{}
	err = json.Unmarshal(answ, &jsonAnsw)
	checkErr(err)

	myGuestBook := GuestBook{
		MyTime:      t.Format("2006-01-02"),
		Pass:        passList,
		PassStrong:  passListStrong,
		YourIP:      jsonAnsw.IP,
		YourCountry: jsonAnsw.Country,
		YourCC:      jsonAnsw.CC,
	}

	err = temp.Execute(writer, myGuestBook)
	checkErr(err)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
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
