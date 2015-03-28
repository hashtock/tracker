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

func (m *mgoCounter) collection(collectionName string) *mgo.Collection {
	lsession := m.session.Copy()
	col := lsession.DB(dbName).C(collectionName)
	return col
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
	query := baseQuery(since, until)

	pipeline := []bson.M{
		showOnlyComplete(m.minumAgeOfCount),

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

	col := m.collection(tagCountCollectionName)
	defer col.Database.Session.Close()

	pipe := col.Pipe(pipeline)
	err = pipe.All(&tagCounts)

	return tagCounts, err
}

func (m *mgoCounter) trendsPipeline(tag string, since, until time.Time, sampling core.Sampling) []bson.M {
	query := baseQuery(since, until)
	if tag != "" {
		query["name"] = tag
	}

	pipeMatch := []bson.M{
		showOnlyComplete(m.minumAgeOfCount),
		bson.M{"$match": query},
		sortByDate(),
	}

	pipeResample := resamplePipeline(sampling)

	pipeGroupResults := []bson.M{
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
				"date":   1,
			},
		},
	}

	pipeline := pipeMatch
	pipeline = append(pipeline, pipeResample...)
	pipeline = append(pipeline, pipeGroupResults...)

	return pipeline
}

func (m *mgoCounter) TagTrends(tag string, since, until time.Time, sampling core.Sampling) (tagCounts core.TagCountTrend, err error) {
	col := m.collection(tagCountCollectionName)
	defer col.Database.Session.Close()

	pipeline := m.trendsPipeline(tag, since, until, sampling)

	pipe := col.Pipe(pipeline)
	err = pipe.One(&tagCounts)
	return
}

func (m *mgoCounter) Trends(since, until time.Time) (tagCounts []core.TagCountTrend, err error) {
	col := m.collection(tagCountCollectionName)
	defer col.Database.Session.Close()

	pipeline := m.trendsPipeline("", since, until, core.SamplingRaw)

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
