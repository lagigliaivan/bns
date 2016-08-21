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

/*
	itemsDescriptions := [...]ItemDescription{
		ItemDescription{ItemId:"fee34fea0d0bb91048a350ab21605380075e95db", Description:"pollo pechuga"},
		ItemDescription{ItemId:"b71c4c694b630b6241ebc3ef92a896fb411f9ae7", Description:"verdura"},
		ItemDescription{ItemId:"99cab8664ebef2109dce6407feb8437da5506b60", Description:"chicle topline"},
		ItemDescription{ItemId:"7cae0f5f4ed22f36d5968118948e19882bd12ad8", Description:"manteca la Seren\u00edsima 130gr"},
		ItemDescription{ItemId:"e44b105fa1d18dc8b7a033cad8a0193bd69cf785", Description:"tarjeta visa"},
		ItemDescription{ItemId:"179c175d205590fedc01ba0f91c757d2a1931626", Description:"vino fuego negro"},
		ItemDescription{ItemId:"b40d5cf745ade72b2256359a7f37067fbc73c715", Description:"queso sancor"},
		ItemDescription{ItemId:"cdbe9b45530dd8adcf9bb66dfd5fdcdfb20da71a", Description:"vino familia gascon"},
		ItemDescription{ItemId:"cf6dc0c32dd646e60cd0313eee7ad99c8e4d940a", Description:"vino tinto elementos"},
		ItemDescription{ItemId:"48319a60b4f1510d2f224bdb466f5e8cb604d46e", Description:"queso cremoso sue\u00f1o"},
		ItemDescription{ItemId:"dce7c45f5eebfac0a4990c5909986f54b7b56924", Description:"jam\u00f3n cocido piamontesa 150gt"},
		ItemDescription{ItemId:"9e9fdf9d3b9b0ca706a23b819e5b54d40f1b199c", Description:"rollo cocina cartabella"},
		ItemDescription{ItemId:"ac89876ede16a5f579ee8c874bd4e80bd18fb3ee", Description:"manteca la caba\u00f1a 200gr"},
		ItemDescription{ItemId:"8000c386381808ee850446ef153e42d1822cc55b", Description:"costeleta almacor"},
		ItemDescription{ItemId:"690b15d4f8d1d4e866cc5250916cb9aafebd18b1", Description:"vino misterio"},
		ItemDescription{ItemId:"b672c3011133a805c41000e42e8f89ba967212f7", Description:"carne tapa de asado"},
		ItemDescription{ItemId:"827c7065857c9958809807a898128ab455fe1161", Description:"lechuga"},
		ItemDescription{ItemId:"13866c94380c18546dee3e3b0468e3767ec8c898", Description:"queso cremoso manfrey"},
		ItemDescription{ItemId:"153a9eb2796aa197a95d61183a6de59622ff5e22", Description:"vinagre favinco"},
	}*/

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