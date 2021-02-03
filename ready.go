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

	clog "github.com/coredns/coredns/plugin/pkg/log"
	ztcentral "github.com/zerotier/go-ztcentral"
)

// Ready :check the connection with TOKEN
func (n Zerotier) Ready() bool {
	c := ztcentral.NewClient(n.Token)
	ctx := context.Background()

	// get list of networks
	_, err := c.GetNetworks(ctx)
	if err != nil {
		clog.Fatalf("Error reading all networks %v", err)
		return false
	}

	return true
}
