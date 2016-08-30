package wordrankingservice

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"strings"
	"os"
	"io/ioutil"
	"encoding/json"
)

const (
	WORD_RANKING_TABLE ="WordRanking"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
/*

type ItemDescription struct {

	ItemId string `json:"itemid"`
	Description string `json:"description"`
}
*/


//This method sets what resources are going to be managed by the router
type WordRankingService struct {

}

func (service WordRankingService) ConfigureRouter(router *mux.Router) {


	routes := Routes{

		Route{
			"get_purchases",
			"POST",
			"/word",
			service.handlePostWordRankingCalculation,
		},
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)

	}
}

func (service WordRankingService) handlePostWordRankingCalculation (w http.ResponseWriter, r *http.Request){

}



func (service WordRankingService) calculateWordRankingFor (userId string){


	file, err := os.Open("./purchases_backup.json") // For read access.

	if err != nil {
		log.Printf("mock data file could not be opened.")
		return
	}


	body, err := ioutil.ReadAll(file)

	purchaseContainer := new (PurchaseContainer)

	if err := json.Unmarshal(body, purchaseContainer); err != nil {

		log.Printf("Error when reading response %s", err)
		return
	}
	

	wordsRanking := make (map[string] int)


	for _, purchase := range purchaseContainer.Purchases {

		for _, item := range purchase.Items {

			words := strings.Split(item.Description, " ")

			for _, word := range words {

				toLow := strings.ToLower(word)
				wordsRanking[toLow] = wordsRanking[toLow] + 1
			}
		}

	}


	for k, v := range wordsRanking {
		log.Printf("word: %s ranking:%d", k, v)
	}


}