package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	tb "github.com/bp1222/tinybeans-api/go"
	"gorm.io/gorm"
)

const (
	KEY_FILE = "./tinybeans.key"
)

type Tinybeans interface {
	Authenticate() error
	GetJournals() (*tb.Journals, error)
	GetJournalEntries(journal, fetchSize int64, last int64) (*tb.Entries, error)
	GetJournalEntriesSince(journal, fetchSize int64, since int64) (*tb.Entries, error)
	GetMostRecentEntry(journal int64) int64
	IsLikedBy(int64, int64) bool
}

type tinybeans struct {
	db        *gorm.DB
	api       *tb.APIClient
	authtoken *string
}

var _ Tinybeans = &tinybeans{}

func InitTinybeans(db *gorm.DB) Tinybeans {
	configuration := tb.NewConfiguration()

	/* Local mitmproxy for debugging
	proxyUrl, _ := url.Parse("http://localhost:8080")
	configuration.HTTPClient = &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}
	*/

	t := tinybeans{
		db:        db,
		api:       tb.NewAPIClient(configuration),
		authtoken: readStoredKey(),
	}

	return &t
}

func saveStoredKey(key *string) error {
	if err := os.WriteFile(KEY_FILE, []byte(*key), fs.FileMode(0700)); err != nil {
		return err
	}
	return nil
}

func readStoredKey() *string {
	data, err := os.ReadFile(KEY_FILE)
	if err != nil {
		return nil
	}

	ret := string(data)
	return &ret
}

func (t *tinybeans) getContext() (context.Context, error) {
	if t.authtoken == nil {
		return nil, fmt.Errorf("getContext: unauthenticated use")
	}

	return context.WithValue(
		context.Background(),
		tb.ContextAPIKeys,
		map[string]tb.APIKey{
			"Authorization": {
				Key: *t.authtoken,
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
		fmt.Fprintf(os.Stderr, "Error when calling `AuthApi.UsersMe``: %v\n", err)
		return false
	}

	return resp.GetStatus() == "ok"
}

func (t *tinybeans) Authenticate() error {
	if t.authtoken != nil && t.checkAuth() {
		return nil
	}

	authenticateRequst := tb.AuthenticateRequst{}

	authenticateRequst.SetClientId("d324d503-0127-4a85-a547-d9f2439ffeae") // Web UI id, sorry analytics.
	authenticateRequst.SetUsername(os.Getenv("TINYBEANS_USERNAME"))
	authenticateRequst.SetPassword(os.Getenv("TINYBEANS_PASSWORD"))

	resp, _, err := t.api.AuthApi.Login(context.Background()).AuthenticateRequst(authenticateRequst).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthApi.Login``: %v\n", err)
		return err
	}

	fmt.Fprintf(os.Stdout, "Response from `AuthApi.Login`: %v\n", resp)

	t.authtoken = resp.AccessToken

	if err := saveStoredKey(t.authtoken); err != nil {
		return fmt.Errorf("unable to save auth-token to file cache")
	}

	return nil
}

func (t *tinybeans) GetJournals() (*tb.Journals, error) {
	ctx, err := t.getContext()
	if err != nil {
		return nil, err
	}

	resp, _, err := t.api.JournalsApi.Journals(ctx).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `JournalsApi.Journals``: %v\n", err)
		return nil, err
	}

	fmt.Fprintf(os.Stdout, "Response from `JournalsApi.Journals`: %v\n", resp)

	return &resp, nil
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
		fmt.Fprintf(os.Stderr, "Error when calling `JournalsApi.JournalEntries``: %v\n", err)
		return nil, err
	}

	fmt.Fprintf(os.Stdout, "Response from `JournalsApi.JournalEntries`: %v\n", resp)

	return &resp, nil
}

func (t *tinybeans) IsLikedBy(entry_id, user_id int64) bool {
	var l Like
	t.db.Where(Like{
		EntryId: entry_id,
		UserId:  user_id,
	}).First(&l)
	return l.EntryId == entry_id
}

func (t *tinybeans) GetMostRecentEntry(journal int64) int64 {
	var e Entry
	t.db.Order("timestamp desc").First(&e)
	return e.Timestamp
}
