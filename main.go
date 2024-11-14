// //todo: handle slash at the end of url

// ● {{base_url}}
// ● {{base_url}}?format=text
// ● {{base_url}}?branch=X
// ● {{base_url}}?year=X
// ● {{base_url}}/:id

package main

import (
	"bufio"
	"encoding/json"
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
	Uid    string `json:"uid"`
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
	}
	return "hyderabad"
}

type FilePath struct {
	Path string
}

func (path *FilePath) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	file, err := os.Open(path.Path)
	if err != nil {
		http.Error(w, `{"error": "File not found", "status": 404}`, http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, `{"error": "File was found but unable to get file info", "status": 500}`, http.StatusInternalServerError)
		return
	}
	fileSize := fileInfo.Size()

	scanner := bufio.NewScanner(file)
	queries := r.URL.Query()
	var ids []string

	if len(queries) == 0 {
		for scanner.Scan() {
			line := scanner.Text()
			ids = append(ids, line)
		}
		json.NewEncoder(w).Encode(Data{Ids: ids})
	} else if queries.Get("format") == "text" {
		fileData := make([]byte, fileSize)
		_, err = file.Read(fileData)
		if err != nil {
			http.Error(w, `{"error": "Error reading file", "status": 500}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(fileData)
	} else if queries.Get("branch") != "" {
		for scanner.Scan() {
			line := scanner.Text()
			if get_brach(line) == queries.Get("branch") {
				ids = append(ids, line)
			}
		}
		json.NewEncoder(w).Encode(Data{Ids: ids})
	} else if queries.Get("year") != "" {
		X, _ := strconv.Atoi(queries.Get("year"))
		for scanner.Scan() {
			line := scanner.Text()
			year, _ := strconv.Atoi(line[0:4])
			if year == 2024+1-X {
				ids = append(ids, line)
			}
		}
		json.NewEncoder(w).Encode(Data{Ids: ids})
	} else {
		http.Error(w, `{"error": "Invalid request", "status": 400}`, http.StatusBadRequest)
	}
}

func (path *FilePath) id_Handler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(path.Path)
	if err != nil {
		http.Error(w, `{"error": "File not found", "status": 404}`, http.StatusNotFound)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	w.Header().Set("Content-Type", "application/json")

	responseSent := false
	for scanner.Scan() {
		line := scanner.Text()
		if line[8:12] == r.URL.Path[1:] {
			type data struct {
				Id details `json:"id"`
			}
			year, _ := strconv.Atoi(line[0:4])
			year = 2025 - year
			branch := get_brach(line)
			campus := get_campus(line)
			email := "f" + line[0:4] + line[8:12] + "@" + campus + ".bits-pilani.ac.in"
			idDetail := details{
				Year:   year,
				Branch: branch,
				Campus: campus,
				Email:  email,
				Id:     line,
				Uid:    line[8:12],
			}
			json.NewEncoder(w).Encode(data{Id: idDetail})
			responseSent = true
			break
		}
	}
	if !responseSent {
		http.Error(w, `{"error": "ID not found", "status": 404}`, http.StatusNotFound)
	}
}

func main() {
	path := FilePath{Path: "./data.txt"}
	mux := http.NewServeMux()

	mux.HandleFunc("/", path.handler)
	mux.HandleFunc("/{id}", path.id_Handler)
	http.ListenAndServe(":8080", mux)
}
