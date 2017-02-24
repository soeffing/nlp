package sparql

import (
	"errors"
	"fmt"
	"github.com/knakk/sparql"
	"time"
)

var (
	dbpediaURL = "http://es.dbpedia.org/sparql"
)

// GetLabelAbstractByTerm takes a term and returns labels and abstracts found .. just testing
// []map[string]string
func GetLabelAbstractByTerm(term string) ([]map[string]string, error) {
	var results []map[string]string

	if term == "" {
		return results, errors.New("term cannot be empty")
	}

	repo, err := sparql.NewRepo(dbpediaURL,
		sparql.Timeout(10*time.Second),
	)

	if err != nil {
		fmt.Println("Error: Cannot connect to sparql endpoint")
		return results, err
	}

	queryString := `
    PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
    PREFIX dbo: <http://dbpedia.org/ontology/>
    SELECT DISTINCT ?s ?abstract ?p ?o ?label
    WHERE { <http://dbpedia.org/resource/%s> ?p ?o ;
              dbo:abstract ?abstract ;
              rdfs:label ?label .
          }
    LIMIT 20
    `

	finalQueryString := fmt.Sprintf(queryString, term)

	res, err := repo.Query(finalQueryString)

	if err != nil {
		fmt.Printf("Error: Cannot query the sparql endpoint %s", err)
		return results, err
	}

	bindings := res.Results.Bindings

	for _, r := range bindings {
		singleResult := make(map[string]string)
		// cycle through triples
		for k, v := range r {
			if k == "abstract" && v.Type == "literal" {
				singleResult[k] = v.Value
			}

			if k == "label" && v.Type == "literal" {
				singleResult[k] = v.Value
			}
		}

		results = append(results, singleResult)
	}
	return results, nil
}
