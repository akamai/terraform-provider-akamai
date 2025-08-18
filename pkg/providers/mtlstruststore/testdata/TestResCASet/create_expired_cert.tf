provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlstruststore_ca_set" "test" {
  name                = "mgw-res-1"
  description         = "Test CA Set for validation"
  allow_insecure_sha1 = false
  version_description = "Initial version for testing"
  certificates = [
    {
      certificate_pem = "-----BEGIN CERTIFICATE-----\nMIIC3jCCAcagAwIBAgIBATANBgkqhkiG9w0BAQsFADAAMB4XDTI0MDcwMTAwMDAwMFoXDTI1MDcwMjEwNDcxNlowADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL1xDWbGQoeGrUimkp7KUnlj1w0+aHNs9QH9sXygMxxis9cNZeBewJv9fL7n2MmFSkgmsAxJ8/90G19cyfWzNnPtO9PF9kBVarPU79CUVPx9D3hJBPIlDKozrbZYy2H4HRbQ41xM9DF4DjXIqX3Lk8YslTf8SOSxgIpQLKVrdvIxSTY3uH+u0E67dtcTcz6Ytop1Z0u4Q7GesC6iUoqWYNNPRGTETN++kTZ1XqVXWoVWML4ffeHpqUqHm/ITY0OKeXcIMTD/lg0zFdMqMYY01Y76Vddgts5utqmgt7qJ6mWlETHpVXNiIn/ooukxCsAgxgfqS/iXEyOvWrmCT/O4rqECAwEAAaNjMGEwHQYDVR0OBBYEFNG9BugMXYKr404m1c4nIIuCbB6VMB8GA1UdIwQYMBaAFNG9BugMXYKr404m1c4nIIuCbB6VMA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgGGMA0GCSqGSIb3DQEBCwUAA4IBAQB+0lz3a78X+Eg0QqLBUJBtNzuVwNMb2Yq3in0s+91QOMymfc5uVl55S/8RxInbTdwD51jkW9akAl3fpQu2PBRwLPJPpewHnWZJgVK/xws0TDJYHe0iWPpNTfQRU5QuciTAp5lwhyyzpamZK2uE76lYTwZ8y6ZHblPDp1JCu6k2soH0YkvTzrKSUJUki70jhVajEFEUZ8S19PQ8+UeEycxn6c629ZPgw87aej8SEbPiY60J1vq+o4px/9HpW9pGeZNMilIvvY9ezDtqERmC4mKbYPNzSFwkYJ5mqG4yUlafITpl6/nvPi9rihMcf06rEhhPye+nztQrscMYcvr7qaDh\n-----END CERTIFICATE-----\n"
      description     = "Test certificate"
    }
  ]
}