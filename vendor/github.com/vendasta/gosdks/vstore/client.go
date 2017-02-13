package vstore

import (
	"golang.org/x/net/context"
	"crypto/x509"
	"crypto/tls"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc"
	"golang.org/x/oauth2/google"
	"github.com/vendasta/gosdks/pb/vstorepb"
)

const (
	prodAddress = "vstore-api-production.vendasta-internal.com:443"
	testAddress = "vstore-api-test.vendasta-internal.com:443"
	demoAddress = "vstore-api-demo.vendasta-internal.com:443"
	localAddress = "vstore-api-test.vendasta-internal.com:443" // Namespacing manages separation of developer-specific data
	internalAddress = "127.0.0.1:10000"
)

// Creates a new vstore client
func newClient(e env, dialOption ...grpc.DialOption) (vstorepb.VStoreClient, vstorepb.VStoreAdminClient, error) {
	dts, err := google.DefaultTokenSource(context.Background(), "https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		return nil, nil, err
	}

	var address string
	var rootCAs *x509.CertPool
	certificates := []tls.Certificate{}

	if (e == Internal) {
		address = internalAddress
		cer, err := tls.X509KeyPair([]byte(LocalCert), []byte(LocalKey)); if err != nil {
			return nil, nil, err
		}
		rootCAs = x509.NewCertPool()
		rootCAs.AppendCertsFromPEM([]byte(LocalCa))
		certificates = append(certificates, cer)
	} else if (e == Local) {
		address = localAddress
	} else if (e == Demo) {
		address = demoAddress
	} else if (e == Test) {
		address = testAddress
	} else if (e == Prod) {
		address = prodAddress
	}

	config := &tls.Config{
		Certificates: certificates,
		RootCAs: rootCAs,
	}
	creds := credentials.NewTLS(config)
	dialOption = append(
		dialOption,
		grpc.WithTransportCredentials(creds),
		grpc.WithBalancer(grpc.RoundRobin(NewPoolResolver(3, &DialSettings{Endpoint: address}))),
		grpc.WithBackoffConfig(grpc.DefaultBackoffConfig),
		grpc.WithPerRPCCredentials(oauth.TokenSource{dts}),
	)
	client, err := grpc.Dial(
		address,
		dialOption...
	); if err != nil {
		return nil, nil, err
	}
	return vstorepb.NewVStoreClient(client), vstorepb.NewVStoreAdminClient(client), nil
}

var (
	//LocalCert is a cert that can be used on a local development env
	LocalCert = `-----BEGIN CERTIFICATE-----
MIIEIzCCAwugAwIBAgIUGHRws8M0+QMre/Cn3qBIFlDhrMUwDQYJKoZIhvcNAQEL
BQAwZzELMAkGA1UEBhMCVVMxDzANBgNVBAgTBk9yZWdvbjERMA8GA1UEBxMIUG9y
dGxhbmQxEzARBgNVBAoTCkt1YmVybmV0ZXMxDTALBgNVBAsTBGdSUEMxEDAOBgNV
BAMTB2dSUEMgQ0EwHhcNMTYwODA5MjAwNjAwWhcNMTcwODA5MjAwNjAwWjByMQsw
CQYDVQQGEwJVUzEPMA0GA1UECBMGT3JlZ29uMREwDwYDVQQHEwhQb3J0bGFuZDET
MBEGA1UEChMKS3ViZXJuZXRlczEOMAwGA1UECxMFSGVsbG8xGjAYBgNVBAMTEWhl
bGxvLmV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
0jKekeB55bfGFgJWlvMgL50paqFZSrzdQDPL/UKJ7Mwbp9pLgwed16WwH6Rj4U97
CSHPEe0uHrBIqcCg+4oOD9C8Gka+sQ2TRJ9a8m4wRVKZ0RemRBxdJUbJrMHCjxkF
vcJ8yHleAo/EsSB3aA0cNvUfe7awydkw6G/Tjgb0gdLIg1xl/uI5HzG4XMtOC8Mg
l93VUr0YEiuraN+hEcQaYZJc2cwuldUE9cpnALpyPoaMkD4ySnvYwXp1OdJrk69t
dtD19NT2tE4YfA+8ZXnIWdYigWSTpyP1uEnUXy5QozsTtRqH6vQVZ1hZqVZJrXLW
Fo35VGULcQZ9bXk5XEXsbwIDAQABo4G7MIG4MA4GA1UdDwEB/wQEAwIFoDATBgNV
HSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBS24IhrS06G
G7AT7ilDueLpLtsXwzAfBgNVHSMEGDAWgBSBT/30bYfYIfyxk84miGq286XMJjBD
BgNVHREEPDA6ghFoZWxsby5leGFtcGxlLmNvbYIfaGVsbG8uZGVmYXVsdC5zdmMu
Y2x1c3Rlci5sb2NhbIcEfwAAATANBgkqhkiG9w0BAQsFAAOCAQEAiTXwLhIdm9GW
erSLo+rXZY36HfVGnzPyfl87CrkmznWflIxV0cc3hNnCFj5QX5nDuQjcFwV7DLbt
l8Z41b/tyv5gfmy3yJmbhxr8SipzIje9J2ZZSrAcA5grtWZBNgFcRRkr3s8lvlzv
rFGDSKcd6MzIIm64CdqKgwrLtKOGkTB/RKWjRmZavUZvjn9hVpsq47JBqvGvDA5d
dagPcW/M9onWbPeCVjP2uXakvO13vnbyfo6VKPtzEJjTOTjLneMOMSaWq/huK7gO
hXSZk4CqXAun1pbUEOh+AOXzwHNLB6DnxsAFW6wKAK7ZUWmPB5QrDOn0auKnuh+6
nQBCUzfayw==
-----END CERTIFICATE-----`

	//LocalKey is a key that can be used on a local development env
	LocalKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0jKekeB55bfGFgJWlvMgL50paqFZSrzdQDPL/UKJ7Mwbp9pL
gwed16WwH6Rj4U97CSHPEe0uHrBIqcCg+4oOD9C8Gka+sQ2TRJ9a8m4wRVKZ0Rem
RBxdJUbJrMHCjxkFvcJ8yHleAo/EsSB3aA0cNvUfe7awydkw6G/Tjgb0gdLIg1xl
/uI5HzG4XMtOC8Mgl93VUr0YEiuraN+hEcQaYZJc2cwuldUE9cpnALpyPoaMkD4y
SnvYwXp1OdJrk69tdtD19NT2tE4YfA+8ZXnIWdYigWSTpyP1uEnUXy5QozsTtRqH
6vQVZ1hZqVZJrXLWFo35VGULcQZ9bXk5XEXsbwIDAQABAoIBAC672GugmBmN+Xmt
PWlEVvSfIbU2eG7YfOyoV4NQhu/iFYgFTeKtD9gBW549Y4OVs0o9fReEP0vNb+pm
DKTAdg3oH9pLvlwJI6QPNh5Oh2byTYailnHwSHlgOrixP7unGZKuKiY8bb1uD5I+
IK4+s/Y4G67a2IWYR3p2WnfqbfRq7MkHRKrA87e2xXS5n3OMTZsv6YcYrbZKmQZ2
yF7F3W6WxoypB/uVvWpy3X1ZxXrgb15nZuF0TR4S5O2k9bUrBg7OQY0UYfD1br/b
iqba6pUIwiWNMbGOJ5iZFt/dew9yesFUUeRmEwkxUNujwIV+k3kHXfuJOIkZa38p
zM2GKhkCgYEA6NuQdJGVyi2z8gerrIE1/VYQL9sYN7vVUKOvy+mkj2foh4DL5FCv
dq6GEJ9UrGZOWev/QHNvz8bZtoVETj51u8Wlsp0D+b8WCJ8IP6GV6e7rB7eB4d+Z
KD+TPR2r7Az+6WMZe68/80oXEJggAvDeJzbG5ObpxBRlAlpXtt7yyAMCgYEA5xaI
thh2v8jqYlr5As0OCjHgZcvv6uEXzYfkN8w3Avtqq32EZha74JqlKXzNAxjJoYF2
h0kjSV0D0Ts0/2kQzoap6aoKhtWkAwzmac0WG+AZPPAFq5l6g2wEKvpiGiXM8vab
q4Cgw7Fo0Kqq8z5lwamY4+vD1ot0DdZNIm7vrCUCgYA5AH6lOnpLitKQ/fW1fc/k
myvNOzn7cryuR9Oh/CjvfgU7HnlLA8FgMSraaNaGeWjWtGHAukF1wHzNJGRrLvkN
JT4BslQlz/Qp2hxfz0Nuh7D7K53c2Cqa4q09ecT7PNct9LdpQqZJ/SoWQtcbQTFw
sgUQRcKV4FQ1tj3go0UVVwKBgQC09ThpIA8db7/a9VI5l0l/Qj9ud5yQWWPCVr+n
0griEu8dC1U6fGLzJyZerpP78NUz26VtmyA+us/acHq35xZ6I4m6qKVFoNambNuh
zi+Z9IrO5UYLckw1zcgVv6xCvYcYW3TbgAZkN/DUNlFX2WzlkmFfWagpwVpH26Db
bfPQ4QKBgQCI2u/NWeLZ7Ey+Mrd4xPsUhL0chK+FeSugjzw/aWQo2jDNPqyQgh54
3cnDKnIyVOido3TBkjFIPCMQhVeSNQGNhMAElbbmJ4t52ffAe+lpHSfIhss/Kj/r
jyRqXkkkajcCGNsC32OcQxu3ufaxoswDqK3DV9hSLnjsX8mFyeKDqw==
-----END RSA PRIVATE KEY-----`

	//LocalCa is a certificate authority than can be used on a local development env
	LocalCa = `-----BEGIN CERTIFICATE-----
MIIDvzCCAqegAwIBAgIUKQ0JtJHBL9KwCQlvJlJe/Ut+NhEwDQYJKoZIhvcNAQEL
BQAwZzELMAkGA1UEBhMCVVMxDzANBgNVBAgTBk9yZWdvbjERMA8GA1UEBxMIUG9y
dGxhbmQxEzARBgNVBAoTCkt1YmVybmV0ZXMxDTALBgNVBAsTBGdSUEMxEDAOBgNV
BAMTB2dSUEMgQ0EwHhcNMTYwODA5MjAwNjAwWhcNMjEwODA4MjAwNjAwWjBnMQsw
CQYDVQQGEwJVUzEPMA0GA1UECBMGT3JlZ29uMREwDwYDVQQHEwhQb3J0bGFuZDET
MBEGA1UEChMKS3ViZXJuZXRlczENMAsGA1UECxMEZ1JQQzEQMA4GA1UEAxMHZ1JQ
QyBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANhFmsqxDyANF5S7
AxkhDU172DJT+hpsfKMNRO6c33XUaWiSfRzw8je6Zdy7RCzgG5I/26mu0nnefb/T
K2XrSlwNETv59L5q/80rAR7MAxlfyifAvLS8cMkbWyskfM/QSA6jq59a1bvag3te
dRIDGc42ZGnuq2FNIxaoRQoblEltRY99dlqFj5mj2pD4c7RGhXgm6c2cQCmAjxYY
3nN0/zom9DSi7Jl+0httZVPLTbvDTtf4KznJ4iJflw1arTds3BRk7wfMz4m538Nk
sXtF7kHwRaP2/SMI8hrDnf2/unRm7Js4FdbwtXH3c3+XjdJkzH8BP4GTuCON+tOf
DtOVZQkCAwEAAaNjMGEwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8w
HQYDVR0OBBYEFIFP/fRth9gh/LGTziaIarbzpcwmMB8GA1UdIwQYMBaAFIFP/fRt
h9gh/LGTziaIarbzpcwmMA0GCSqGSIb3DQEBCwUAA4IBAQC0GSUzT2EcIJV2voXs
JED3DmG36vlo8vEeR4+gSFWLdw9UogwhfJoyiLYA4XG0g7GMVdAA0l3x+XhjFRh3
Gg0P/Xy+ga/z3AFFjoFC0YsMX2O4iGdh/VhbQcAUnd4ZU34Ap1AE+s2ekuOG13yA
xrO+3JSVbx+1cWg0yPvPIcKc1j1Yu8k+rV1r6t5a8gNTDtqNPqjGiTXPMN+Ecjkv
K5IVRwQQ6TVpyIcqW8xUGIZYiXMvC1BW1ClNy1WT+m9uCS2ZSiswUorTPcb2BTFf
broDjL8ptKtY1qBOBkfcMuv/5ACnW+aEvjD6aLMMvd/w4jhTp2ca5pVNTbta/IZ4
pV3m
-----END CERTIFICATE-----`
)
