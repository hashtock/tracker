package core

import (
	"time"
)

type Tag struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
}

type TagCount struct {
	Name  string    `bson:"name,omitempty" json:"name,omitempty"`
	Date  time.Time `bson:"date,omitempty" json:"-"`
	Count int       `bson:"count,omitempty" json:"count,omitempty"`
}

type Count struct {
	Date  time.Time `bson:"date,omitempty" json:"date"`
	Count int       `bson:"count,omitempty" json:"count"`
}

type TagCountTrend struct {
	Name   string  `bson:"name,omitempty" json:"name,omitempty"`
	Counts []Count `bson:"counts,omitempty" json:"counts"`
}
