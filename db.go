package main

import (
	"github.com/asdine/storm"
	"go.etcd.io/bbolt"
	"log"
)

type DB struct {
	*storm.DB
}

func newDB(path string) *DB {
	db, err := storm.Open(path)

	if err != nil {
		log.Fatalf("%s %s", err, path)
	}
	return &DB{db}
}

func (db *DB) bucketInit(data interface{}) {
	err := db.Init(data)
	switch err {
	case bbolt.ErrBucketExists:
		db.ReIndex(data)
	case nil:
		return
	default:
		log.Println(err)
	}
}

func (db *DB) init() {
	db.bucketInit(&Exceptions{})
	db.bucketInit(&Subtitles{})
}
