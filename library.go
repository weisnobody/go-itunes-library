package itunes

import (
    "encoding/xml"
    "fmt"
    "time"
)

// Library represents an entire iTunes library data structure
type Library struct {
    MajorVersion        int
    MinorVersion        int
    ApplicationVersion  string
    Date                time.Time
    Features            int
    ShowContentRatings  bool
    LibraryPersistentID string
    MusicFolder         string
    Tracks              []Track
    Playlists           []interface{}
}

// UnmarshalXML is a custom unmarshaller function for itunes library xml format
func (lib *Library) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

    for {
        token, err := decoder.Token()
        if err != nil {
            return err
        }

        switch t := token.(type) {

        case xml.StartElement:

            if t.Name.Local == "key" {

                key, err := readOpenTagAsString(decoder)
                if err != nil {
                    return err
                }

                resolveKeyOnStruct(lib, key, decoder)

            } else {

                fmt.Printf("skip %s\n", t.Name.Local)
                _ = decoder.Skip()

            }

        }

    }

}
