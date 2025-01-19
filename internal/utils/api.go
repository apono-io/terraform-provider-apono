package utils

import (
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"net/http"
)

func GetDiagnosticsForApiError(err error, actionType string, objectName string, objectId string) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	var errorMessagePrefix string
	if objectId != "" {
		errorMessagePrefix = fmt.Sprintf("Failed to %s %s with id %s", actionType, objectName, objectId)
	} else {
		errorMessagePrefix = fmt.Sprintf("Failed to %s %s", actionType, objectName)
	}

	if apiError, ok := err.(*apono.GenericOpenAPIError); ok {
		diagnostics.AddError("Client Error", fmt.Sprintf("%s, error: %s, body: %s", errorMessagePrefix, apiError.Error(), string(apiError.Body())))
	} else if apiError, ok := err.(*aponoapi.GenericOpenAPIError); ok {
		diagnostics.AddError("Client Error", fmt.Sprintf("%s, error: %s, body: %s", errorMessagePrefix, apiError.Error(), string(apiError.Body())))
	} else {
		diagnostics.AddError("Client Error", fmt.Sprintf("%s: %s", errorMessagePrefix, err.Error()))
	}

	return diagnostics
}

func IsAponoApiNotFoundError(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}
