// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"time"
)

var html = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>Hello Universe</title>
</head>
<body>
  <h3>Hello Universe</h3>
  <p>Hostname: %s</p>
</body>
</html>
`

func httpHandler(w http.ResponseWriter, r *http.Request) {
	format := "%s - [%s] \"%s %s %s\" %s\n"
	fmt.Printf(format, hostname, time.Now().Format(time.RFC1123),
		r.Method, r.URL.Path, r.Proto, r.UserAgent())
	fmt.Fprintf(w, html, hostname)
}
