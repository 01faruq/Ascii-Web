package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var tem *template.Template
var err error

func home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello Word !")

	if r.URL.Path != "/" {
		http.Error(w, "404 This Page does not exist", http.StatusNotFound)
		return
	}
	tem.Execute(w, nil)

}

func Route(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello Word !")
	if r.Method != "POST" {
		http.Error(w, "405 : Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	text := r.FormValue("text")
	banner := r.FormValue("banner")
	ban, err := Readfile(banner)
	if err != nil {
		http.Error(w, "500: Server Error", http.StatusInternalServerError)
		return
	}

	art, err := GenerateASCII(text, ban)
	if err != nil {
		http.Error(w, "400: You typed a non-ascii character", http.StatusBadRequest)
		return
	}

	tem.Execute(w, art)

}

func main() {
	tem, _ = template.ParseFiles("my.html")
	http.HandleFunc("/", home)
	http.HandleFunc("/ascii_art", Route)

	fmt.Println("server is running on http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func Readfile(bannerName string) ([]string, error) {
	data, err := os.ReadFile(bannerName)
	if err != nil {
		return nil, err
	}
	str := string(data)
	return strings.Split(str, "\n"), nil

}

func GenerateASCII(text string, sli []string) (string, error) {

	Text := strings.ReplaceAll(text, "\\n", "\n")
	texts := strings.Split(Text, "\n")

	var Art strings.Builder
	for index, word := range texts {
		if word == "" && index == len(texts)-1 {
			break
		}
		if word == "" {
			Art.WriteString("\n")
			continue
		}

		for line := 1; line <= 8; line++ {
			for _, char := range word {
				// check Non Ascii character
				if char < 32 || char > 126 {
					return "", fmt.Errorf("invalid character: %q", char)
				}
				startIndex := int(char-32) * 9

				Art.WriteString(sli[startIndex+line])

			}

			Art.WriteString("\n")
		}
	}
	return Art.String(), nil
}
