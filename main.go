package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/caarlos0/env/v6"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	_ "go.uber.org/automaxprocs"
)

var wg sync.WaitGroup

var (
	reqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dnsapi_processed_req_total",
		Help: "The total number of requests to resolve",
	})
)

type config struct {
	Host  string `env:"HOST" envDefault:"0.0.0.0"`
	Port  int    `env:"PORT" envDefault:"8080"`
	Debug bool   `env:"DEBUG" envDefault:false`
}

func resolve(domain string) []string {
	// resolve domain and return ips
	var ans []string

	domain = strings.Trim(domain, " ")
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", domain)
	if err != nil {
		log.Info(err)
	}
	for _, ip := range ips {
		ip := ip.String() + "/32"
		ans = append(ans, ip)
	}
	reqProcessed.Inc()
	return ans
}

func dnsResolve(w http.ResponseWriter, r *http.Request) {
	// Get domains in `item` and resolve them concurrently
	// curl --data "item=ya.ru,mail.ru,  google.com" -v localhost:8080/dns-resolve
	// {"Nets":["87.250.250.242/32","217.69.139.200/32","94.100.180.201/32",
	// "94.100.180.200/32","217.69.139.202/32","64.233.162.100/32",
	// "64.233.162.102/32","64.233.162.101/32","64.233.162.139/32",
	// "64.233.162.138/32","64.233.162.113/32"]}

	var ips []string

	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 - Just use POST!"))

	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Info(err)
			return
		}
		domains := r.FormValue("item")
		if len(domains) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - We don't found values in your request!"))
			return
		}
		s := strings.Split(domains, ",")
		for _, domain := range s {
			res := make(chan []string)
			wg.Add(1)
			go func(d string, res chan []string) {
				defer wg.Done()
				domains := resolve(d)
				res <- domains
				close(res)
			}(domain, res)
			ips = append(ips, <-res...)
		}

	}

	js, err := json.Marshal(map[string][]string{"Nets": ips})
	if err != nil {
		http.Error(w, "Can't Marshal json", http.StatusInternalServerError)
		log.Info(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetFormatter(&log.JSONFormatter{})

	http.HandleFunc("/dns-resolve", dnsResolve)
	http.Handle("/metrics", promhttp.Handler())

	addr := cfg.Host + ":" + fmt.Sprint(cfg.Port)
	http.ListenAndServe(addr, nil)

}
