package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/carlosokumu/dubbedapi/controllers"
	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("Fetching data  from api")

	//Initialize the database
	//path := "root:<abc>@tcp(localhost:3306)/test"
	//path := "root:abc@tcp(127.0.0.1:3306)/test"
	dsn := "host=localhost user=postgres password=Agent047 dbname=postgres port=2768 sslmode=disable TimeZone=Asia/Shanghai"
	database.Connect(dsn)
	database.Migrate()

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
	router := initRouter()
	router.Run(":8080")

	fmt.Printf("API Response as struct %+v\n", responseObject)

}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/tradex")
	{

		api.POST("/user/register", controllers.RegisterUser)

	}
	return router
}
