package main

import (
	"log"

	"github.com/jinzhu/gorm"
)

func testutilGetDB() *gorm.DB {
	db, err := gorm.Open("mysql", "all:@/biocleaner?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}
