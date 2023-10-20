package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "homework/app"
    "homework/handler"
    "log"
    "net/http"
    "os"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }

    h := handler.Handler{
        Service: app.NewService(),
    }

    http.HandleFunc("/api/device", func(writer http.ResponseWriter, request *http.Request) {
        switch request.Method {
        case "GET":
            h.GetDevice(writer, request)
        case "POST":
            h.CreateDevice(writer, request)
        case "PUT":
            h.UpdateDevice(writer, request)
        case "DELETE":
            h.DeleteDevice(writer, request)
        }
    })

    url := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
    log.Fatal(http.ListenAndServe(url, nil))
}
