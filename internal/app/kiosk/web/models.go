// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package web

// Errors encapsulates the HTTP error response model for 4xx and 5xx statuses.
type Errors struct {
	Errors []Error `json:"errors"`
}

// Error encapsulates a specific error that contains the error details.
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
