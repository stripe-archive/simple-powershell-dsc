package dsc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRegistration(t *testing.T) {
	regKey := `f65e1a0c-46b0-424c-a6a5-c3701aef32e5`

	testCases := []struct {
		Authz   string
		MsDate  string
		Body    string
		Success bool
	}{
		// Test values from documentation here:
		//     https://github.com/PowerShellOrg/tug/blob/master/references/regkey-authorization.md
		{
			Authz:   `Shared SM095lQD5iEVzrToxnyuuoDAYfX2zA23YoZsZlZDyFU=`,
			MsDate:  `2016-12-21T23:43:48.4718366Z`,
			Body:    `{"AgentInformation":{"LCMVersion":"2.0","NodeName":"EC2AMAZ-VT1I874","IPAddress":"10.50.1.9;127.0.0.1;fe80::288e:6e98:1555:55e9%6;::2000:0:0:0;::1;::2000:0:0:0"},"ConfigurationNames":["ClientConfig2"],"RegistrationInformation":{"CertificateInformation":{"FriendlyName":"DSC-OaaS Client Authentication","Issuer":"CN=http://10.50.1.5:8080/PSDSCPullServer.svc","NotAfter":"2017-12-21T11:40:36.0000000-05:00","NotBefore":"2016-12-21T16:30:36.0000000-05:00","Subject":"CN=http://10.50.1.5:8080/PSDSCPullServer.svc","PublicKey":"U3lzdGVtLlNlY3VyaXR5LkNyeXB0b2dyYXBoeS5YNTA5Q2VydGlmaWNhdGVzLlB1YmxpY0tleQ==","Thumbprint":"AC5849ACDB6DD19FD79B6ACA2D077E71CEE31C4F","Version":3},"RegistrationMessageType":"ConfigurationRepository"}}`,
			Success: true,
		},

		// As above, but with a bad Authorization header
		{
			Authz:   `Shared AA095lQD5iEVzrToxnyuuoDAYfX2zA23YoZsZlZDyFU=`,
			MsDate:  `2016-12-21T23:43:48.4718366Z`,
			Body:    `{"AgentInformation":{"LCMVersion":"2.0","NodeName":"EC2AMAZ-VT1I874","IPAddress":"10.50.1.9;127.0.0.1;fe80::288e:6e98:1555:55e9%6;::2000:0:0:0;::1;::2000:0:0:0"},"ConfigurationNames":["ClientConfig2"],"RegistrationInformation":{"CertificateInformation":{"FriendlyName":"DSC-OaaS Client Authentication","Issuer":"CN=http://10.50.1.5:8080/PSDSCPullServer.svc","NotAfter":"2017-12-21T11:40:36.0000000-05:00","NotBefore":"2016-12-21T16:30:36.0000000-05:00","Subject":"CN=http://10.50.1.5:8080/PSDSCPullServer.svc","PublicKey":"U3lzdGVtLlNlY3VyaXR5LkNyeXB0b2dyYXBoeS5YNTA5Q2VydGlmaWNhdGVzLlB1YmxpY0tleQ==","Thumbprint":"AC5849ACDB6DD19FD79B6ACA2D077E71CEE31C4F","Version":3},"RegistrationMessageType":"ConfigurationRepository"}}`,
			Success: false,
		},

		// As above, but with a bad x-ms-date header
		{
			Authz:   `Shared SM095lQD5iEVzrToxnyuuoDAYfX2zA23YoZsZlZDyFU=`,
			MsDate:  `2006-12-21T23:43:48.4718366Z`,
			Body:    `{"AgentInformation":{"LCMVersion":"2.0","NodeName":"EC2AMAZ-VT1I874","IPAddress":"10.50.1.9;127.0.0.1;fe80::288e:6e98:1555:55e9%6;::2000:0:0:0;::1;::2000:0:0:0"},"ConfigurationNames":["ClientConfig2"],"RegistrationInformation":{"CertificateInformation":{"FriendlyName":"DSC-OaaS Client Authentication","Issuer":"CN=http://10.50.1.5:8080/PSDSCPullServer.svc","NotAfter":"2017-12-21T11:40:36.0000000-05:00","NotBefore":"2016-12-21T16:30:36.0000000-05:00","Subject":"CN=http://10.50.1.5:8080/PSDSCPullServer.svc","PublicKey":"U3lzdGVtLlNlY3VyaXR5LkNyeXB0b2dyYXBoeS5YNTA5Q2VydGlmaWNhdGVzLlB1YmxpY0tleQ==","Thumbprint":"AC5849ACDB6DD19FD79B6ACA2D077E71CEE31C4F","Version":3},"RegistrationMessageType":"ConfigurationRepository"}}`,
			Success: false,
		},

		// As above, but with a bad body
		{
			Authz:   `Shared SM095lQD5iEVzrToxnyuuoDAYfX2zA23YoZsZlZDyFU=`,
			MsDate:  `2016-12-21T23:43:48.4718366Z`,
			Body:    `{"AgentInformation":{"LCMVersion":"3.0","NodeName":"EC2AMAZ-VT1I874","IPAddress":"10.50.1.9;127.0.0.1;fe80::288e:6e98:1555:55e9%6;::2000:0:0:0;::1;::2000:0:0:0"},"ConfigurationNames":["ClientConfig2"],"RegistrationInformation":{"CertificateInformation":{"FriendlyName":"DSC-OaaS Client Authentication","Issuer":"CN=http://10.50.1.5:8080/PSDSCPullServer.svc","NotAfter":"2017-12-21T11:40:36.0000000-05:00","NotBefore":"2016-12-21T16:30:36.0000000-05:00","Subject":"CN=http://10.50.1.5:8080/PSDSCPullServer.svc","PublicKey":"U3lzdGVtLlNlY3VyaXR5LkNyeXB0b2dyYXBoeS5YNTA5Q2VydGlmaWNhdGVzLlB1YmxpY0tleQ==","Thumbprint":"AC5849ACDB6DD19FD79B6ACA2D077E71CEE31C4F","Version":3},"RegistrationMessageType":"ConfigurationRepository"}}`,
			Success: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc=%d", i), func(t *testing.T) {
			req, _ := http.NewRequest("PUT", "/something", strings.NewReader(tc.Body))
			req.Header.Set("Authorization", tc.Authz)
			req.Header.Set("x-ms-date", tc.MsDate)

			t.Logf("%d: headers = %+v", i, req.Header)

			called := false
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			})

			resp := httptest.NewRecorder()

			mware := checkRegistration([]string{regKey})
			wrappedHandler := mware(handler)
			wrappedHandler.ServeHTTP(resp, req)

			if tc.Success {
				assert.True(t, called)
				assert.Equal(t, http.StatusOK, resp.Code)
			} else {
				assert.False(t, called)
				assert.Equal(t, http.StatusUnauthorized, resp.Code)
			}
		})
	}
}
