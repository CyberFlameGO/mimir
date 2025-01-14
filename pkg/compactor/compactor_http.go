// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/compactor/compactor_http.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package compactor

import (
	"html/template"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"

	util_log "github.com/grafana/mimir/pkg/util/log"
)

var (
	compactorStatusPageTemplate = template.Must(template.New("main").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Compactor Ring</title>
		</head>
		<body>
			<h1>Compactor Ring</h1>
			<p>{{ .Message }}</p>
		</body>
	</html>`))
)

func writeMessage(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	err := compactorStatusPageTemplate.Execute(w, struct {
		Message string
	}{Message: message})

	if err != nil {
		level.Error(util_log.Logger).Log("msg", "unable to serve compactor ring page", "err", err)
	}
}

func (c *MultitenantCompactor) RingHandler(w http.ResponseWriter, req *http.Request) {
	if c.State() != services.Running {
		// we cannot read the ring before MultitenantCompactor is in Running state,
		// because that would lead to race condition.
		writeMessage(w, "Compactor is not running yet.")
		return
	}

	c.ring.ServeHTTP(w, req)
}
