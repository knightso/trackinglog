package trackinglog

import (
	"appengine"
	"appengine/datastore"
	"appengine/delay"
	"code.google.com/p/go-uuid/uuid"
	"github.com/knightso/base/errors"
	"github.com/knightso/base/gae/ds"
	"net/http"
)

const KIND_USER_AGENT = "UserAgent"
const KIND_TRACKING_LOG = "TrackingLog"

type UserAgent struct {
	UserAgent string `json:"userAgent"`
	ds.Meta
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

func init() {
	http.HandleFunc("/tracking/", saveTracking)
}

func saveTracking(w http.ResponseWriter, r *http.Request) {
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
