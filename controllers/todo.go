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

type TestCollectionStruct struct{
	Item	string	`json:"item"`
	Qty 	int		`json:"qty"`
	Size 	interface{} 	`json:"size"`
	Status 	string 		`json:"status"`
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
var deviceCollections,serviceCollection,userCollection,testCollection *mongo.Collection

func UserCollection(c *mongo.Database){
	userCollection = c.Collection("user")
}

func DeviceCollection(c *mongo.Database) {
	deviceCollections = c.Collection("device")
}

func ServiceCollection(c *mongo.Database){
	serviceCollection = c.Collection("service")
}

func TestCollection(c *mongo.Database){
	testCollection = c.Collection("inventory_query_top")
}


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

	result, err := deviceCollections.InsertOne(context.TODO(), newTodo)

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
		"data" : result.InsertedID,
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

	result, err := serviceCollection.InsertOne(context.TODO(), newTodo)

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
		"data" : result.InsertedID,
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
	//folders := []Device{}
	c.BindJSON(&userIds)


	fmt.Println("User ID : ",userIds)
	//pipelineResult := make([]OrderStatusTotal, 0)
	/*pipeline := make([]bson.M, 0)

	groupStage := bson.M{
		"$group": bson.M{
		"user": primitive.ObjectID(userIds.ID),
		"total": bson.M{"$sum": 1},
		},
	}

	matchStage := bson.M{
		"$match": bson.M{
			"user": primitive.ObjectID(userIds.ID),
		},
	}

	pipeline = append(pipeline, matchStage,groupStage)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	data, err := deviceCollections.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err.Error())
		fmt.Errorf("failed to execute aggregation %s", err.Error())
		return
	}*/

	query := []bson.M{{
		"$lookup": bson.M{ // lookup the documents table here
		  "from":         "user",
		  "localField":   "user",
		  "foreignField": "_id",
		  "as":           "user",
		}},
		{"$match": bson.M{
		  //"level": bson.M{"$lte": user.Level},
		  "user": userIds.ID,
	  }}}

	  //lookupStage := bson.D{{"$lookup", bson.D{{"from", "user"}, {"localField", "user"}, {"foreignField", "_id"}, {"as", "user"}}}}
	  
	  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	  data, err  := deviceCollections.Aggregate(ctx, query)
	  if err != nil {
		log.Println(err.Error())
		fmt.Errorf("failed to execute aggregation %s", err.Error())
		return
	}
	//err1 := data.All(&folders)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Services",
		"data" : data,
	})
	return
}

func TestInsert(c *gin.Context) {

		// Start Example 6

		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
		}

		result, err := testCollection.InsertMany(context.Background(), docs)
		
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
			"message": "Inserted Successfully",
			"data" : result.InsertedIDs,
		})
		return

}


func TestQueryOne(c *gin.Context) {

	todos := []TestCollectionStruct{}

	var testStruct TestCollectionStruct
	//folders := []Device{}
	c.BindJSON(&testStruct)
	fmt.Println("test Struct : ",testStruct.Status)
	cursor, err := testCollection.Find(
		context.Background(),
		bson.D{{"status", testStruct.Status}},
	)

	if err != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	for cursor.Next(context.TODO()) {
		var todo TestCollectionStruct
		cursor.Decode(&todo)
		todos = append(todos, todo)
		}
	// End Example 9

	//require.NoError(t, err)
	//requireCursorLength(t, cursor, 2)

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Queried Successfully",
		"data" : todos,
	})

	return

}