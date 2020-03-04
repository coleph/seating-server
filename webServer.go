package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var seatingJson = []byte{}
var People []Person
var nameSearched = ""

// for figuring out who has been searched for
var seatedAt = ""

// for figuring out who is sitting with the person who has been searched

var TableNumber = 1
var Place1 = 0
var Place2 = 9

type Person struct {
	Firstname string
	Lastname  string
	Table     string
	Table2nd  string
}

func handler(w http.ResponseWriter, r *http.Request) {

	var output []string

	// Let's print the info
	fmt.Println("Incoming Request: ")
	fmt.Println("Method: ", r.Method, " ", r.URL)

	name := strings.Split(r.URL.Path, "/")[2]
	//borrowed from Dylan's code- get a name from the URL
	name, _ = url.PathUnescape(name)

	nameSearched = name

	// Collect all the header keys
	headerKeys := make([]string, len(r.Header))
	i := 0
	for k := range r.Header {
		headerKeys[i] = k
		i++
	}
	// Show Client Headers
	for _, line := range headerKeys {
		fmt.Println("  > ", line, ":", r.Header.Get(line))
	}

	for index := 0; index < 290; index++ {
		if People[index].Firstname == nameSearched {
			People[index].Table = seatedAt
		}

	}
	for index := 0; index < 290; index++ {
		if People[index].Table == seatedAt {
			output = append(People[index].Firstname)

		}

	}
	// Answer the Client request
	seatingJson, _ = json.Marshal(output)
	tableJSON := "{\"Table\": \"" + seatedAt + "\", \"Names\": " + string(seatingJson) + "}"
	// borrowed from Dylan, puts people seated with you in the same JSON as your table number
	// Answer the Client request
	fmt.Fprintf(w, string(seatingJson))
	// not sure what w means
}

func main() {

	csvFile, _ := os.Open("list.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		People = append(People, Person{
			Firstname: line[0],
			Lastname:  line[1],
		})
	}

	// randimization credit to youbasic.org go tutorial
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(People), func(i, j int) { People[i], People[j] = People[j], People[i] })
	// assign 8 people at a time a table number 1-31, int is converted to a string for Table
	for count := 0; count <= 30; count++ {

		for index := Place1; index < Place2; index++ {
			TableNumberString := strconv.FormatInt(int64(TableNumber), 10)
			People[index].Table = TableNumberString
		}
		// change the position of people[] modified
		TableNumber = TableNumber + 1
		Place1 = Place1 + 8
		Place2 = Place2 + 8

	}
	for index := 248; index < 279; index++ {
		People[index].Table = "Waiter"
	}
	for index := 279; index < 290; index++ {
		People[index].Table = "KC"
	}

	fmt.Println(People)
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":80", nil))

}

//http://localhost/
