package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	guuid "github.com/google/uuid"
)

type Todo struct {
	ID        string    `json:"id"`
	ProductId string	`json:"p_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Completed string    `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID          string `json:"id"`
	//TodoId     string `bson:"podcast,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    string  `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}

// DATABASE INSTANCE
var collection,proCollection *mongo.Collection


func TodoCollection(c *mongo.Database) {
	collection = c.Collection("todo")
}

func ProductCollection(c *mongo.Database){
	proCollection = c.Collection("pingpong")
}

func GetAllProducts(c *gin.Context) {
	pro := []Product{}
	cursor, err := proCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	// Iterate through the returned cursor.
    for cursor.Next(context.TODO()) {
				var prod Product
        cursor.Decode(&prod)
        pro = append(pro, prod)
		}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Products",
		"data":    pro,
	})
	return
}

func GetAllProductTodo(c *gin.Context) {
	id := c.Param("proId")
	pro := []Product{}

	matchStage := bson.D{{"$match", bson.D{{"p_id", id}}}}
	groupStage := bson.D{{"$group", bson.D{{"id", "$product" }}}}
	cursor, err := proCollection.Find(context.TODO(), bson.M{})

	showInfoCursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		panic(err)
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(context.TODO(), &showsWithInfo); err != nil {
		panic(err)
	}
	fmt.Println(showsWithInfo)

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	// Iterate through the returned cursor.
	
    for cursor.Next(context.TODO()) {
				var prod Product
        cursor.Decode(&prod)
        pro = append(pro, prod)
		}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Products",
		"data":    pro,
	})
	return
}

func CreateProd(c *gin.Context) {
	var pro Product
	c.BindJSON(&pro)
	fmt.Println("pro : ",pro)
	title := pro.Title
	body := pro.Description
	duration := pro.Duration
	id := guuid.New().String()
	fmt.Println("Id : ",id)
	fmt.Println("Desc : ",body)


	newTodo := Product{
		ID: id,
		Title:     title,
		Description:      body,
		Duration: duration,
		CreatedAt: time.Now(),
	}

	_, err := proCollection.InsertOne(context.TODO(), newTodo)

	if err != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Product created Successfully",
	})
	return
}


func GetAllTodos(c *gin.Context) {
	todos := []Todo{}
	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	// Iterate through the returned cursor.
    for cursor.Next(context.TODO()) {
				var todo Todo
        cursor.Decode(&todo)
        todos = append(todos, todo)
		}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Todos",
		"data":    todos,
	})
	return
}

func CreateTodo(c *gin.Context) {
	var todo Todo
	c.BindJSON(&todo)
	fmt.Println("todo : ",todo)
	title := todo.Title
	body := todo.Body
	proId := todo.ProductId
	completed := todo.Completed
	id := guuid.New().String()
	fmt.Println("Id : ",id)


	newTodo := Todo{
		ID: id,
		ProductId : proId,
		Title:     title,
		Body:      body,
		Completed: completed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), newTodo)

	if err != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Todo created Successfully",
	})
	return
}

func GetSingleTodo(c *gin.Context) {
	todoId := c.Param("todoId")

	todo := Todo{}
	err := collection.FindOne(context.TODO(), bson.M{"id": todoId}).Decode(&todo)
	if err != nil {
			log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Todo",
		"data": todo,
	})
	return
}

func EditTodo(c *gin.Context) {
	todoId := c.Param("todoId")
	var todo Todo
	c.BindJSON(&todo)
	completed := todo.Completed

	newData := bson.M{
            "$set": bson.M{
            "completed":       completed,
            "updated_at": time.Now(),
            },
        }

	_, err := collection.UpdateOne(context.TODO(), bson.M{"id": todoId}, newData)
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"message":  "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Todo Edited Successfully",
	})
	return
}

func DeleteTodo(c *gin.Context) {
todoId := c.Param("todoId")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": todoId})
	if err != nil {
		log.Printf("Error while deleting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully",
	})
	return
}