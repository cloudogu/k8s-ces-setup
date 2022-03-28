package core

import "fmt"

// ResourceError wraps the original error and takes additional arguments to identify a K8s resource by kind, version and name
type ResourceError struct {
	err           error
	wrapperErrMsg string
	kind          string
	apiVersion    string
	resourceName  string
}

// NewResourceError creates a custom k8s error that identifies the resource by kind, version and name.
func NewResourceError(err error, wrapperErrMsg, kind, apiVersion, resourceName string) *ResourceError {
	return &ResourceError{
		err:           err,
		wrapperErrMsg: wrapperErrMsg,
		kind:          kind,
		apiVersion:    apiVersion,
		resourceName:  resourceName,
	}
}

func (e *ResourceError) Error() string {
	return fmt.Sprintf("%s (resource %s/%s/%s): %+v", e.wrapperErrMsg, e.kind, e.apiVersion, e.resourceName, e.err)
}

// Unwrap returns the original error.
func (e *ResourceError) Unwrap() error {
	return e.err
}