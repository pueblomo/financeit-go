package expenses

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"main.go/conf"
	"main.go/database"
)

type Expense  struct {
	Id int64 `json:"id"`
	Date time.Time `json:"date" validate:"required"`
	Description string `json:"description" validate:"required,max=60"`
	Price int `json:"price" validate:"required,number,min=0"`
	CategoryId int64 `json:"categoryId" validate:"required,min=0"`
}

func getExpensesFromRows(rows pgx.Rows) []Expense {
	var expenses []Expense
	for rows.Next() {
		var id int64
		var date time.Time
		var desc string
		var price int
		var catId int64
		rows.Scan(&id, &desc, &date, &price, &catId)
		expenses = append(expenses, Expense{id, date,desc,price,catId})
	}
	return expenses
}

func PostExpense(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	var newExp Expense
	if err:=ctx.BodyParser(&newExp); err!=nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Post Expense: %v\n",newExp)

	_, err := database.Conn.Exec(context.Background(), "INSERT INTO expense (date, description, price, categoryid) VALUES ($1, $2, $3, $4)", newExp.Date,newExp.Description,newExp.Price,newExp.CategoryId)
	if err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func GetExpenses(ctx *fiber.Ctx) error {
	log.Println("Get Expenses")
	id := ctx.Params("catId")
	if id == "" {
		log.Println("can't get id from path")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Get Expenses with id: %v", id)

	rows, err := database.Conn.Query(context.Background(), "SELECT * FROM expense WHERE categoryId = $1 order by date", id)
	defer rows.Close()
	conf.CheckError(err)

	return ctx.Status(fiber.StatusOK).JSON(getExpensesFromRows(rows))
}