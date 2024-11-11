package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	// "reflect"
)

type Data struct{
	Ids []string `json:"ids"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.URL)
	file, err := os.Open("./data.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	var ids []string
	for scanner.Scan(){
		line := scanner.Text()
		ids = append(ids, line)
	}
	data := Data{Ids: ids}

	format := r.URL.Query().Get("format")
    // Encoding
	if format == "text" {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Println(format)
	}else{
		w.Header().Set("Content-Type", "application/json")
	}
	json.NewEncoder(w).Encode(data)
	if err != nil {  
        http.Error(w, err.Error(), http.StatusInternalServerError)  
        return 
   }

}

func main() {

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

    http.HandleFunc("/", handler)
    // http.HandleFunc("/format=text", handler)
    http.ListenAndServe(":8080", nil)

 }