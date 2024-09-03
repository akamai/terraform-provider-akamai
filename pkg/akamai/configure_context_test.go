package akamai

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_validateRetryConfiguration(t *testing.T) {

	tests := map[string]struct {
		name    string
		args    contextConfig
		wantErr bool
		errMsg  string
	}{
		"valid values": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: 0,
				retryWaitMax: 0,
			},

			wantErr: false,
		},
		"invalid values - negative number of retries": {
			args: contextConfig{
				retryMax:     -3,
				retryWaitMin: 0,
				retryWaitMax: 0,
			},

			wantErr: true,
			errMsg:  "wrong retry values: maximum number of retries (-3), minimum retry wait time (0s), maximum retry wait time (0s) cannot be negative",
		},
		"invalid values - negative min wait time": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: -1,
				retryWaitMax: 0,
			},

			wantErr: true,
			errMsg:  "wrong retry values: maximum number of retries (0), minimum retry wait time (-1ns), maximum retry wait time (0s) cannot be negative",
		},
		"invalid values - negative max wait time": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: 0,
				retryWaitMax: -1,
			},
			wantErr: true,
			errMsg:  "wrong retry values: maximum number of retries (0), minimum retry wait time (0s), maximum retry wait time (-1ns) cannot be negative",
		},
		"invalid values - min wait time cannot be higher than max wait time": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: 1,
				retryWaitMax: 0,
			},
			wantErr: true,
			errMsg:  "wrong retry values: maximum retry wait time (0s) cannot be lower than minimum retry wait time (1ns)",
		},
		"invalid values - too many retries": {
			args: contextConfig{
				retryMax:     51,
				retryWaitMin: 0,
				retryWaitMax: 0,
			},
			wantErr: true,
			errMsg:  "wrong retry values: too many retries, maximum number of retries (51) cannot be higher than 50",
		},
		"invalid values - retry time too long (retryWaitMin)": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: time.Hour * 25,
				retryWaitMax: time.Hour * 26, // needs to be higher than retryWaitMin
			},
			wantErr: true,
			errMsg:  "wrong retry values: retry wait time too long, minimum retry wait time (25h0m0s) cannot be higher than 24h0m0s or maximum retry wait time (26h0m0s) cannot be higher than 24h0m0s",
		},
		"invalid values - retry time too long (retryWaitMax)": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: 1,
				retryWaitMax: time.Hour * 25,
			},
			wantErr: true,
			errMsg:  "wrong retry values: retry wait time too long, minimum retry wait time (1ns) cannot be higher than 24h0m0s or maximum retry wait time (25h0m0s) cannot be higher than 24h0m0s",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateRetryConfiguration(tt.args)
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
