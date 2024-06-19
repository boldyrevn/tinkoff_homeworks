package main

import (
    "flag"
    "github.com/boldyrevn/servconf"
    "homework/internal/adapters/maprepo"
    "homework/internal/app"
    "homework/internal/ports/httpvanilla"
    "log"
    "net/http"
)

func main() {
    path := flag.String("conf", "", "path to config file")
    flag.Parse()
    if *path == "" {
        log.Fatal("config file is not specified")
    }

    conf, err := servconf.LoadYamlConfig(*path)
    if err != nil {
        log.Fatal("error during reading config file")
    }

    r := http.NewServeMux()
    httpvanilla.AppRouter(r, app.NewUseCase(maprepo.NewRepository()))

    s := http.Server{
        Addr:         conf.Addr,
        ReadTimeout:  conf.ReadTimeout,
        WriteTimeout: conf.WriteTimeout,
        Handler:      r,
    }
    log.Fatal(s.ListenAndServe())
}
