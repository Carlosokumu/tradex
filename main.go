package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/carlosokumu/dubbedapi/models"
)

func main() {

	fmt.Println("Fetching data  from api")

	//Initialize the database
	//path := "root:<abc>@tcp(localhost:3306)/test"
	// path := "root:abc@tcp(127.0.0.1:3306)/test"
	// database.Connect(path)
	// //database.Connect(path)
	// database.Migrate()

	//Create a client
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.spotware.com/connect/tradingaccounts?access_token=Ct2rFyZKl7-tSXWgkXJxrScJYdMTR-sdrVc9AGDoTzw", nil)

	if err != nil {
		fmt.Println(err)
	}

	//Set  headers to the requests
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	//Use the client to make the requests with the given [configurations]
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject models.Response

	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Println(responseObject.Data[0].AccountID)
	// fmt.Println(responseObject.AccountNumber)

	fmt.Printf("API Response as struct %+v\n", responseObject)
}
