package produce_api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ProduceList struct {
	List []Produce `json:"Produce"`
}

// getaAllHandler returns the JSON of all produce items
func (store *Store) getAllHandler(c *gin.Context) {
	c.JSON(http.StatusOK, store.ProduceItems)
}

// getProduceHandler returns the JSON of one or more produce items based on URL paramaters
// expects a 'code' parameter containing the produce code
func (store *Store) getProduceHandler(c *gin.Context) {
	params := c.Request.URL.Query()

	// create a new slice to hold all of the returned produce items
	var foundProduce []Produce
	for _, code := range params["code"] {
		index, produceItem := store.FindProduce(code)
		if index >= 0 { // if no error returned, then assume the product is valid
			foundProduce = append(foundProduce, produceItem)
		}
	}

	if len(foundProduce) > 0 { // return produce data if any was found
		c.JSON(http.StatusOK, foundProduce)
	} else {
		// return that the request was processed, but no data was found
		c.JSON(http.StatusNoContent, gin.H{"Warning": "No produce items matched the provided produce code(s)."})
	}
}

// addProduceHandler handles the POST request from a client for adding a produce item to the internal DB
func (store *Store) addProduceHandler(c *gin.Context) {
	var list ProduceList
	c.BindJSON(&list) // bind JSON body to Produce struct
	fmt.Println(list)

	fmt.Println("list", list)
	for _, produce := range list.List {
		fmt.Println(produce)
		err := store.AddProduce(produce) // add the new produce to the db
		if err != nil {
			c.String(http.StatusBadRequest, "The item(s) have not been added, please ensure the format of the produce item is corrected.")
		}
	}

	c.String(http.StatusAccepted, "The item(s) have been added")
}

// deleteProduceHandler handles the DELETE request when a client requests to delete a produce item
func (store *Store) deleteProduceHandler(c *gin.Context) {
	params := c.Request.URL.Query()
	targetCode := params["Produce Code"]

	// if the internal slice is empty, no point in continuing
	if len(store.ProduceItems) == 0 {
		c.JSON(http.StatusNoContent, "There are currently no produce items.")
		return
	}

	// confirm the passed in code is of a valid format
	// make a temporary Produce struct to pass into IsValid
	err := IsValid(Produce{Name: "Test", ProduceCode: targetCode[0], Price: 1.00})
	if err != nil {
		c.JSON(http.StatusBadRequest, "the inputted codes is of an invalid format")
		return
	}

	preLength := len(store.ProduceItems)
	store.ProduceItems, _ = store.RemoveProduce(targetCode[0]) // delete one item

	// check if the internal slice was adjusted at all
	if preLength == len(store.ProduceItems) {
		c.JSON(http.StatusNoContent, "The item was not found, so no deletion occurred.")
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "Delete successful", "produceList": store.ProduceItems})
	}
}

func APIMain() {
	router := gin.Default()
	router.Use(cors.Default())

	store := CreateStore()         // create store struct model for use in the API
	store.PopulateDefaultProduce() // populate default produce items as specified in the specifications

	// API GET endpoints
	router.GET("/produce/getall", store.getAllHandler)
	router.GET("/produce/getitem", store.getProduceHandler)

	// API POST endpoints
	router.POST("/produce/add", store.addProduceHandler)
	router.DELETE("/produce/delete", store.deleteProduceHandler)

	router.Run()
}
