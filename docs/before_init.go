// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import "os"

var is_testnet = os.Getenv("TESTNET")

func init() {
	if is_testnet != "" {
		SwaggerInfo.BasePath = "/test"
	}
}