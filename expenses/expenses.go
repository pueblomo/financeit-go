package expenses

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"main.go/database"
)

type Expense  struct {
	Id int64 `json:"id"`
	Date time.Time `json:"date"`
	Description string `json:"description"`
	Price json.Number `json:"price"`
	CategoryId int64 `json:"categoryId"`
}

func getExpensesFromRows(rows pgx.Rows) []Expense {
	var expenses []Expense
	for rows.Next() {
		var id int64
		var date time.Time
		var desc string
		var price json.Number
		var catId int64
		rows.Scan(&id, &date, &desc, &price, &catId)
		expenses = append(expenses, Expense{id, date,desc,price,catId})
	}
	return expenses
}

func PostExpense(ctx *fiber.Ctx) error {
	var newExp Expense
	if err:=ctx.BodyParser(&newExp); err!=nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	log.Printf("Post Expense: %v\n",newExp)

	_, err := database.Conn.Exec(context.Background(), "INSERT INTO expense (date, description, price, categoryid) VALUES ($1, $2, $3, $4)", newExp.Date,newExp.Description,newExp.Price,newExp.CategoryId)
	if err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}