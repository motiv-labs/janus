package api_test

import (
	"encoding/json"
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstanceOfDefinition(t *testing.T) {
	instance := api.NewDefinition()

	assert.IsType(t, &api.Definition{}, instance)
	assert.True(t, instance.Active)
}

func TestSuccessfulValidation(t *testing.T) {
	instance := api.NewDefinition()
	instance.Name = "Test"
	instance.Proxy.ListenPath = "/"
	instance.Proxy.Upstreams = &proxy.Upstreams{
		Balancing: "roundrobin",
		Targets: []*proxy.Target{
			{Target: "http:/example.com"},
		},
	}

	isValid, err := instance.Validate()
	require.NoError(t, err)
	assert.True(t, isValid)
}

func TestFailedValidation(t *testing.T) {
	instance := api.NewDefinition()
	isValid, err := instance.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestNameValidation(t *testing.T) {
	instanceSimple := api.NewDefinition()
	instanceSimple.Name = "simple"
	instanceSimple.Proxy.ListenPath = "/"
	isValid, err := instanceSimple.Validate()

	require.NoError(t, err)
	require.True(t, isValid)

	instanceDash := api.NewDefinition()
	instanceDash.Name = "with-dash-and-123"
	instanceDash.Proxy.ListenPath = "/"
	isValid, err = instanceDash.Validate()

	require.NoError(t, err)
	require.True(t, isValid)

	instanceBadSymbol := api.NewDefinition()
	instanceBadSymbol.Name = "test~"
	instanceBadSymbol.Proxy.ListenPath = "/"
	isValid, err = instanceBadSymbol.Validate()

	require.Error(t, err)
	require.False(t, isValid)
}

func TestConfiguration_EqualsTo(t *testing.T) {
	def11 := api.NewDefinition()
	def12 := api.NewDefinition()
	def13 := api.NewDefinition()
	def14 := api.NewDefinition()

	def21 := api.NewDefinition()
	def22 := api.NewDefinition()
	def23 := api.NewDefinition()
	def24 := api.NewDefinition()

	require.NoError(t, json.Unmarshal([]byte(def1), &def11))
	require.NoError(t, json.Unmarshal([]byte(def2), &def12))
	require.NoError(t, json.Unmarshal([]byte(def3), &def13))
	require.NoError(t, json.Unmarshal([]byte(def4), &def14))

	require.NoError(t, json.Unmarshal([]byte(def1), &def21))
	require.NoError(t, json.Unmarshal([]byte(def2), &def22))
	require.NoError(t, json.Unmarshal([]byte(def3), &def23))
	require.NoError(t, json.Unmarshal([]byte(def4), &def24))

	c1 := &api.Configuration{
		Definitions: []*api.Definition{def11, def12, def13, def14},
	}
	c2 := &api.Configuration{
		Definitions: []*api.Definition{def21, def22, def23, def24},
	}

	assert.True(t, c1.EqualsTo(c2))
	assert.True(t, c2.EqualsTo(c1))
}

const (
	def1 = `{
    "name" : "users",
    "active" : true,
    "proxy" : {
        "preserve_host" : false,
        "listen_path" : "/users/*",
        "upstreams" : {
            "balancing" : "weight",
            "targets" : [
                {
                    "target" : "http://localhost:8000/users",
                    "weight" : 0
                },
                {
                    "target" : "http://auth-service.live-k8s.hellofresh.io/users",
                    "weight" : 100
                }
            ]
        },
        "insecure_skip_verify" : false,
        "strip_path" : true,
        "append_path" : false,
        "enable_load_balancing" : false,
        "methods" : [
            "ALL"
        ],
        "hosts" : []
    },
    "plugins" : [
        {
            "name" : "cors",
            "enabled" : true,
            "config" : {
                "request_headers" : [
                    "Accept",
                    "Accept-CH",
                    "Accept-Charset",
                    "Accept-Datetime",
                    "Accept-Encoding",
                    "Accept-Ext",
                    "Accept-Features",
                    "Accept-Language",
                    "Accept-Params",
                    "Accept-Ranges",
                    "Access-Control-Allow-Credentials",
                    "Access-Control-Allow-Headers",
                    "Access-Control-Allow-Methods",
                    "Access-Control-Allow-Origin",
                    "Access-Control-Expose-Headers",
                    "Access-Control-Max-Age",
                    "Access-Control-Request-Headers",
                    "Access-Control-Request-Method",
                    "Age",
                    "Allow",
                    "Alternates",
                    "Authentication-Info",
                    "Authorization",
                    "C-Ext",
                    "C-Man",
                    "C-Opt",
                    "C-PEP",
                    "C-PEP-Info",
                    "CONNECT",
                    "Cache-Control",
                    "Compliance",
                    "Connection",
                    "Content-Base",
                    "Content-Disposition",
                    "Content-Encoding",
                    "Content-ID",
                    "Content-Language",
                    "Content-Length",
                    "Content-Location",
                    "Content-MD5",
                    "Content-Range",
                    "Content-Script-Type",
                    "Content-Security-Policy",
                    "Content-Style-Type",
                    "Content-Transfer-Encoding",
                    "Content-Type",
                    "Content-Version",
                    "Cookie",
                    "Cost",
                    "DAV",
                    "DELETE",
                    "DNT",
                    "DPR",
                    "Date",
                    "Default-Style",
                    "Delta-Base",
                    "Depth",
                    "Derived-From",
                    "Destination",
                    "Differential-ID",
                    "Digest",
                    "ETag",
                    "Expect",
                    "Expires",
                    "Ext",
                    "From",
                    "GET",
                    "GetProfile",
                    "HEAD",
                    "HTTP-date",
                    "Host",
                    "IM",
                    "If",
                    "If-Match",
                    "If-Modified-Since",
                    "If-None-Match",
                    "If-Range",
                    "If-Unmodified-Since",
                    "Keep-Alive",
                    "Label",
                    "Last-Event-ID",
                    "Last-Modified",
                    "Link",
                    "Location",
                    "Lock-Token",
                    "MIME-Version",
                    "Man",
                    "Max-Forwards",
                    "Media-Range",
                    "Message-ID",
                    "Meter",
                    "Negotiate",
                    "Non-Compliance",
                    "OPTION",
                    "OPTIONS",
                    "OWS",
                    "Opt",
                    "Optional",
                    "Ordering-Type",
                    "Origin",
                    "Overwrite",
                    "P3P",
                    "PEP",
                    "PICS-Label",
                    "POST",
                    "PUT",
                    "Pep-Info",
                    "Permanent",
                    "Position",
                    "Pragma",
                    "ProfileObject",
                    "Protocol",
                    "Protocol-Query",
                    "Protocol-Request",
                    "Proxy-Authenticate",
                    "Proxy-Authentication-Info",
                    "Proxy-Authorization",
                    "Proxy-Features",
                    "Proxy-Instruction",
                    "Public",
                    "RWS",
                    "Range",
                    "Referer",
                    "Refresh",
                    "Resolution-Hint",
                    "Resolver-Location",
                    "Retry-After",
                    "Safe",
                    "Sec-Websocket-Extensions",
                    "Sec-Websocket-Key",
                    "Sec-Websocket-Origin",
                    "Sec-Websocket-Protocol",
                    "Sec-Websocket-Version",
                    "Security-Scheme",
                    "Server",
                    "Set-Cookie",
                    "Set-Cookie2",
                    "SetProfile",
                    "SoapAction",
                    "Status",
                    "Status-URI",
                    "Strict-Transport-Security",
                    "SubOK",
                    "Subst",
                    "Surrogate-Capability",
                    "Surrogate-Control",
                    "TCN",
                    "TE",
                    "TRACE",
                    "Timeout",
                    "Title",
                    "Trailer",
                    "Transfer-Encoding",
                    "UA-Color",
                    "UA-Media",
                    "UA-Pixels",
                    "UA-Resolution",
                    "UA-Windowpixels",
                    "URI",
                    "Upgrade",
                    "User-Agent",
                    "Variant-Vary",
                    "Vary",
                    "Version",
                    "Via",
                    "Viewport-Width",
                    "WWW-Authenticate",
                    "Want-Digest",
                    "Warning",
                    "Width",
                    "X-Content-Duration",
                    "X-Content-Security-Policy",
                    "X-Content-Type-Options",
                    "X-CustomHeader",
                    "X-DNSPrefetch-Control",
                    "X-Forwarded-For",
                    "X-Forwarded-Port",
                    "X-Forwarded-Proto",
                    "X-Frame-Options",
                    "X-Modified",
                    "X-OTHER",
                    "X-PING",
                    "X-PINGOTHER",
                    "X-Powered-By",
                    "X-Requested-With"
                ],
                "domains" : [
                    "*"
                ],
                "exposed_headers" : [
                    "X-Debug-Token",
                    "X-Debug-Token-Link"
                ],
                "methods" : [
                    "CONNECT",
                    "DEBUG",
                    "DELETE",
                    "DONE",
                    "GET",
                    "HEAD",
                    "HTTP",
                    "HTTP/0.9",
                    "HTTP/1.0",
                    "HTTP/1.1",
                    "HTTP/2",
                    "OPTIONS",
                    "ORIGIN",
                    "ORIGINS",
                    "PATCH",
                    "POST",
                    "PUT",
                    "QUIC",
                    "REST",
                    "SESSION",
                    "SHOULD",
                    "SPDY",
                    "TRACE",
                    "TRACK"
                ]
            }
        },
        {
            "name" : "rate_limit",
            "enabled" : false,
            "config" : {
                "limit" : "50-S",
                "policy" : "local"
            }
        },
        {
            "name" : "oauth2",
            "enabled" : true,
            "config" : {
                "server_name" : "live oauth server"
            }
        },
        {
            "name" : "compression",
            "enabled" : false,
            "config" : {}
        }
    ],
    "health_check" : {
        "url" : "",
        "timeout" : 0
    }
}`
	def2 = `{
    "name" : "csi-users",
    "active" : true,
    "domain" : "",
    "proxy" : {
        "listen_path" : "/csi/user/*",
        "append_listen_path" : false,
        "target_list" : [],
        "methods" : [
            "ALL"
        ],
        "preserve_host" : false,
        "strip_path" : true,
        "upstreams" : {
            "balancing" : "roundrobin",
            "targets" : [
                {
                    "target" : "http://localhost:5004/user"
                }
            ]
        }
    },
    "allowed_ips" : [],
    "oauth_server_name" : "live oauth server",
    "plugins" : [
        {
            "name" : "cors",
            "enabled" : true,
            "config" : {
                "domains" : [
                    "*"
                ],
                "methods" : [
                    "CONNECT",
                    "DEBUG",
                    "DELETE",
                    "DONE",
                    "GET",
                    "HEAD",
                    "HTTP",
                    "HTTP/0.9",
                    "HTTP/1.0",
                    "HTTP/1.1",
                    "HTTP/2",
                    "OPTIONS",
                    "ORIGIN",
                    "ORIGINS",
                    "PATCH",
                    "POST",
                    "PUT",
                    "QUIC",
                    "REST",
                    "SESSION",
                    "SHOULD",
                    "SPDY",
                    "TRACE",
                    "TRACK"
                ],
                "request_headers" : [
                    "Accept",
                    "Accept-CH",
                    "Accept-Charset",
                    "Accept-Datetime",
                    "Accept-Encoding",
                    "Accept-Ext",
                    "Accept-Features",
                    "Accept-Language",
                    "Accept-Params",
                    "Accept-Ranges",
                    "Access-Control-Allow-Credentials",
                    "Access-Control-Allow-Headers",
                    "Access-Control-Allow-Methods",
                    "Access-Control-Allow-Origin",
                    "Access-Control-Expose-Headers",
                    "Access-Control-Max-Age",
                    "Access-Control-Request-Headers",
                    "Access-Control-Request-Method",
                    "Age",
                    "Allow",
                    "Alternates",
                    "Authentication-Info",
                    "Authorization",
                    "C-Ext",
                    "C-Man",
                    "C-Opt",
                    "C-PEP",
                    "C-PEP-Info",
                    "CONNECT",
                    "Cache-Control",
                    "Compliance",
                    "Connection",
                    "Content-Base",
                    "Content-Disposition",
                    "Content-Encoding",
                    "Content-ID",
                    "Content-Language",
                    "Content-Length",
                    "Content-Location",
                    "Content-MD5",
                    "Content-Range",
                    "Content-Script-Type",
                    "Content-Security-Policy",
                    "Content-Style-Type",
                    "Content-Transfer-Encoding",
                    "Content-Type",
                    "Content-Version",
                    "Cookie",
                    "Cost",
                    "DAV",
                    "DELETE",
                    "DNT",
                    "DPR",
                    "Date",
                    "Default-Style",
                    "Delta-Base",
                    "Depth",
                    "Derived-From",
                    "Destination",
                    "Differential-ID",
                    "Digest",
                    "ETag",
                    "Expect",
                    "Expires",
                    "Ext",
                    "From",
                    "GET",
                    "GetProfile",
                    "HEAD",
                    "HTTP-date",
                    "Host",
                    "IM",
                    "If",
                    "If-Match",
                    "If-Modified-Since",
                    "If-None-Match",
                    "If-Range",
                    "If-Unmodified-Since",
                    "Keep-Alive",
                    "Label",
                    "Last-Event-ID",
                    "Last-Modified",
                    "Link",
                    "Location",
                    "Lock-Token",
                    "MIME-Version",
                    "Man",
                    "Max-Forwards",
                    "Media-Range",
                    "Message-ID",
                    "Meter",
                    "Negotiate",
                    "Non-Compliance",
                    "OPTION",
                    "OPTIONS",
                    "OWS",
                    "Opt",
                    "Optional",
                    "Ordering-Type",
                    "Origin",
                    "Overwrite",
                    "P3P",
                    "PEP",
                    "PICS-Label",
                    "POST",
                    "PUT",
                    "Pep-Info",
                    "Permanent",
                    "Position",
                    "Pragma",
                    "ProfileObject",
                    "Protocol",
                    "Protocol-Query",
                    "Protocol-Request",
                    "Proxy-Authenticate",
                    "Proxy-Authentication-Info",
                    "Proxy-Authorization",
                    "Proxy-Features",
                    "Proxy-Instruction",
                    "Public",
                    "RWS",
                    "Range",
                    "Referer",
                    "Refresh",
                    "Resolution-Hint",
                    "Resolver-Location",
                    "Retry-After",
                    "Safe",
                    "Sec-Websocket-Extensions",
                    "Sec-Websocket-Key",
                    "Sec-Websocket-Origin",
                    "Sec-Websocket-Protocol",
                    "Sec-Websocket-Version",
                    "Security-Scheme",
                    "Server",
                    "Set-Cookie",
                    "Set-Cookie2",
                    "SetProfile",
                    "SoapAction",
                    "Status",
                    "Status-URI",
                    "Strict-Transport-Security",
                    "SubOK",
                    "Subst",
                    "Surrogate-Capability",
                    "Surrogate-Control",
                    "TCN",
                    "TE",
                    "TRACE",
                    "Timeout",
                    "Title",
                    "Trailer",
                    "Transfer-Encoding",
                    "UA-Color",
                    "UA-Media",
                    "UA-Pixels",
                    "UA-Resolution",
                    "UA-Windowpixels",
                    "URI",
                    "Upgrade",
                    "User-Agent",
                    "Variant-Vary",
                    "Vary",
                    "Version",
                    "Via",
                    "Viewport-Width",
                    "WWW-Authenticate",
                    "Want-Digest",
                    "Warning",
                    "Width",
                    "X-Content-Duration",
                    "X-Content-Security-Policy",
                    "X-Content-Type-Options",
                    "X-CustomHeader",
                    "X-DNSPrefetch-Control",
                    "X-Forwarded-For",
                    "X-Forwarded-Port",
                    "X-Forwarded-Proto",
                    "X-Frame-Options",
                    "X-Modified",
                    "X-OTHER",
                    "X-PING",
                    "X-PINGOTHER",
                    "X-Powered-By",
                    "X-Requested-With"
                ],
                "exposed_headers" : [
                    "X-Debug-Token",
                    "X-Debug-Token-Link"
                ]
            }
        },
        {
            "name" : "rate_limit",
            "enabled" : false,
            "config" : {
                "limit" : "50-S",
                "policy" : "local"
            }
        },
        {
            "name" : "oauth2",
            "enabled" : true,
            "config" : {
                "server_name" : "live oauth server"
            }
        },
        {
            "name" : "compression",
            "enabled" : false
        }
    ]
}`
	def3 = `{
    "name" : "sku-mapper",
    "active" : true,
    "proxy" : {
        "preserve_host" : false,
        "listen_path" : "/scm/recipes/*",
        "upstreams" : {
            "balancing" : "roundrobin",
            "targets" : [
                {
                    "target" : "http://sku-mapper.live-k8s.hellofresh.io/recipes",
                    "weight" : 50
                }
            ]
        },
        "insecure_skip_verify" : false,
        "strip_path" : true,
        "append_path" : false,
        "enable_load_balancing" : false,
        "methods" : [
            "ALL"
        ],
        "hosts" : []
    },
    "plugins" : [
        {
            "name" : "cors",
            "enabled" : true,
            "config" : {
                "request_headers" : [
                    "Accept",
                    "Accept-CH",
                    "Accept-Charset",
                    "Accept-Datetime",
                    "Accept-Encoding",
                    "Accept-Ext",
                    "Accept-Features",
                    "Accept-Language",
                    "Accept-Params",
                    "Accept-Ranges",
                    "Access-Control-Allow-Credentials",
                    "Access-Control-Allow-Headers",
                    "Access-Control-Allow-Methods",
                    "Access-Control-Allow-Origin",
                    "Access-Control-Expose-Headers",
                    "Access-Control-Max-Age",
                    "Access-Control-Request-Headers",
                    "Access-Control-Request-Method",
                    "Age",
                    "Allow",
                    "Alternates",
                    "Authentication-Info",
                    "Authorization",
                    "C-Ext",
                    "C-Man",
                    "C-Opt",
                    "C-PEP",
                    "C-PEP-Info",
                    "CONNECT",
                    "Cache-Control",
                    "Compliance",
                    "Connection",
                    "Content-Base",
                    "Content-Disposition",
                    "Content-Encoding",
                    "Content-ID",
                    "Content-Language",
                    "Content-Length",
                    "Content-Location",
                    "Content-MD5",
                    "Content-Range",
                    "Content-Script-Type",
                    "Content-Security-Policy",
                    "Content-Style-Type",
                    "Content-Transfer-Encoding",
                    "Content-Type",
                    "Content-Version",
                    "Cookie",
                    "Cost",
                    "DAV",
                    "DELETE",
                    "DNT",
                    "DPR",
                    "Date",
                    "Default-Style",
                    "Delta-Base",
                    "Depth",
                    "Derived-From",
                    "Destination",
                    "Differential-ID",
                    "Digest",
                    "ETag",
                    "Expect",
                    "Expires",
                    "Ext",
                    "From",
                    "GET",
                    "GetProfile",
                    "HEAD",
                    "HTTP-date",
                    "Host",
                    "IM",
                    "If",
                    "If-Match",
                    "If-Modified-Since",
                    "If-None-Match",
                    "If-Range",
                    "If-Unmodified-Since",
                    "Keep-Alive",
                    "Label",
                    "Last-Event-ID",
                    "Last-Modified",
                    "Link",
                    "Location",
                    "Lock-Token",
                    "MIME-Version",
                    "Man",
                    "Max-Forwards",
                    "Media-Range",
                    "Message-ID",
                    "Meter",
                    "Negotiate",
                    "Non-Compliance",
                    "OPTION",
                    "OPTIONS",
                    "OWS",
                    "Opt",
                    "Optional",
                    "Ordering-Type",
                    "Origin",
                    "Overwrite",
                    "P3P",
                    "PEP",
                    "PICS-Label",
                    "POST",
                    "PUT",
                    "Pep-Info",
                    "Permanent",
                    "Position",
                    "Pragma",
                    "ProfileObject",
                    "Protocol",
                    "Protocol-Query",
                    "Protocol-Request",
                    "Proxy-Authenticate",
                    "Proxy-Authentication-Info",
                    "Proxy-Authorization",
                    "Proxy-Features",
                    "Proxy-Instruction",
                    "Public",
                    "RWS",
                    "Range",
                    "Referer",
                    "Refresh",
                    "Resolution-Hint",
                    "Resolver-Location",
                    "Retry-After",
                    "Safe",
                    "Sec-Websocket-Extensions",
                    "Sec-Websocket-Key",
                    "Sec-Websocket-Origin",
                    "Sec-Websocket-Protocol",
                    "Sec-Websocket-Version",
                    "Security-Scheme",
                    "Server",
                    "Set-Cookie",
                    "Set-Cookie2",
                    "SetProfile",
                    "SoapAction",
                    "Status",
                    "Status-URI",
                    "Strict-Transport-Security",
                    "SubOK",
                    "Subst",
                    "Surrogate-Capability",
                    "Surrogate-Control",
                    "TCN",
                    "TE",
                    "TRACE",
                    "Timeout",
                    "Title",
                    "Trailer",
                    "Transfer-Encoding",
                    "UA-Color",
                    "UA-Media",
                    "UA-Pixels",
                    "UA-Resolution",
                    "UA-Windowpixels",
                    "URI",
                    "Upgrade",
                    "User-Agent",
                    "Variant-Vary",
                    "Vary",
                    "Version",
                    "Via",
                    "Viewport-Width",
                    "WWW-Authenticate",
                    "Want-Digest",
                    "Warning",
                    "Width",
                    "X-Content-Duration",
                    "X-Content-Security-Policy",
                    "X-Content-Type-Options",
                    "X-CustomHeader",
                    "X-DNSPrefetch-Control",
                    "X-Forwarded-For",
                    "X-Forwarded-Port",
                    "X-Forwarded-Proto",
                    "X-Frame-Options",
                    "X-Modified",
                    "X-OTHER",
                    "X-PING",
                    "X-PINGOTHER",
                    "X-Powered-By",
                    "X-Requested-With"
                ],
                "domains" : [
                    "*"
                ],
                "exposed_headers" : [
                    "X-Debug-Token",
                    "X-Debug-Token-Link"
                ],
                "methods" : [
                    "CONNECT",
                    "DEBUG",
                    "DELETE",
                    "DONE",
                    "GET",
                    "HEAD",
                    "HTTP",
                    "HTTP/0.9",
                    "HTTP/1.0",
                    "HTTP/1.1",
                    "HTTP/2",
                    "OPTIONS",
                    "ORIGIN",
                    "ORIGINS",
                    "PATCH",
                    "POST",
                    "PUT",
                    "QUIC",
                    "REST",
                    "SESSION",
                    "SHOULD",
                    "SPDY",
                    "TRACE",
                    "TRACK"
                ]
            }
        },
        {
            "name" : "rate_limit",
            "enabled" : false,
            "config" : {
                "limit" : "50-S",
                "policy" : "local"
            }
        },
        {
            "name" : "oauth2",
            "enabled" : true,
            "config" : {
                "server_name" : "live oauth server"
            }
        },
        {
            "name" : "compression",
            "enabled" : false,
            "config" : {}
        }
    ],
    "health_check" : {
        "url" : "",
        "timeout" : 0
    }
}`
	def4 = `{
    "name" : "Bacchus-wine-service",
    "active" : true,
    "proxy" : {
        "preserve_host" : false,
        "listen_path" : "/wineclub/*",
        "upstreams" : {
            "balancing" : "roundrobin",
            "targets" : [
                {
                    "target" : "http://bacchus-service.staging-k8s.hellofresh.io",
                    "weight" : 0
                }
            ]
        },
        "insecure_skip_verify" : false,
        "strip_path" : true,
        "append_path" : false,
        "enable_load_balancing" : false,
        "methods" : [
            "ALL"
        ],
        "hosts" : []
    },
    "plugins" : [
        {
            "name" : "cors",
            "enabled" : true,
            "config" : {
                "domains" : [
                    "*"
                ],
                "exposed_headers" : [
                    "X-Debug-Token",
                    "X-Debug-Token-Link"
                ],
                "methods" : [
                    "CONNECT",
                    "DEBUG",
                    "DELETE",
                    "DONE",
                    "GET",
                    "HEAD",
                    "HTTP",
                    "HTTP/0.9",
                    "HTTP/1.0",
                    "HTTP/1.1",
                    "HTTP/2",
                    "OPTIONS",
                    "ORIGIN",
                    "ORIGINS",
                    "PATCH",
                    "POST",
                    "PUT",
                    "QUIC",
                    "REST",
                    "SESSION",
                    "SHOULD",
                    "SPDY",
                    "TRACE",
                    "TRACK"
                ],
                "request_headers" : [
                    "Accept",
                    "Accept-CH",
                    "Accept-Charset",
                    "Accept-Datetime",
                    "Accept-Encoding",
                    "Accept-Ext",
                    "Accept-Features",
                    "Accept-Language",
                    "Accept-Params",
                    "Accept-Ranges",
                    "Access-Control-Allow-Credentials",
                    "Access-Control-Allow-Headers",
                    "Access-Control-Allow-Methods",
                    "Access-Control-Allow-Origin",
                    "Access-Control-Expose-Headers",
                    "Access-Control-Max-Age",
                    "Access-Control-Request-Headers",
                    "Access-Control-Request-Method",
                    "Age",
                    "Allow",
                    "Alternates",
                    "Authentication-Info",
                    "Authorization",
                    "C-Ext",
                    "C-Man",
                    "C-Opt",
                    "C-PEP",
                    "C-PEP-Info",
                    "CONNECT",
                    "Cache-Control",
                    "Compliance",
                    "Connection",
                    "Content-Base",
                    "Content-Disposition",
                    "Content-Encoding",
                    "Content-ID",
                    "Content-Language",
                    "Content-Length",
                    "Content-Location",
                    "Content-MD5",
                    "Content-Range",
                    "Content-Script-Type",
                    "Content-Security-Policy",
                    "Content-Style-Type",
                    "Content-Transfer-Encoding",
                    "Content-Type",
                    "Content-Version",
                    "Cookie",
                    "Cost",
                    "DAV",
                    "DELETE",
                    "DNT",
                    "DPR",
                    "Date",
                    "Default-Style",
                    "Delta-Base",
                    "Depth",
                    "Derived-From",
                    "Destination",
                    "Differential-ID",
                    "Digest",
                    "ETag",
                    "Expect",
                    "Expires",
                    "Ext",
                    "From",
                    "GET",
                    "GetProfile",
                    "HEAD",
                    "HTTP-date",
                    "Host",
                    "IM",
                    "If",
                    "If-Match",
                    "If-Modified-Since",
                    "If-None-Match",
                    "If-Range",
                    "If-Unmodified-Since",
                    "Keep-Alive",
                    "Label",
                    "Last-Event-ID",
                    "Last-Modified",
                    "Link",
                    "Location",
                    "Lock-Token",
                    "MIME-Version",
                    "Man",
                    "Max-Forwards",
                    "Media-Range",
                    "Message-ID",
                    "Meter",
                    "Negotiate",
                    "Non-Compliance",
                    "OPTION",
                    "OPTIONS",
                    "OWS",
                    "Opt",
                    "Optional",
                    "Ordering-Type",
                    "Origin",
                    "Overwrite",
                    "P3P",
                    "PEP",
                    "PICS-Label",
                    "POST",
                    "PUT",
                    "Pep-Info",
                    "Permanent",
                    "Position",
                    "Pragma",
                    "ProfileObject",
                    "Protocol",
                    "Protocol-Query",
                    "Protocol-Request",
                    "Proxy-Authenticate",
                    "Proxy-Authentication-Info",
                    "Proxy-Authorization",
                    "Proxy-Features",
                    "Proxy-Instruction",
                    "Public",
                    "RWS",
                    "Range",
                    "Referer",
                    "Refresh",
                    "Resolution-Hint",
                    "Resolver-Location",
                    "Retry-After",
                    "Safe",
                    "Sec-Websocket-Extensions",
                    "Sec-Websocket-Key",
                    "Sec-Websocket-Origin",
                    "Sec-Websocket-Protocol",
                    "Sec-Websocket-Version",
                    "Security-Scheme",
                    "Server",
                    "Set-Cookie",
                    "Set-Cookie2",
                    "SetProfile",
                    "SoapAction",
                    "Status",
                    "Status-URI",
                    "Strict-Transport-Security",
                    "SubOK",
                    "Subst",
                    "Surrogate-Capability",
                    "Surrogate-Control",
                    "TCN",
                    "TE",
                    "TRACE",
                    "Timeout",
                    "Title",
                    "Trailer",
                    "Transfer-Encoding",
                    "UA-Color",
                    "UA-Media",
                    "UA-Pixels",
                    "UA-Resolution",
                    "UA-Windowpixels",
                    "URI",
                    "Upgrade",
                    "User-Agent",
                    "Variant-Vary",
                    "Vary",
                    "Version",
                    "Via",
                    "Viewport-Width",
                    "WWW-Authenticate",
                    "Want-Digest",
                    "Warning",
                    "Width",
                    "X-Content-Duration",
                    "X-Content-Security-Policy",
                    "X-Content-Type-Options",
                    "X-CustomHeader",
                    "X-DNSPrefetch-Control",
                    "X-Forwarded-For",
                    "X-Forwarded-Port",
                    "X-Forwarded-Proto",
                    "X-Frame-Options",
                    "X-Modified",
                    "X-OTHER",
                    "X-PING",
                    "X-PINGOTHER",
                    "X-Powered-By",
                    "X-Requested-With"
                ]
            }
        },
        {
            "name" : "rate_limit",
            "enabled" : false,
            "config" : {
                "limit" : "50-S",
                "policy" : "local"
            }
        },
        {
            "name" : "oauth2",
            "enabled" : true,
            "config" : {
                "server_name" : "live oauth server"
            }
        },
        {
            "name" : "compression",
            "enabled" : false,
            "config" : {}
        }
    ],
    "health_check" : {
        "url" : "",
        "timeout" : 0
    }
}`
)
