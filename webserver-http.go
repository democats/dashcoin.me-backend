package main

import (
    "log"
    "net/http"
)

func redirectTLS(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://api.dashcoin.me:443"+r.RequestURI, http.StatusMovedPermanently)
}

func main() {
    http.Handle("/.well-known/acme-challenge/", http.FileServer(http.FileSystem(http.Dir("/var/tmp/letsencrypt/"))))

    if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
        log.Fatalf("ListenAndServe error: %v", err)
    }
}
