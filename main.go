package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const apiUrl = "http://run.shareprogress.org/api/v1/buttons/analytics"

var apiKey = os.Getenv("API_KEY")

type PostData struct {
	ButtonId string `json:"button_id"`
}

func main() {
	http.HandleFunc("/fetch", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		notFound(w, r, http.StatusNotFound)
	} else {
		decoder := json.NewDecoder(r.Body)
		var button PostData
		err := decoder.Decode(&button)

		if err != nil {
			panic(err)
		}

		log.Println(button.ButtonId)
		getButtonFromSp(button.ButtonId)
		w.WriteHeader(http.StatusOK)
	}
}

func writeDataToRecord(data string) {
	fmt.Printf("%s\n", string(data))
}

func getButtonFromSp(buttonId string) {
	go func(button string) {
		resp, _ := http.PostForm(apiUrl, url.Values{"key": {apiKey}, "id": {button}})
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Printf("ERROR! %s\n", err)
		}

		log.Printf("Data Received: %d %s\n", resp.StatusCode, string(contents))
	}(buttonId)
}

func notFound(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}
