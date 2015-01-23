package main

import (
    "fmt"
    "time"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

const (
    DATABASE             = "Tags"
    TAG_COLLECTION       = "Tag"
    TAG_COUNT_COLLECTION = "TagCount"
)

var session *mgo.Session = nil

type Tag struct {
    Name string `bson:"name,omitempty"`
}

type TagCount struct {
    Name  string    `bson:"name,omitempty"`
    Date  time.Time `bson:"date,omitempty"`
    Count int       `bson:"count,omitempty"`
}

func (t *Tag) Hashed() string {
    return "#" + t.Name
}

func startSession(dbUrl string) error {
    msession, err := mgo.Dial(dbUrl)
    if err != nil {
        return err
    }

    session = msession
    return nil
}

func AddTagsToTrack(hashtags []string) error {
    if len(hashtags) == 0 {
        return nil
    }

    lsession := session.Copy()
    col := lsession.DB(DATABASE).C(TAG_COLLECTION)
    tags := make([]interface{}, len(hashtags))
    for i, tagName := range hashtags {
        tags[i] = Tag{Name: tagName}
    }

    return col.Insert(tags...)
}

func AddTagCounts(tagCounts []TagCount) error {
    lsession := session.Copy()
    col := lsession.DB(DATABASE).C(TAG_COUNT_COLLECTION)

    tagCountsInt := make([]interface{}, len(tagCounts))
    for i, tag := range tagCounts {
        tagCountsInt[i] = tag
    }

    return col.Insert(tagCountsInt...)
}

func GetTagsToTrack() (tags []Tag) {
    lsession := session.Copy()
    col := lsession.DB(DATABASE).C(TAG_COLLECTION)
    tags = make([]Tag, 0)
    col.Find(nil).Sort("name").All(&tags)
    return
}

func GetTagCountFor(delta time.Duration) []TagCount {
    since := time.Now().Add(-delta)

    pipeline := []bson.M{
        bson.M{
            "$match": bson.M{
                "count": bson.M{"$gt": 0},
                "date":  bson.M{"$gt": since},
            },
        },

        bson.M{
            "$group": bson.M{
                "_id":   "$name",
                "count": bson.M{"$sum": "$count"},
            },
        },

        bson.M{"$sort": bson.M{"count": -1}},

        bson.M{
            "$project": bson.M{
                "_id":   0,
                "name":  "$_id",
                "count": 1,
            },
        },
    }

    tagCounts := make([]TagCount, 0)

    lsession := session.Copy()
    col := lsession.DB(DATABASE).C(TAG_COUNT_COLLECTION)
    pipe := col.Pipe(pipeline)

    if err := pipe.All(&tagCounts); err != nil {
        fmt.Println("Could not count tags.", err.Error())
    }

    return tagCounts
}

func DropAll() error {
    lsession := session.Copy()
    return lsession.DB(DATABASE).DropDatabase()
}

func DropCollection(collection string) error {
    lsession := session.Copy()
    return lsession.DB(DATABASE).C(collection).DropCollection()
}
