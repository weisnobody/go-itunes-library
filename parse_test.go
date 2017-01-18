package ituneslib_test

import (
    "testing"

    "github.com/rydrman/go-itunes-library"
)

func ParseFileTest(*testing.T) {
    itunes.ParseFile("~/Downloads/iTunes Music Library.xml")
}
