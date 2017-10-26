package study

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}
func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 8, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

type ByArtist []*Track

func (x ByArtist) Len() int {
	return len(x)
}
func (x ByArtist) Less(i, j int) bool {
	return x[i].Artist < x[j].Artist
}
func (x ByArtist) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}
func SortByArtistTest() {
	tras := ByArtist(tracks)
	sort.Sort(tras)
	//sort.Sort(sort.Reverse(tras))//逆序
	printTracks(tras)
}

//==========综合排序========
type CustomSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (c CustomSort) Len() int {
	return len(c.t)
}
func (c CustomSort) Less(i, j int) bool {
	return c.less(c.t[i], c.t[j])
}
func (c CustomSort) Swap(i, j int) {
	c.t[i], c.t[j] = c.t[j], c.t[i]
}
func cusLess(x, y *Track) bool {
	if x.Title != y.Title {
		return x.Title < y.Title
	}
	if x.Year != y.Year {
		return x.Year < y.Year
	}
	if x.Length != y.Length {
		return x.Length < y.Length
	}
	return false
}
func SortCustomTest() {
	tras := CustomSort{tracks, cusLess}
	sort.Sort(tras)
	fmt.Println(sort.IsSorted(tras)) //判断是否排过序
	printTracks(tras.t)
}
