package trackinglog

import (
	"appengine"
	"appengine/datastore"
	"appengine/delay"
	"encoding/json"
	"code.google.com/p/go-uuid/uuid"
	"github.com/knightso/base/errors"
	"github.com/knightso/base/gae/ds"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"net/http"
)

const KIND_USER_AGENT = "UserAgent"
const KIND_TRACKING_LOG = "TrackingLog"

func init() {
	http.Handle("/", goji.DefaultMux)
	goji.Get("/tracking/*", saveTracking)
	goji.Get("/api/useragents", queryUserAgents)
	goji.Get("/api/useragents/:key", getUserAgent)
	goji.Get("/api/useragents/:key/trackinglogs", queryTrackingLogs)
}

type UserAgent struct {
	UserAgent string `json:"userAgent"`
	ds.Meta
}

type JsonUserAgent struct {
	*UserAgent
	Key *datastore.Key `json:"key"`
}

func NewUserAgentKey(c appengine.Context, ua string) *datastore.Key {
	return datastore.NewKey(c, KIND_USER_AGENT, ua, 0, nil)
}

type TrackingLog struct {
	UserAgent string `json:"userAgent"`
	URL       string `json:"url"`
	ds.Meta
}

func NewTrackingLogKey(c appengine.Context, uuid string) *datastore.Key {
	return datastore.NewKey(c, KIND_TRACKING_LOG, uuid, 0, nil)
}

func saveTracking(gojic web.C, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var t TrackingLog
	t.Key = NewTrackingLogKey(c, uuid.New())
	t.UserAgent = r.UserAgent()
	t.URL = r.URL.String()

	DoSaveTracking.Call(c, &t)
}

var DoSaveTracking = delay.Func("DoSaveTracking", func(c appengine.Context, t *TrackingLog) error {

	err := datastore.RunInTransaction(c, func(c appengine.Context) error {

		uaKey := NewUserAgentKey(c, t.UserAgent)

		var ua UserAgent
		if err := ds.Get(c, uaKey, &ua); err != nil {
			if errors.Root(err) == datastore.ErrNoSuchEntity {
				ua.Key = uaKey
				ua.UserAgent = t.UserAgent
			} else {
				return errors.WrapOr(err)
			}
		}
		if err := ds.Put(c, &ua); err != nil {
			return errors.WrapOr(err)
		}

		if err := ds.Put(c, t); err != nil {
			return errors.WrapOr(err)
		}

		return nil
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		return errors.WrapOr(err)
	}
	return nil
})

func getUserAgent(gojic web.C, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	uaKey, err := datastore.DecodeKey(gojic.URLParams["key"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ua UserAgent
	if err := ds.Get(c, uaKey, &ua); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jua := JsonUserAgent {
		UserAgent: &ua,
		Key: ua.Key,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsn, err := json.MarshalIndent(jua, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsn); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func queryUserAgents(gojic web.C, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery(KIND_USER_AGENT).Order("-UpdatedAt")

	var uas []*UserAgent
	if err := ds.ExecuteQuery(c, q, &uas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	juas := make([]*JsonUserAgent, len(uas))
	for i, ua := range uas {
		jua := JsonUserAgent{
			UserAgent: ua,
			Key: ua.Key,
		}
		juas[i] = &jua
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsn, err := json.MarshalIndent(juas, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsn); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func queryTrackingLogs(gojic web.C, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	uaKey, err := datastore.DecodeKey(gojic.URLParams["key"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ua UserAgent
	if err := ds.Get(c, uaKey, &ua); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := datastore.NewQuery(KIND_TRACKING_LOG).Filter("UserAgent =", ua.UserAgent).Order("-CreatedAt").Limit(500)

	var logs []*TrackingLog
	if err := ds.ExecuteQuery(c, q, &logs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsn, err := json.MarshalIndent(logs, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsn); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
