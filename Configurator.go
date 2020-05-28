package main

import (
    "github.com/BurntSushi/toml"
    "io/ioutil"
    "log"
)

type Config struct {
    Username string
    Password string
    IP       string
    PORT     string
}

var adminAuthUser JSONLogin
var IP string
var PORT string

func readConfig() {
    configBytes, err := ioutil.ReadFile(configPath + configName)
    if err != nil {
        log.Fatal("ERROR: Failed to read config file:\n" + err.Error())
    }
    tomlData := string(configBytes)
    var conf Config
    if _, err := toml.Decode(tomlData, &conf); err != nil {
        log.Fatal("ERROR: Failed to parse config file:\n " + err.Error())
    }

    adminAuthUser.Username = conf.Username
    adminAuthUser.Password = conf.Password
    IP = conf.IP
    PORT = conf.PORT
}
