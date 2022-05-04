package categories

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"main.go/aggregates"
	"main.go/conf"
	"main.go/database"
	"main.go/model"
)

var validate = validator.New()

func getCategoriesFromRows(rows pgx.Rows) []model.Category {
	var categories []model.Category 
	for rows.Next() {
		var id int64
		var name string
		var icon string
		rows.Scan(&id, &name, &icon)
		categories = append(categories, model.Category {id,name,icon, int64(getAmount(id))})
	}
	return categories
}

func GetCategories(ctx *fiber.Ctx) error{
	log.Println("Get Categories")
	rows, err := database.Conn.Query(context.Background(),"SELECT * FROM category order by name")
	defer rows.Close()
	conf.CheckError(err)

	return ctx.Status(fiber.StatusOK).JSON(getCategoriesFromRows(rows))
}

func PostCategorie(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	var newCat model.Category 
	log.Println("Post category")
	if err:= ctx.BodyParser(&newCat); err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if err:= validate.Struct(newCat); err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	log.Printf("Post Category : %v\n", newCat)

	_, err := database.Conn.Exec(context.Background(),"INSERT INTO category (name, icon) VALUES ($1, $2)", newCat.Name, newCat.Icon)
	if err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	go aggregates.SendCategory(newCat)

	return ctx.SendStatus(fiber.StatusCreated)
}


func getAmount(id int64) int {
	t := time.Now()
	firstday := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	lastday := firstday.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	firstdaySplit := strings.Split(firstday.String(), " ")
	lastdaySplit := strings.Split(lastday.String(), " ")
	rows, err := database.Conn.Query(context.Background(),"SELECT SUM (price) FROM expense WHERE categoryid = $1 and date >= $2 and date <= $3", strconv.FormatInt(id,10), firstdaySplit[0] + " " + firstdaySplit[1], lastdaySplit[0] + " " + lastdaySplit[1])
	conf.CheckError(err)
	var amount int
	for rows.Next() {
		rows.Scan(&amount)
	}

	return amount
}
