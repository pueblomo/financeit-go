package model

import "time"

type Category struct {
	Id int64 `json:"id"`
	Name string `json:"name" validate:"required,max=60"`
	Icon string `json:"icon" validate:"required,max=20"`
	Amount int64 `json:"amount"`
}

type Expense  struct {
	Id int64 `json:"id"`
	Date time.Time `json:"date" validate:"required"`
	Description string `json:"description" validate:"required,max=60"`
	Price int `json:"price" validate:"required,number,min=0"`
	CategoryId int64 `json:"categoryId" validate:"required,min=0"`
}