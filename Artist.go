package itunes

import (
	"fmt"
	"time"
)

// Artist represents an Artist (not a real iTunes object)
type Artist struct {
	Name                string	// AlbumArtist (or Artist)
	YearFirst           int
	YearLast            int
	TotalTime           int
	DateModified        time.Time
	DateAdded           time.Time
	Rating              int
	AlbumRatings        []int
	Albums              map[string]*Album
	//PodCasts            []*Track  // likely should be a map of maps, rather than individuals
	//TVShows             []*Track
	//Movies              []*Track
	//MusicVideos         []*Track
	Songs                 []*Track
	//Tracks              []*Track // not sure about this
	Genres              map[string]bool
	PlayCount           int
	PlayCountAll        int
	PlayDate            int
	PlayDateUTC         time.Time
	CountPodcasts       int  // counts should be a map?
	CountMusicVideos    int
	CountMovies         int
	CountTVShows        int
	CountSongs          int
	CountTracks         int
	SortName            string
	Loved               int
	Liked               int
	Disliked            int
	AlbumLoved          int
}

func (t *Artist) String() string {

	return fmt.Sprintf("Artist: %s [%v albums, %v songs, %v podcasts, %v music videos, %v movies, %v tv shows]", t.Name, len(t.Albums), t.CountSongs, t.CountPodcasts, t.CountMusicVideos, t.CountMovies, t.CountTVShows)

}