package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func routeLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[Request URI]: ", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("authorization")
		if !strings.HasPrefix(authHeader, "Awoo ") {
			res := Message{
				Status:  "failed",
				Message: "Invalid token",
			}
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(res)
			return
		}
		auth, err := getUserAuth(authHeader[5:])
		if err != nil {
			res := Message{
				Status:  "failed",
				Message: "Bad request. Wrong token or invalid token",
			}
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(res)
			return
		}
		expiry := auth.LastAccess.Add(15 * time.Minute)
		if compareTime(time.Now(), expiry) == -1 {
			next.ServeHTTP(w, r)
		} else {
			res := Message{
				Status:  "failed",
				Message: "Unauthorized. Please login again.",
			}
			setJSONHeader(w).WriteHeader(401)
			json.NewEncoder(w).Encode(res)
			return
		}

	})
}
