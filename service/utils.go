/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func initBoolParam(args map[string]string, envName, attrName string, defValue bool) bool {
	if envName != "" {
		if v, exists := os.LookupEnv(envName); exists && v != "" {
			if r, err := strconv.ParseBool(v); err == nil {
				return r
			}
		}
	}

	if attrName != "" {
		if v, exists := args[attrName]; exists && v != "" {
			if r, err := strconv.ParseBool(v); err == nil {
				return r
			}
		}
	}

	return defValue
}

func initIntParam(args map[string]string, envName, attrName string, defValue int64) int64 {
	if envName != "" {
		if v, exists := os.LookupEnv(envName); exists && v != "" {
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				return i
			}
		}
	}

	if attrName != "" {
		if v, exists := args[attrName]; exists && v != "" {
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				return i
			}
		}
	}

	return defValue
}

func initStringParam(args map[string]string, envName, attrName, defValue string) string {

	if v, exists := os.LookupEnv(envName); exists && v != "" {
		return v
	}

	if v, exists := args[attrName]; exists && v != "" {
		return v
	}

	return defValue
}

func buildRequest(r *http.Request) (req *DBRequest, err error) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &req)

	return req, err
}

func mapArgs(args []string) map[string]string {
	var p interface{} = nil
	mappedArgs := map[string]string{}

	for _, s := range args {
		if strings.HasPrefix(s, "-") || strings.HasPrefix(s, "--") {
			p = s
			mappedArgs[p.(string)] = "true"
		} else {
			if p != nil {
				mappedArgs[p.(string)] = s
				p = nil
			}
		}
	}

	return mappedArgs
}

func checkMethodAllowed(r *http.Request, methods []string) error {
	if len(methods) == 0 {
		return nil
	}

	b := false
	m := r.Method

	for _, method := range HTTPMethods {
		if m == method {
			b = true
		}
	}

	if !b {
		return fmt.Errorf("no such method %v in HTTP methods list", m)
	}

	b = false

	for _, method := range methods {
		if m == method {
			b = true
		}
	}

	if !b {
		return fmt.Errorf("method %v not allowed", m)
	}

	return nil
}

func buildKey(pkey map[string]interface{}, ckey map[string]interface{}) (string, error) {
	key := ""

	for _, v := range pkey {
		key += v.(string)
	}

	for _, v := range ckey {
		key += v.(string)
	}

	return key, nil
}
