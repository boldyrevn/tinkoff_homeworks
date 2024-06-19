package httpvanilla

import (
    "homework/internal/app"
    "net/http"
)

func AppRouter(r *http.ServeMux, uc app.UseCase) {
    r.HandleFunc("/api/device", func(writer http.ResponseWriter, request *http.Request) {
        switch request.Method {
        case "GET":
            GetDevice(uc)(writer, request)
        case "POST":
            CreateDevice(uc)(writer, request)
        case "PUT":
            UpdateDevice(uc)(writer, request)
        case "DELETE":
            DeleteDevice(uc)(writer, request)
        }
    })
}
