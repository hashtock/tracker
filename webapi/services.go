package webapi

import (
	"net/http"

	"github.com/hashtock/tracker/core"
)

const UserContextKey = "user"

type counterService struct {
	counter    core.CountReaderWritter
	serializer Serializer
}

func (c *counterService) allTags(rw http.ResponseWriter, req *http.Request) {
	tags, err := c.counter.Tags()
	if err != nil {
		c.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	c.serializer.JSON(rw, http.StatusOK, tags)
}

func (c *counterService) addTag(rw http.ResponseWriter, req *http.Request) {
	if adminErr := isAdmin(req); adminErr != nil {
		c.serializer.JSON(rw, http.StatusForbidden, adminErr.Error())
		return
	}

	name := req.URL.Query().Get(":name")

	if err := c.counter.AddTag(name); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusCreated)
	}
}

func (c *counterService) counts(rw http.ResponseWriter, req *http.Request) {
	since, until, err := parseQuery(req.URL.Query())
	if err != nil {
		c.serializer.JSON(rw, http.StatusBadRequest, err.Error())
		return
	}

	counts, err := c.counter.Counts(since, until)
	if err != nil {
		c.serializer.JSON(rw, http.StatusBadRequest, err.Error())
		return
	}

	c.serializer.JSON(rw, http.StatusOK, counts)
}

func (c *counterService) trends(rw http.ResponseWriter, req *http.Request) {
	since, until, err := parseQuery(req.URL.Query())
	if err != nil {
		c.serializer.JSON(rw, http.StatusBadRequest, err.Error())
		return
	}

	trends, err := c.counter.Trends(since, until)
	if err != nil {
		c.serializer.JSON(rw, http.StatusBadRequest, err.Error())
		return
	}

	mapped := make(map[string][]core.Count, len(trends))
	for _, trend := range trends {
		mapped[trend.Name] = trend.Counts
	}

	c.serializer.JSON(rw, http.StatusOK, mapped)
}

func (c *counterService) tagTrends(rw http.ResponseWriter, req *http.Request) {
	tag := req.URL.Query().Get(":name")

	if tag == "" {
		c.serializer.JSON(rw, http.StatusBadRequest, nil)
		return
	}

	since, until, err := parseQuery(req.URL.Query())
	if err != nil {
		c.serializer.JSON(rw, http.StatusBadRequest, err.Error())
		return
	}

	sampling, err := getSamplingFromQuery(req.URL.Query())
	if err != nil {
		c.serializer.JSON(rw, http.StatusBadRequest, err.Error())
		return
	}

	trends, err := c.counter.TagTrends(tag, since, until, sampling)
	if err != nil {
		c.serializer.JSON(rw, http.StatusNotFound, err.Error())
		return
	}

	c.serializer.JSON(rw, http.StatusOK, trends.Counts)
}
