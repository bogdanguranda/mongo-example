package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	log "github.com/Sirupsen/logrus"
)

type Genre string

const (
	ACTION Genre = "Action"
	HORROR Genre = "Horror"
	SCI_FI Genre = "Science Fiction"
)

type Movie struct {
	Name        string
	Genre       Genre
	ReleaseYear int
	Rating      float32
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Panicf("Could not connnect to MongoDB! Error was: %v", err.Error())
		return
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	initializeDB(session)
	collection := session.DB("test").C("movies")

	var result []Movie

	err = collection.Find(bson.M{"$and": []interface{}{bson.M{"genre": SCI_FI}, bson.M{"rating": bson.M{"$gt": 8.5}}}}).All(&result)
	if err != nil {
		log.Warnf("Error when querying MongoDB. Error was: %v", err.Error())
		return
	}

	for _, movie := range result {
		log.Infof("Found movie with name: %v, genre: %v, release year: %v, rating: %v", movie.Name, movie.Genre, movie.ReleaseYear, movie.Rating)
	}
}

func initializeDB(session *mgo.Session) {
	c := session.DB("test").C("movies")
	c.DropCollection()
	c.Insert(&Movie{Name: "Star Wars: Episode IV - A New Hope", Genre: SCI_FI, ReleaseYear: 1977, Rating: 8.7},
		&Movie{Name: "Star Wars: Episode V - Empire Strikes Back", Genre: SCI_FI, ReleaseYear: 1980, Rating: 8.8},
		&Movie{Name: "Star Wars: Episode VI - Return of the Jedi", Genre: SCI_FI, ReleaseYear: 1983, Rating: 8.4},
		&Movie{Name: "Ip Man 2", Genre: ACTION, ReleaseYear: 2010, Rating: 7.6},
		&Movie{Name: "The Ring", Genre: HORROR, ReleaseYear: 2002, Rating: 7.1})
}