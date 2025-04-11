package frontend

import "embed"

// TemplateFS holds all embedded templates and can be used like a regular file system.
// The templates are embedded in the binary at compile time.
//
//go:embed templates/**/*.html assets/css/*.css assets/js/*.js
var TemplateFS embed.FS
