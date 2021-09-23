package main

import (
	"errors"
	"log"
	"net/http"
)

func NewAuthMiddleware(endpoint string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := Refresh(w, r, endpoint)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func Refresh(w http.ResponseWriter, r *http.Request, endpoint string) error {
	client := &http.Client{}
	authReq, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	for _, cookie := range r.Cookies() {
		authReq.AddCookie(cookie)
	}
	authReq.Header.Set("Accept", "application/json")
	authResp, err := client.Do(authReq)
	if err != nil {
		return err
	}
	log.Printf("auth response: %d", authResp.StatusCode)
	for _, cookie := range authResp.Cookies() {
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
	}
	if authResp.StatusCode != http.StatusOK {
		return errors.New("unauthorized")
	}
	return nil
}
