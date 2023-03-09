package models

import (
	"database/sql"
	"sync"
	"time"

	"github.com/golang-sql/civil"
)

const (
	RetryFibonacciTimeFactor = 10 * time.Millisecond
	RetryFibonacciAmount     = 7
)

// Offer ...
type Offer struct {
	ItemID       string `json:"item_id"`
	ShortWebName sql.NullString
	Attributes   []*OfferAttribute
	Level        int
	Category1    sql.NullString
	Category2    sql.NullString
	Category3    sql.NullString
	Category4    sql.NullString
	Category5    sql.NullString
	Category6    sql.NullString
	Category7    sql.NullString
	Category8    sql.NullString
	Category9    sql.NullString
	Category10   sql.NullString
}

// OfferAttribute ...
type OfferAttribute struct {
	ItemID sql.NullString
	Key    sql.NullString
	Value  sql.NullString
}

// OfferImage ...
type OfferImage struct {
	ItemID    sql.NullString
	ImageLink sql.NullString
}

// OffersSyncMap ...
type OffersSyncMap struct {
	mx sync.RWMutex
	M  map[string]*Offer
}

// YourStruct ...
type YourStruct struct {
	ItemID                    string              `json:"item_id"`
	UpdatedDate               time.Time           `json:"updated_date"`
	Model                     string              `json:"model"`
	Items                     []string            `json:"items"`
	AdvancedParams            map[string][]string `json:"advanced_params"`
	AdvancedParamsCreatedDate time.Time           `json:"advanced_params_created_date"`
	VisenzeResponseErrors     []string            `json:"visenze_response_errors"`
}

// Task ...
type Task struct {
	Name        string
	CronPattern string
	IsActive    bool
	Func        func() error
}

type TemplateRequest struct {
	ServerFrom           string   `json:"server_from"`
	ServerFromUser       string   `json:"server_from_user"`
	ServerFromPass       string   `json:"server_from_pass"`
	ServerFromVersion    int      `json:"server_from_version"`
	ServerTo             string   `json:"server_to"`
	ServerToUser         string   `json:"server_to_user"`
	ServerToPass         string   `json:"server_to_pass"`
	ServerToVersion      int      `json:"server_to_version"`
	IsForce              bool     `json:"is_force"`
	IsUnsafeMem          bool     `json:"is_unsafe_mem"`
	IndicesForProcessRaw []string `json:"indexes_to_transfer"`
	IndicesForProcess    []string `json:"-"`
}

type TemplateResponse struct {
	Msg                 string   `json:"msg"`
	IndicesSpent        []string `json:"indices_spent"`
	IndicesInProcessing []string `json:"indices_in_processing"`
	IndicesForProcess   []string `json:"indices_for_process"`
}

type ScrollIdent struct {
	IndexNameFrom string    `json:"index_name_from" bson:"index_name_from"`
	IndexNameTo   string    `json:"index_name_to"   bson:"index_name_to"`
	ServerFrom    string    `json:"server_from"     bson:"server_from"`
	ServerTo      string    `json:"server_to"       bson:"server_to"`
	ScrollID      string    `json:"scroll_id"       bson:"scroll_id"`
	UpsertDate    time.Time `json:"upsert_date"     bson:"upsert_date"`
	IsDone        bool      `json:"is_done"         bson:"is_done"`
}

type QueryOptions struct {
	Query   string
	Timeout time.Duration
}

type CredentialsDB struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
}

type StatisticsPostgres struct {
	AdvertID        string
	AdvertName      string
	AdvertiserName  string
	CampaignName    string
	Views           int64
	Clicks          int64
	CTR             float64
	Orders          int64
	CR              float64
	GMV             int64
	ActionDate      civil.Date
	AdvertStartDate civil.Date
	AdvertEndDate   civil.Date
}
