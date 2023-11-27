/*
Copyright 2022 Upbound Inc.
*/

package config

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/upjet/pkg/config"
)

const (
	// ErrFmtNoAttribute is an error string for not-found attributes
	ErrFmtNoAttribute = `"attribute not found: %s`
	// ErrFmtUnexpectedType is an error string for attribute map values of unexpected type
	ErrFmtUnexpectedType = `unexpected type for attribute %s: Expecting a string`
)

func region(parameters map[string]any) (string, error) {
	region, ok := parameters["region"]
	if !ok {
		return "", errors.Errorf(ErrFmtNoAttribute, "region")
	}
	regionStr, ok := region.(string)
	if !ok {
		return "", errors.Errorf(ErrFmtUnexpectedType, "region")
	}
	return regionStr, nil
}

func domain(parameters map[string]any) (string, error) {
	domain, ok := parameters["domain"]
	if !ok {
		return "", errors.Errorf(ErrFmtNoAttribute, "domain")
	}
	domainStr, ok := domain.(string)
	if !ok {
		return "", errors.Errorf(ErrFmtUnexpectedType, "domain")
	}
	return domainStr, nil
}

var credentialsFromProvider = config.ExternalName{
	SetIdentifierArgumentFn: config.NopSetIdentifierArgument,
	GetExternalNameFn:       config.IDAsExternalName,
	GetIDFn: func(ctx context.Context, externalName string, parameters map[string]any, providerConfig map[string]any) (string, error) {
		region, err := region(parameters)
		if err != nil {
			return region, err
		}
		domain, err := domain(parameters)
		if err != nil {
			return domain, err
		}

		return fmt.Sprintf("%s:%s@%s", region, externalName, domain), nil
	},
	DisableNameInitializer: true,
}

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// Import requires using a randomly generated ID from provider: nl-2e21sda
	"mailgun_domain_credential": credentialsFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
