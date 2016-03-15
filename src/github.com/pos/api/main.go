// @APIVersion 1.0.0
// @APITitle Teamwork Desk
// @APIDescription Bend Teamwork Desk to your will using these read and write endpoints
// @Contact support@teamwork.com
// @TermsOfServiceUrl https://www.teamwork.com/termsofservice
// @License BSD
// @LicenseUrl http://opensource.org/licenses/BSD-2-Clause
package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"
//	"flag"
)

func main() {

	/*var (
	        staticContent = flag.String("staticPath", "../../../swaggerui", "Path to folder with Swagger UI")
	        //apiurl = flag.String("api", "http://127.0.0.1", "The base path URI of the API service")
	)*/
	router := mux.NewRouter();

	db := infrastructure.CatalogDB{}
	service := NewService(db)

	router.HandleFunc("/catalog/products/{id}", service.HandleRequest)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../../../swaggerui/")))
	http.Handle("/", router)
	router.Methods("GET", "PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}