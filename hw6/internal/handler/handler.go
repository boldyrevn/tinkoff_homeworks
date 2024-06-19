package handler

import (
    "encoding/json"
    "errors"
    "homework/internal/device"
    "homework/internal/model"
    "log"
    "net/http"
)

type Handler struct {
    UseCase device.UseCase
}

func writeAnswer(w http.ResponseWriter, a any, code int) {
    m, _ := json.MarshalIndent(a, "", "  ")
    w.WriteHeader(code)
    _, err := w.Write(m)
    if err != nil {
        log.Println(err)
    }
}

func writeMessage(w http.ResponseWriter, s string) {
    writeAnswer(w, model.Message{Text: s}, http.StatusOK)
}

func writeError(w http.ResponseWriter, err error, code int) {
    log.Printf("Error: %s %v", err.Error(), code)
    writeAnswer(w, model.Message{Text: err.Error()}, code)
}

func decodeDevice(r *http.Request) (model.Device, error) {
    dc := json.NewDecoder(r.Body)
    var dev model.Device
    err := dc.Decode(&dev)
    if err != nil {
        return model.Device{}, errors.New("wrong json format")
    }
    return dev, nil
}

func (h Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
    num := r.URL.Query().Get("num")
    if num == "" {
        writeError(w, errors.New("query must contain `num` parameter"), http.StatusBadRequest)
        return
    }
    d, err := h.UseCase.GetDevice(num)
    if err != nil {
        writeError(w, err, http.StatusNotFound)
        return
    }
    writeAnswer(w, d, http.StatusOK)
}

func (h Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
    dev, err := decodeDevice(r)
    if err != nil {
        writeError(w, err, http.StatusBadRequest)
        return
    }
    err = h.UseCase.CreateDevice(dev)
    if err != nil {
        writeError(w, err, http.StatusBadRequest)
        return
    }
    writeMessage(w, "device was created successfully")
}

func (h Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
    num := r.URL.Query().Get("num")
    if num == "" {
        writeError(w, errors.New("query must contain `num` parameter"), http.StatusBadRequest)
        return
    }
    err := h.UseCase.DeleteDevice(num)
    if err != nil {
        writeError(w, err, http.StatusNotFound)
        return
    }
    writeMessage(w, "device was deleted successfully")
}

func (h Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
    dev, err := decodeDevice(r)
    if err != nil {
        writeError(w, err, http.StatusBadRequest)
        return
    }
    err = h.UseCase.UpdateDevice(dev)
    if err != nil {
        writeError(w, err, http.StatusNotFound)
        return
    }
    writeMessage(w, "device was updated successfully")
}
