package tinybeans

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"

	"github.com/bp1222/photoSync/database"
	tb "github.com/bp1222/tinybeans-api/go-client"
	log "github.com/sirupsen/logrus"
)

const (
	KEY_FILE = "./tinybeans.key"
)

type Tinybeans interface {
	GetJournalEntries(journal, fetchSize int64, last int64) (*tb.Entries, error)
	GetJournalEntriesSince(journal, fetchSize int64, since int64) (*tb.Entries, error)
}

type tinybeans struct {
	api       *tb.APIClient
	db        database.Database
	authtoken string
}

func InitTinybeans(db database.Database) Tinybeans {
	configuration := tb.NewConfiguration()

	if os.Getenv("DEBUG") == "true" {
		proxyUrl, _ := url.Parse("http://localhost:8080")
		configuration.HTTPClient = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		}
	}

	t := tinybeans{
		db:        db,
		api:       tb.NewAPIClient(configuration),
		authtoken: readStoredKey(),
	}

	if err := t.authenticate(); err != nil {
		log.Fatal("unable to authenticate to tinybeans")
	}

	return &t
}

func saveStoredKey(key string) error {
	if err := os.WriteFile(KEY_FILE, []byte(key), fs.FileMode(0700)); err != nil {
		return err
	}
	return nil
}

func readStoredKey() string {
	data, err := os.ReadFile(KEY_FILE)
	if err != nil {
		return ""
	}

	return string(data)
}

func (t *tinybeans) getContext() (context.Context, error) {
	if t.authtoken == "" {
		return nil, fmt.Errorf("getContext: unauthenticated use")
	}

	return context.WithValue(
		context.Background(),
		tb.ContextAPIKeys,
		map[string]tb.APIKey{
			"ApiKeyAuth": {
				Key: t.authtoken,
			},
		},
	), nil
}

func (t *tinybeans) checkAuth() bool {
	ctx, err := t.getContext()
	if err != nil {
		return false
	}

	resp, _, err := t.api.AuthApi.UsersMe(ctx).Execute()
	if err != nil {
		log.Fatal("Error when calling `AuthApi.UsersMe``: %v", err)
		return false
	}

	return resp.GetStatus() == "ok"
}

func (t *tinybeans) authenticate() error {
	if t.authtoken != "" && t.checkAuth() {
		log.Infof("Existing AuthToken is still Valid")
		return nil
	}
	log.Infof("Existing AuthToken is Invalid, acquiring new one")

	authenticateRequest := tb.AuthenticateRequst{}

	authenticateRequest.SetClientId("d324d503-0127-4a85-a547-d9f2439ffeae") // Web UI id, sorry analytics.
	authenticateRequest.SetUsername(os.Getenv("TINYBEANS_USERNAME"))
	authenticateRequest.SetPassword(os.Getenv("TINYBEANS_PASSWORD"))

	resp, _, err := t.api.AuthApi.Login(context.Background()).AuthenticateRequst(authenticateRequest).Execute()
	if err != nil {
		log.Fatal("Error when calling `AuthApi.Login``: %v", err)
		return err
	}

	log.Infof("Response from `AuthApi.Login`: %v", resp)

	if err := saveStoredKey(resp.GetAccessToken()); err != nil {
		log.Fatal("unable to save auth-token to file cache")
	}

	t.authtoken = resp.GetAccessToken()
	return nil
}

type JournalEntryCallback func(tb.ApiJournalEntriesRequest) tb.ApiJournalEntriesRequest

func (t *tinybeans) GetJournalEntries(journal, fetchSize int64, last int64) (*tb.Entries, error) {
	return t.getJournalEntries(journal, fetchSize, func(req tb.ApiJournalEntriesRequest) tb.ApiJournalEntriesRequest {
		return req.Last(last)
	})
}

func (t *tinybeans) GetJournalEntriesSince(journal, fetchSize int64, since int64) (*tb.Entries, error) {
	return t.getJournalEntries(journal, fetchSize, func(req tb.ApiJournalEntriesRequest) tb.ApiJournalEntriesRequest {
		return req.Since(since)
	})
}

func (t *tinybeans) getJournalEntries(journal, fetchSize int64, cb JournalEntryCallback) (*tb.Entries, error) {
	ctx, err := t.getContext()
	if err != nil {
		return nil, err
	}

	req := t.api.JournalsApi.JournalEntries(ctx, journal).FetchSize(fetchSize).IdsOnly(1)

	if cb != nil {
		req = cb(req)
	}

	resp, _, err := req.Execute()
	if err != nil {
		log.Fatalf("Error when calling `JournalsApi.JournalEntries``: %v", err)
		return nil, err
	}

	log.Infof("Response from `JournalsApi.JournalEntries`: (%d) new entries to process", len(resp.Entries))
	if num, has := resp.GetNumEntriesRemainingOk(); has && *num > 0 {
		log.Infof("\t And there are (%d) more to process", *num)
	}

	return resp, nil
}
