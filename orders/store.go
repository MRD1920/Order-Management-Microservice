package main

import "context"

type store struct {
	// Add monogDB instance here.
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(context.Context) error {
	return nil
}
