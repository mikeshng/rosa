//go:build windows
// +build windows

/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains the function that returns the trusted CA certificates for Windows. This is
// needed because currently Go doesn't know how to load the Windows trusted CA store. See the
// following issues for more information:
//
//	https://github.com/golang/go/issues/16736
//	https://github.com/golang/go/issues/18609

package internal

import (
	"crypto/x509"
)

// loadSystemCAs loads the certificates of the CAs that we will trust. Currently this uses a fixed
// set of CA certificates, which is obviusly going to break in the future, but there is not much we
// can do (or know to do) till Go learns to read the Windows CA trust store.
func loadSystemCAs() (pool *x509.CertPool, err error) {
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(ssoCA1)
	pool.AppendCertsFromPEM(ssoCA2)
	pool.AppendCertsFromPEM(apiCA1)
	pool.AppendCertsFromPEM(apiCA2)
	return
}

// The SSO certificates has been obtained with the following command:
//
//	$ openssl s_client -connect sso.redhat.com:443 -showcerts
var ssoCA1 = []byte(`
-----BEGIN CERTIFICATE-----
MIIGbTCCBVWgAwIBAgIQDW8O1qMpxFU/O4crQErx7jANBgkqhkiG9w0BAQsFADB1
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMTQwMgYDVQQDEytEaWdpQ2VydCBTSEEyIEV4dGVuZGVk
IFZhbGlkYXRpb24gU2VydmVyIENBMB4XDTI0MDUwNzAwMDAwMFoXDTI1MDQyMDIz
NTk1OVowgcoxEzARBgsrBgEEAYI3PAIBAxMCVVMxGTAXBgsrBgEEAYI3PAIBAhMI
RGVsYXdhcmUxHTAbBgNVBA8MFFByaXZhdGUgT3JnYW5pemF0aW9uMRAwDgYDVQQF
EwcyOTQ1NDM2MQswCQYDVQQGEwJVUzEXMBUGA1UECBMOTm9ydGggQ2Fyb2xpbmEx
EDAOBgNVBAcTB1JhbGVpZ2gxFjAUBgNVBAoTDVJlZCBIYXQsIEluYy4xFzAVBgNV
BAMTDnNzby5yZWRoYXQuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwl7+
NobUYJIGiHRZf/B6psCRdWWWJWP5M6yC/3SKnvUVKSReDeaZ7hYNPcTZASsC3dOM
N54wvIxa/xDYFUiEzKOCA2wwggNoMB8GA1UdIwQYMBaAFD3TUKXWoK3u80pgCmXT
IdT4+NYPMB0GA1UdDgQWBBQ6jEVDz+ylBZm0aTC5G5tPtp/wLzAZBgNVHREEEjAQ
gg5zc28ucmVkaGF0LmNvbTBKBgNVHSAEQzBBMAsGCWCGSAGG/WwCATAyBgVngQwB
ATApMCcGCCsGAQUFBwIBFhtodHRwOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwDgYD
VR0PAQH/BAQDAgOIMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjB1BgNV
HR8EbjBsMDSgMqAwhi5odHRwOi8vY3JsMy5kaWdpY2VydC5jb20vc2hhMi1ldi1z
ZXJ2ZXItZzMuY3JsMDSgMqAwhi5odHRwOi8vY3JsNC5kaWdpY2VydC5jb20vc2hh
Mi1ldi1zZXJ2ZXItZzMuY3JsMIGIBggrBgEFBQcBAQR8MHowJAYIKwYBBQUHMAGG
GGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBSBggrBgEFBQcwAoZGaHR0cDovL2Nh
Y2VydHMuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0U0hBMkV4dGVuZGVkVmFsaWRhdGlv
blNlcnZlckNBLmNydDAMBgNVHRMBAf8EAjAAMIIBfgYKKwYBBAHWeQIEAgSCAW4E
ggFqAWgAdgDPEVbu1S58r/OHW9lpLpvpGnFnSrAX7KwB0lt3zsw7CAAAAY9Sx23S
AAAEAwBHMEUCIQD3ni/fGMOjKCzCIuvyUcBJMpjM7XKfrkHR+MxvFuqXmgIgLoss
OAc1y40putPNdg82Piu8igOIq3TeeO7zTI7Wc3oAdgB9WR4S4XgqexxhZ3xe/fjQ
h1wUoE6VnrkDL9kOjC55uAAAAY9Sx22SAAAEAwBHMEUCIEjQ/8mEJiKBTbRRkjZ/
10q2dMCCInzE0J8VV7FxbOOfAiEApydkGC8PM0f6azEGtKbi7vjhbGhtIfU46iSy
qUn08pUAdgDm0jFjQHeMwRBBBtdxuc7B0kD2loSG+7qHMh39HjeOUAAAAY9Sx22p
AAAEAwBHMEUCIGKoyNjFuFU8ScXufev+vO0bcmY588FqDQQT3XH54T+0AiEA43Pg
oMLz7rpPedYihEh8GzirmFpNg9GYOPh9QqWzGP8wDQYJKoZIhvcNAQELBQADggEB
AM9QyKcfZUj4oyWb0vyFcjSs30HYzBPsCvDkFaGs3ypVWg0+CodsUzR1JmN4d3bG
UFiBc/5/68GH1bVG74aa8y6hoHXfF2SbF8SJrHX6Qm8ZkuSvkPj9AKGnw1a3dwyH
utu92dI/4J8DSrVV6Wu1Puyvx+iWUuo/XLsnvqtwOGcMT00wG6NYMd/30CE6OB7M
4ONXa/j4Lnk9aOL4zkk4OM//2bwoQP+/9T66SUF7ACFpVBOP2GEs5w7TKGT5Wi9m
Qe8lkOoe4hqpMKmdj4wlkvfI4W5mVgUrqK6NAhl2gMhcJqPk7gzcZRtAQ1jyvbTf
HN9/ze44odnpb4zW5lMbcuw=
-----END CERTIFICATE-----
`)

var ssoCA2 = []byte(`
-----BEGIN CERTIFICATE-----
MIIEtjCCA56gAwIBAgIQDHmpRLCMEZUgkmFf4msdgzANBgkqhkiG9w0BAQsFADBs
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSswKQYDVQQDEyJEaWdpQ2VydCBIaWdoIEFzc3VyYW5j
ZSBFViBSb290IENBMB4XDTEzMTAyMjEyMDAwMFoXDTI4MTAyMjEyMDAwMFowdTEL
MAkGA1UEBhMCVVMxFTATBgNVBAoTDERpZ2lDZXJ0IEluYzEZMBcGA1UECxMQd3d3
LmRpZ2ljZXJ0LmNvbTE0MDIGA1UEAxMrRGlnaUNlcnQgU0hBMiBFeHRlbmRlZCBW
YWxpZGF0aW9uIFNlcnZlciBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBANdTpARR+JmmFkhLZyeqk0nQOe0MsLAAh/FnKIaFjI5j2ryxQDji0/XspQUY
uD0+xZkXMuwYjPrxDKZkIYXLBxA0sFKIKx9om9KxjxKws9LniB8f7zh3VFNfgHk/
LhqqqB5LKw2rt2O5Nbd9FLxZS99RStKh4gzikIKHaq7q12TWmFXo/a8aUGxUvBHy
/Urynbt/DvTVvo4WiRJV2MBxNO723C3sxIclho3YIeSwTQyJ3DkmF93215SF2AQh
cJ1vb/9cuhnhRctWVyh+HA1BV6q3uCe7seT6Ku8hI3UarS2bhjWMnHe1c63YlC3k
8wyd7sFOYn4XwHGeLN7x+RAoGTMCAwEAAaOCAUkwggFFMBIGA1UdEwEB/wQIMAYB
Af8CAQAwDgYDVR0PAQH/BAQDAgGGMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEF
BQcDAjA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRp
Z2ljZXJ0LmNvbTBLBgNVHR8ERDBCMECgPqA8hjpodHRwOi8vY3JsNC5kaWdpY2Vy
dC5jb20vRGlnaUNlcnRIaWdoQXNzdXJhbmNlRVZSb290Q0EuY3JsMD0GA1UdIAQ2
MDQwMgYEVR0gADAqMCgGCCsGAQUFBwIBFhxodHRwczovL3d3dy5kaWdpY2VydC5j
b20vQ1BTMB0GA1UdDgQWBBQ901Cl1qCt7vNKYApl0yHU+PjWDzAfBgNVHSMEGDAW
gBSxPsNpA/i/RwHUmCYaCALvY2QrwzANBgkqhkiG9w0BAQsFAAOCAQEAnbbQkIbh
hgLtxaDwNBx0wY12zIYKqPBKikLWP8ipTa18CK3mtlC4ohpNiAexKSHc59rGPCHg
4xFJcKx6HQGkyhE6V6t9VypAdP3THYUYUN9XR3WhfVUgLkc3UHKMf4Ib0mKPLQNa
2sPIoc4sUqIAY+tzunHISScjl2SFnjgOrWNoPLpSgVh5oywM395t6zHyuqB8bPEs
1OG9d4Q3A84ytciagRpKkk47RpqF/oOi+Z6Mo8wNXrM9zwR4jxQUezKcxwCmXMS1
oVWNWlZopCJwqjyBcdmdqEU79OX2olHdx3ti6G8MdOu42vi/hw15UJGQmxg7kVkn
8TUoE6smftX3eg==
-----END CERTIFICATE-----
`)

// The API certificates have been obtained with the following command:
//
//	$ openssl s_client -connect api.openshift.com:443 -showcerts
var apiCA1 = []byte(`
-----BEGIN CERTIFICATE-----
MIIE7zCCA9egAwIBAgISAx+yWIEh5U0//4KPNfNy1o92MA0GCSqGSIb3DQEBCwUA
MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
EwJSMzAeFw0yNDA1MjIxNTM1MzlaFw0yNDA4MjAxNTM1MzhaMBwxGjAYBgNVBAMT
EWFwaS5vcGVuc2hpZnQuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEA4sRpb1jErC7X3jEyG3oPnfWiq0us2YlgY/Xjs6t77gw7GE1dpDNGedgZQ6GG
FtunP8f20CFjHRpv8sSqC4AulDygQBIPhXOSeQ3oABVOcMm0qz3AsHUoA60jnA0N
3oQbMWdbdhTvonrx9t2XfIlE1zfai4BHdRTqJVzzZx7LSs0lxrc3/xMFz6668OsJ
saBLUe7eU4q5xhalGfRgENnkSqrlS16xGt71d/uTKL5epS2kd9v75TErkcQtcCSY
4loXhAUuxinya5Gfql86xw4yt6gPS6F0/It9SLs6u8P3uYA3DW3zV5TeAqqNDiUd
qHAjiGps21Una1L/utIk2P7bkQIDAQABo4ICEzCCAg8wDgYDVR0PAQH/BAQDAgWg
MB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0G
A1UdDgQWBBQ9k6sKfs7FfMkeGYjf6Eg/fe/N/DAfBgNVHSMEGDAWgBQULrMXt1hW
y65QCUDmH6+dixTCxjBVBggrBgEFBQcBAQRJMEcwIQYIKwYBBQUHMAGGFWh0dHA6
Ly9yMy5vLmxlbmNyLm9yZzAiBggrBgEFBQcwAoYWaHR0cDovL3IzLmkubGVuY3Iu
b3JnLzAcBgNVHREEFTATghFhcGkub3BlbnNoaWZ0LmNvbTATBgNVHSAEDDAKMAgG
BmeBDAECATCCAQQGCisGAQQB1nkCBAIEgfUEgfIA8AB2AD8XS0/XIkdYlB1lHIS+
DRLtkDd/H4Vq68G/KIXs+GRuAAABj6EoibkAAAQDAEcwRQIgaLG6upSKmj/dca1A
Qee5L3/bK/QyMfnyuhWRdcu197UCIQCNhW5cA5jhsFm6A9s6ejLnnCdLUuRHx6xT
/OoHwxa1AAB2AN/hVuuqBa+1nA+GcY2owDJOrlbZbqf1pWoB0cE7vlJcAAABj6Eo
i7QAAAQDAEcwRQIgeUduQpuLAgRKEbUM9SKsdc9OusXvVwlomDrkNtXo91wCIQDp
8kGl+ZjGygLJ2FYj1/wLzE7jsK4fAIA/1/rFC427jTANBgkqhkiG9w0BAQsFAAOC
AQEABifzI44dyzY8z9JfYrFIGuubxzHD9Op3XLVnr+WclVpBHebyD3oYjvN5ILXU
9ndOyxNs7mqvvL4cqhDhd37bUKOxHbnGUNSiS1UY0VH5kseuudfaYnzXw1JpbeIw
tevTkxkaO2KRVaGynDrhywxUabS3S+RwfNDTGf43v2Kj/cZeqXGy6z2TihzRn7U3
PrL3UdMYtZUkbNi70HFMXgemCbYE0lzU7EGxjVicxoRWuSQ/EHfVBCzAQm0Gy2Om
/AmNVPwoea6TkXuwM1GHnRlt+N3GQgNqrnC+QxIzCOb+A6IvFr8rd+zb4R1K0ngN
315oufYhHZrYyQ11NtyDz1v84Q==
-----END CERTIFICATE-----
`)

var apiCA2 = []byte(`
-----BEGIN CERTIFICATE-----
MIIFFjCCAv6gAwIBAgIRAJErCErPDBinU/bWLiWnX1owDQYJKoZIhvcNAQELBQAw
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjAwOTA0MDAwMDAw
WhcNMjUwOTE1MTYwMDAwWjAyMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg
RW5jcnlwdDELMAkGA1UEAxMCUjMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQC7AhUozPaglNMPEuyNVZLD+ILxmaZ6QoinXSaqtSu5xUyxr45r+XXIo9cP
R5QUVTVXjJ6oojkZ9YI8QqlObvU7wy7bjcCwXPNZOOftz2nwWgsbvsCUJCWH+jdx
sxPnHKzhm+/b5DtFUkWWqcFTzjTIUu61ru2P3mBw4qVUq7ZtDpelQDRrK9O8Zutm
NHz6a4uPVymZ+DAXXbpyb/uBxa3Shlg9F8fnCbvxK/eG3MHacV3URuPMrSXBiLxg
Z3Vms/EY96Jc5lP/Ooi2R6X/ExjqmAl3P51T+c8B5fWmcBcUr2Ok/5mzk53cU6cG
/kiFHaFpriV1uxPMUgP17VGhi9sVAgMBAAGjggEIMIIBBDAOBgNVHQ8BAf8EBAMC
AYYwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMBMBIGA1UdEwEB/wQIMAYB
Af8CAQAwHQYDVR0OBBYEFBQusxe3WFbLrlAJQOYfr52LFMLGMB8GA1UdIwQYMBaA
FHm0WeZ7tuXkAXOACIjIGlj26ZtuMDIGCCsGAQUFBwEBBCYwJDAiBggrBgEFBQcw
AoYWaHR0cDovL3gxLmkubGVuY3Iub3JnLzAnBgNVHR8EIDAeMBygGqAYhhZodHRw
Oi8veDEuYy5sZW5jci5vcmcvMCIGA1UdIAQbMBkwCAYGZ4EMAQIBMA0GCysGAQQB
gt8TAQEBMA0GCSqGSIb3DQEBCwUAA4ICAQCFyk5HPqP3hUSFvNVneLKYY611TR6W
PTNlclQtgaDqw+34IL9fzLdwALduO/ZelN7kIJ+m74uyA+eitRY8kc607TkC53wl
ikfmZW4/RvTZ8M6UK+5UzhK8jCdLuMGYL6KvzXGRSgi3yLgjewQtCPkIVz6D2QQz
CkcheAmCJ8MqyJu5zlzyZMjAvnnAT45tRAxekrsu94sQ4egdRCnbWSDtY7kh+BIm
lJNXoB1lBMEKIq4QDUOXoRgffuDghje1WrG9ML+Hbisq/yFOGwXD9RiX8F6sw6W4
avAuvDszue5L3sz85K+EC4Y/wFVDNvZo4TYXao6Z0f+lQKc0t8DQYzk1OXVu8rp2
yJMC6alLbBfODALZvYH7n7do1AZls4I9d1P4jnkDrQoxB3UqQ9hVl3LEKQ73xF1O
yK5GhDDX8oVfGKF5u+decIsH4YaTw7mP3GFxJSqv3+0lUFJoi5Lc5da149p90Ids
hCExroL1+7mryIkXPeFM5TgO9r0rvZaBFOvV2z0gp35Z0+L4WPlbuEjN/lxPFin+
HlUjr8gRsI3qfJOQFy/9rKIJR0Y/8Omwt/8oTWgy1mdeHmmjk7j1nYsvC9JSQ6Zv
MldlTTKB3zhThV1+XWYp6rjd5JW1zbVWEkLNxE7GJThEUG3szgBVGP7pSWTUTsqX
nLRbwHOoq7hHwg==
-----END CERTIFICATE-----
`)
