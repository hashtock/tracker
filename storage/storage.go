package storage

import (
    "fmt"
    "log"
    "time"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "github.com/hashtock/tracker/conf"
)

const (
    DATABASE             = "Tags"
    TAG_COLLECTION       = "Tag"
    TAG_COUNT_COLLECTION = "TagCount"
)

var session *mgo.Session = nil

type Tag struct {
    Name string `bson:"name,omitempty" json:"name,omitempty"`
}

type TagCount struct {
    Name  string    `bson:"name,omitempty" json:"name,omitempty"`
    Date  time.Time `bson:"date,omitempty" json:"-"`
    Count int       `bson:"count,omitempty" json:"count,omitempty"`
}

func init() {
    cfg := conf.GetConfig()

    if err := startSession(cfg.General.DB); err != nil {
        log.Fatalln("Could not connect to DB.", err.Error())
    }
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

    for _, tagName := range hashtags {
        tag := Tag{Name: tagName}
        if _, err := col.Upsert(tag, tag); err != nil {
            return err
        }
    }

    return nil
}

func AddTagCounts(tagCounts []TagCount) error {
    if len(tagCounts) == 0 {
        return nil
    }

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

func GetTagCountForLast(delta time.Duration) []TagCount {
    since := time.Now().Add(-delta)

    return GetTagCount(since, time.Time{})
}

func GetTagCount(since, until time.Time) []TagCount {
    query := bson.M{
        "count": bson.M{"$gt": 0},
        "date": bson.M{
            "$gte": since,
            "$lt":  until,
        },
    }

    if since.IsZero() && until.IsZero() {
        delete(query, "date")
    } else if since.IsZero() {
        delete(query["date"].(bson.M), "$gte")
    } else if until.IsZero() {
        delete(query["date"].(bson.M), "$lt")
    }

    pipeline := []bson.M{
        bson.M{"$match": query},

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
