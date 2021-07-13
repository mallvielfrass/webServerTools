package wst

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mallvielfrass/fmc"
)

type parseQ struct {
	Key, Value string
}

func parseQuery(query string) []parseQ {
	if query != "" {
		var p []parseQ
		arrayQuery := strings.Split(query, "&")
		for _, item := range arrayQuery {
			res := strings.Split(item, "=")
			p = append(p, parseQ{Key: res[0], Value: res[1]})
		}
		return p
	}
	return []parseQ{}
}
func colorQUery(p []parseQ) string {
	var str string
	parse := ", "
	lp := len(p) - 1
	for i, item := range p {
		if i == lp {
			parse = ""
		}
		str += fmt.Sprintf("#gbt'#wbt%s#gbt'#ybt:#gbt'#wbt%s#gbt'#ybt%s", item.Key, item.Value, parse)
	}
	return "#bbt[" + str + "#bbt]"
}
func MiddlewareURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmc.Printfln("[#ybtMiddlewareURL#RRR]: \n\t#bbtURL#RRR: [#gbt%s#RRR] \n\t#bbtMethod#RRR: [#gbt%s#RRR] \n\t#bbtTime#RRR: [#gbt%s#RRR]", r.URL.Path, r.Method, time.Now().Format(time.RFC1123))
		//fmc.Printfln("#gbturl#ybt: #gbt[ #bbt%s #gbt]", r.URL.Path)
		query := parseQuery(r.URL.RawQuery)
		if len(query) != 0 {
			fmc.Printfln(colorQUery(query))
		}

		next.ServeHTTP(w, r)
	})
}
func MiddlewareJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
func MiddlewareAllowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}
