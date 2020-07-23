/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package restapi

import (
	"fmt"

	"github.com/trustbloc/hub-auth/pkg/restapi/operation"
)

// New returns new controller instance.
func New(config *operation.Config) (*Controller, error) {
	var allHandlers []operation.Handler

	rpService, err := operation.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize hub-auth-rest operations: %w", err)
	}

	allHandlers = append(allHandlers, rpService.GetRESTHandlers()...)

	return &Controller{handlers: allHandlers}, nil
}

// Controller contains handlers for controller.
type Controller struct {
	handlers []operation.Handler
}

// GetOperations returns all controller endpoints.
func (c *Controller) GetOperations() []operation.Handler {
	return c.handlers
}
