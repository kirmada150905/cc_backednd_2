package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Data struct{
	Ids []string `json:"ids"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open("./data.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	var ids []string
	for scanner.Scan(){
		line := scanner.Text()
		fmt.Println(line)
		ids = append(ids, line)
	}
	data := Data{Ids: ids}

    // Encoding
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

	if err != nil {  
        http.Error(w, err.Error(), http.StatusInternalServerError)  
        return 
   }
	// jsonStr, err := json.Marshal(data)
	// w.Write(jsonStr)
}

func main() {

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)

 }

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// type Person struct {
// 	Name string `json:"name"`
// 	Age  int    `json:"age"`
// }

// func handler(w http.ResponseWriter, r *http.Request) { 
//     person := Person{  Name: "John",  Age: 30, } 

//     // Encoding - One step
//     jsonStr, err := json.Marshal(person) 

//     if err != nil {  
//         http.Error(w, err.Error(), http.StatusInternalServerError)  
//         return 
//     } 
// 	fmt.Println(person)
//     w.Write(jsonStr)
// }

// func main() {
// 	http.HandleFunc("/", handler)
// 	http.ListenAndServe(":8080", nil)
// }
