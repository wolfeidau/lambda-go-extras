// Package middleware implements a middleware chaining solution for Go based lambdas.
package middleware

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
)

// FieldMap used to pass in a list of attribute value pairs.
type FieldMap map[string]any

// Middleware A constructor for a a piece of middleware.
// Some middleware use this constructor out of the box,
// so in most cases you can just pass somepackage.New.
type Middleware func(next lambda.Handler) lambda.Handler

// Chain is a middleware chain.
type Chain struct {
	middlewares []Middleware
}

// New creates a new chain.
func New(middlewares ...Middleware) Chain {
	return Chain{append([]Middleware(nil), middlewares...)}
}

// Use append one or more middleware(s) onto the existing chain.
func (c *Chain) Use(middlewares ...Middleware) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// Then chains the middleware and returns the final lambda.Handler.
//
//	New(m1, m2, m3).Then(h)
//
// is equivalent to:
//
//	m1(m2(m3(h)))
//
// When the request comes in, it will be passed to m1, then m2, then m3
// and finally, the given handler
// (assuming every middleware calls the following one).
func (c Chain) Then(h lambda.Handler) lambda.Handler {
	for i := range c.middlewares {
		h = c.middlewares[len(c.middlewares)-1-i](h)
	}

	return h
}

// ThenFunc works identically to Then, but takes
// a HandlerFunc instead of a Handler.
func (c Chain) ThenFunc(fn lambdaextras.HandlerFunc) lambda.Handler {
	return c.Then(fn)
}
