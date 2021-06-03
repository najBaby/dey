package service

import (
	"context"
	"log"

	"deyforyou/dey/schema"
)

// var expression = regexp.MustCompile(`(http(s|)://.*.mp4)`)

type MovieServiceServer struct {
	NewMovie     []*schema.Movie
	PopularMovie []*schema.Movie
	schema.UnimplementedMovieServiceServer
}

// NewMovieServiceServer is
func NewMovieServiceServer() *MovieServiceServer {
	return &MovieServiceServer{
		NewMovie:                        make([]*schema.Movie, 0),
		PopularMovie:                    make([]*schema.Movie, 0),
		UnimplementedMovieServiceServer: schema.UnimplementedMovieServiceServer{},
	}
}

// ListMovies is
func (ass *MovieServiceServer) ListMovies(
	context context.Context,
	request *schema.ListMoviesRequest,
) (*schema.ListMoviesResponse, error) {
	var movies []*schema.Movie

	streamComplet3 := NewStreamComplet3()
	if query := request.GetQuery(); query != "" {
		log.Println(query)
		movies = streamComplet3.Search(query)
	} else {
		movies = streamComplet3.NewsMovies(request.NextPage)
	}
	return &schema.ListMoviesResponse{Movies: movies}, nil
}

func (ass *MovieServiceServer) Movie(
	context context.Context,
	request *schema.MovieRequest,
) (*schema.MovieResponse, error) {

	streamComplet3 := NewStreamComplet3()
	movie := streamComplet3.Movie(request.GetSource())
	return &schema.MovieResponse{Movie: movie}, nil
}

type Movie interface {
	MovieFilm
	MovieSerie
}

type MovieFilm interface {
	Film()
	Films()
}

type MovieSerie interface {
	Serie()
	Series()
}
