package main

import (
	"fmt"
	"io"
	//"bytes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"flag"
)

var dsnFront = "root:dbpassword@tcp("
var dsnBack = ":3306)/aq?parseTime=true"
var db *sql.DB

var debug bool

type Reading struct {
    	
	Ts string
	Lat float32
	Lon float32
	Co float32
	Pm float32
	Hum float32
	Temp float32
}

func upload(w http.ResponseWriter, r *http.Request){
	text := r.URL.Query()["text"][0]
	if debug {fmt.Println("add request for text="+text)}
	
	res, err:= db.Exec("INSERT INTO test VALUES('"+text+"')")
	if err != nil{
		log.Print(err)
	} else {
		rows, _ := res.RowsAffected()
		result := text + " entered into DB! " + strconv.FormatInt(rows, 10) + " rows affected."		
		io.WriteString(w, result)
		fmt.Println(result)		
	}
}

func jsonTest(w http.ResponseWriter, r *http.Request){
	//test json decoding and encoding

	// initialize a Reading to hold the values
	var reading Reading
	// STUB OUT a new json upload
	b := []byte(`{"lat":30.1234,"lon":70.4321, "co":0.1, "pm":0.17, "hum":0.8, "temp":30}`)
	// decode the json string into the reading object
	json.Unmarshal(b, &reading)
	// print the contents of reading to make sure they came through
	fmt.Fprint(w, reading)

}

/*
* The server attaches handlers and listens for REST requests on port 80
 */
func main() {		
	var err error
	//Grab server addresses from command line args	
	// Usage: -debug=true or -debug=false
	debugPtr := flag.Bool("debug", false, "Turn console output on or off")
	
    flag.Parse()

    debug = *debugPtr    //-debug=?
        
	dsn := dsnFront + flag.Args()[0] + dsnBack			
	//Open an MySQL connection to the RDS
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err) //Couldn't open the database
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err) //Couldn't ping the database
	} else {
		fmt.Println("DB open!") //Ok, the database is open
	}	
	
	//Attach handlers

	http.HandleFunc("/json", jsonTest)
	http.HandleFunc("/upload", upload)
	port := ":8000"
	fmt.Println("Datasink starting... on port", port)
	log.Fatal(http.ListenAndServe(port, nil))


}