package sql

import (
	_ "embed"
)

//go:embed insert_event.sql
var InsertEvent string

//go:embed get_link.sql
var GetLink string
