package utils

import (
    "errors"
    "strconv"
    "strings"
)

var OctetsCountError = errors.New("wrong count of octets")
var IncorrectOctetNumber = errors.New("incorrect number in bit octet")
var UnfilledError = errors.New("all fields must be filled")

func ValidateIP(addr string) error {
    parts := strings.Split(addr, ".")
    if len(parts) != 4 {
        return OctetsCountError
    }
    for _, part := range parts {
        n, err := strconv.Atoi(part)
        if err != nil || !(0 <= n && n <= 255) {
            return IncorrectOctetNumber
        }
    }
    return nil
}
