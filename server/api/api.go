package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"server/db"
	"server/models"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// var userCollection *mongo.Collection = db.GetCollection("users")
var validate = validator.New()
var psql = db.DB.GetDB()

func SignInHandler(c *fiber.Ctx) error {

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body structure
	if err := c.BodyParser(&user); err != nil {
		return Response(c, err.Error(), http.StatusBadRequest, "Body Parser Error", "X-Correlation-Id")
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return Response(c, validationErr.Error(), http.StatusBadRequest, "Validation Error", "X-Correlation-Id")
	}

	{ //Check if the user already exists
		//Check id
		var count int
		err := psql.QueryRow("SELECT COUNT(*) FROM users WHERE user_id = $1", user.UserID).Scan(&count)
		if err != nil {
			log.Panicln(err)
			return Response(c, err.Error(), http.StatusInternalServerError, "Check UserID Error", "X-Correlation-Id")
		}
		if count > 0 {
			log.Panicln(err)
			return Response(c, "User already exists", http.StatusConflict, "User already exists", "X-Correlation-Id")
		}
		//Check email
		err = psql.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&count)
		defer cancel()
		if err != nil {
			log.Panicln(err)
			return Response(c, err.Error(), http.StatusInternalServerError, "Check Email Error", "X-Correlation-Id")
		}
		if count > 0 {
			log.Panicln(err)
			return Response(c, "User already exists", http.StatusConflict, "User already exists", "X-Correlation-Id")
		}
	}

	//insert the new user field
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Password = HashPassword(user.Password)
	log.Println("p: ", user.Password)
	user.Token, user.RefreshToken = GenerateTokens(user.Email, user.FirstName, user.LastName, user.UserID)

	err := psql.QueryRow(
		`INSERT INTO users (user_id,first_name, last_name, email, password, token, refresh_token, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.UserID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Token,
		user.RefreshToken,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		log.Panic(err)
		return Response(c, err.Error(), http.StatusInternalServerError, "DB Insert Error", "X-Correlation-Id")
	}

	return Response(c, user.ID, http.StatusCreated, "User Created Successfully", "X-Correlation-Id")
}

func LogInHandler(c *fiber.Ctx) error {
	var user models.User

	//validate the request body structure
	if err := c.BodyParser(&user); err != nil {
		return Response(c, err.Error(), http.StatusBadRequest, "Body Parser Error", "X-Correlation-Id")
	}
	password := user.Password

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := psql.QueryRow("SELECT * FROM users WHERE user_id = $1;", user.UserID).Scan(
		&user.ID,
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Token,
		&user.RefreshToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		log.Panicln(err)
		return Response(c, err.Error(), http.StatusInternalServerError, "DB Query Error", "X-Correlation-Id")
	}
	// TODO: errore in caso non trovi nulla nel db
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		return c.Status(http.StatusNotFound).JSON(models.ErrorResponse{
	// 			Meta: models.Meta{
	// 				Status:        http.StatusNotFound,
	// 				Message:       "no data found",
	// 				TimeStamp:     time.Now(),
	// 				CorrelationId: "X-Correlation-Id"},
	// 			Errors: &fiber.Map{"data": "User not found"}})
	// 	}
	// 	return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{
	// 		Meta: models.Meta{
	// 			Status:        http.StatusInternalServerError,
	// 			Message:       "error finding data process",
	// 			TimeStamp:     time.Now(),
	// 			CorrelationId: "X-Correlation-Id"},
	// 		Errors: &fiber.Map{"data": err.Error()}})
	// }

	if VerifyPassword(password, user.Password) {

		// update tokens
		token, refresh_token := GenerateTokens(user.Email, user.FirstName, user.LastName, user.UserID)
		update_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		if _, err := psql.Query("UPDATE users SET token = $1, refresh_token = $2, updated_at = $3 WHERE user_id = $4;", token, refresh_token, update_at, user.UserID); err != nil {
			log.Panicln(err)
			return Response(c, err.Error(), http.StatusInternalServerError, "DB Update tokens Error", "X-Correlation-Id")
		}

		return Response(c, &fiber.Map{"token": user.Token, "user_id": user.UserID}, http.StatusOK, "OK", "X-Correlation-Id")
	}

	log.Fatal("Authentication Failed: ID or Password is incorrect")
	return Response(c, "Authentication Failed: ID or Password is incorrect", http.StatusUnauthorized, "Authentication Failed", "X-Correlation-Id")
}

// func GetUsersHandler(c *fiber.Ctx) error {
// 	var checkToken models.User
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	token := c.GetReqHeaders()["Token"]
// 	if token == "" {
// 		return c.Status(http.StatusUnauthorized).JSON(models.ErrorResponse{
// 			Meta: models.Meta{
// 				Status:        http.StatusNotFound,
// 				Message:       "Unauthorized",
// 				TimeStamp:     time.Now(),
// 				CorrelationId: "X-Correlation-Id"},
// 			Errors: &fiber.Map{"data": "Unauthorized"}})
// 	}

// 	err := userCollection.FindOne(ctx, bson.M{"token": token}).Decode(&checkToken)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return c.Status(http.StatusNotFound).JSON(models.ErrorResponse{
// 				Meta: models.Meta{
// 					Status:        http.StatusNotFound,
// 					Message:       "no data found",
// 					TimeStamp:     time.Now(),
// 					CorrelationId: "X-Correlation-Id"},
// 				Errors: &fiber.Map{"data": err.Error()}})
// 		}
// 		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{
// 			Meta: models.Meta{
// 				Status:        http.StatusInternalServerError,
// 				Message:       "error finding data process",
// 				TimeStamp:     time.Now(),
// 				CorrelationId: "X-Correlation-Id"},
// 			Errors: &fiber.Map{"data": err.Error()}})
// 	}

// 	if checkToken.UserType == "ADMIN" {

// 		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
// 		if err != nil || recordPerPage < 1 {
// 			recordPerPage = 10
// 		}
// 		page, err := strconv.Atoi(c.Query("page"))
// 		if err != nil || page < 1 {
// 			page = 1
// 		}
// 		startIndex, err := strconv.Atoi(c.Query("startIndex"))
// 		if err != nil || startIndex < 1 {
// 			startIndex = (page - 1) * recordPerPage
// 		}

// 		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
// 		groupStage := bson.D{{Key: "$group", Value: bson.D{
// 			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
// 			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
// 			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
// 		}}}
// 		projectStage := bson.D{{Key: "$project", Value: bson.D{
// 			{Key: "_id", Value: 0},
// 			{Key: "total_count", Value: 1},
// 			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"data", startIndex, recordPerPage}}}},
// 		}}}

// 		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
// 			matchStage,
// 			groupStage,
// 			projectStage,
// 		})
// 		if err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{
// 				Meta: models.Meta{
// 					Status:        http.StatusInternalServerError,
// 					Message:       "error finding data process",
// 					TimeStamp:     time.Now(),
// 					CorrelationId: "X-Correlation-Id"},
// 				Errors: &fiber.Map{"data": err.Error()}})
// 		}

// 		var users []bson.M
// 		if err := result.All(ctx, &users); err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{
// 				Meta: models.Meta{
// 					Status:        http.StatusInternalServerError,
// 					Message:       "error finding data process",
// 					TimeStamp:     time.Now(),
// 					CorrelationId: "X-Correlation-Id"},
// 				Errors: &fiber.Map{"data": err.Error()}})
// 		}

// 		return c.Status(http.StatusOK).JSON(models.SuccessResponse{
// 			Meta: models.Meta{
// 				Status:        http.StatusOK,
// 				Message:       "OK",
// 				TimeStamp:     time.Now(),
// 				CorrelationId: "X-Correlation-Id"},
// 			Result: users})

// 	}

// 	return c.Status(http.StatusUnauthorized).JSON(models.ErrorResponse{
// 		Meta: models.Meta{
// 			Status:        http.StatusUnauthorized,
// 			Message:       "Unauthorized",
// 			TimeStamp:     time.Now(),
// 			CorrelationId: "X-Correlation-Id"},
// 		Errors: &fiber.Map{"data": "Unauthorized"}})
// }

func GetUserHandler(c *fiber.Ctx) error {

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// var result models.User
	// token := c.GetReqHeaders()["Token"]
	// user_id := c.Params("user_id")
	// if token == "" {
	// 	log.Panic("token is required")
	// 	return Response(c, "Unauthorized", http.StatusUnauthorized, "Unauthorized", "X-Correlation-Id")
	// }

	// err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		return c.Status(http.StatusNotFound).JSON(models.ErrorResponse{
	// 			Meta: models.Meta{
	// 				Status:        http.StatusNotFound,
	// 				Message:       "no data found",
	// 				TimeStamp:     time.Now(),
	// 				CorrelationId: "X-Correlation-Id"},
	// 			Errors: &fiber.Map{"data": err.Error()}})
	// 	}
	// 	return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{
	// 		Meta: models.Meta{
	// 			Status:        http.StatusInternalServerError,
	// 			Message:       "error finding data process",
	// 			TimeStamp:     time.Now(),
	// 			CorrelationId: "X-Correlation-Id"},
	// 		Errors: &fiber.Map{"data": err.Error()}})
	// }

	// if result.UserType == "USER" && result.Token == token {
	// 	return c.Status(http.StatusOK).JSON(models.SuccessResponse{
	// 		Meta: models.Meta{
	// 			Status:        http.StatusOK,
	// 			Message:       "OK",
	// 			TimeStamp:     time.Now(),
	// 			CorrelationId: "X-Correlation-Id"},
	// 		Result: result})
	// }
	// if result.UserType == "ADMIN" {
	// 	return c.Status(http.StatusOK).JSON(models.SuccessResponse{
	// 		Meta: models.Meta{
	// 			Status:        http.StatusOK,
	// 			Message:       "OK",
	// 			TimeStamp:     time.Now(),
	// 			CorrelationId: "X-Correlation-Id"},
	// 		Result: result})
	// }
	// return c.Status(http.StatusUnauthorized).JSON(models.ErrorResponse{
	// 	Meta: models.Meta{
	// 		Status:        http.StatusUnauthorized,
	// 		Message:       "Unauthorized",
	// 		TimeStamp:     time.Now(),
	// 		CorrelationId: "X-Correlation-Id"},
	// 	Errors: &fiber.Map{"data": "Unauthorized"}})
	return nil
}
