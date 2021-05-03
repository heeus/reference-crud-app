/*
 * Copyright (c) 2019-present Heeus authors
 */

package main

import "github.com/heeus/reference-crud-app/service"

func main() {
	s := service.Service{}

	if err := s.Init(); err != nil {
		panic(err)
	}

	s.Start()
}
