package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Movie struct {
	Id          string
	Title       string
	Year        uint
	Rating      float32
	Runtime     uint
	Genre       []string
	Director    string
	Stars       []string
	VoteCount   uint
	ReleaseDate string
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {

	movie, err := os.ReadFile("./input.json")
	checkErr(err)

	movieList := []Movie{}
	err = json.Unmarshal(movie, &movieList)
	checkErr(err)

	//sortByYear(&movieList)
	//sortByRating(&movieList)
	//filteredMov := filterByYear(1999, movieList)
	filteredMov := filterByYears(1999, 2014, movieList)
	fmt.Println(filteredMov)

}

// SORTS

// sortByYear sorts the movielist in acending order according to year
func sortByYear(store *[]Movie) {
	sort.Slice(*store, func(i, j int) bool {
		if (*store)[i].Year < (*store)[j].Year {
			return true
		}
		return false
	})
}

// sortByRating sorts the movielist in acending order according to rating
func sortByRating(store *[]Movie) {
	sort.Slice(*store, func(i int, j int) bool {
		if (*store)[i].Rating < (*store)[j].Rating {
			return true
		}
		return false
	})
}

// sortByRating sorts the movielist in acending order according to rating
func sortByVotes(store *[]Movie) {
	sort.Slice(*store, func(i int, j int) bool {
		if float32((*store)[i].VoteCount) < float32((*store)[j].VoteCount) {
			return true
		}
		return false
	})
}

// sortByRating sorts the movielist in acending order according to rating
func sortByTitle(store *[]Movie) {
	sort.Slice(*store, func(i int, j int) bool {
		if (*store)[i].Title < (*store)[j].Title {
			return true
		}
		return false
	})
}

// Filters

// filterByYear will filter all the movies according to year passed in
func filterByYear(year uint, store []Movie) []Movie {
	filteredMovies := []Movie{}
	for _, mov := range store {
		if mov.Year == year {
			filteredMovies = append(filteredMovies, mov)
		}
	}
	return filteredMovies
}

// filterByYear will filter all the movies according to year passed in
func filterByYears(yearStart uint, yearEnd uint, store []Movie) []Movie {
	filteredMovies := []Movie{}
	for _, mov := range store {
		if mov.Year <= yearEnd && mov.Year >= yearStart {
			filteredMovies = append(filteredMovies, mov)
		}
	}
	return filteredMovies
}

func filterByRating(rating float32, store []Movie) []Movie {
	filteredList := []Movie{}
	for _, mov := range store {
		if mov.Rating >= rating {
			filteredList = append(filteredList, mov)
		}
	}
	return filteredList
}

func filterByDirector(director string, store []Movie) []Movie {
	filteredList := []Movie{}
	for _, mov := range store {
		if strings.ToLower(mov.Director) == strings.ToLower(director) {
			filteredList = append(filteredList, mov)
		}
	}
	return filteredList
}

// velvet := &Movie{
// 	Id:    "123112",
// 	Title: "Blue Velvet",
// 	Year:  "1986",
// }
//
// encodedMovies, err := json.Marshal(velvet)
// checkErr(err)
//
//fmt.Println(string(encodedMovies))
