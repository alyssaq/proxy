# Simple proxy server example

Usage:

    $ go run proxy.go -proxyURL=<http://api.something.com> -port=<port>

port `9000` is used by default if not specified.

The proxy server will proxy calls from `http://api.something.com` to `localhost:9000`

### Use cases
You need to modify an API but have no influence to do so.

For example:

* Your web app is on a https endpoint and the API is on http.
* You wish to add more headers to the API (such as cache-control)

You can host the proxy server on a https endpoint, make https calls and it'll return results from the http API.
