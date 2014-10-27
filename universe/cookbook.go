// Author:: Jonathan Hartman (<j@p4nt5.com>)
//
// Copyright (C) 2014, Jonathan Hartman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package goulash implements an API client for the Chef Supermarket.

This file implements a struct for a cookbook as described by a Berkshelf-style
universe endpoint, e.g.

https://supermarket.getchef.com/universe =>

{
	"chef": {
		"0.12.0": {
			"location_type": "opscode",
			"location_path": "https://supermarket.getchef.com/api/v1",
			"download_url": "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.12.0/download",
			"dependencies": {
				"runit":">= 0.0.0",
				"couchdb":">= 0.0.0",
				...
			}
		},
		"0.20.0": {
			"location_type": "opscode",
			"location_path": "https://supermarket.getchef.com/api/v1",
			"download_url": "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.20.0/download",
			"dependencies": {
				"zlib":">= 0.0.0",
				"xml": ">= 0.0.0",
				...
			}
		},
		...
	},
	...
*/
package universe

import (
	"reflect"
)

// Cookbook is just a map of version strings to Version structs
type Cookbook struct {
	Name     string
	Versions map[string]*CookbookVersion
}

// NewCookbook generates an empty Cookbook struct.
func NewCookbook() (c *Cookbook) {
	c = new(Cookbook)
	c.Versions = map[string]*CookbookVersion{}
	return
}

// Empty checks whether a Cookbook struct has been populated with anything or
// still holds all the base defaults.
func (c *Cookbook) Empty() (empty bool) {
	empty = true
	if c == nil {
		return
	}
	r := reflect.ValueOf(c).Elem()
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		switch f.Kind() {
		case reflect.String:
			if f.String() != "" {
				empty = false
				break
			}
		case reflect.Map:
			for _, k := range f.MapKeys() {
				method := f.MapIndex(k).MethodByName("Empty")
				if !method.Call([]reflect.Value{})[0].Bool() {
					empty = false
					return
				}
			}
		}
	}
	return
}

// Equals implements an equality test for a Cookbook.
func (c1 *Cookbook) Equals(c2 *Cookbook) (res bool) {
	res = reflect.DeepEqual(c1, c2)
	return
}

// Diff returns any attributes that have changed from one Cookbook struct to
// another.
func (c1 *Cookbook) Diff(c2 *Cookbook) (pos, neg *Cookbook) {
	if c1.Equals(c2) {
		return
	}
	pos = NewCookbook()
	neg = NewCookbook()

	if c1.Name != c2.Name {
		pos.Name = c2.Name
		neg.Name = c1.Name
	}
	for k, _ := range c1.Versions {
		if c2.Versions[k] == nil {
			neg.Versions[k] = c1.Versions[k]
		} else if !c1.Versions[k].Equals(c2.Versions[k]) {
			pos.Versions[k], neg.Versions[k] = c1.Versions[k].Diff(c2.Versions[k])
		}
	}
	for k, _ := range c2.Versions {
		if c1.Versions[k] == nil {
			pos.Versions[k] = c2.Versions[k]
		}
	}
	if pos.Empty() {
		pos = nil
	}
	if neg.Empty() {
		neg = nil
	}
	return
}
