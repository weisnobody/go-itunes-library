package itunes

import (
    "encoding/xml"
    "fmt"
)

// Playlist represents an entire iTunes Playlist data structure
type Playlist struct {
    Master               bool
    PlaylistID           int
    PlaylistPersistentID string
    AllItems             bool
    Visible              bool
    Protected            bool
    Name                 string
    Description          string
    SmartInfo            []byte
    SmartCriteria        []byte
    ParentPersistentID   string
    DistinguishedKind    int
    Kind                 string
    PurchasedMusic       bool
    Disliked             bool
    Music                bool
    Movies               bool
    TVShows              bool
    Podcasts             bool
    ITunesU              bool
    Audiobooks           bool
    Books                bool
    Folder               bool
    PlaylistItems        []*Track
}

func (p *Playlist) String() string {

    return fmt.Sprintf("Playlist: %s (%d tracks)", p.Name, len(p.PlaylistItems))

}

func PlaylistDistinguished(kind int) string {
    newKind := fmt.Sprintf("%v", kind)
    
    kinds := map[int]string{
        0: "",
        2: "Movies",
        3: "TV Shows",
        4: "Music",
        5: "Audiobooks",
        10: "Podcasts",
        17: "Voice Memos",
        19: "Purchased",
        65: "Downloaded",
        66: "Downloaded",
        67: "Downloaded",
    }

    if _, exists := kinds[kind]; exists {
        newKind = kinds[kind]
    }
    
    if false && newKind != fmt.Sprintf("%v", kind) {
        fmt.Println(fmt.Sprintf("Cleaning Distinguished: %s (%s)", kind, newKind))
    }

    return newKind

}

// UnmarshalXML is a custom unmarshaller function for itunes Track xml format
func (p *Playlist) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

    for {
        token, err := decoder.Token()
        if err != nil {
            return err
        }

        switch tok := token.(type) {

        case xml.StartElement:

            if tok.Name.Local == "key" {

                key, err := readOpenTagAsString(decoder)
                if err != nil {
                    return err
                }

                err = resolveKeyOnStruct(p, key, decoder)
                if err != nil {
                    return err
                }

            } else {

                fmt.Printf("skip %s (Playlist)\n", tok.Name.Local)
                _ = decoder.Skip()

            }

        case xml.EndElement:

            if start.End() == tok {
                return nil
            }

            return NewInvalidFormatError(
                fmt.Sprintf("Unexpected end element in playlist: %s", tok.Name.Local))

        }

    }

}
