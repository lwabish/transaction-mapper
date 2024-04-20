package util

import (
	"log"
	"strconv"
	"strings"
)

func ParseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println(err)
	}
	return f
}
