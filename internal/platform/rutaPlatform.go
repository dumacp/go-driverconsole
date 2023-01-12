package platform

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

const (
	RutasPlatformName = "ROUTES"
)

type Route struct {
	ID              string `json:"id"`
	DivipolCode     string `json:"divipolCode"`
	AuthorityCode   string `json:"authorityCode"`
	Code            string `json:"code"`
	CompanyID       string `json:"companyId"`
	Name            string `json:"name"`
	OrganizationID  string `json:"organizationId"`
	EmpresaID       string `json:"empresa"`
	Itinerarys      []string
	StartPoints     map[string][2]float64
	EndPoints       map[string][2]float64
	InitTimeService time.Time
}

func (r *Route) Save(db *bbolt.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if err := db.Update(func(tx *bbolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(RutasPlatformName))
		if err != nil {
			return err
		}
		data, err := json.Marshal(r)
		if err != nil {
			return err
		}
		if err := bk.Put([]byte(r.ID), data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
