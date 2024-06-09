package itunes

import (
	"fmt"
	"time"
	"strings"
	"regexp"
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

	return fmt.Sprintf("Artist: %s (%s; a: %v; m: %v) [%v albums, %v songs, %v podcasts, %v music videos, %v movies, %v tv shows]", t.Name, t.SortName, t.DateAdded, t.DateAdded, len(t.Albums), t.CountSongs, t.CountPodcasts, t.CountMusicVideos, t.CountMovies, t.CountTVShows)

}

func cleanArtist(name string) string {
	newName := name
	if strings.Contains(name, " Feat. ") {
		re := regexp.MustCompile(` Feat. .*$`)
		newName = re.ReplaceAllString(name, "")
	}

	if strings.ToLower(newName) == "thirty seconds to mars" {
		newName = "30 Seconds To Mars"
	} else if strings.ToLower(newName) == "pink" {
		newName = "P!nk"
	} else if strings.ToLower(newName) == "30 seconds to mars" && newName != "30 Seconds To Mars" {
		newName = "30 Seconds To Mars"
	} else if newName == "Andrea McArdle & Ensemble" || newName == "Andrea McArdle, Orphans" {
		newName = "Andrea McArdle"
	} else if newName == "Biggest Loser Workout Mix" {
		newName = "Biggest Loser"
	} else if newName == "Debbie Byrne (Fantine) / Garry Morris (Valjean)" {
		newName = "Debbie Byrne (Fantine)"
	} else if newName == "Florence & The Machine" {
		newName = "Florence + The Machine"
	} else if newName == "Jay-Z & Linkin Park" {
		newName = "Linkin Park & Jay-Z"
	} else if newName == "Lost Prophets" {
		newName = "Lostprophets"
	} else if newName == "Me First and The Gimme Gimmes" {
		newName = "Me First & The Gimme Gimmes"
	} else if newName == "MTV's 120 Minutes" {
		newName = "MTV"
	} else if newName == "Puddle of Mud" {
		newName = "Puddle Of Mudd"
	} else if newName == "Queensryche" {
		newName = "Queensrÿche"
	} else if newName == "various" || newName == "Various Artist" {
		newName = "Various Artists"
	} else if newName == "Vega, Suzanne" {
		newName = "Suzanne Vega"
	} else if newName == "The Arcade Fire" {
		newName = "Arcade Fire"
	} else if newName == "Art Of Noise" {
		newName = "The Art of Noise"
	} else if newName == "Black Eyed Peas" {
		newName = "The Black Eyed Peas"
	} else if newName == "Cars" {
		newName = "The Cars"
	} else if newName == "Chemical Brothers" {
		newName = "The Chemical Brothers"
	} else if newName == "Corrs" {
		newName = "The Corrs"
	} else if newName == "Fray" {
		newName = "The Fray"
	} else if newName == "KLF" {
		newName = "The KLF"
	} else if newName == "Moody Blues" {
		newName = "The Moody Blues"
	} else if newName == "Pussycat Dolls" {
		newName = "The Pussycat Dolls"
	} else if newName == "Refreshments" {
		newName = "The Refreshments"
	} else if newName == "blink-182" {
		newName = "Blink-182"
	} else if newName == "Dada" {
		newName = "dada"
	}
	
	if false && name != newName {
		fmt.Println(fmt.Sprintf("Cleaning: %s (%s)", name, newName))
	}

	return newName
}