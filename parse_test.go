package itunes_test

import (
    "testing"

    "github.com/rydrman/go-itunes-library"
)

func TestParseFile(t *testing.T) {

    _, err := itunes.ParseFile("/home/Downloads/rbottriell/iTunes Music Library.xml")
    if nil != err {
        t.Fatal(err)
    }

    t.Fail()

}
