/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import "fmt"

//NopDriver s.e.
type NopDriver struct {
	logger *Logger
}

//Name s.e.
func (d *NopDriver) Name() string {
	return "Nop mode"
}

//Info s.e.
func (d *NopDriver) Info() string {
	return "Nop mode is On"
}

//Init s.e.
func (d *NopDriver) Init(args map[string]string) error {
	fmt.Println("nop mode active")
	return nil
}

//Free s.e.
func (d *NopDriver) Free() error {
	fmt.Println("nop mode stoped")
	return nil
}

//Read s.e.
func (d *NopDriver) Clean(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Read s.e.
func (d *NopDriver) Read(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Insert s.e.
func (d *NopDriver) Insert(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Update s.e.
func (d *NopDriver) Update(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Scan s.e.
func (d *NopDriver) Scan(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}

//Delete s.e.
func (d *NopDriver) Delete(r *DBRequest) *DBResponse {
	return &DBResponse{Status: 200}
}
