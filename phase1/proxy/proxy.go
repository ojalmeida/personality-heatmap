package proxy

import (
	"context"
	"encoding/json"
	"github.com/elazarl/goproxy"
	"github.com/ojalmeida/personality-heatmap/phase1/data"
	"github.com/ojalmeida/personality-heatmap/phase1/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	NodeName string
	FakeGPS  models.FakeLocation
}

var ProxyConfig Config

var server *http.Server
var proxy = goproxy.NewProxyHttpServer()

func init() {

	file, err := os.OpenFile("/tmp/proxy.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 660)

	if err != nil {

		panic(err.Error())

	}

	log.SetOutput(file)

	err = setCustomCertificate()

	if err != nil {
		panic(err.Error())
	}

	proxy.Verbose = false
	proxy.Logger = log.Default()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	// Modify geolocation
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {

		if resp.Request.Host == "www.googleapis.com" {

			originalResponseBody, err := ioutil.ReadAll(resp.Body)

			if err != nil {

				return resp

			}

			modifiedResponseBody, err := json.Marshal(ProxyConfig.FakeGPS)

			resp.Body = ioutil.NopCloser(strings.NewReader(string(modifiedResponseBody)))

			log.Println("Geolocation request intercepted")
			log.Println("Original response: ", string(originalResponseBody))
			log.Println("Modified response: ", string(modifiedResponseBody))

			resp.StatusCode = 200
			resp.Status = "200 OK"

		}

		return resp

	})

	// Store API Token
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {

		if resp.Request.Host == "api.gotinder.com" {

			xAuthToken := resp.Request.Header.Get("X-Auth-Token")

			if xAuthToken != "" {

				for _, node := range data.Data.Nodes {

					if node.Name == ProxyConfig.NodeName {

						if node.APIToken == "" {

							log.Printf("X-Auth-Token: %s stored for node %s\n", xAuthToken, node.Name)
							node.APIToken = xAuthToken

							break
						}
					}

				}

			}

		}

		return resp

	})

}

func Start() {

	server = &http.Server{
		Addr:    ":8888",
		Handler: proxy,
	}

	log.Println("Starting proxy")

	go func() {

		err := server.ListenAndServe()

		if err != nil {

			log.Println(err.Error())

		}

	}()

}

func Stop() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {

		panic(err.Error())

	}

}
