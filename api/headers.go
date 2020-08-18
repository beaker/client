package api

const (
	// HeaderAuthor can be set to an account's API token to set that account as
	// the author of a created resource. If omitted, the author defaults to the
	// request's authenticated user as specified in the Authorization header.
	HeaderAuthor = "Beaker-Author"

	// HeaderVersion can be set on any request to specify the client's version.
	HeaderVersion = "Beaker-Version"
)
