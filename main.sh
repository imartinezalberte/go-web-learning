#!/bin/bash
#

reflex -r '\.go$' -s -- sh -c 'go mod tidy && go run ./cmd/web -addr=:8080'

# Go http.FileServer allows partial content(The range header)
# curl -i -H "Range: bytes=100-199" --output - http://localhost:8080/static/img/logo.png
#
