package storage

import (
    "fmt"
    "log"
    "time"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "github.com/hashtock/tracker/conf"
    "github.com/hashtock/tracker/core"
)

const (
    DATABASE             = "Tags"
    TAG_COLLECTION       = "Tag"
    TAG_COUNT_COLLECTION = "TagCount"
)

var session *mgo.Session = nil

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
    defer lsession.Close()
    col := lsession.DB(DATABASE).C(TAG_COLLECTION)

    for _, tagName := range hashtags {
        tag := core.Tag{Name: tagName}
        if _, err := col.Upsert(tag, tag); err != nil {
            return err
        }
    }

    return nil
}

func AddTagCounts(tagCounts []core.TagCount) error {
    if len(tagCounts) == 0 {
        return nil
    }

    lsession := session.Copy()
    defer lsession.Close()
    col := lsession.DB(DATABASE).C(TAG_COUNT_COLLECTION)

    var lastErr error = nil
    for _, tag := range tagCounts {
        selector := core.TagCount{
            Name: tag.Name,
            Date: tag.Date,
        }

        update_with := bson.M{
            "$inc": bson.M{"count": tag.Count},
        }
        _, err := col.Upsert(selector, update_with)

        if err != nil {
            lastErr = err
        }
    }

    return lastErr
}

func GetTagsToTrack() (tags []core.Tag) {
    lsession := session.Copy()
    defer lsession.Close()
    col := lsession.DB(DATABASE).C(TAG_COLLECTION)
    tags = make([]core.Tag, 0)
    col.Find(nil).Sort("name").All(&tags)
    return
}

func GetTagCountForLast(delta time.Duration) []core.TagCount {
    since := time.Now().Add(-delta)

    return GetTagCount(since, time.Time{})
}

func GetTagDetailedCountForLast(delta time.Duration) []core.TagCountTrend {
    since := time.Now().Add(-delta)

    return GetTagCountDetailed(since, time.Time{})
}

func GetTagCountDetailed(since, until time.Time) []core.TagCountTrend {
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

        bson.M{"$sort": bson.M{"date": 1}},

        bson.M{
            "$group": bson.M{
                "_id": "$name",
                "counts": bson.M{
                    "$push": bson.M{
                        "date":  "$date",
                        "count": "$count",
                    },
                },
            },
        },

        bson.M{"$sort": bson.M{"_id": 1}},

        bson.M{
            "$project": bson.M{
                "_id":    0,
                "name":   "$_id",
                "counts": 1,
            },
        },
    }

    tagCounts := make([]core.TagCountTrend, 0)

    lsession := session.Copy()
    defer lsession.Close()

    col := lsession.DB(DATABASE).C(TAG_COUNT_COLLECTION)
    pipe := col.Pipe(pipeline)

    if err := pipe.All(&tagCounts); err != nil {
        fmt.Println("Could not count tags.", err.Error())
    }

    return tagCounts
}

func GetTagCount(since, until time.Time) []core.TagCount {
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

    tagCounts := make([]core.TagCount, 0)

    lsession := session.Copy()
    defer lsession.Close()

    col := lsession.DB(DATABASE).C(TAG_COUNT_COLLECTION)
    pipe := col.Pipe(pipeline)

    if err := pipe.All(&tagCounts); err != nil {
        fmt.Println("Could not count tags.", err.Error())
    }

    return tagCounts
}

func DropAll() error {
    lsession := session.Copy()
    defer lsession.Close()

    return lsession.DB(DATABASE).DropDatabase()
}

func DropCollection(collection string) error {
    lsession := session.Copy()
    defer lsession.Close()

    return lsession.DB(DATABASE).C(collection).DropCollection()
}
