// Package ctrl contains components for performing CRUD operations on the database. 
// Functions are tied to HTTP Handlers in main/main.go.
package ctrl

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"tpgl_inventory_srvr/model"
	"tpgl_inventory_srvr/util"
	"tpgl_inventory_srvr/view"
)

// DbMap maps data from disk in memory
var DbMap = make(model.Database)

var mu = &sync.Mutex{}

// Home displays home page
func Home(w http.ResponseWriter, r *http.Request) {
	if err := view.TemplExe(view.TmplMap["home"], w, nil); err != nil {
		util.FailLog(err)
		fmt.Fprintf(w, "html template failed to execute: %s", err)
		fmt.Printf("html template failed to execute: %s", err)
		return
	}
	return
}

// ReadOpList Reads all items in database from memory and returns dataset to view
func ReadOpList(w http.ResponseWriter, r *http.Request) {
	if err := view.TemplExe(view.TmplMap["list"], w, DbMap); err != nil {
		util.FailLog(err)
		fmt.Fprintf(w, "html template failed to execute: %s", err)
		fmt.Printf("html template failed to execute: %s", err)
		return
	}
	fmt.Println("'list' called; template executed")
	return
}

// ReadOpPrice Reads price for specified item
func ReadOpPrice(w http.ResponseWriter, r *http.Request) {
	name, _ := view.GetFields(r)
	item, ok := DbMap[name] // check if item exists

	// handle if named item doesn't exist
	if !ok {
		data := view.HTMLData{Name: name, Price: ""}
		if err := view.TemplExe(view.TmplMap["priceFail"], w, data); err != nil {
			util.FailLog(err)
			fmt.Fprintf(w, "html template failed to execute: %s", err)
			fmt.Printf("html template failed to execute: %s", err)
			return
		}
		fmt.Printf("'price' call failed: non-existent item: '%s'\n", name)
		return
	}

	// display price for specified item to user
	data := view.HTMLData{Name: name, Price: item.Price.String()}
	if err := view.TemplExe(view.TmplMap["price"], w, data); err != nil {
		util.FailLog(err)
		fmt.Fprintf(w, "html template failed to execute: %s", err)
		fmt.Printf("html template failed to execute: %s", err)
		return
	}
	fmt.Printf("'price' called for '%s'\n", name)
	return
}

// CreUpOp controls logic for Create/Update operations on database
func CreUpOp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("'update' called")
	name, price := view.GetFields(r)

	// check price is set
	if price == "" {
		uErr := "error: price not set"
		if err := view.TemplExe(view.TmplMap["updateFail"], w, uErr); err != nil {
			util.FailLog(err)
			fmt.Fprintf(w, "html template failed to execute: %s", err)
			fmt.Printf("html template failed to execute: %s", err)
			return
		}
		fmt.Println("'update' call failed: price not set; fail template executed")
		return
	}

	// convert string from URL to float64 and check value is numeric
	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		uErr := "error: price must be numerical value"
		if err := view.TemplExe(view.TmplMap["updateFail"], w, uErr); err != nil {
			util.FailLog(err)
			fmt.Fprintf(w, "html template failed to execute: %s", err)
			fmt.Printf("html template failed to execute: %s", err)
			return
		}
		fmt.Println("'update' call failed: price set to non-numerical value; fail template executed")
		return
	}

	// verify price is >= 0
	if p < 0 {
		uErr := "error: price must be greater than or equal to 0"
		if err := view.TemplExe(view.TmplMap["updateFail"], w, uErr); err != nil {
			util.FailLog(err)
			fmt.Fprintf(w, "html template failed to execute: %s", err)
			fmt.Printf("html template failed to execute: %s", err)
			return
		}
		fmt.Println("'update' call failed: price less than 0; fail template executed")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Create/Update transaction
	err = model.CreUpTx(name, p)
	if err != nil {
		if tErr := view.TemplExe(view.TmplMap["dbFail"], w, err); tErr != nil {
			util.FailLog(err)
			fmt.Fprintf(w, "html template failed to execute: %s", err)
			fmt.Printf("html template failed to execute: %s", err)
			return
		}
		util.FailLog(err)
		return
	}

	// update in memory
	DbMap[name] = model.Product{name, model.Dollars(p)} // convert float to dollars type
	fmt.Println("'update': in-memory map updated")

	// execute template
	data := view.HTMLData{Name: name, Price: DbMap[name].Price.String()}
	if err := view.TemplExe(view.TmplMap["update"], w, data); err != nil {
		util.FailLog(err)
		fmt.Fprintf(w, "html template failed to execute: %s", err)
		fmt.Printf("html template failed to execute: %s", err)
		return
	}
	fmt.Println("'delete': template executed")

	return
}

// DelOp controls logic for Delete operations on database
func DelOp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("'delete' called")
	name, _ := view.GetFields(r)

	// check item exists in db
	_, ok := DbMap[name]
	if !ok {
		if err := view.TemplExe(view.TmplMap["deleteFail"], w, name); err != nil {
			util.FailLog(err)
			fmt.Fprintf(w, "html template failed to execute: %s", err)
			fmt.Printf("html template failed to execute: %s", err)
			return
		}
		fmt.Printf("'delete' call failed: item '%s' does not exist; fail template executed\n", name)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Delete transaction
	err := model.DelTx(name)
	if err != nil {
		if tErr := view.TemplExe(view.TmplMap["dbFail"], w, err); tErr != nil {
			util.FailLog(err)
			fmt.Fprintf(w, "html template failed to execute: %s", err)
			fmt.Printf("html template failed to execute: %s", err)
		}
		util.FailLog(err)
		return
	}

	// delete from in-memory map
	delete(DbMap, name)
	fmt.Println("'delete': in-memory map updated")

	// execute template
	if err := view.TemplExe(view.TmplMap["delete"], w, name); err != nil {
		util.FailLog(err)
		fmt.Fprintf(w, "html template failed to execute: %s", err)
		fmt.Printf("html template failed to execute: %s", err)
		return
	}
	fmt.Println("'delete': template executed")

	return
}
