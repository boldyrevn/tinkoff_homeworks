package handler

import (
    "encoding/json"
    "errors"
    "homework/app"
    "homework/model"
    "log"
    "net/http"
    "strconv"
    "strings"
)

type Handler struct {
    Service app.Service
}

func writeAnswer(w http.ResponseWriter, a any, code int) {
    m, _ := json.MarshalIndent(a, "", "   ")
    w.WriteHeader(code)
    _, err := w.Write(m)
    if err != nil {
        log.Println(err)
    }
}

func writeSuccess(w http.ResponseWriter, s string) {
    writeAnswer(w, model.Message{Text: s}, http.StatusOK)
}

func writeError(w http.ResponseWriter, err error) {
    writeAnswer(w, model.Message{Text: err.Error()}, http.StatusBadRequest)
}

func checkIP(addr string) bool {
    parts := strings.Split(addr, ".")
    for _, part := range parts {
        n, err := strconv.Atoi(part)
        if err != nil || !(0 <= n && n <= 255) {
            return false
        }
    }
    return true
}

func decodeDevice(r *http.Request) (model.Device, error) {
    dc := json.NewDecoder(r.Body)
    var dev model.Device
    err := dc.Decode(&dev)
    if err != nil {
        return model.Device{}, errors.New("wrong json format")
    } else if dev.IP == "" || dev.SerialNum == "" || dev.Model == "" {
        return model.Device{}, errors.New("all fields must be filled")
    } else if !checkIP(dev.IP) {
        return model.Device{}, errors.New("wrong IP address format")
    }
    return dev, nil
}

func (h Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
    num := r.URL.Query().Get("num")
    if num == "" {
        writeError(w, errors.New("query must contain `num` parameter"))
        return
    }
    d, err := h.Service.GetDevice(num)
    if err != nil {
        writeError(w, err)
        return
    }
    writeAnswer(w, d, http.StatusOK)
}

func (h Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
    dev, err := decodeDevice(r)
    if err != nil {
        writeError(w, err)
        return
    }
    err = h.Service.CreateDevice(dev)
    if err != nil {
        writeError(w, err)
        return
    }
    writeSuccess(w, "device was created successfully")
}

func (h Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
    num := r.URL.Query().Get("num")
    if num == "" {
        writeError(w, errors.New("query must contain `num` parameter"))
        return
    }
    err := h.Service.DeleteDevice(num)
    if err != nil {
        writeError(w, err)
        return
    }
    writeSuccess(w, "device was deleted successfully")
}

func (h Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
    dev, err := decodeDevice(r)
    if err != nil {
        writeError(w, err)
        return
    }
    err = h.Service.UpdateDevice(dev)
    if err != nil {
        writeError(w, err)
        return
    }
    writeSuccess(w, "device was updated successfully")
}
