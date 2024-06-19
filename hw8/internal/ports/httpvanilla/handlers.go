package httpvanilla

import (
    "encoding/json"
    "errors"
    "fmt"
    "homework/internal/app"
    "homework/internal/device"
    "log"
    "net/http"
)

func writeAnswer(w http.ResponseWriter, a any, code int) {
    m, _ := json.Marshal(a)
    w.WriteHeader(code)
    _, err := w.Write(m)
    if err != nil {
        log.Println(err)
    }
}

func writeOK(w http.ResponseWriter) {
    writeAnswer(w, Message{Text: ""}, http.StatusOK)
}

func writeError(w http.ResponseWriter, err error, code int) {
    writeAnswer(w, Message{Text: err.Error()}, code)
}

func decodeDevice(r *http.Request) (device.Device, error) {
    dc := json.NewDecoder(r.Body)
    var dev device.Device
    err := dc.Decode(&dev)
    if err != nil {
        return device.Device{}, fmt.Errorf("wrong json format: %s", err.Error())
    }
    return dev, nil
}

func GetDevice(uc app.UseCase) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        num := request.URL.Query().Get("num")
        if num == "" {
            writeError(writer, errors.New("query must contain `num` parameter"), http.StatusBadRequest)
            return
        }
        d, err := uc.GetDevice(num)
        if err != nil {
            writeError(writer, err, http.StatusNotFound)
            return
        }
        writeAnswer(writer, d, http.StatusOK)
    }
}

func CreateDevice(uc app.UseCase) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        dev, err := decodeDevice(request)
        if err != nil {
            writeError(writer, err, http.StatusBadRequest)
            return
        }
        err = uc.CreateDevice(dev)
        if err != nil {
            writeError(writer, err, http.StatusBadRequest)
            return
        }
        writeOK(writer)
    }
}

func DeleteDevice(uc app.UseCase) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        num := request.URL.Query().Get("num")
        if num == "" {
            writeError(writer, errors.New("query must contain `num` parameter"), http.StatusBadRequest)
            return
        }
        err := uc.DeleteDevice(num)
        if err != nil {
            writeError(writer, err, http.StatusNotFound)
            return
        }
        writeOK(writer)
    }
}

func UpdateDevice(uc app.UseCase) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        dev, err := decodeDevice(request)
        if err != nil {
            writeError(writer, err, http.StatusBadRequest)
            return
        }
        err = uc.UpdateDevice(dev)
        if err != nil {
            writeError(writer, err, http.StatusNotFound)
            return
        }
        writeOK(writer)
    }
}
