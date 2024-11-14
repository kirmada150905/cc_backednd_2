package main

import (
	"bufio"
	"encoding/json"
	"fmt"

	// "log"
	"net/http"
	"os"
	"strconv"
)

type Data struct {
	Ids []string `json:"ids"`
}

type details struct {
	Year   int    `json:"year"`
	Branch string `json:"branch"`
	Campus string `json:"campus"`
	Email  string `json:"email"`
	Id     string `json:"id"`
	Uid    string `json:uid`
}

func get_brach(id string) string {
	dict := map[string]string{
		"A1":         "chemical",
		"A2":         "civil",
		"A3":         "eee",
		"A4":         "mech",
		"A5":         "pharma",
		"A7":         "cs",
		"A8":         "eni",
		"AA":         "ece",
		"AB":         "Manu",
		"D2":         "genstudies",
		"B1":         "bio",
		"B2":         "chem",
		"B3":         "eco",
		"B4":         "math",
		"B5":         "phy",
		"chemical":   "A1",
		"civil":      "A2",
		"eee":        "A3",
		"mech":       "A4",
		"pharma":     "A5",
		"cs":         "A7",
		"eni":        "A8",
		"ece":        "AA",
		"Manu":       "AB",
		"genstudies": "D2",
		"bio":        "B1",
		"chem":       "B2",
		"eco":        "B3",
		"math":       "B4",
		"phy":        "B5",
	}

	return dict[id[4:6]]
}
func get_campus(id string) string {
	if id[12] == 'P' {
		return "pilani"
	}
	if id[12] == 'G' {
		return "goa"
	} else {
		return "hyderabad"

	}

}

type FilePath struct {
	Path string
}

func (path *FilePath) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	file, err := os.Open(path.Path) // For read access.
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	fileInfo, err := file.Stat()
	if err != nil {
		w.Write([]byte("Unable to get file info"))
		// w.Write({"Unable to get file info": http.StatusInternalServerError})
		return
	}
	fileSize := fileInfo.Size()

	scanner := bufio.NewScanner(file)
	queries := r.URL.Query()
	var ids []string //array for storing id's

	//{{base_url}}
	if len(queries) == 0 {
		for scanner.Scan() {
			line := scanner.Text()
			ids = append(ids, line)
		}

		dataStruct := Data{Ids: ids} //struct for storing id's
		//encoding dataStruct into JSON
		jsonStr, err := json.Marshal(dataStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonStr)
	} else if queries.Get("format") == "text" { //{{base_url}}?format=text

		//getting file size

		file_data := make([]byte, fileSize)
		_, err = file.Read(file_data)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write(file_data)

	} else if queries.Get("branch") != "" { //{{base_url}}?branch=X
		for scanner.Scan() {
			line := scanner.Text()
			if get_brach(line) == queries.Get("branch") {
				ids = append(ids, line)
			}
		}

		dataStruct := Data{Ids: ids} //struct for storing id's

		//encoding dataStruct into JSON
		jsonStr, err := json.Marshal(dataStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonStr)
	} else if queries.Get("year") != "" { //{{base_url}}?year=X
		X, _ := strconv.Atoi(queries.Get("year"))
		for scanner.Scan() {
			line := scanner.Text()
			year, _ := strconv.Atoi(line[0:4]) //slicing oart of the line(id) that represents year
			if year == 2024+1-X {
				ids = append(ids, line)
				fmt.Println(line)
			}
		}

		dataStruct := Data{Ids: ids} //struct for storing id's

		//encoding dataStruct into JSON
		jsonStr, err := json.Marshal(dataStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonStr)
	} else {
		w.Write([]byte("error"))
	}

}

func (path *FilePath) id_Handler(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open(path.Path) // For read access.

	scanner := bufio.NewScanner(file)
	w.Header().Set("Content-Type", "application/json")

	response_sent := false
	for scanner.Scan() {
		line := scanner.Text()
		if line[8:12] == r.URL.Path[1:] {
			type data struct {
				Id details `json:"id"`
			}
			year, _ := strconv.Atoi(line[0:4])
			year = 2025 - year
			brach := get_brach(line)
			campus := get_campus(line)
			email := "f" + line[0:4] + line[8:12] + "@" + campus + ".bits-pilani.ac.in"
			id := line
			uid := line[8:12]
			id_deail := details{Year: year, Branch: brach, Campus: campus, Email: email, Id: id, Uid: uid}
			response := data{Id: id_deail}
			jsonStr, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(jsonStr)
			response_sent = true
			break
		}
	}
	if !response_sent {
		w.Write([]byte("id not found"))
	}
}

func main() {

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	path := FilePath{Path: "./data.txt"}
	mux := http.NewServeMux()

	mux.HandleFunc("/", path.handler)
	mux.HandleFunc("/{id}", path.id_Handler)
	http.ListenAndServe(":8080", mux)

}

//todo: handle slash at the end of url

// ● {{base_url}}
// ● {{base_url}}?format=text
// ● {{base_url}}?branch=X
// ● {{base_url}}?year=X
// ● {{base_url}}/:id
