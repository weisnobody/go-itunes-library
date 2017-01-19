package itunes

import (
    "encoding/xml"
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "time"
)

// readOpenTagAsString reads the current open tag in decoder and returns the
// string upon success
func readOpenTagAsString(decoder *xml.Decoder) (string, error) {

    value := ""
    for {

        token, err := decoder.Token()
        if nil != err {
            return "", err
        }

        switch t := token.(type) {

        case xml.CharData:
            value = string(t)

        case xml.EndElement:
            return value, nil

        default:
            return "", ErrInvalidFormat

        }
    }

}

// resolveKeyOnStruct resolves the given key by reading the expected data from decoder
func resolveKeyOnStruct(item interface{}, key string, decoder *xml.Decoder) error {

    token, err := decoder.Token()
    if err != nil {
        return err
    }

    // get information about posiible matching
    // struct field
    fieldName := strings.Replace(key, " ", "", -1)
    field := reflect.ValueOf(item).Elem().FieldByName(fieldName)
    if !field.IsValid() {
        fmt.Printf("skipping missing field: %s\n", fieldName)
        _ = decoder.Skip()
        return nil
    }

    // convert the type based on the field type
    switch field.Interface().(type) {

    case string:
        v, err := readOpenTagAsString(decoder)
        if nil != err {
            return err
        }
        field.SetString(v)

    case int:
        v, err := readOpenTagAsString(decoder)
        if nil != err {
            return err
        }
        i, err := strconv.Atoi(v)
        if nil != err {
            return err
        }
        field.SetInt(int64(i))

    case bool:
        v := token.(xml.StartElement).Name.Local
        field.SetBool(v == "true")

    case time.Time:
        v, err := readOpenTagAsString(decoder)
        if nil != err {
            return err
        }
        t, err := time.Parse(time.RFC3339, v)
        if nil != err {
            return err
        }
        field.Set(reflect.ValueOf(t))

    case []Track:

        start, err := findNextStartElement(decoder, token)
        if nil != err {
            return err
        }

        tracks, err := decodeTracks(decoder, start)
        if err != nil {
            return err
        }
        field.Set(reflect.ValueOf(tracks))

    default:
        fmt.Printf("unknown field type for %s\n", key)
        _ = decoder.Skip()

    }

    return nil

}

func findNextStartElement(decoder *xml.Decoder, begin xml.Token) (xml.StartElement, error) {

    token := begin
    var err error

    for {

        switch t := token.(type) {

        case xml.StartElement:

            return t, nil

        default:

            token, err = decoder.Token()
            if err != nil {
                return xml.StartElement{}, err
            }

        }

    }

}

func decodeTracks(decoder *xml.Decoder, start xml.StartElement) ([]Track, error) {

    tracks := make([]Track, 0)

    for {
        token, err := decoder.Token()
        if err != nil {
            return nil, err
        }

        switch t := token.(type) {

        case xml.StartElement:

            if t.Name.Local == "key" {
                continue
            }

            if t.Name.Local == "dict" {

                track := Track{}
                track.UnmarshalXML(decoder, t)
                tracks = append(tracks, track)
                fmt.Println(track.Name)

            } else {

                return nil, ErrInvalidFormat

            }

        case xml.EndElement:
            if t == start.End() {

                return tracks, nil

            }

        default:
            return nil, ErrInvalidFormat

        }

    }

}
