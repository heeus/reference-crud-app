/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import "fmt"

//LightDriver s.e.
type LightDriver struct {
	logger *Logger
}

//Name s.e.
func (d *LightDriver) Name() string {
	return "Light driver"
}

//Info s.e.
func (d *LightDriver) Info() string {
	return "Light driver"
}

//Init s.e.
func (d *LightDriver) Init(args map[string]string) error {
	fmt.Println("light driver initialized")
	return nil
}

//Free s.e.
func (d *LightDriver) Free() error {
	fmt.Println("light driver freed")
	return nil
}

//Read s.e.
func (d *LightDriver) Clean(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Read s.e.
func (d *LightDriver) Read(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Insert s.e.
func (d *LightDriver) Insert(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Update s.e.
func (d *LightDriver) Update(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Scan s.e.
func (d *LightDriver) Scan(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Delete s.e.
func (d *LightDriver) Delete(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}
