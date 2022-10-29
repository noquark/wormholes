package sql

import (
	_ "embed"
)

//go:embed update_link.sql
var Update string

//go:embed get_link.sql
var Get string

//go:embed delete_link.sql
var Delete string

//go:embed insert_link.sql
var Insert string
