module github.com/akamai/AkamaiOPEN-edgegrid-golang

require (
	github.com/go-ini/ini v1.44.0
	github.com/google/go-querystring v1.0.0
	github.com/google/uuid v1.1.1
	github.com/h2non/gock v0.0.0-00010101000000-000000000000
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	github.com/stretchr/testify v1.3.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	gopkg.in/h2non/gock.v1 v1.0.15
	gopkg.in/ini.v1 v1.44.0 // indirect
)

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14

go 1.13
