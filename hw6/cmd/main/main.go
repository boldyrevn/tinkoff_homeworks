package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "homework/internal/device"
    "homework/internal/handler"
    "log"
    "net/http"
    "os"
)

func main() {
    err := godotenv.Load("homework/.env")
    if err != nil {
        log.Fatal(err)
    }

    h := handler.Handler{
        UseCase: device.NewUseCase(device.NewRepository()),
    }

    http.HandleFunc("/api/device", func(writer http.ResponseWriter, request *http.Request) {
        switch request.Method {
        case http.MethodGet:
            h.GetDevice(writer, request)
        case http.MethodPost:
            h.CreateDevice(writer, request)
        case http.MethodPut:
            h.UpdateDevice(writer, request)
        case http.MethodDelete:
            h.DeleteDevice(writer, request)
        }
    })

    url := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
    log.Fatal(http.ListenAndServe(url, nil))
}
