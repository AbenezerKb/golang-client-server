package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type C2BPaymentQueryResult struct {
	Password      string `xml:"Password"`
	Xmlstring     string `xml:"Xmlstring"`
	Text          string `xml:",chardata"`
	ResultCode    string `xml:"ResultCode"`
	ResultDesc    string `xml:"ResultDesc"`
	TransID       string `xml:"TransID"`
	BillRefNumber string `xml:"BillRefNumber"`
	UtilityName   string `xml:"UtilityName"`
	CustomerName  string `xml:"CustomerName"`
	Amount        string `xml:"Amount"`
}

type XML1 struct {
	XMLName xml.Name              `xml:"xml"`
	Text    string                `xml:",chardata"`
	C2      C2BPaymentQueryResult `xml:"C2BPaymentQueryResult"`
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	location := XML1{}

	jsn, err := ioutil.ReadAll(r.Body)

	//checking if the server is able to read the data and if not send bad request 400
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("Error reading the body", err)
	}

	// unmarshaling the JSON data
	err = json.Unmarshal(jsn, &location)

	//checking if unmarshaling has error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("Decoding error: ", err)
	}

	//printing the recieved password and xmlstring
	fmt.Printf("Password: %v\n", location.C2.Password)
	fmt.Printf("xmlstring: %v\n", location.C2.Xmlstring)

	password := "1234567"
	hash := md5.Sum([]byte(password))

	//check for the password equality
	if hex.EncodeToString(hash[:]) == location.C2.Password {

		//create and write the  xml  to file
		os.Create("C:\\Users\\Administrator\\Desktop\\Web\\code\\third\\success.xml")
		writer, err := os.OpenFile("C:\\Users\\Administrator\\Desktop\\Web\\code\\third\\success.xml", os.O_RDWR|os.O_APPEND, 0660)

		//if saving to file failed it will send status of 500
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("error opening file")
			return
		}
		encoder := xml.NewEncoder(writer)
		er := encoder.Encode(location)

		//if saving to file failed it will send status of 500
		if er != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(er)
			return
		}
		fmt.Println("password equality success, success.xml file created and payload data written")
	} else {
		//create and write the  xml  to file
		os.Create("C:\\Users\\Administrator\\Desktop\\Web\\code\\third\\failed.xml")
		writer, err := os.OpenFile("C:\\Users\\Administrator\\Desktop\\Web\\code\\third\\failed.xml", os.O_RDWR|os.O_APPEND, 0660)

		//if saving to file failed it will send status of 500
		if err != nil {
			fmt.Println("error opening file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		encoder := xml.NewEncoder(writer)
		er := encoder.Encode(location)

		//if saving to file failed it will send status of 500
		if er != nil {
			fmt.Println(er)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("password equality failed, failed.xml file created and payload data written")
	}

}
func server() {

	http.HandleFunc("/", dataHandler)
	http.ListenAndServe(":8088", nil)
}

func client() {

	//the password and xmlstring variables
	password := "1234567"
	xmlstring := "the xml string"

	// encrypting the password
	hash := md5.Sum([]byte(password))
	hex.EncodeToString(hash[:])

	// initailizing C2BPayment with the password and xmlstring variables xml data
	cc := C2BPaymentQueryResult{
		Xmlstring: xmlstring,
		Password:  hex.EncodeToString(hash[:]),
	}

	// initializing the xml with C2BPayment
	xm := XML1{
		C2: cc,
	}

	//marshaling into json
	locJson, err := json.Marshal(xm)

	//sending the json xml data
	req, err := http.NewRequest("POST", "http://localhost:8088", bytes.NewBuffer(locJson))
	req.Header.Set("Content-Type", "application/json")

	//setting time out for the client
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	//check if the status code of the response is different than 200 and print with error
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unable to send the data, Status code: %v", resp.StatusCode)
		return
	}

	//print the status code
	fmt.Println("Response: ", resp.Status)
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("Response: ", string(body))
	resp.Body.Close()
}

func main() {

	//Run the server in different go routine
	go server()

	//Run the client in the main

	client()

}
