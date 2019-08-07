package model

import (
	"fmt"

	"github.com/boltdb/bolt"
	"tpgl_inventory_srvr/util"
)

// Database stores entries in memory
type Database map[string]Product // key string == product ID

// Product stores inventory item data as struct
type Product struct {
	Name  string
	Price Dollars
}

// Dollars represents dollar value of items in DB and includes String() method
type Dollars float32

func (d Dollars) String() string { return fmt.Sprintf("$%.2f", d) }

// MainTx initializes Database map and loads offline contents into memory
func MainTx(dbMap Database) error {
	// create/load existing database into in-memory map
	db, err := bolt.Open("db/inventory.db", 0755, nil)
	if err != nil {
		fmt.Println("'main': FATAL: 'db/inventory.db' failed to open")
		return fmt.Errorf("'main': FATAL: 'db/inventory.db' failed to open")
	}
	defer db.Close()
	fmt.Println("'main': 'db/inventory.db' opened")

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("inventory"))
		if err != nil {
			fmt.Println("'main': FATAL: 'db/inventory.db': 'inventory' bucket failed to open")
			return fmt.Errorf("'main': FATAL: 'db/inventory.db': 'inventory' bucket failed to open: %v", err)
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			dbMap[string(k)] = Product{string(k), Dollars(util.BytesToFloat64(v))}
		}
		fmt.Println("'main': database contents loaded into memory")
		return nil
	}); err != nil {
		fmt.Println("'main': FATAL: 'db/inventory.db': 'inventory' bucket failed to open")
		return fmt.Errorf("FATAL: 'db/inventory.db': 'inventory' bucket failed to open: %v", err)
	}

	fmt.Println("'main': 'db/inventory.db' closed")
	return nil
}

// DelTx creates an Update transaction to delete specified item from offline database
func DelTx(name string) error {
	// open db to initialize Delete transaction
	oldb, err := bolt.Open("db/inventory.db", 0755, nil) // offline database
	if err != nil {
		fmt.Println("'delete': 'db/inventory.db' failed to open; fail template executed")
		return fmt.Errorf("db/inventory.db' failed to open: %v", err)
	}
	fmt.Println("'delete': 'db/inventory.db' opened")
	defer oldb.Close()

	// Delete transaction
	if err := oldb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("inventory"))
		if err != nil {
			fmt.Println("'delete': 'db/inventory.db': 'inventory' bucket failed to open")
			return fmt.Errorf("offline database could not be opened; try again\n%v", err)
		}
		if err := b.Delete([]byte(name)); err != nil {
			fmt.Println("'delete': 'db/inventory.db': delete operation failed")
			return fmt.Errorf("could not delete; try again\n%v", err)
		}
		fmt.Printf("'delete': entry deleted from disk: '%s'\n", name)
		return nil
	}); err != nil {
		fmt.Println("'delete': delete operation failed; fail template executed")
		return err
	}

	fmt.Println("'delete': 'db/inventory.db' closed")
	return nil
}

// CreUpTx creates update transaction to create/update offline db entries with specified values
func CreUpTx(name string, price float64) error {
	// open db
	oldb, err := bolt.Open("db/inventory.db", 0755, nil) // offline database
	if err != nil {
		fmt.Println("'update': 'db/inventory.db' failed to open; fail template executed")
		return fmt.Errorf("db/inventory.db' failed to open: %v", err)
	}
	defer oldb.Close()
	fmt.Println("'update': 'db/inventory.db' opened")

	// Create/Update transaction
	if err := oldb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("inventory"))
		if err != nil {
			fmt.Println("'update': 'db/inventory.db': 'inventory' bucket failed to open")
			return fmt.Errorf("offline database could not be opened:\n%v", err)
		}
		if err := b.Put([]byte(name), util.Float64ToBytes((price))); err != nil { // serialize k,v
			fmt.Printf("'update': 'db/inventory.db': '%s': '%v' failed to store\n", name, price)
			return fmt.Errorf("could not update:\n%v", err)
		}
		fmt.Printf("'update': entry stored: name: '%s', price: '%v'\n", name, price)
		return nil
	}); err != nil {
		fmt.Println("'update': data store failed; fail template executed")
		return err
	}
	fmt.Println("'update': 'db/inventory.db' closed")
	return nil
}
