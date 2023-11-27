package mailgun

import (
	"github.com/crossplane/upjet/pkg/config"
)

// Configure  database resource
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mailgun_domain_credential", func(r *config.Resource) {
		r.UseAsync = true

		r.ShortGroup = "credential.domain.mailgun"
	})
}
