package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/centrifugal/gocent"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte("secret")
var dbConn *gorm.DB
var cent *gocent.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	dbConn = db()

	if err != nil {
		log.Fatal(err)
	}

	cent = gocent.New(gocent.Config{
		Addr: "http://116.193.190.203:8091/api",
		Key:  "ba3a0cc4-3093-4d9a-aa0f-9f45f5728905",
	})

	test, err := cent.Info(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(test)

	jwt, err := GenerateJWT(User{Id: 1, Username: "test", Password: "test"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jwt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},                                       // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow specific methods
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // Set to true if you want to allow credentials
		MaxAge:           300,  // Maximum value for preflight requests
	}))

	r.Post("/register", register)

	r.Post("/login", login)

	r.Get("/generate-qr", generateQr)

	r.Post("/get-qr-session", getQrSession)

	r.Route("/", func(r chi.Router) {
		r.Use(JwtMiddleware)

		r.Post("/scan-qr", scanQr)
	})

	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		type Test struct {
			ChannelName string `json:"channel_name"`
		}

		req := Test{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})

	fmt.Println()
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), r)
}

// register request
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func register(w http.ResponseWriter, r *http.Request) {
	// Create an instance of RegisterRequest
	req := RegisterRequest{}

	// Decode the JSON body into the struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// Handle error (e.g., invalid JSON)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log or use the data (just for demonstration)
	fmt.Printf("Received data: Username: %s, Password: %s\n", req.Username, req.Password)

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		// Handle error (e.g., password generation failed)
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		return
	}

	// Log or use the password (just for demonstration)
	fmt.Printf("Generated password: %s\n", password)

	// Create a new user in the database
	user := User{
		Username: req.Username,
		Password: string(password),
	}

	err = dbConn.Create(&user).Error
	if err != nil {
		// Handle error (e.g., user creation failed)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Log or use the user (just for demonstration)
	fmt.Printf("Created user: %+v\n", user)

	// Create a response
	response := map[string]string{
		"message": "User created successfully",
	}

	// Encode the response as JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// Handle error (e.g., JSON encoding failed)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	// Create an instance of LoginRequest
	req := LoginRequest{}

	// Decode the JSON body into the struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// Handle error (e.g., invalid JSON)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log or use the data (just for demonstration)
	fmt.Printf("Received data: Username: %s, Password: %s\n", req.Username, req.Password)

	// Find the user in the database
	var user User
	err = dbConn.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		// Handle error (e.g., user not found)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Log or use the user (just for demonstration)
	fmt.Printf("Found user: %+v\n", user)

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Handle error (e.g., invalid password)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Log or use the user (just for demonstration)
	fmt.Printf("Verified user: %+v\n", user)

	// Generate a JWT token
	token, err := GenerateJWT(user)
	if err != nil {
		// Handle error (e.g., JWT generation failed)
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	// Log or use the token (just for demonstration)
	fmt.Printf("Generated token: %s\n", token)

	// Create a response
	response := map[string]string{
		"token": token,
	}

	// Encode the response as JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// Handle error (e.g., JSON encoding failed)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type GenerateQrResponse struct {
	Uuid           string `json:"uuid"`
	WebsocketToken string `json:"websocket_token"`
	ChannelName    string `json:"channel_name"`
}

func generateQr(w http.ResponseWriter, r *http.Request) {
	// Generate a UUID
	tempUuid := uuid.New().String()

	qrscan := Qrscan{
		Uuid:       tempUuid,
		UserId:     null.Int64{},
		IsValid:    false,
		ValidUntil: time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05 -0700"),
	}

	err := dbConn.Create(&qrscan).Error
	if err != nil {
		// Handle error (e.g., user creation failed)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Log or use the user (just for demonstration)
	fmt.Printf("Created user: %+v\n", qrscan)

	token := GenerateJWTForWebsocket(fmt.Sprintf("public:%v", qrscan.Id))

	// Create a response
	response := GenerateQrResponse{
		Uuid:           tempUuid,
		WebsocketToken: token,
		ChannelName:    fmt.Sprintf("public:%v", qrscan.Id),
	}

	// Encode the response as JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// Handle error (e.g., JSON encoding failed)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type ScanQrRequest struct {
	Uuid string `json:"uuid"`
}

func scanQr(w http.ResponseWriter, r *http.Request) {
	// Handle the request here
	ctx := r.Context()

	// get uuid from request
	var req ScanQrRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// get user data from middleware
	user, ok := ctx.Value("user").(User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// find qrscan by uuid
	qrscan := Qrscan{}
	err = dbConn.Where("uuid = ?", req.Uuid).First(&qrscan).Error
	if err != nil {
		http.Error(w, "Qrscan not found", http.StatusNotFound)
		return
	}

	// update qrscan
	qrscan.UserId = null.IntFrom(int64(user.Id))
	qrscan.IsValid = true

	err = dbConn.Save(&qrscan).Error
	if err != nil {
		http.Error(w, "Failed to update qrscan", http.StatusInternalServerError)
		return
	}

	// Log or use the user (just for demonstration)
	fmt.Printf("Found user: %+v\n", user)
	channelName := fmt.Sprintf("public:%v", qrscan.Id)

	err = cent.Publish(ctx, channelName, []byte(`{"event": "qrscan"}`))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response
	response := map[string]string{
		"uuid": req.Uuid,
	}

	// Encode the response as JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// Handle error (e.g., JSON encoding failed)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type GetQrSessionRequest struct {
	Uuid string `json:"uuid"`
}

func getQrSession(w http.ResponseWriter, r *http.Request) {
	// get uuid from request
	var req GetQrSessionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// get qrscam data
	qrData := Qrscan{}
	err = dbConn.Where("uuid = ?", req.Uuid).First(&qrData).Error
	if err != nil {
		http.Error(w, "Qrscan not found", http.StatusNotFound)
		return
	}

	if !qrData.IsValid {
		http.Error(w, "Qrscan not valid", http.StatusNotFound)
		return
	}

	// get user data
	userId := qrData.UserId.Int64
	user := User{}
	err = dbConn.Where("id = ?", userId).First(&user).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// delete qrscan
	err = dbConn.Delete(&qrData).Error
	if err != nil {
		http.Error(w, "Failed to delete qrscan", http.StatusInternalServerError)
		return
	}

	// Log or use the user (just for demonstration)
	fmt.Printf("Found user: %+v\n", user)

	// generate new JWT
	token, err := GenerateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	// Create a response
	response := map[string]string{
		"token": token,
	}

	// Encode the response as JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// Handle error (e.g., JSON encoding failed)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
