// Copyright 2013 Matthew Baird
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elastigo

import (
	"encoding/json"
)

type AliasOption struct {
	Actions []AliasAction `json:"actions"`
}

type AliasAction map[string]AliasActionItem

type AliasActionItem struct {
	Alias string `json:"alias"`
	Index string `json:"index"`
}

type Aliases map[string]interface{}

func (c *Conn) PutAliases(oldIndex, newIndex, alias string) (BaseResponse, error) {
	var retval BaseResponse

	actions := make([]AliasAction, 0)
	aliasOption := AliasOption{
		Actions: actions,
	}

	if len(oldIndex) > 0 {
		action := AliasAction{
			"remove": AliasActionItem{
				Alias: alias,
				Index: oldIndex,
			},
		}
		actions = append(actions, action)
	}

	if len(newIndex) > 0 {
		action := AliasAction{
			"add": AliasActionItem{
				Alias: alias,
				Index: newIndex,
			},
		}
		actions = append(actions, action)
	}

	requestBody, err := json.Marshal(aliasOption)
	if err != nil {
		return retval, err
	}

	body, err := c.DoCommand("POST", "/_aliases", nil, requestBody)
	if err != nil {
		return retval, err
	}

	jsonErr := json.Unmarshal(body, &retval)
	if jsonErr != nil {
		return retval, jsonErr
	}

	return retval, nil
}
