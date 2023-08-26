package xapikey

import "net/http"



func CheckApiKey(next http.Handler) http.Handler {
	availableKey :=  map[string]bool {
		"qwerty123": true,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !availableKey[r.Header.Get("X-Auth-Key")] {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w,r)
	})
}