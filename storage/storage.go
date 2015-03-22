package storage

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/hashtock/tracker/core"
)

const (
	dbName                 = "Tags"
	tagCollectionName      = "Tag"
	tagCountCollectionName = "TagCount"
)

type mgoCounter struct {
	session         *mgo.Session
	db              string
	minumAgeOfCount time.Duration
}

func NewMongoCounter(dbURL string, minumAgeOfCount time.Duration) (*mgoCounter, error) {
	msession, err := mgo.Dial(dbURL)
	if err != nil {
		return nil, err
	}

	return &mgoCounter{
		db:              dbURL,
		session:         msession,
		minumAgeOfCount: minumAgeOfCount,
	}, nil
}

func (m *mgoCounter) showOnlyComplete() bson.M {
	return bson.M{
		"$match": bson.M{
			"date": bson.M{
				"$lt": time.Now().Add(-m.minumAgeOfCount),
			},
		},
	}
}

func (m *mgoCounter) Tags() (tags []core.Tag, err error) {
	lsession := m.session.Copy()
	defer lsession.Close()
	col := lsession.DB(dbName).C(tagCollectionName)
	tags = make([]core.Tag, 0)
	col.Find(nil).Sort("name").All(&tags)
	return
}

func (m *mgoCounter) Counts(since, until time.Time) (tagCounts []core.TagCount, err error) {
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
		m.showOnlyComplete(),

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

	lsession := m.session.Copy()
	defer lsession.Close()

	col := lsession.DB(dbName).C(tagCountCollectionName)
	pipe := col.Pipe(pipeline)
	err = pipe.All(&tagCounts)

	return tagCounts, err
}

func (m *mgoCounter) Trends(since, until time.Time) (tagCounts []core.TagCountTrend, err error) {
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
		m.showOnlyComplete(),

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

	lsession := m.session.Copy()
	defer lsession.Close()

	col := lsession.DB(dbName).C(tagCountCollectionName)
	pipe := col.Pipe(pipeline)
	err = pipe.All(&tagCounts)

	return tagCounts, err
}

func (m *mgoCounter) AddTag(tagName string) error {
	if tagName == "" {
		return nil
	}

	lsession := m.session.Copy()
	defer lsession.Close()
	col := lsession.DB(dbName).C(tagCollectionName)

	tag := core.Tag{Name: tagName}
	if _, err := col.Upsert(tag, tag); err != nil {
		return err
	}

	return nil
}

func (m *mgoCounter) AddTagCounts(tagCounts []core.TagCount) error {
	if len(tagCounts) == 0 {
		return nil
	}

	lsession := m.session.Copy()
	defer lsession.Close()
	col := lsession.DB(dbName).C(tagCountCollectionName)

	var lastErr error
	for _, tag := range tagCounts {
		selector := core.TagCount{
			Name: tag.Name,
			Date: tag.Date,
		}

		updateWith := bson.M{
			"$inc": bson.M{"count": tag.Count},
		}
		_, err := col.Upsert(selector, updateWith)

		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (m *mgoCounter) RemoveAll() error {
	lsession := m.session.Copy()
	defer lsession.Close()

	return lsession.DB(dbName).DropDatabase()
}

func (m *mgoCounter) RemoveCounts() error {
	lsession := m.session.Copy()
	defer lsession.Close()

	return lsession.DB(dbName).C(tagCountCollectionName).DropCollection()
}
