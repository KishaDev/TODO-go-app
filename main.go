package main

import (
	"encoding/json" // package to encode and decode JSON
	"log"           // package for logging errors
	"net/http"      // package for handling http requests and responses

	// package to convert string to int

	"github.com/gin-gonic/gin"                // package for creating RESTful APIs
	"github.com/gorilla/mux"                  // package for handling router and requests
	"github.com/jinzhu/gorm"                  // package for interacting with databases
	_ "github.com/jinzhu/gorm/dialects/mysql" // MySQL dialect for GORM
)

var db *gorm.DB
var err error

type Todo struct {
	gorm.Model        // This embeds the default fields (ID, CreatedAt, UpdatedAt, DeletedAt) in the struct
	ID         int    `json:"id"`        // the ID generated - has to be incremental
	Action     string `json:"action"`    // The title of the task
	Completed  bool   `json:"completed"` // Whether the task is completed or not
}

// Get all todos
//
//	func getTodos(c *gin.Context) {
//		var todos []Todo
//
// db.Find(&todos)
// c.JSON(http.StatusOK, todos)
// }
func getTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	db.Find(&todos)
	json.NewEncoder(w).Encode(todos)
}

// Get a todo by ID
func getTodo(c *gin.Context) {
	log.Println("Received GET request for todos")
	var todo Todo
	id := c.Params.ByName("id")
	db.First(&todo, id)
	if todo.ID == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		log.Printf("Todo with ID %s not found", id)
	} else {
		c.JSON(http.StatusOK, todo)
	}
}

// Create a new todo
//
//	func createTodo(c *gin.Context) {
//		var todo Todo
//		err := json.NewDecoder(c.Request.Body).Decode(&todo)
//		if err != nil {
//			log.Printf("Error decoding request body: %v", err)
//			c.AbortWithStatus(http.StatusBadRequest)
//			return
//		}
//		db.Create(&todo)
//		c.JSON(http.StatusOK, todo)
//	}
func createTodo(w http.ResponseWriter, r *http.Request) {
	var create_todo Todo
	err := json.NewDecoder(r.Body).Decode(&create_todo) // decode the request body into a todo struct
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db.Create(&create_todo)                //create the todo record in the database
	json.NewEncoder(w).Encode(create_todo) //encode the newly created todo and send it back as the respose
}

// Update a todo by ID
func updateTodo(c *gin.Context) {
	var todo Todo
	id := c.Params.ByName("id")
	db.First(&todo, id)
	if todo.ID == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		log.Printf("Todo with ID %s not found", id)
		return
	}
	err := json.NewDecoder(c.Request.Body).Decode(&todo)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.Save(&todo)
	c.JSON(http.StatusOK, todo)
}

// Delete a todo by ID
func deleteTodo(c *gin.Context) {
	var todo Todo
	id := c.Params.ByName("id")
	db.First(&todo, id)
	if todo.ID == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		log.Printf("Todo with ID %s not found", id)
		return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, todo)
}

func main() {
	//initialize database connection
	log.Println(("Starting server"))
	//TODO
	db, err = gorm.Open("mysql", "DB:PASSWORD/todo_app")
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Server open..")
	defer db.Close()
	log.Println("Server serving...")
	//migrate the schema
	db.AutoMigrate(&Todo{})

	//create a new router using the Gorilla mux router
	router := mux.NewRouter()
	log.Println("Mux router initialized")

	//define endpoints
	router.HandleFunc("/todos", getTodos).Methods("GET")
	log.Println("Got my Todos..")
	router.HandleFunc("/todos", createTodo).Methods("POST")
	log.Println("Posting my Todos..")
	//router.GET("/todos", getTodos)
	//router.GET("/todos/:id", getTodo)
	//router.POST("/todos", createTodo)
	//router.PUT("/todos/:id", updateTodo)
	//router.DELETE("/todos/:id", deleteTodo)
	///run the server
	//router.Run(":8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}