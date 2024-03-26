package main

import "time"

type ErrorResponse struct {
	Error string `json:"error`
}

type User struct {
	ID 				int64 		`json:"id"`	
	Email 			string 		`json:"email"`	
	FirstName 		string 		`json:"firstName"`	
	LastName 		string 		`json:"lastName"`	
	Password 		string  	`json:"password"`
	CreatedAt 		time.Time 	`json:"createdAt"`	
}

type Task struct {
	ID 				int64 		`json:"id"`	
	Name 			string 		`json:"name"`	
	Status 			string 		`json:"status"`	
	ProjectID 		int64 		`json:"projectID"`	
	AssignedToID 	int64 		`json:"assignedTo"`	
	CreatedAt 		time.Time 	`json:"createdAt"`	
}

type Project struct {
	ID 				int64 		`json:"id"`	
	Name 			string 		`json:"name"`	
	CreatedAt 		time.Time 	`json:"createdAt"`	
}