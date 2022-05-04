package expenses

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"main.go/aggregates"
	"main.go/conf"
	"main.go/database"
	"main.go/model"
)

func getExpensesFromRows(rows pgx.Rows) []model.Expense {
	var expenses []model.Expense 
	for rows.Next() {
		var id int64
		var date time.Time
		var desc string
		var price int
		var catId int64
		rows.Scan(&id, &desc, &date, &price, &catId)
		expenses = append(expenses, model.Expense{id, date,desc,price,catId})
	}
	return expenses
}

func PostExpense(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	var newExp model.Expense
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

	go aggregates.SendExpense(newExp)

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
	t := time.Now()
	firstday := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	lastday := firstday.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	firstdaySplit := strings.Split(firstday.String(), " ")
	lastdaySplit := strings.Split(lastday.String(), " ")

	rows, err := database.Conn.Query(context.Background(), "SELECT * FROM expense WHERE categoryId = $1 and date >= $2 and date <= $3 order by date desc", id, firstdaySplit[0] + " " + firstdaySplit[1], lastdaySplit[0] + " " + lastdaySplit[1])
	defer rows.Close()
	conf.CheckError(err)

	return ctx.Status(fiber.StatusOK).JSON(getExpensesFromRows(rows))
}