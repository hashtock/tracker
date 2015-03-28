package storage

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/hashtock/tracker/core"
)

func showOnlyComplete(minumAgeOfCount time.Duration) bson.M {
	return bson.M{
		"$match": bson.M{
			"date": bson.M{
				"$lt": time.Now().Add(-minumAgeOfCount),
			},
		},
	}
}

func baseQuery(since, until time.Time) bson.M {
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

	return query
}

func sortByDate() bson.M {
	return bson.M{"$sort": bson.M{"date": 1}}
}

func resamplePipeline(sampling core.Sampling) []bson.M {

	seconds := sampling.Duration().Seconds()
	log.Println("Sampling:", int(sampling))
	log.Println("Seconds:", seconds)

	if seconds == 0 {
		return []bson.M{}
	}

	unix := time.Unix(0, 0)
	dateDiff := bson.M{"$subtract": []interface{}{"$date", unix}}

	// Round date to by number of seconds
	// Mongo does not allow to manipule dates easily
	// so some extra work was required
	roundDate := bson.M{
		"$add": []interface{}{
			bson.M{
				"$subtract": []bson.M{
					dateDiff,
					bson.M{
						"$mod": []interface{}{
							bson.M{
								"$subtract": []interface{}{"$date", unix},
							},
							seconds * 1000,
						},
					},
				},
			},
			unix,
		},
	}

	pipeResample := []bson.M{
		bson.M{
			"$project": bson.M{
				"name":  1,
				"count": 1,
				"date":  roundDate,
			},
		},

		// Resample using new date
		bson.M{
			"$group": bson.M{
				"_id":   bson.M{"name": "$name", "date": "$date"},
				"name":  bson.M{"$first": "$name"},
				"date":  bson.M{"$first": "$date"},
				"count": bson.M{"$avg": "$count"},
			},
		},

		// Keep data sorted
		sortByDate(),
	}

	return pipeResample
}
