package kubernetes

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/client-go/transport"

	"github.com/hashicorp/terraform-provider-kubernetes/authplugin"
)

func BuildWrappers(d *schema.ResourceData) []transport.WrapperFunc {
	wrappers := make([]transport.WrapperFunc, 0)

	// Add logger wrapper if log level meet debug requirement
	if logging.IsDebugOrHigher() {
		log.Printf("[TRACE] Creating logger wrappers")
		wrappers = append(wrappers, debugLogWrapper)
	}

	// Add Huawei Cloud authentication wrapper
	hwAuthParams, err := BuildHWAuthParameters(d)
	if err == nil && hwAuthParams.HasRequiredAttributes() {
		log.Printf("[TRACE] Creating Huawei Cloud authentication wrappers")
		wrappers = append(wrappers, func(rt http.RoundTripper) http.RoundTripper {
			return authplugin.CreateHWAuthWrapperByAuthParams(hwAuthParams, rt)
		})
	}

	return wrappers
}

func BuildHWAuthParameters(d *schema.ResourceData) (authplugin.HWAuthParameters, error) {
	// The configuration only valid when the following combinations are given in the
	// same time:
	// 1. hw_access_key, hw_secret_key, hw_project_id
	// 2. hw_access_key, hw_secret_key, hw_project_id, hw_security_token
	authParameters := authplugin.HWAuthParameters{}
	if ak, ok := d.Get(authplugin.AccessKeyConfiguration).(string); ok {
		authParameters.AccessKey = ak
	}
	if sk, ok := d.Get(authplugin.SecretKeyConfiguration).(string); ok {
		authParameters.SecretKey = sk
	}
	if projectId, ok := d.Get(authplugin.ProjectIdConfiguration).(string); ok {
		authParameters.ProjectId = projectId
	}
	if securityToken, ok := d.Get(authplugin.SecurityTokenConfiguration).(string); ok {
		authParameters.SecurityToken = securityToken
	}

	if err := authParameters.Valid(); err != nil {
		return authplugin.HWAuthParameters{}, err
	}

	return authParameters, nil
}

func debugLogWrapper(rt http.RoundTripper) http.RoundTripper {
	return logging.NewTransport("kubernetes", rt)
}
