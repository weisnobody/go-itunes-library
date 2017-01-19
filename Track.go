package itunes

import (
    "encoding/xml"
    "fmt"
    "time"
)

// Track represents an entire iTunes Track data structure
type Track struct {
    TrackID             int
    Size                int
    TotalTime           int
    DiscNumber          int
    DiscCount           int
    TrackNumber         int
    Year                int
    DateModified        time.Time
    DateAdded           time.Time
    BitRate             int
    SampleRate          int
    PlayCount           int
    PlayDate            int
    PlayDateUTC         time.Time
    SkipCount           int
    SkipDate            time.Time
    Rating              int
    AlbumRating         int
    AlbumRatingComputed bool
    PersistentID        string
    TrackType           string
    FileFolderCount     int
    LibraryFolderCount  int
    Name                string
    Artist              string
    AlbumArtist         string
    Composer            string
    Album               string
    Genre               string
    Kind                string
    Location            string
}

// UnmarshalXML is a custom unmarshaller function for itunes Track xml format
func (lib *Track) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

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
