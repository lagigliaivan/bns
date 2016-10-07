package wordrankingservice

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"strings"
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
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


	file, err := os.Open("./items.json") // For read access.

	if err != nil {
		log.Printf("mock data file could not be opened.")
		return
	}


	body, err := ioutil.ReadAll(file)

	items := new ([]ItemDescription)

	if err := json.Unmarshal(body, &items); err != nil {

		log.Printf("Error when reading response %s", err)
		return
	}
	

	wordsRanking := make (map[string] int)


	for _, item := range *items {
		itemRepetition, _ := strconv.Atoi(item.Quantity)
		words := strings.Split(item.Description, " ")

		for _, word := range words {

			toLowCases := strings.ToLower(word)
			if len(toLowCases) > 2 {
				wordsRanking[toLowCases] = wordsRanking[toLowCases] + (1 * itemRepetition)
			}
		}


	}


	for k, v := range wordsRanking {
		log.Printf("word: %s ranking:%d", k, v)
	}


}