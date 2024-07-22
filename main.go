package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func findFileByExt(root, ext string) ([]string, error) {
	var a []string
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

func iterCsvFile(filePath string, condition func([]string) bool) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	csvReader.LazyQuotes = true

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to parse CSV row in %s: %w", filePath, err)
		}
		if !condition(record) {
			break
		}
	}
	return nil
}

type City struct {
	Name       string `json:"name,omitempty"`
	Population int    `json:"population,omitempty"`
	Country    string `json:"country,omitempty"`
}

type findCityWithPopulationResult struct {
	result []*City
}

const (
	indexName         = 1
	indexFeatureClass = 6
	indexFeatureCode  = 7
	indexCountryCode  = 8
	indexPopulation   = 14
	indexTZ           = 17
)

func (f *findCityWithPopulationResult) checkCity(record []string, population int) {
	// skip if not city
	if record[indexFeatureClass] != "P" {
		return
	}

	myPopulation, err := strconv.Atoi(record[indexPopulation])
	if err != nil {
		log.Println("error while converting population to int: ", record[indexPopulation])
		return
	}
	c := &City{Name: strings.TrimSpace(record[indexName]), Population: myPopulation, Country: strings.TrimSpace(record[indexCountryCode])}

	if len(f.result) == 0 {
		f.result = append(f.result, c)
		return
	}

	currentDiff := population - f.result[0].Population
	if currentDiff < 0 {
		currentDiff = -currentDiff
	}

	myDiff := population - myPopulation
	if myDiff < 0 {
		myDiff = -myDiff
	}

	if myDiff < currentDiff {
		f.result = []*City{c}
	}

	if myDiff == currentDiff {
		f.result = append(f.result, c)
	}

}

func handleByPopulation(w http.ResponseWriter, r *http.Request) {
	handleError := func(err error, code int) {
		w.WriteHeader(code)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		log.Println("error while handling request: ", err)
	}

	population := filepath.Base(r.URL.Path)
	prefix := r.URL.Path[:len(r.URL.Path)-len(population)]
	if prefix != "/api/by-population/" {
		handleError(fmt.Errorf("Population is required, example: /api/by-population/1234"), http.StatusBadRequest)
		return
	}

	p, err := strconv.Atoi(population)
	if err != nil {
		handleError(fmt.Errorf("Population should be a number: %w", err), http.StatusBadRequest)
		return
	}

	if p < 0 {
		handleError(fmt.Errorf("Population should be a positive number: %d", p), http.StatusBadRequest)
		return
	}

	f, err := findCityWithPopulation(p)
	if err != nil {
		handleError(fmt.Errorf("error while finding city with population: %w", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"cities": f}); err != nil {
		log.Printf("error while encoding response: %s", err)
		return
	}

}

const indexBody = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="color-scheme" content="light dark" />
    <title>Find city by population</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css" />
  </head>
  <body>
  <main class="container">
   <h2>Find city with closest population</h2>
   <form action="?" method="get">
    <label for="population">Population:</label>
    <input type="number" id="population" name="population"><br/><br/>
    <input type="submit" value="Submit">
   </form>
   <br/><br/>
   %s
  </main>
  </body>
</html>
`

var dataFiles []string

func findCityWithPopulation(population int) ([]*City, error) {
	f := &findCityWithPopulationResult{}

	for _, s := range dataFiles {
		if len(filepath.Base(s)) != 6 {
			log.Println("Skipped as file name is not 6 characters long: ", s)
			continue
		}

		err := iterCsvFile(s, func(record []string) bool {
			f.checkCity(record, population)
			return true
		})
		if err != nil {
			return nil, fmt.Errorf("error while processing file: %s: %w", s, err)
		}
	}

	return f.result, nil
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var result string
	if population := r.URL.Query().Get("population"); population != "" {
		p, err := strconv.Atoi(population)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Population should be a number"))
			return
		}
		if p < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Population should be a positive number"))
			return
		}

		f, err := findCityWithPopulation(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal server error"))
			log.Println("error while finding city with population: ", err)
			return
		}

		result = `
      <h2>Cities with closest population:</h2>
      <table>
      <thead>
       <tr>
        <th>City</th><th>Population</th><th>Country</th>
       </tr>
      </thead>
      <tbody>`

		for _, c := range f {
			result += fmt.Sprintf(`      <tr>
       <td>%s</td>
       <td>%d</td>
       <td>%s</td>
      </tr>`, c.Name, c.Population, c.Country)
		}

		result += `      </tbody>
    </table>`
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(fmt.Sprintf(indexBody, result)))
}

func main() {
	var err error

	// find all data files with .txt extension and store in global variable
	dataFiles, err = findFileByExt(".", ".txt")
	if err != nil {
		log.Fatal("error while finding files: ", err)
	}

	// expose web service
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", handleIndex)
	mux.HandleFunc("/api/by-population/", handleByPopulation)
	log.Fatal(http.ListenAndServe(":8081", mux))

}
