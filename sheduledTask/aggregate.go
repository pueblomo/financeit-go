package sheduledTask

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"main.go/aggregates"
	"main.go/database"
	"main.go/model"
)

func SendToAggregator(){
	log.Println("Sheduled Task send failed categories again")
	rows, err := database.Conn.Query(context.Background(),"SELECT categoryid, name, icon FROM caterror")
	if err != nil {
		log.Println(err)
		return
	}

	categories := getCategoriesFromRows(rows)
	test, err := database.Conn.Exec(context.Background(), "DELETE FROM caterror")
	log.Println(test)
	log.Printf("Sending %v categories", len(categories))
	for _,cat := range categories {
		aggregates.SendCategory(cat)
	} 
	log.Println("Finished sending")

	log.Println("Sheduled Task send failed expenses again")
	rows, err = database.Conn.Query(context.Background(),"SELECT expenseid, description, date, price, categoryid FROM experror")
	if err != nil {
		log.Println(err)
		return
	}

	expenses := getExpensesFromRows(rows)
	test, err = database.Conn.Exec(context.Background(), "DELETE FROM experror")
	log.Println(test)
	log.Printf("Sending %v expenses", len(categories))
	for _,exp := range expenses {
		aggregates.SendExpense(exp)
	} 
	log.Println("Finished sending")
}

func getCategoriesFromRows(rows pgx.Rows) []model.Category {
	var categories []model.Category 
	for rows.Next() {
		var id int64
		var name string
		var icon string
		rows.Scan(&id, &name, &icon)
		categories = append(categories, model.Category {id,name,icon, int64(0)})
	}
	return categories
}

func getExpensesFromRows(rows pgx.Rows) []model.Expense {
	var expenses []model.Expense 
	for rows.Next() {
		var id int64
		var description string
		var date time.Time
		var price int
		var categoryId int64
		rows.Scan(&id, &description, &date, &price, &categoryId)
		expenses = append(expenses, model.Expense {id,date,description,price,categoryId})
	}
	return expenses
}