// ● {{base_url}}
// ● {{base_url}}?format=text
// ● {{base_url}}?branch=X
// ● {{base_url}}?year=X
// ● {{base_url}}/:id

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Data struct {
	Ids []string `json:"ids"`
}

type ErrorStruct struct {
	// E []byte `json:error`
	Error string `json:"error"`
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
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorStruct{Error: "Data not not found"})
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorStruct{Error: "File was found but unable to get file info"})
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
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorStruct{Error: "Error reading file"})
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
		if len(ids) != 0 {
			json.NewEncoder(w).Encode(Data{Ids: ids})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorStruct{Error: "No data with requested branch"})
		}
	} else if queries.Get("year") != "" {
		X, _ := strconv.Atoi(queries.Get("year"))
		for scanner.Scan() {
			line := scanner.Text()
			year, _ := strconv.Atoi(line[0:4])
			if year == 2024+1-X {
				ids = append(ids, line)
			}
		}
		if len(ids) != 0 {
			json.NewEncoder(w).Encode(Data{Ids: ids})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorStruct{Error: "No data with requested year"})
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorStruct{Error: "Invalid Request"})
	}
}

func (path *FilePath) id_Handler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	_, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorStruct{Error: "did you forget ? before query?"})
		return
	}
	file, err := os.Open(path.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorStruct{Error: "File was not found"})
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
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorStruct{Error: "id not found"})
	}
}

func main() {
	var filePath string
	// fmt.Print("Enter the file path: ")
	// fmt.Scanln(&filePath)
	for{
		fmt.Print("Enter the file path: ")
		fmt.Scanln(&filePath)

		// Check if the file exists
		_, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("File not found! Please try again.")
			} else {
				fmt.Println("Error accessing the file:", err)
			}
		} else {
			// If file exists, break out of the loop
			break
		}
	}
	// path := FilePath{Path: "./data.txt"}
	path := FilePath{Path: filePath}
	mux := http.NewServeMux()

	mux.HandleFunc("/", path.handler)
	mux.HandleFunc("/{id}", path.id_Handler)
	fmt.Println("Server up and running at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
