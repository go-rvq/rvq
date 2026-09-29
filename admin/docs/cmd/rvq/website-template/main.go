package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-rvq/rvq/admin/docs/cmd/rvq/website-template/admin"
	"github.com/go-rvq/rvq/x/osenv"
)

var (
	port       = osenv.Get("RVQ_PORT", "The port to serve on", "9001")
	publishURL = osenv.Get("RVQ_PUBLISH_URL", "Publish Target URL", "")
)

func main() {
	// CMS server

	cmsMux := admin.InitApp()
	cmsServer := &http.Server{
		Addr:    ":" + port,
		Handler: cmsMux,
	}
	go cmsServer.ListenAndServe()
	fmt.Println("CMS Served at http://localhost:" + port + "/admin")

	// Publish server
	u, _ := url.Parse(publishURL)
	publishPort := u.Port()
	if publishPort == "" {
		publishPort = "9001"
	}
	publishMux := http.FileServer(http.Dir(admin.PublishDir))
	publishServer := &http.Server{
		Addr:    ":" + publishPort,
		Handler: publishMux,
	}
	fmt.Println("Publish Served at http://localhost:" + publishPort)
	log.Fatal(publishServer.ListenAndServe())
}
