package main

import (
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

func (ra resultSlice) Len() int {
	return len(ra)
}

func (ra resultSlice) Swap(i int, j int) {
	ra[i], ra[j] = ra[j], ra[i]
}

func (ra resultSlice) Less(i int, j int) bool {
	return ra[i].score > ra[j].score
}

func main() {
	rt := template.New("Result Template")

	rt.Funcs(template.FuncMap{"Intersperse": strings.Join})

	_, err := rt.ParseFiles(resultsTemplatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

	http.HandleFunc("/opensearch.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, openSearchPath)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, indexPath)
		case "POST":
			q := r.FormValue("q")

			resultCh := make(chan result)
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

			results := make([]result, 0, len(mergedResults))

			for _, result := range mergedResults {
				results = append(results, result)
			}

			sort.Sort(resultSlice(results))

			err := rt.ExecuteTemplate(w, resultsTemplatePath, results)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		}
	})

	err = http.ListenAndServeTLS(port, tlsCertPath, tlsKeyPath, nil)
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
}
