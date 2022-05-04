package aggregates

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"main.go/conf"
	"main.go/database"
	"main.go/model"
)

func SendCategory(cat model.Category ){
	cat.Id = getIdForCategory(cat)
	log.Printf("Send %v to aggregator\n", cat)
	catJson,err := json.Marshal(cat)
	if err != nil {
		saveToErrorTable(cat, err)
		return
	}
	resp, err := http.Post("http://"+conf.GetAggregatorUrl()+":7080/categories", "application/json",bytes.NewBuffer(catJson))
	if err != nil {
		saveToErrorTable(cat, err)
		return
	}
	if resp.StatusCode != 201 {
		saveToErrorTable(cat,fmt.Errorf("Response Code: %v", resp.StatusCode))
	}
}

func getIdForCategory(cat model.Category) int64{
	rows, err:= database.Conn.Query(context.Background(),"SELECT id FROM category where name = $1", cat.Name)
	defer rows.Close()
	conf.CheckError(err)
	var id int64
	if rows.Next() {
		rows.Scan(&id)
	}
	return id
}

func saveToErrorTable(cat model.Category, err error){
	log.Println(err)
	log.Println("Save to error table")
	test, err := database.Conn.Exec(context.Background(),"INSERT INTO caterror (categoryid, name, icon) VALUES ($1, $2, $3)", cat.Id, cat.Name, cat.Icon)
	log.Println(test)
	conf.CheckError(err)
}

func SendExpense(exp model.Expense) {
	exp.Id = getIdForExpense(exp)
	log.Printf("Send %v to aggregator", exp)
	expJson, err := json.Marshal(exp)
	if err != nil {
		saveToExpErrorTable(exp,err)
		return
	}
	resp, err := http.Post("http://"+conf.GetAggregatorUrl()+":7080/expenses", "application/json",bytes.NewBuffer(expJson))
	if err != nil {
		saveToExpErrorTable(exp, err)
		return
	}
	if resp.StatusCode != 201 {
		saveToExpErrorTable(exp, fmt.Errorf("Response Code: %v", resp.StatusCode))
	}
	
}

func getIdForExpense(exp model.Expense) int64 {
	rows, err := database.Conn.Query(context.Background(), "SELECT id from expense where description = $1 and date = $2", exp.Description, exp.Date)
	defer rows.Close()
	conf.CheckError(err)
	var id int64
	if rows.Next() {
		rows.Scan(&id)
	}
	return id
}

func saveToExpErrorTable(exp model.Expense, err error) {
	log.Println(err)
	log.Panicln("Save to expense error table")
	test, err := database.Conn.Exec(context.Background(), "INSERT INTO experror (expenseid, description, date, price, categoryid) VALUES ($1, $2, $3, $4, $5)", exp.Id, exp.Description, exp.Date, exp.Price, exp.CategoryId)
	log.Println(test)
	conf.CheckError(err)
}