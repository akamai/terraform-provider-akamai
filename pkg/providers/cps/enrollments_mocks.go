package cps

import "github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/cps"

func mockLetsEncryptChallenges() *cps.Change {
	allowedInput := cps.AllowedInput{
		Info:              "/cps/v2/enrollments/27552/changes/2848126/input/info/lets-encrypt-challenges",
		RequiredToProceed: false,
		Type:              "lets-encrypt-challenges",
		Update:            "/cps/v2/enrollments/27552/changes/2848126/input/update/lets-encrypt-challenges-completed",
	}
	statusInfo := cps.StatusInfo{
		DeploymentSchedule: nil,
		Description:        "Coordinating domain name validation.",
		Error:              nil,
		State:              "suspended",
		Status:             "coodinate-domain-validation",
	}
	change := &cps.Change{
		AllowedInput: []cps.AllowedInput{allowedInput},
		StatusInfo:   &statusInfo,
	}
	return change
}

func mockThirdPartyCSRChallenges() *cps.Change {
	allowedInput := cps.AllowedInput{
		Info:              "/cps/v2/enrollments/10002/changes/10002/input/info/third-party-csr",
		RequiredToProceed: true,
		Type:              "third-party-certificate",
		Update:            "/cps/v2/enrollments/10002/changes/10002/input/update/third-party-cert-and-trust-chain",
	}
	statusInfo := cps.StatusInfo{
		DeploymentSchedule: nil,
		Description:        "Waiting for you to upload and submit your third party certificate and trust chain.",
		Error:              nil,
		State:              "awaiting-input",
		Status:             "wait-upload-third-party",
	}
	change := &cps.Change{
		AllowedInput: []cps.AllowedInput{allowedInput},
		StatusInfo:   &statusInfo,
	}
	return change
}

func mockEVChallenges() *cps.Change {
	allowedInput := cps.AllowedInput{
		Info:              "/cps/v2/enrollments/10002/changes/10002/input/info/ev",
		RequiredToProceed: true,
		Type:              "ev",
		Update:            "/cps/v2/enrollments/10002/changes/10002/input/update/ev-and-trust-chain",
	}
	statusInfo := cps.StatusInfo{
		DeploymentSchedule: nil,
		Description:        "Waiting for you to upload and submit your ev certificate and trust chain.",
		Error:              nil,
		State:              "awaiting-input",
		Status:             "wait-upload-ev",
	}
	change := &cps.Change{
		AllowedInput: []cps.AllowedInput{allowedInput},
		StatusInfo:   &statusInfo,
	}
	return change
}

func mockEmptyChanges() *cps.Change {
	allowedInput := cps.AllowedInput{}
	statusInfo := cps.StatusInfo{}
	change := &cps.Change{
		AllowedInput: []cps.AllowedInput{allowedInput},
		StatusInfo:   &statusInfo,
	}
	return change
}

func mockDVArray() *cps.DVArray {
	challenges0 := cps.Challenge{
		Error:             "",
		FullPath:          "http://TestFullPath",
		RedirectFullPath:  "http://TestRedirectFullPath.com",
		ResponseBody:      "TestResponseBody",
		Status:            "pending",
		Token:             "TestToken123",
		Type:              "http-01",
		ValidationRecords: nil,
	}
	challenges1 := cps.Challenge{
		Error:             "",
		FullPath:          "TestFullPath",
		RedirectFullPath:  "",
		ResponseBody:      "TestResponseBody",
		Status:            "pending",
		Token:             "TestToken123",
		Type:              "dns-01",
		ValidationRecords: nil,
	}
	dv := cps.DV{
		Challenges:         []cps.Challenge{challenges0, challenges1},
		Domain:             "TestDomain",
		Error:              "The domain TestDomain is not ready for HTTP validation.",
		Expires:            "2022-07-25T10:17:44Z",
		RequestTimestamp:   "2022-07-18T10:17:44Z",
		Status:             "Awaiting user",
		ValidatedTimestamp: "2022-07-19T09:35:29Z",
		ValidationStatus:   "DATA_NOT_READY",
	}
	dvArray := &cps.DVArray{DV: []cps.DV{dv}}
	return dvArray
}

func mockThirdPartyCSRDVArray() *cps.DVArray {
	challenges0 := cps.Challenge{
		Error:             "",
		FullPath:          "http://testFullPath.com",
		RedirectFullPath:  "http://testRedirectFullPath.com",
		ResponseBody:      "TestResponseBody",
		Status:            "pending",
		Token:             "TestToken",
		Type:              "third-party",
		ValidationRecords: nil,
	}
	dv := cps.DV{
		Challenges:         []cps.Challenge{challenges0},
		Domain:             "test.com",
		Error:              "2022-07-25T10:17:44Z",
		Expires:            "2022-07-25T10:17:44Z",
		RequestTimestamp:   "2022-07-25T10:17:44Z",
		Status:             "completed",
		ValidatedTimestamp: "2022-07-25T10:17:44Z",
		ValidationStatus:   "done",
	}
	dvArray := &cps.DVArray{DV: []cps.DV{dv}}
	return dvArray
}

func mockEVDVArray() *cps.DVArray {
	challenges0 := cps.Challenge{
		Error:             "",
		FullPath:          "http://testFullPath.com",
		RedirectFullPath:  "http://testRedirectFullPath.com",
		ResponseBody:      "TestResponseBody",
		Status:            "pending",
		Token:             "TestToken",
		Type:              "ev",
		ValidationRecords: nil,
	}
	dv := cps.DV{
		Challenges:         []cps.Challenge{challenges0},
		Domain:             "test.com",
		Error:              "2022-07-25T10:17:44Z",
		Expires:            "2022-07-25T10:17:44Z",
		RequestTimestamp:   "2022-07-25T10:17:44Z",
		Status:             "completed",
		ValidatedTimestamp: "2022-07-25T10:17:44Z",
		ValidationStatus:   "done",
	}
	dvArray := &cps.DVArray{DV: []cps.DV{dv}}
	return dvArray
}

func mockEmptyDVArray() *cps.DVArray {
	dv := cps.DV{
		Challenges:         []cps.Challenge{},
		Domain:             "",
		Error:              "",
		Expires:            "",
		RequestTimestamp:   "",
		Status:             "",
		ValidatedTimestamp: "",
		ValidationStatus:   "",
	}
	dvArray := &cps.DVArray{DV: []cps.DV{dv}}
	return dvArray
}
