package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	_ "go.uber.org/automaxprocs"
)

func Test_resolve(t *testing.T) {
	type args struct {
		domain string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "resolve one.one.one.one.",
			args: args{domain: "one.one.one.one."},
			want: []string{"1.0.0.1/32", "1.1.1.1/32"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolve(tt.args.domain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dnsResolve(t *testing.T) {

	req := httptest.NewRequest(http.MethodPost, "/dns-resolve", strings.NewReader("item=one.one.one.one"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dnsResolve)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"Nets":["1.0.0.1/32","1.1.1.1/32"]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
