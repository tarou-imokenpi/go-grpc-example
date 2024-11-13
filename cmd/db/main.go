package main

import (
	"github.com/dgraph-io/badger/v4"
	"log"
)

func main() {
	opts := badger.DefaultOptions("/tmp/badger")
	opts.IndexCacheSize = 100 << 20 // 100 MB

	// open DB
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//// 読み取り専用トランザクション
	//err = db.View(func(txn *badger.Txn) error {
	//	// Your code here…
	//	return nil
	//})
	//
	//// 読み取り/書き込みトランザクション
	//err = db.Update(func(txn *badger.Txn) error {
	//	// Your code here…
	//	return nil
	//})

	err = entryByte(db, []byte("key"), []byte("value"))
	if err != nil {
		log.Fatal(err)
	}
}

func entryByte(db *badger.DB, key []byte, value []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, value)
		err := txn.SetEntry(e)
		return err
	})
}
