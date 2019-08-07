// Package main initializes server and in-memory map loaded from offline database, 
// and registers handlers to controller functions.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gpl/ch7/exercises/e7.12/ctrl"
	"gpl/ch7/exercises/e7.12/model"
	"gpl/ch7/exercises/e7.12/util"
)

func main() {
	// create db directory and log directory/file if non-existent
	if _, err := os.Stat("db"); os.IsNotExist(err) {
		os.Mkdir("db", 0744)
		fmt.Println("'main': 'db' directory created")
	}
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", 0744)
		fmt.Println("'main': 'log' directory created")
	}
	if _, err := os.Stat("log/fail_log.log"); os.IsNotExist(err) {
		os.Create("log/fail_log.log")
		fmt.Println("'main':'log/fail_log.log' created")
	}

	// create/open database and load contents in memory
	if err := model.MainTx(ctrl.DbMap); err != nil {
		util.FailLog(err)
		os.Exit(1)
	}

	// create server and handler functions
	fmt.Println("initializing server...")
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         "localhost:8000",
		Handler:      nil,
	}
	fmt.Printf("server address: %v\nread timeout: %v\nwrite timeout: %v\n",
		srv.Addr, srv.ReadTimeout, srv.WriteTimeout)
	fmt.Printf("listening at: '%v'...\n", srv.Addr)

	http.HandleFunc("/home", ctrl.Home)
	http.HandleFunc("/list", ctrl.ReadOpList)
	http.HandleFunc("/price", ctrl.ReadOpPrice)
	http.HandleFunc("/update", ctrl.CreUpOp)
	http.HandleFunc("/delete", ctrl.DelOp)
	log.Fatal(srv.ListenAndServe())
}
