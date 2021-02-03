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
	"context"
	"fmt"
	"time"

	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/imkira/go-ttlmap"

	ztcentral "github.com/zerotier/go-ztcentral"
)

type Record struct {
	Address  string `json:"address"`
	HostName string `json:"dns_name,omitempty"`
}

type RecordsList struct {
	Records []Record `json:"results"`
}

var localCache = ttlmap.New(nil)

func query(url, token, dns_name string, duration time.Duration) []string {
	var addresses []string

	item, err := localCache.Get(dns_name)
	if err == nil {
		clog.Debug(fmt.Sprintf("Found in local cache %s", dns_name))
		return item.Value().([]string)
	} else {
		c := ztcentral.NewClient(token)
		ctx := context.Background()

		// get list of networks
		networks, err := c.GetNetworks(ctx)
		if err != nil {
			clog.Fatalf("Error reading all networks %v", err)
		}

		// print networks and members
		for _, n := range networks {
			members, err := c.GetMembers(ctx, n.ID)
			if err != nil {
				clog.Fatalf("Error reading network info %v", err)
			}

			for _, m := range members {
				if m.Name == dns_name {
					clog.Debugf("Found %s recorde: %s\t%s\n", dns_name, m.Name, m.Config.IPAssignments)
					addresses = append(addresses, m.Config.IPAssignments...)
				}
			}
		}

		if len(addresses) == 0 {
			clog.Debugf("Recored %s not found.", dns_name)
			return nil
		} else {
			localCache.Set(dns_name, ttlmap.NewItem(addresses, ttlmap.WithTTL(duration)), nil)
			return addresses
		}
	}
}
