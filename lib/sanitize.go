// Copyright 2013 Matthew Baird, Dimitri Roche
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
	"regexp"
	"strings"
)

// Sanitize Elastic Search query string input.
// Inspired by https://github.com/lanetix/node-elasticsearch-sanitize
// Removed escaping of white space because Elasticsearch complained
func Sanitize(input string) string {
	charReg := regexp.MustCompile("[\\*\\+\\-=~><\"\\?^\\${}\\(\\)\\:\\!\\/[\\]]+")

	output := charReg.ReplaceAllString(input, "\\$0")
	output = strings.Replace(output, "||", "\\||", -1)
	output = strings.Replace(output, "&&", "\\&&", -1)
	output = strings.Replace(output, "AND", "\\A\\N\\D", -1)
	output = strings.Replace(output, "OR", "\\O\\R", -1)
	output = strings.Replace(output, "NOT", "\\N\\O\\T", -1)
	return output
}
