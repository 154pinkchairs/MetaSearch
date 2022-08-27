package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)

type resultSlice []result

func mergeResults(r result, s result) (t result) {
	t.Title = r.Title
	t.Link = r.Link
	t.Description = r.Description
	t.SearchEngines = append(r.SearchEngines, s.SearchEngines...)
	t.score = r.score + s.score

	return
}

//declare an imgres struct
type imgres struct {
	Link        string
	Description string
	Thumbnail   string
	score       float64
}

type imgResultSlice []imgres

func mergeImgResults(r imgres, s imgres) (t imgres) {
	t.Link = r.Link
	t.Description = r.Description
	t.Thumbnail = r.Thumbnail
	t.score = r.score + s.score

	return
}

//declare a function to get the image results from the search engines

func (ra resultSlice) Len() int {
	return len(ra)
}

func (ra resultSlice) Swap(i int, j int) {
	ra[i], ra[j] = ra[j], ra[i]
}

func (ra resultSlice) Less(i int, j int) bool {
	return ra[i].score > ra[j].score
}

type results struct {
	Query   string
	Results resultSlice
}

func listenAndServe(addr string, handler http.Handler) {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	server.ListenAndServe()
}

func main() {
	rt := template.New("Result Template")
	searchEngines := []searchEngine{duckduckgo{}, google{}, bing{}, duckduckgoimages{}}
	rt.Funcs(template.FuncMap{"Intersperse": strings.Join})
	resultsTemplatePath := "results.html"
	openSearchPath := "opensearch.xml"
	faviconPath := "favicon.ico"
	indexPath := "index.html"

	_, err := rt.ParseFiles(resultsTemplatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

	http.HandleFunc("/opensearch.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, openSearchPath)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, indexPath)
		case "POST":
			q := r.FormValue("q")

			resultCh := make(chan result, len(searchEngines))
			var wg sync.WaitGroup

			wg.Add(len(searchEngines))

			go func() {
				wg.Wait()
				close(resultCh)
			}()

			for _, se := range searchEngines {
				go func(se searchEngine) {
					se.search(q, resultCh)
					wg.Done()
				}(se)
			}

			mergedResults := make(map[string]result)

			for r := range resultCh {
				if existingResult, exists := mergedResults[r.Link]; exists {
					mergedResults[r.Link] = mergeResults(r, existingResult)
				} else {
					mergedResults[r.Link] = r
				}
			}

			var results results
			results.Query = q
			results.Results = make([]result, 0, len(mergedResults))

			for _, result := range mergedResults {
				results.Results = append(results.Results, result)
			}

			sort.Sort(resultSlice(results.Results))

			err := rt.ExecuteTemplate(w, resultsTemplatePath, results)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		}
	})

	//run listenAndServe on port 8080, reading the arguments from config.json
	//open the config.json file
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	//read the config.json file and pass the values to the listenAndServe function
	decoder := json.NewDecoder(configFile)
	configuration := struct {
		Addr string
	}{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println(err)
	}
	//close the config.json file
	configFile.Close()
	//create a Http handler interface
	handler := http.NewServeMux()
	//run the listenAndServe function
	listenAndServe(configuration.Addr, handler)

}
