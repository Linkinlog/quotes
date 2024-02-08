package assets

import (
	"embed"
)

//go:embed *
var assets embed.FS

func Files() embed.FS {
	return assets
}

// Howdy detective, nothing fun to see here.
// See it all(and more) at https://github.com/Linkinlog/quotes
