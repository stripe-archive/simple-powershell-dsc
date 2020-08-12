# gateway

This is a forked and lightly-patched version of https://github.com/apex/gateway
that makes a single change: always base64-encoding the response body so that
it's easier to have a single API Gateway fronting this entire API.

The original README can be found at Readme.md.orig, and the license file at LICENSE
