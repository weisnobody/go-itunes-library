package itunes

import (
    "encoding/xml"
    "io/ioutil"
    "os"
)

// ParseFile parses the file at the given filepath
// as an itunes library file
func ParseFile(filename string) (*Library, error) {

    libraryFile, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer libraryFile.Close()

    bytes, err := ioutil.ReadAll(libraryFile)
    if err != nil {
        return nil, err
    }

    return ParseBytes(bytes)

}

// ParseBytes parses the given byte slice as an library xml file
func ParseBytes(bytes []byte) (*Library, error) {

    lib := &Library{}

    err := xml.Unmarshal(bytes, lib)
    if nil != err {
        return nil, err
    }

    return lib, nil

}
