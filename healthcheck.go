package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func StartHealthCheck(appCtx AppSettings) {

	ticker := time.NewTicker(time.Duration(appCtx.PingInterval) * time.Millisecond)

	for range ticker.C {
		endpoints := strings.Split(appCtx.Endpoints, ",")
		for _, endpoint := range endpoints {
			go checkTarget(endpoint, appCtx.PingInterval-10)
		}
	}
}

func checkTarget(endpoint string, timeout int) {
	begin := time.Now()
	status := 1
	isValid := false
	repl := strings.NewReplacer("http://", "", "https://", "")
	target := repl.Replace(endpoint)
	if target == "" {
		return
	}
	transport := &http.Transport{
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 0,
	}
	t := time.Duration(time.Duration(timeout) * time.Millisecond)
	transport.ResponseHeaderTimeout = t
	client := &http.Client{
		Transport: transport,
		Timeout:   t,
	}

	r, err := client.Get(endpoint + "/ping")
	if err != nil {
		http_healthcheck_total.WithLabelValues(target, strconv.Itoa(status)).Inc()
		log.Println(err.Error())
		return
	}
	defer r.Body.Close()

	status = r.StatusCode
	for i := 200; i <= 299; i++ {
		if r.StatusCode == i {
			isValid = true
			break
		}
	}

	if !isValid {
		http_healthcheck_total.WithLabelValues(target, strconv.Itoa(status)).Inc()
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status = 2
		http_healthcheck_total.WithLabelValues(target, strconv.Itoa(status)).Inc()
		return
	}

	if len(body) > 0 {
		strBody := strings.ToLower(string(body[:]))
		isValid, _ = regexp.MatchString(strings.ToLower("pong"), strBody)
	} else {
		status = 3
	}

	took := time.Since(begin)
	http_healthcheck_total.WithLabelValues(target, strconv.Itoa(status)).Inc()
	http_healthcheck_latency.WithLabelValues(target).Observe(took.Seconds())
}
