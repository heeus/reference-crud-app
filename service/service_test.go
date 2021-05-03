/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import (
	"testing"
)

func Test_Multyinit(t *testing.T) {
	{
		go initCassandraDriver(t)
		go initCassandraDriver(t)
	}
}

func initCassandraDriver(t *testing.T) {
	t.Log("Initiating cassandra driver")

	d := CasandraDriver{}
	args := map[string]string{}

	err := d.Init(args)

	if err != nil {
		t.Fatal(err)
	}
}
