package itunes_test

import (
    "testing"

    "github.com/rydrman/go-itunes-library"
)

func TestParseFile(t *testing.T) {

    lib, err := itunes.ParseFile("test_data/test_library_1.xml")
    if nil != err {
        t.Fatal(err)
    }

    t.Logf("Library loaded successfully!\n")
    t.Logf(" - found %d tracks and %d playlists", len(lib.Tracks), len(lib.Playlists))

}
