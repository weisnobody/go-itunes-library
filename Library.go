package itunes

import (
    "encoding/xml"
    "fmt"
    "io"
    "time"
    "strings"
    "strconv"
    "regexp"
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
    Music               []*Track
    MusicVideos         []int
    TVShows             []int
    Movies              []*Track
    Podcasts            []int
    Tracks              []*Track
    tracksByID          map[int]*Track
    Playlists           []*Playlist
    Artists             map[string]*Artist
    Albums              map[string]*Album
}

// String returns a nice string representation of this Library
func (lib *Library) String() string {
    return fmt.Sprintf("iTunes Library: %d Tracks, %d Playlists, [%d music, %d movies, %d tv, %d music videos, %d podcasts, %d total; %d Artists, %d Albums] (%s)",
        len(lib.Tracks), 
        len(lib.Playlists), 
        len(lib.Music),
        len(lib.Movies),
        len(lib.TVShows),
        len(lib.MusicVideos),
        len(lib.Podcasts),
        len(lib.tracksByID), 
        len(lib.Artists),
        len(lib.Albums),
        lib.LibraryPersistentID)
}

func (lib *Library) len() int {
    return len(lib.tracksByID)
}


func (lib *Library) getTrackByID(trackID int) *Track {
    return lib.tracksByID[trackID]
}

func (lib *Library) GetNumTracksX() int {
    return len(lib.tracksByID)
}

func getNumTracks(lib Library) int {
    return len(lib.tracksByID)
}


// sortTracksByID sorts the tracks in the Tracks slice by id into
// the tracksByID slice for easy lookup later
func (lib *Library) sortTracksByID() {

    byID := make(map[int]*Track)

    //MusicByID := []int{}
    Music := []*Track{}
    //MoviesByID := []int{}
    Movies := []*Track{}
    MusicVideosByID := []int{}
    TVShowsByID := []int{}
    PodcastsByID := []int{}
    
    for _, t := range lib.Tracks {

        if t.Album == "Voice Memos" && t.Genre == "Voice Memo" {
            continue
        }
        
        byID[t.TrackID] = t
        
        trackKind := "music"
        // 
        if t.Movie {
            // MoviesByID = append(MoviesByID, t.TrackID)
            //Movies = append(Movies, t)
            trackKind = "movie"
        } else if t.TVShow {
            //TVShowsByID = append(TVShowsByID, t.TrackID)
            trackKind = "tvshow"
        } else if t.Podcast {
            //PodcastsByID = append(PodcastsByID, t.TrackID)
            trackKind = "podcast"
        } else if t.MusicVideo {
            //MusicVideosByID = append(MusicVideosByID, t.TrackID)
            trackKind = "musicvideo"
        } else if strings.Contains(t.Location, "/Movies/") {
            trackKind = "movie"
        } else if strings.Contains(t.Location, "/TV%20Shows/") {
            trackKind = "tvshow"
        } else if t.Album == "Voice Memos" && t.Genre == "Voice Memo" {
            trackKind = "voice_memo"
        } else if t.TotalTime > 7200000 { // time greater than 2 hours
            trackKind = "audiobook"
        
        
        
        } else {
            // MusicByID = append(MusicByID, t.TrackID)
            //Music = append(Music, t)
            trackKind = "music"
        }
        t.TrackKind = trackKind
        
        doArtistsFor := map[string]bool {
            "music": true,
        }
        t.NameOriginal = t.Name
        if strings.ToLower(t.Name) == "[blank]" || strings.ToLower(t.Name) == "[hidden track]" {
            t.Name = fmt.Sprintf("%s-t:%v", t.Name, t.TrackNumber)
        } else if t.PersistentID == "6D073C0ED894D8DB" && t.Name == "Everlong" {
            t.Name = "Everlong (Acoustic)"
        } else if t.PersistentID == "CD0C69A4B5AF8E48" && t.Name == "Tear You Apart" {
            t.Name = "Tear You Apart (Bonus Track)"
        }
        
        
        //if _, containsKind := doArtistsFor[trackKind]; containsKind {
        if doArtistsFor[trackKind] {
            // pull together "artist" info
            
            if strings.Contains(t.Album, "(Original Motion Picture Soundtrack)") {
                t.AlbumArtist = "Soundtrack"
            }
            
            AlbumArtist := strings.TrimSpace(t.Artist)
            t.ArtistOriginal = AlbumArtist  // save original after trimming
            AlbumArtist = cleanArtist(AlbumArtist)
            t.Artist = AlbumArtist

            AlbumArtistSort := AlbumArtist
            if t.SortArtist != "" {
                //fmt.Println("Has SortArtist")
                AlbumArtistSort = strings.TrimSpace(t.SortArtist)
            } else {
                // TODO create ArtistSort
                AlbumArtistSort = strings.TrimSpace(t.Artist)
            }
            
            //fmt.Println(fmt.Sprintf("1] AA: %s, AAS: %s, SA: %s, SAA: %s", AlbumArtist, AlbumArtistSort, t.SortArtist, t.SortAlbumArtist))
            
            if t.AlbumArtist != "" {
                AlbumArtist = strings.TrimSpace(t.AlbumArtist)
                t.AlbumArtistOriginal = AlbumArtist
                AlbumArtist = cleanArtist(AlbumArtist)
                if t.SortAlbumArtist != "" {
                    //fmt.Println("Has SortAlbumArtist")
                    AlbumArtistSort = strings.TrimSpace(t.SortAlbumArtist)
                } else {
                    // TODO create ArtistSort
                    AlbumArtistSort = strings.TrimSpace(t.AlbumArtist)
                }
            }
            t.AlbumArtist = AlbumArtist
            
            
            //fmt.Println(fmt.Sprintf("2] AA: %s, AAS: %s, SA: %s, SAA: %s", AlbumArtist, AlbumArtistSort, t.SortArtist, t.SortAlbumArtist))
            
            curArtist := &Artist{}
            // for the Artists key, should likely create a "sanitized" version (all lowercase would be a start
            artistKey := strings.ToLower(AlbumArtist)
            if _, isMapContainsKey := lib.Artists[artistKey]; isMapContainsKey {
                curArtist = lib.Artists[artistKey]
            } else {
                curArtist.Name = AlbumArtist
            }
            
            //if AlbumArtist == "2008 - Disk 1" {
            //fmt.Println(fmt.Sprintf("3] AA: %s, AAS: %s, SA: %s, SAA: %s", AlbumArtist, AlbumArtistSort, t.SortArtist, t.SortAlbumArtist))
                
            if curArtist.SortName == "" || (curArtist.SortName == AlbumArtist && AlbumArtistSort != AlbumArtist) {
                if AlbumArtistSort == "various" {
                    AlbumArtistSort = "Various Artists"
                }
                //fmt.Println("Setting SortName")
                curArtist.SortName = AlbumArtistSort
            //} else {
            //    fmt.Println(fmt.Sprintf("%s == %s && %s != %s", curArtist.SortName, AlbumArtist, AlbumArtistSort, AlbumArtist))
            }

            // Should we collect Song Artist / Composer info at artist level?
            
            if curArtist.YearFirst == 0 || curArtist.YearFirst > t.Year {
                curArtist.YearFirst = t.Year
            }
            if curArtist.YearLast < t.Year {
                curArtist.YearLast = t.Year
            }
            curArtist.TotalTime += t.TotalTime
            if curArtist.DateModified.IsZero()  || curArtist.DateModified.Before(t.DateModified) {
                curArtist.DateModified = t.DateModified
            }
            //if curArtist.DateModified.IsZero() {
            //    fmt.Println(fmt.Sprintf("t: %v, c: %v", t.DateModified, curArtist.DateModified))
            //}


            if curArtist.DateAdded.IsZero() || curArtist.DateAdded.After(t.DateAdded) {
                curArtist.DateAdded = t.DateAdded
            }
            //if curArtist.DateAdded.IsZero() {
            //    fmt.Println(fmt.Sprintf("t: %v, c: %v", t.DateAdded, curArtist.DateAdded))
            //}

            // need to rethink all the rating stuff, how to aggregate songs / album ratings up through the tree
            if curArtist.Rating < t.Rating {
                // this should really be a combination of all Ratings / Loves, Likes and Dislikes
                curArtist.Rating = t.Rating
            }
            if t.Rating > 0 {
                curArtist.AlbumRatings = append(curArtist.AlbumRatings, t.Rating)
            }
            curArtist.Songs = append(curArtist.Songs, t)
            if t.Genre != "" {
                if strings.Contains(t.Genre, ",") {
                    genres := strings.Split(t.Genre, ",")
                    for _, genre := range genres {
                        genre = strings.TrimSpace(genre)
                        if _, containsGenre := curArtist.Genres[genre]; containsGenre {
                            curArtist.Genres[genre] = true
                        }
                    }
                } else if _, containsGenre := curArtist.Genres[strings.TrimSpace(t.Genre)]; containsGenre {
                    curArtist.Genres[strings.TrimSpace(t.Genre)] = true
                }
            }
            
            // increment counts
            // TODO this should be a map most like (Counts := map[string][int])
            if t.Movie {
                curArtist.CountMovies ++
            } else if t.TVShow {
                curArtist.CountTVShows ++
            } else if t.Podcast {
                curArtist.CountPodcasts ++
            } else if t.MusicVideo {
                curArtist.CountMusicVideos ++
            } else {
                curArtist.CountSongs ++
            }
            curArtist.CountTracks ++

            if t.Loved {
                curArtist.Loved ++
            }
            if t.Liked {
                curArtist.Liked ++
            }
            if t.Disliked {
                curArtist.Disliked ++
            }
            if t.AlbumLoved {
                curArtist.AlbumLoved ++
            }
            
            // pause artist stuff here to do album

            //
            // pull together "Album" info
            
            t.AlbumOriginal = t.Album
            AlbumName := strings.TrimSpace(t.Album)
            t.AlbumOriginal = AlbumName  // save Original after trimming
            AlbumName = cleanAlbum(AlbumName)
            t.Album = AlbumName

            AlbumSort := strings.TrimSpace(t.Album)
            if t.SortAlbum != "" {
                AlbumSort = strings.TrimSpace(t.SortAlbum)
            } else {
                // TODO create AlbumSort
                AlbumSort = strings.TrimSpace(t.Album)
                AlbumSort = cleanAlbum(AlbumSort)
            }
            t.SortAlbum = AlbumSort
            
            if t.DiscNumber == 0 {
                reDisk := regexp.MustCompile(`(?i)Dis[ck] \d*$`) 
                diskFound := reDisk.MatchString(t.AlbumOriginal)
                reDisk2 := regexp.MustCompile(`(?i)Dis[ck] (\d*)`) 
                diskFound2 := reDisk2.FindStringSubmatch(t.AlbumOriginal)
                if diskFound {
                    reNum := regexp.MustCompile(`\d*$`) 
                    discNum, _ := strconv.Atoi(reNum.FindString(t.AlbumOriginal))
                    t.DiscNumber = discNum
                } else if len(diskFound2) > 0 {
                    t.DiscNumber, _ = strconv.Atoi(diskFound2[0])
                } else if strings.Contains(strings.ToLower(t.AlbumOriginal), "disk one") {
                    t.DiscNumber = 1
                } else if strings.Contains(strings.ToLower(t.AlbumOriginal), "disk two") {
                    t.DiscNumber = 2
                } else if strings.Contains(strings.ToLower(t.AlbumOriginal), "disk three") {
                    t.DiscNumber = 3
                } else if strings.Contains(strings.ToLower(t.AlbumOriginal), "disk a") {
                    t.DiscNumber = 1
                } else if strings.Contains(strings.ToLower(t.AlbumOriginal), "disk b") {
                    t.DiscNumber = 2
                }
            }

            curAlbum := &Album{}
            // for the Album key, should likely create a "sanitized" version (all lowercase would be a start, removing deluxe, etc)
            albumKey := fmt.Sprintf("%s##%s", artistKey, strings.ToLower(AlbumName))
            if _, isMapContainsKey := lib.Albums[albumKey]; isMapContainsKey {
                curAlbum = lib.Albums[albumKey]
            } else {
                curAlbum.Name = AlbumName
            }
            if curAlbum.SortName == AlbumName && AlbumSort != AlbumArtist {
                curAlbum.SortName = AlbumSort
            }
            if curAlbum.SortArtist == AlbumName && AlbumSort != AlbumArtist {
                curAlbum.SortArtist = AlbumSort
            }
            if curAlbum.Artist == "" && AlbumArtist != "" {
                curAlbum.Artist = AlbumArtist
            }
            if curAlbum.SortArtist == AlbumArtist && AlbumArtistSort != AlbumArtist {
                curAlbum.SortArtist = AlbumArtistSort
            }
            
            if curAlbum.Name == "Classical Marches" {
                fmt.Println("CM: ", t.Album, t.Artist, curAlbum.Artist, AlbumName)
            }

            // Should we collect Song Artist / Composer info at album level?
            
            if curAlbum.Year == 0 || curAlbum.Year > t.Year {
                curAlbum.Year = t.Year
            }
            //itunes/Library.go:246:49: invalid operation: t.ReleaseDate < curAlbum.ReleaseDate (operator < not defined on struct)
            //if curAlbum.ReleaseDate.IsZero() || t.ReleaseDate < curAlbum.ReleaseDate {
            if curAlbum.ReleaseDate.IsZero() {
                curAlbum.ReleaseDate = t.ReleaseDate
            }
            if curAlbum.TrackCount < t.TrackCount {
                curAlbum.TrackCount = t.TrackCount
            }

            curAlbum.TracksHave ++
            //TracksCount we don't know yet
            
            curAlbum.TotalTime += t.TotalTime
            if curAlbum.DateModified.IsZero() || curAlbum.DateModified.Before(t.DateModified) {
                curAlbum.DateModified = t.DateModified
            }
            if curAlbum.DateAdded.IsZero() || curAlbum.DateAdded.After(t.DateAdded) {
                curAlbum.DateAdded = t.DateAdded
            }
            if t.PlayCount > curAlbum.PlayCount {
                curAlbum.PlayCount = t.PlayCount
            }
            curAlbum.PlayCountAll += t.PlayCount
            if curAlbum.PlayDate < t.PlayDate {
                curAlbum.PlayDate = t.PlayDate
            }
            if curAlbum.PlayDateUTC.IsZero() || curAlbum.PlayDateUTC.Before(t.PlayDateUTC) {
                curAlbum.PlayDateUTC = t.PlayDateUTC
            }
            
            // need to rethink all the rating stuff, how to aggregate songs / album ratings up through the tree
            if curAlbum.Rating < t.Rating {
                // this should really be a combination of all Ratings / Loves, Likes and Dislikes
                curAlbum.Rating = t.Rating
            }
            if t.Rating > 0 {
                curAlbum.SongRatings = append(curAlbum.SongRatings, t.Rating)
            }
            if t.AlbumRating > curAlbum.AlbumRating {
                curAlbum.AlbumRating = t.AlbumRating
            }
            if t.AlbumRatingComputed {
                curAlbum.AlbumRatingComputed ++
            }
            
            // if doing multiple types, need to adjust this to put it into the correct bucket
            // likely should be a map of maps
            curAlbum.Songs = append(curAlbum.Songs, t)
            if t.Genre != "" {
                if strings.Contains(t.Genre, ",") {
                    genres := strings.Split(t.Genre, ",")
                    for _, genre := range genres {
                        genre = strings.TrimSpace(genre)
                        if _, containsGenre := curAlbum.Genres[genre]; containsGenre {
                            curAlbum.Genres[genre] = true
                        }
                    }
                } else if _, containsGenre := curAlbum.Genres[strings.TrimSpace(t.Genre)]; containsGenre {
                    curAlbum.Genres[strings.TrimSpace(t.Genre)] = true
                }
            }
            if t.ArtworkCount > curAlbum.ArtworkCount {
                curAlbum.ArtworkCount = t.ArtworkCount
            }
            if t.Compilation {
                curAlbum.Compilation = true
            }

            // increment counts
            // TODO this should be a map most like (Counts := map[string][int])
            if t.Movie {
                curAlbum.CountMovies ++
            } else if t.TVShow {
                curAlbum.CountTVShows ++
            } else if t.Podcast {
                curAlbum.CountPodcasts ++
            } else if t.MusicVideo {
                curAlbum.CountMusicVideos ++
            } else {
                curAlbum.CountSongs ++
            }
            curAlbum.CountTracks ++
            
            if t.Matched {
                curAlbum.Matched ++
            }
            if t.Unplayed {
                curAlbum.Unplayed ++
            }
            if t.Purchased {
                curAlbum.Purchased ++
            }
            if t.Explicit {
                curAlbum.Explicit ++
            }
            if t.Clean {
                curAlbum.Clean ++
            }
            
            if t.Loved {
                curAlbum.Loved ++
            }
            if t.Liked {
                curAlbum.Liked ++
            }
            if t.Disliked {
                curAlbum.Disliked ++
            }
            if t.AlbumLoved {
                curAlbum.AlbumLoved ++
            }

            if t.Protected {
                curAlbum.Protected ++
            }
            if t.PartOfGaplessAlbum {
                curAlbum.GaplessAlbum = true
            }
            

            // save Album to Library
            
            if len(lib.Albums) == 0 {
                newAlbums := make(map[string]*Album)
                newAlbums[albumKey] = curAlbum
                lib.Albums = newAlbums
            } else {
                lib.Albums[albumKey] = curAlbum
            }
            
            // save Album to Artists
            
            if len(curArtist.Albums) == 0 {
                newAlbums := make(map[string]*Album)
                newAlbums[albumKey] = curAlbum
                curArtist.Albums = newAlbums
            } else {
                curArtist.Albums[albumKey] = curAlbum
            }

            // done with album    
            
            // save Artist to Library
            
            if len(lib.Artists) == 0 {
                newArtists := make(map[string]*Artist)
                newArtists[artistKey] = curArtist
                lib.Artists = newArtists
            } else {
                lib.Artists[artistKey] = curArtist
            }

        }
        
        byID[t.TrackID] = t
        
        if trackKind == "movie" {
            Movies = append(Movies, t)
        } else if trackKind == "tvshow" {
            TVShowsByID = append(TVShowsByID, t.TrackID)
        } else if trackKind == "podcast" {
            PodcastsByID = append(PodcastsByID, t.TrackID)
        } else if trackKind == "musicvideo" {
            MusicVideosByID = append(MusicVideosByID, t.TrackID)
        } else if trackKind == "music"{
            Music = append(Music, t)
        }

    }
    

    lib.tracksByID = byID
    
    lib.TVShows = TVShowsByID
    lib.Podcasts = PodcastsByID
    //lib.Movies = MoviesByID
    lib.Movies = Movies
    lib.MusicVideos = MusicVideosByID
    //lib.Music = MusicByID
    lib.Music = Music

}

// mapAllPlaylistTracks updates the tracks in each playlist to ensure
// that they map to the valid track struct in the Library.Tracks slice
func (lib *Library) mapAllPlaylistTracks() {

    for _, p := range lib.Playlists {

        for i, t := range p.PlaylistItems {

            p.PlaylistItems[i] = lib.tracksByID[t.TrackID]

        }

    }

}

// UnmarshalXML is a custom unmarshaller function for itunes library xml format
func (lib *Library) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

    linenum := 0
    for {
        token, err := decoder.Token()
        if err == io.EOF {
            lib.sortTracksByID()
            lib.mapAllPlaylistTracks()

            fmt.Println(lib.GetNumTracksX())

            return nil
        }
        if err != nil {
            return err
        }
        linenum ++
        //fmt.Printf("Library/UnmarshalXML: %i", linenum)

        switch t := token.(type) {

        case xml.StartElement:

            if t.Name.Local == "key" {

                key, err := readOpenTagAsString(decoder)
                if err != nil {
                    return err
                }

                err = resolveKeyOnStruct(lib, key, decoder)
                if err != nil {
                    return err
                }

            } else {

                fmt.Printf("skip %s (Library)\n", t.Name.Local)
                err = decoder.Skip()
                if err != nil {
                    return err
                }

            }

        }

    }

}
