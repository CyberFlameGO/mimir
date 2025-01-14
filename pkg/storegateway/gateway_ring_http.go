// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/storegateway/gateway_http.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package storegateway

import (
	"net/http"
	"text/template"

	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"

	util_log "github.com/grafana/mimir/pkg/util/log"
)

var (
	statusPageTemplate = template.Must(template.New("main").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Store Gateway Ring</title>
		</head>
		<body>
			<h1>Store Gateway Ring</h1>
			<p>{{ .Message }}</p>
		</body>
	</html>`))
)

func writeMessage(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	err := statusPageTemplate.Execute(w, struct {
		Message string
	}{Message: message})

	if err != nil {
		level.Error(util_log.Logger).Log("msg", "unable to serve store gateway ring page", "err", err)
	}
}

func (c *StoreGateway) RingHandler(w http.ResponseWriter, req *http.Request) {
	if c.State() != services.Running {
		// we cannot read the ring before the store gateway is in Running state,
		// because that would lead to race condition.
		writeMessage(w, "Store gateway is not running yet.")
		return
	}

	c.ring.ServeHTTP(w, req)
}
