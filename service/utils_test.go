/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mapArgs(t *testing.T) {
	{
		s := []string{"-nop", "--ll", "1", "-todo", "-fruit", "banana", "-doit", "no", "wrong"}

		args := mapArgs(s)

		assert.Equal(t, len(args), 5)
		assert.Equal(t, args["-nop"], "true")
		assert.Equal(t, args["--ll"], "1")
		assert.Equal(t, args["-todo"], "true")
		assert.Equal(t, args["-fruit"], "banana")
		assert.Equal(t, args["-doit"], "no")
	}
}
