package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Request Recieved:", r.Method, r.URL.Path)
	fmt.Println(w, "OK")
}

func main(){
	port := os.Getenv("PORT")
	if port == ""{
		port ="8080"
}

mux := http.NewServeMux()
mux.HandleFunc("/healthz", healthCheckHandler)

server := &http.Server{
	Addr: ":" + port,
	Handler: mux,
	ReadTimeout: 5 * time.Second,
	WriteTimeout: 10 * time.Second,
	IdleTimeout: 120 * time.Second,
}

log.Printf("Server running on port %s",port)

if err := server.ListenAndServe(); err != nil {
	log.Fatalf("Server failed:%v",err)
}

}

