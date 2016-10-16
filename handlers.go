package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func homeIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	app, ok := ctx.Value("app").(AppContext)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		rnd.JSON(w, http.StatusOK, app)
	} else {
		rnd.HTML(w, http.StatusOK, "home", app)
	}
}

func homePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func homePanic(w http.ResponseWriter, r *http.Request) {
	panic("ERROR!!!")
}

func eventIngestHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var p Payload
	err := decoder.Decode(&p)
	if err != nil {
		log.Println(err.Error())
		response := Status{
			Status:  "400",
			Message: "malformed data",
		}
		rnd.JSON(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	app, _ := ctx.Value("app").(AppContext)
	if app.Role == "proxy" && app.Endpoints != "" {
		endpoints := strings.Split(app.Endpoints, ",")
		for _, endpoint := range endpoints {
			err := redirectPaylod(p, endpoint+"/ingest/data")
			if err != nil {
				response := Status{
					Status:  "502",
					Message: "endpoint " + endpoint + " error: " + err.Error(),
				}
				rnd.JSON(w, http.StatusBadGateway, response)
			}
		}
	}
}

func redirectPaylod(p Payload, url string) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(p)

	r, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}

func dataIngestHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var p Payload
	err := decoder.Decode(&p)
	if err != nil {
		log.Println(err.Error())
		response := Status{
			Status:  "400",
			Message: "malformed data",
		}
		rnd.JSON(w, http.StatusBadRequest, response)
		return
	}
}
func emptyHandler(w http.ResponseWriter, r *http.Request) {

}
