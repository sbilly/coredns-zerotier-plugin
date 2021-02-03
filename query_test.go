// Copyright 2020 Oz Tiram <oz.tiram@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zerotier

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestQuery(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://example.org/api/ipam/ip-addresses/").MatchParams(
		map[string]string{"dns_name": "my_host"}).Reply(
		200).BodyString(
		`{"count":1, "results":[{"address": ["10.0.0.1", "10.0.0.2"], "dns_name": "my_host"}]}`)

	want := []string{"10.0.0.1", "10.0.0.2"}
	got := query("https://example.org/api/ipam/ip-addresses", "mytoken", "my_host", time.Millisecond*100)

	if testEq(want, got) != true {
		t.Fatalf("Expected %s but got %s", want, got)
	}
}

func TestNoSuchHost(t *testing.T) {
	var want []string

	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://example.org/api/ipam/ip-addresses/").MatchParams(
		map[string]string{"dns_name": "NoSuchHost"}).Reply(
		200).BodyString(`{"count":0,"next":null,"previous":null,"results":[]}`)

	want = nil
	got := query("https://example.org/api/ipam/ip-addresses", "mytoken", "NoSuchHost", time.Millisecond*100)
	if testEq(want, got) != true {
		t.Fatalf("Expected empty string but got %s", got)
	}

}

func TestLocalCache(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://example.org/api/ipam/ip-addresses/").MatchParams(
		map[string]string{"dns_name": "my_host"}).Reply(
		200).BodyString(
		`{"count":1, "results":[{"address": ["10.0.0.1", "10.0.0.2"], "dns_name": "my_host"}]}`)

	ip_address := ""

	got := query("https://example.org/api/ipam/ip-addresses", "mytoken", "my_host", time.Millisecond*100)

	item, err := localCache.Get("my_host")
	if err == nil {
		ip_address = item.Value().(string)
	}

	assert.Equal(t, got, ip_address, "local cache item didn't match")

}

func TestLocalCacheExpiration(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://example.org/api/ipam/ip-addresses/").MatchParams(
		map[string]string{"dns_name": "my_host"}).Reply(
		200).BodyString(
		`{"count":1, "results":[{"address": "10.0.0.2/25", "dns_name": "my_host"}]}`)

	query("https://example.org/api/ipam/ip-addresses", "mytoken", "my_host", time.Millisecond*100)
	<-time.After(101 * time.Millisecond)
	item, err := localCache.Get("my_host")
	if err != nil {
		t.Fatalf("Expected errors, but got: %v", item)
	}
}

func testEq(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
