package itunes

import (
	"fmt"
	"time"
)

// Album represents an Album / group of tracks (not a real iTunes object)
type Album struct {
	Name                string // album name
	Artist              string // Album Artist
	Genres              map[string]bool
	Year                int
	ReleaseDate         time.Time  // should there be a RDFirst and RDLast?
	DateModified        time.Time
	TotalTime           int
	TracksHave          int
	TracksCount         int
	DateAdded           time.Time
	PlayCount           int
	PlayCountAll        int
	PlayDate            int
	PlayDateUTC         time.Time
	Rating              int
	SongRatings         []int
	AlbumRating         int
	AlbumRatingComputed int
	//PodCasts            []*Track // likely should be a map of maps, rather than individuals
	//TVShows             []*Track
	//Movies              []*Track
	//MusicVideos         []*Track
	Songs                 []*Track
	//Tracks              []*Track // not sure about this
	ArtworkCount        int
	Compilation         bool
	CountPodcasts       int  // Counts should be a map?
	CountMusicVideos    int
	CountMovies         int
	CountTVShows        int
	CountSongs          int
	CountTracks         int
	Matched             int    // possibly related to iTunes Match?
	//Podcast             bool  // these can be done via COunts
	//MusicVideo          bool
	//Movie               bool
	//TVShow              bool
	Unplayed            int
	SortName            string
	SortArtist          string
	Purchased           int
	Explicit            int
	Clean               int
	Loved               int
	Liked               int
	Disliked            int
	AlbumLoved          int
	Protected           int
	GaplessAlbum        bool
	//Season              string  // not applicablele to song albums
	//ContentRating       map[string]int // not applicablele to song albums
}

func (t *Album) String() string {

	//return fmt.Sprintf("Album: %s by %s [%v songs, %v podcasts, %v music videos, %v movies, %v tv shows]", t.Name, t.Artist, t.CountSongs, t.CountPodcasts, t.CountMusicVideos, t.CountMovies, t.CountTVShows)
	return fmt.Sprintf("Artist: %s Album: %s [%v songs, %v podcasts, %v music videos, %v movies, %v tv shows]", t.Artist, t.Name, t.CountSongs, t.CountPodcasts, t.CountMusicVideos, t.CountMovies, t.CountTVShows)

}