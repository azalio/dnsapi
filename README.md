# Resolve domain(s) and return in json

## Main

```bash
curl -s --data "item=ya.ru,mail.ru,  google.com" -v localhost:8080/dns-resolve | jq

{
  "Nets": [
    "87.250.250.242/32",
    "94.100.180.201/32",
    "94.100.180.200/32",
    "217.69.139.202/32",
    "217.69.139.200/32",
    "64.233.162.100/32",
    "64.233.162.113/32",
    "64.233.162.101/32",
    "64.233.162.139/32",
    "64.233.162.138/32",
    "64.233.162.102/32"
  ]
}
```

## Config

You can set env variables

```golang
type config struct {
    Host  string `env:"HOST" envDefault:"0.0.0.0"`
    Port  int    `env:"PORT" envDefault:"8080"`
    Debug bool   `env:"DEBUG" envDefault:false`
}
```

## K8s manifests

```make
make make_manifests
```

## Metrics

```bash
curl -s localhost:8080/metrics | grep dnsapi_processed_req_total
# HELP dnsapi_processed_req_total The total number of requests to resolve
# TYPE dnsapi_processed_req_total counter
dnsapi_processed_req_total 0
```
