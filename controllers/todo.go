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
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	//guuid "github.com/google/uuid"
)

/*

response.Header().Set("Content-Type","application/json")
  var user User
  var dbUser User
  json.NewDecoder(request.Body).Decode(&user)
  collection:= client.Database("GODB").Collection("user")
  ctx,_ := context.WithTimeout(context.Background(),10*time.Second)
  err:= collection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&dbUser)

  if err!=nil{
	  response.WriteHeader(http.StatusInternalServerError)
	  response.Write([]byte(`{"message":"`+err.Error()+`"}`))
	  return
  }
  userPass:= []byte(user.Password)
  dbPass:= []byte(dbUser.Password)

  passErr:= bcrypt.CompareHashAndPassword(dbPass, userPass)

  if passErr != nil{
	  log.Println(passErr)
	  response.Write([]byte(`{"response":"Wrong Password!"}`))
	  return
  }

*/

type GetDeviceUser struct{
	ID		primitive.ObjectID 	`bson:"id"`
}

type User struct {
	ID 		primitive.ObjectID `bson:"_id,omitempty"`
	Username	string 	`json:"username"`
	Password 	string 	`json:"password"`
	Address		string	`json:"address"`
	Email		string	`json:"email"`
	Phone_number	string	`json:"phone_number"`
	Is_Active	bool	`json:"is_active"`
}

type Device struct {
	//ID        string    `json:"id"`
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	User	primitive.ObjectID `bson:"user,omitempty"`
	ProductId string	`json:"productid"`
	DeviceName     string    `json:"deviceName"`
	//Owner      string    `json:"owner"`
	Brand string    `json:"brand"`
	Warranty	int		`json:"Warranty"`
	ManufactureDate time.Time `json:"manufactureDate"`
	PurchaseDate time.Time `json:"purchaseDate"`
	Is_Active	bool 	`json:"is_active"`
}

type Service struct {
	//ID          string `json:"id"`
	//TodoId     string `bson:"podcast,omitempty"`
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Device		primitive.ObjectID `bson:"device,omitempty"`
	User		primitive.ObjectID `bson:"user,omitempty"`
	ProductId	string		`json:"productid"`
	DeviceName       string  `json:"deviceName"`
	Issue string  `json:"issue"`
	ReceivedDate    time.Time  `json:"receivedDate"`
	ReturnedDate time.Time `json:"returnedDate"`
	Is_Active	bool 	`json:"is_active"`
	IS_Closed	bool	`json:"is_closed"`
}

// DATABASE INSTANCE
var deviceCollections,serviceCollection,userCollection *mongo.Collection

func UserCollection(c *mongo.Database){
	userCollection = c.Collection("user")
}

func DeviceCollection(c *mongo.Database) {
	deviceCollections = c.Collection("device")
}

func ServiceCollection(c *mongo.Database){
	serviceCollection = c.Collection("service")
}

/*func GetAllProducts(c *gin.Context) {
	pro := []Service{}
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
	pro := []Device{}
	fmt.Println("ID : ",id)
	pipe := []bson.M{{"$match": bson.M{"productid":id}}}
	//resp := []bson.M{}
	showInfoCursor, err := collections.Aggregate(context.TODO(),pipe)
	if err != nil {
		panic(err)
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(context.TODO(), &showsWithInfo); err != nil {
		panic(err)
	}
	fmt.Println("Show With Info : ",showsWithInfo[0]["body"])
	//fmt.Println(resp)
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
		var prod Device
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
*/

func getHash(pwd []byte) string {        
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)          
    if err != nil {
       log.Println(err)
    }
    return string(hash)
}

func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	fmt.Println("pro : ",user)

	name := user.Username
	password := getHash([]byte(user.Password))
	address := user.Address
	email := user.Email
	phone := user.Phone_number
	is_Active := true

	/*uniqEmail, err1 := userCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	if err1 != nil{
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Already Registered User",
			"data" : err1,
		})
		return
	}*/
	
	newUser := User{
		Username : name,
		Password : password,
		Address	: address,
		Email	: email,
		Phone_number : phone,
		Is_Active : is_Active,
	}

	_, err := userCollection.InsertOne(context.TODO(), newUser)

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
		"data" : newUser,
	})
	return
}

func CreateProd(c *gin.Context) {
	var pro Device
	c.BindJSON(&pro)
	fmt.Println("pro : ",pro)

	productId := pro.ProductId
	deviceName := pro.DeviceName
	user := pro.User
	brand := pro.Brand
	warranty := pro.Warranty
	manufactureDate := pro.ManufactureDate
	purchaseDate := pro.PurchaseDate
	//id := guuid.New().String()
	fmt.Println("Id : ",productId)
	fmt.Println("Desc : ",deviceName)


	newTodo := Device{
		ProductId: productId,
		DeviceName:     deviceName,
		User:      user,
		Brand: brand,
		Warranty: warranty,
		ManufactureDate : manufactureDate,
		PurchaseDate : purchaseDate,
		Is_Active : true,
	}

	_, err := deviceCollections.InsertOne(context.TODO(), newTodo)

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

func CreateService(c *gin.Context) {
	var todo Service

	c.BindJSON(&todo)
	fmt.Println("todo : ",todo)
	device := todo.Device
	productId := todo.ProductId
	user := todo.User
	deviceName := todo.DeviceName
	issue := todo.Issue
	receivedDate := todo.ReceivedDate
	//returnedDate := 
	fmt.Println("Id : ",productId)


	newTodo := Service{
		Device : device,
		ProductId : productId,
		User : user,
		DeviceName : deviceName,
		Issue : issue,
		ReceivedDate : receivedDate,
		Is_Active : true,
		IS_Closed : false,
		//ReturnedDate : returnedDate,
	}

	_, err := serviceCollection.InsertOne(context.TODO(), newTodo)

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

/*func GetAllTodos(c *gin.Context) {
	todos := []Todo{}
	cursor, err := collections.Find(context.TODO(), bson.M{})

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

	_, err := collections.InsertOne(context.TODO(), newTodo)

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
	err := collections.FindOne(context.TODO(), bson.M{"id": todoId}).Decode(&todo)
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

	_, err := collections.UpdateOne(context.TODO(), bson.M{"id": todoId}, newData)
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

	_, err := collections.DeleteOne(context.TODO(), bson.M{"id": todoId})
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
}*/

func GetAllDevices(c *gin.Context) {
	todos := []Device{}
	cursor, err := deviceCollections.Find(context.TODO(), bson.M{})

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
				var todo Device
        cursor.Decode(&todo)
        todos = append(todos, todo)
		}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Devices",
		"data":    todos,
	})
	return
}


func GetAllServices(c *gin.Context) {
	todos := []Service{}
	cursor, err := serviceCollection.Find(context.TODO(), bson.M{})

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
				var todo Service
        cursor.Decode(&todo)
        todos = append(todos, todo)
		}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Services",
		"data":    todos,
	})
	return
}

func GetDeviceByUser(c *gin.Context) {

	var userIds GetDeviceUser

	c.BindJSON(&userIds)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "user"}, {"localField", "user"}, {"foreignField", "_id"}, {"as", "user"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$user"}, {"preserveNullAndEmptyArrays", false}}}}

	showLoadedCursor, err := deviceCollections.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	if err != nil {
		panic(err)
	}
	var showsLoaded []bson.M
	if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
		panic(err)
	}
	fmt.Println(showsLoaded)

	/*todos := []Device{}
	cursor, err := serviceCollection.Find(context.TODO(), bson.M{})

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
				var todo Service
        cursor.Decode(&todo)
        todos = append(todos, todo)
		}
		*/
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Services",
		"data":    showsLoaded,
	})
	return
}


func GetDeviceByUserID(c *gin.Context) {

	var userIds GetDeviceUser

	c.BindJSON(&userIds)

	fmt.Println("User ID : ",userIds)
	//pipelineResult := make([]OrderStatusTotal, 0)
	pipeline := make([]bson.M, 0)

	groupStage := bson.M{
		"$group": bson.M{
		"user": userIds,
		"total": bson.M{"$sum": 1},
		},
	}

	matchStage := bson.M{
		"$match": bson.M{
			"user": userIds,
		},
	}

	pipeline = append(pipeline, matchStage,groupStage)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	data, err := deviceCollections.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err.Error())
		fmt.Errorf("failed to execute aggregation %s", err.Error())
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Services",
		"data" : data,
	})
	return
}