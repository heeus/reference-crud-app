/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import (
	"encoding/json"
)

//DBDriver s.e.
type DBDriver interface {
	Init(args map[string]string) error
	Free() error
	Clean(r *DBRequest) *DBResponse
	Read(r *DBRequest) *DBResponse
	Insert(r *DBRequest) *DBResponse
	Update(r *DBRequest) *DBResponse
	Scan(r *DBRequest) *DBResponse
	Delete(r *DBRequest) *DBResponse
	Name() string
	Info() string
}

//ViewView  s.e.
type ViewView struct {
	ViewType     string
	PartitionKey map[string]interface{}
	ClusterKey   map[string]interface{}
}

//ViewMod s.e.
type ViewMod struct {
	ViewView
	Values map[string]interface{}
}

//DBRequest s.e.
type DBRequest struct {
	Partition int64
	ViewViews []ViewView
	ViewMods  []ViewMod
}

//DBResponse s.e.
type DBResponse struct {
	Status  int64
	Error   string
	Records []*Record
}

//Record s.e.
type Record struct {
	Key     string
	Values  map[string]interface{}
	Version int
}

func (r *DBResponse) stringify() []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		return []byte("unable to marshal response: " + err.Error())
	}
	return bytes
}
