package cps

import (
	"encoding/json"
	"fmt"
	"time"

	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Enrollments represents an enrollment
//
// API Docs: https://developer.akamai.com/api/core_features/certificate_provisioning_system/v2.html#enrollments
type Enrollment struct {
	client.Resource
	AdminContact              *Contact              `json:"adminContact"`
	CertificateChainType      *string               `json:"certificateChainType"`
	CertificateType           CertificateType       `json:"certificateType"`
	CertificateSigningRequest *CSR                  `json:"csr"`
	ChangeManagement          bool                  `json:"changeManagement"`
	EnableMultiStacked        bool                  `json:"enableMultiStackedCertificates"`
	Location                  *string               `json:"location"`
	MaxAllowedSans            *int                  `json:"maxAllowedSanNames"`
	MaxAllowedWildcardSans    *int                  `json:"maxAllowedWildcardSanNames"`
	NetworkConfiguration      *NetworkConfiguration `json:"networkConfiguration"`
	Organization              *Organization         `json:"org"`
	PendingChanges            *[]string             `json:"pendingChanges"`
	RegistrationAuthority     RegistrationAuthority `json:"ra"`
	SignatureAuthority        *SHA                  `json:"signatureAlgorithm"`
	TechContact               *Contact              `json:"techContact"`
	ThirdParty                *ThirdParty           `json:"thirdParty"`
	ValidationType            ValidationType        `json:"validationType"`
}

type CreateEnrollmentQueryParams struct {
	ContractID      string
	DeployNotAfter  *string
	DeployNotBefore *string
}

type ListEnrollmentsQueryParams struct {
	ContractID string
}

type CreateEnrollmentResponse struct {
	Location string   `json:"enrollment"`
	Changes  []string `json:"changes"`
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

// Create an Enrollment on CPS
//
//
// API Docs: https://developer.akamai.com/api/core_features/certificate_provisioning_system/v2.html#5aaa335c
// Endpoint: POST /cps/v2/enrollments{?contractId,deploy-not-after,deploy-not-before}
func (enrollment *Enrollment) Create(params CreateEnrollmentQueryParams) (*CreateEnrollmentResponse, error) {
	var request = fmt.Sprintf(
		"/cps/v2/enrollments?contractId=%s",
		params.ContractID,
	)

	if params.DeployNotAfter != nil {
		request = fmt.Sprintf(
			"%s&deploy-not-after=%s",
			request,
			*params.DeployNotAfter,
		)
	}

	if params.DeployNotBefore != nil {
		request = fmt.Sprintf(
			"%s&deploy-not-before=%s",
			request,
			*params.DeployNotBefore,
		)
	}

	req, err := newRequest(
		"POST",
		request,
		enrollment,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response CreateEnrollmentResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Get an enrollment by location
//
//
// API Docs: https://developer.akamai.com/api/core_features/certificate_provisioning_system/v2.html#getasingleenrollment
// Endpoint: POST /cps/v2/enrollments/{enrollmentId}
func GetEnrollment(location string) (*Enrollment, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		location,
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.akamai.cps.enrollment.v7+json")

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response Enrollment
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func ListEnrollments(params ListEnrollmentsQueryParams) ([]Enrollment, error) {
	var enrollments []Enrollment

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/cps/v2/enrollments?contractId={%s}",
			params.ContractID,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, enrollments); err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (enrollment *Enrollment) Exists(enrollments []Enrollment) bool {
	for _, e := range enrollments {
		if e.CertificateSigningRequest.CommonName == enrollment.CertificateSigningRequest.CommonName {
			return true
		}
	}

	return false
}

// CreateEnrollment wraps enrollment.Create to accept json
func CreateEnrollment(data []byte, params CreateEnrollmentQueryParams) (*CreateEnrollmentResponse, error) {
	var enrollment Enrollment
	if err := json.Unmarshal(data, &enrollment); err != nil {
		return nil, err
	}

	return enrollment.Create(params)
}
