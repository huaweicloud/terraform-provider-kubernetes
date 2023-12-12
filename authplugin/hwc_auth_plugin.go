package authplugin

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-provider-kubernetes/authplugin/signer"
)

const (
	AccessKeyConfiguration     = "hw_access_key"
	SecretKeyConfiguration     = "hw_secret_key"
	ProjectIdConfiguration     = "hw_project_id"
	SecurityTokenConfiguration = "hw_security_token"

	ProjectIdHeaderKey     = "X-Project-Id"
	SecurityTokenHeaderKey = "X-Security-Token"
)

// CreateHWAuthWrapperByAuthParams to create Wrapper
func CreateHWAuthWrapperByAuthParams(params HWAuthParameters, rt http.RoundTripper) http.RoundTripper {
	if !params.HasRequiredAttributes() {
		log.Printf("[TRACE] Do not use huawei auth")
		return rt
	}

	log.Printf("[TRACE] Use huawei auth wrapper for request")

	return &huaweiAuthDecorator{parameters: params, next: rt}
}

type HWAuthParameters struct {
	//AccessKey Huawei Cloud access key, could be temporary
	AccessKey string
	//SecretKey Huawei Cloud secret key, could be temporary
	SecretKey string
	//ProjectId Huawei Cloud ProjectId
	ProjectId string
	//SecurityToken Huawei Cloud Key, could be temporary
	SecurityToken string
}

func (p HWAuthParameters) Valid() error {
	if p.HasRequiredAttributes() {
		return nil
	}

	if len(p.AccessKey) != 0 || len(p.SecretKey) != 0 || len(p.ProjectId) != 0 || len(p.SecurityToken) != 0 {
		return fmt.Errorf("huawei cloud authentication requires hw_access_key, hw_secret_key, " +
			"hw_project_id set in the same time")
	}

	return nil
}

func (p HWAuthParameters) HasRequiredAttributes() bool {
	if len(p.AccessKey) != 0 && len(p.SecretKey) != 0 && len(p.ProjectId) != 0 {
		return true
	}

	return false
}

type huaweiAuthDecorator struct {
	parameters HWAuthParameters
	next       http.RoundTripper
}

func (d *huaweiAuthDecorator) RoundTrip(req *http.Request) (*http.Response, error) {
	err := d.decorateRequestHeader(req)
	if err != nil {
		log.Printf("[ERROR] Decorate request header failed: %s", err)

		return nil, err
	}

	err = d.signRequest(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[TRACE] Calling huawei authentication roundTrip")

	return d.next.RoundTrip(req)
}

func (d huaweiAuthDecorator) decorateRequestHeader(req *http.Request) error {
	if len(d.parameters.ProjectId) != 0 {
		req.Header.Add(ProjectIdHeaderKey, d.parameters.ProjectId)
	} else {
		return errors.New(ProjectIdConfiguration + " must be set for authentication")
	}

	if len(d.parameters.SecurityToken) != 0 {
		req.Header.Add(SecurityTokenHeaderKey, d.parameters.SecurityToken)
	}

	return nil
}

func (d huaweiAuthDecorator) signRequest(req *http.Request) error {
	var key, secret []byte

	if len(d.parameters.AccessKey) != 0 {
		key = []byte(d.parameters.AccessKey)
	} else {
		return errors.New(AccessKeyConfiguration + " must be set for authentication")
	}

	if len(d.parameters.SecretKey) != 0 {
		secret = []byte(d.parameters.SecretKey)
	} else {
		return errors.New(SecretKeyConfiguration + " must be set for authencation")
	}

	return createSigner(key, secret).Sign(req)
}

func createSigner(ak []byte, sk []byte) *signer.Signer {
	return &signer.Signer{
		Key:    string(ak),
		Secret: string(sk),
	}
}
