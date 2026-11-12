# go-shopify

The new home of Conversio's Shopify Go library.

**Note**: The library does not have implementations of all Shopify resources, but it is being used in production and should be stable for usage. PRs for new resources and endpoints are welcome, or you can simply implement some yourself as-you-go. See the section "Using your own models" for more info.

[![Build Status](https://github.com/bold-commerce/go-shopify/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/bold-commerce/go-shopify/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/bold-commerce/go-shopify/branch/master/graph/badge.svg)](https://codecov.io/gh/bold-commerce/go-shopify) [![Join the chat at https://gitter.im/bold-commerce/go-shopify](https://badges.gitter.im/bold-commerce/go-shopify.svg)](https://gitter.im/bold-commerce/go-shopify?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Supported Go Versions

The Go versions tested by CI are listed in the [test workflow](.github/workflows/tests.yml).
Older versions are not covered by CI.

## Install v4

```console
$ go get github.com/bold-commerce/go-shopify/v4
```

## Use v4

```go
import "github.com/bold-commerce/go-shopify/v4"
```

This gives you access to the `goshopify` package.

## Install v3

```console
$ go get github.com/bold-commerce/go-shopify/v3
```

## Use v3

```go
import "github.com/bold-commerce/go-shopify/v3"
```

This gives you access to the `goshopify` package.

## Install v2

```console
$ go get github.com/bold-commerce/go-shopify
```

## Use v2

```go
import "github.com/bold-commerce/go-shopify"
```

This gives you access to the `goshopify` package.

#### Oauth

If you don't have an access token yet, you can obtain one with the OAuth flow.
The handler fragments below show the library calls. Generate a random, single-use
state value for each authorization attempt and store it in the user's session.
In the callback, verify that the returned state matches the stored value before
exchanging the code. The fixed state below is only a placeholder.

```go
// Create an app somewhere.
var app = goshopify.App{
    ApiKey: "abcd",
    ApiSecret: "efgh",
    RedirectUrl: "https://example.com/shopify/callback",
    Scope: "read_products,read_orders",
}

// Create an oauth-authorize url for the app and redirect to it.
// In some request handler, you probably want something like this:
func MyHandler(w http.ResponseWriter, r *http.Request) {
    shopName := r.URL.Query().Get("shop")
    state := "nonce"
    authUrl, err := app.AuthorizeUrl(shopName, state)
    if err != nil {
        http.Error(w, "Invalid shop", http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, authUrl, http.StatusFound)
}

// Fetch a permanent access token in the callback
func MyCallbackHandler(w http.ResponseWriter, r *http.Request) {
    // Validate the single-use state against the user's session here.
    // Check that the callback signature is valid
    if ok, _ := app.VerifyAuthorizationURL(r.URL); !ok {
        http.Error(w, "Invalid Signature", http.StatusUnauthorized)
        return
    }

    query := r.URL.Query()
    shopName := query.Get("shop")
    code := query.Get("code")
    ctx := r.Context()
    token, err := app.GetAccessToken(ctx, shopName, code)
    if err != nil {
        http.Error(w, "Token exchange failed", http.StatusBadGateway)
        return
    }

    // Do something with the token, like store it in a DB.
}
```

##### Expiring offline access tokens

Shopify is replacing permanent offline access tokens with
[expiring offline access tokens](https://shopify.dev/docs/apps/build/authentication-authorization/access-tokens/offline-access-tokens)
for public apps. An expiring token comes with a refresh token:

```go
// In the callback, fetch an expiring token instead of a permanent one.
token, err := app.GetExpiringAccessToken(ctx, shopName, code)
// Store token.AccessToken, token.RefreshToken and both lifetimes
// (token.ExpiresIn and token.RefreshTokenExpiresIn, in seconds).

// Before the access token expires, refresh it. The refresh token is rotated,
// so store the returned token again.
token, err = app.RefreshAccessToken(ctx, shopName, token.RefreshToken)

// Exchange a permanent token that was stored earlier for an expiring one.
// Shopify invalidates the permanent token.
token, err = app.MigrateToExpiringAccessToken(ctx, shopName, permanentToken)
```

#### Api calls with a token

With a valid access token, you can make API calls like this. These snippets use
the v4 package. Import `context` and use a request context when available, or
`context.Background()` for a standalone call.

```go
// Create an app somewhere.
app := goshopify.App{
    ApiKey: "abcd",
    ApiSecret: "efgh",
    RedirectUrl: "https://example.com/shopify/callback",
    Scope: "read_products",
}

// Create a new API client
client, err := goshopify.NewClient(app, "shopname", "token")

// Fetch the number of products.
numProducts, err := client.Product.Count(context.Background(), nil)
```

#### Private App Auth

Private Shopify apps use basic authentication and do not require going through the OAuth flow. Here is an example:

```go
// Create an app somewhere.
app := goshopify.App{
	ApiKey: "apikey",
	Password: "apipassword",
}

// Create a new API client (notice the token parameter is the empty string)
client, err := goshopify.NewClient(app, "shopname", "")

// Fetch the number of products.
numProducts, err := client.Product.Count(context.Background(), nil)
```

### Client Options

When creating a client there are configuration options you can pass to NewClient. Simply use the last variadic param and
pass in the built in options or create your own and manipulate the client. See [options.go](https://github.com/bold-commerce/go-shopify/blob/master/options.go)
for more details.

#### WithVersion

Read more details on the [Shopify API Versioning](https://shopify.dev/docs/api/usage/versioning)
to understand the format and release schedules. You can use `WithVersion` to specify a specific version
of the API. If you do not use this option you will be defaulted to the oldest stable API.

```go
client, err := goshopify.NewClient(app, "shopname", "token", goshopify.WithVersion("2026-10"))
```

#### WithRetry

Shopify [Rate Limits](https://shopify.dev/docs/api/admin-rest/usage/rate-limits) their API and if this happens to you they
will send a back off (usually 2s) to tell you to retry your request. To support this functionality seamlessly within
the client a `WithRetry` option exists where you can pass an `int` of how many times you wish to retry per-request
before returning an error. `WithRetry` additionally supports retrying HTTP503 errors.

```go
client, err := goshopify.NewClient(app, "shopname", "", goshopify.WithRetry(3))
```

Read recorded rate-limit information with `client.RateLimits()`. This method
replaces the `client.RateLimits` field and returns an independent snapshot that
is safe to read while other requests are in progress. Changes to the snapshot
do not change the client.

#### Query options

Most API functions take an options `interface{}` as parameter. You can use one
from the library or create your own. For example, to fetch the number of
products created after January 1, 2016, you can do:

```go
// Create standard CountOptions
date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
options := goshopify.CountOptions{CreatedAtMin: date}

// Use the options when calling the API.
numProducts, err := client.Product.Count(context.Background(), options)
```

The options are parsed with Google's
[go-querystring](https://github.com/google/go-querystring) library so you can
use custom options like this:

```go
// Create custom options for the orders.
// Notice the `url:"status"` tag
options := struct {
    Status string `url:"status"`
}{"any"}

// Fetch the order count for orders with status="any"
orderCount, err := client.Order.Count(context.Background(), options)
```

#### Using your own models

Not all endpoints are implemented right now. In those case, feel free to
implement them and make a PR, or you can create your own struct for the data
and use the client's `Get`, `Post`, `Put`, or `Delete` helpers. These helpers add
the configured API path prefix. Pass a resource path such as `webhooks.json`,
without `admin/` or `admin/api/<version>/`.

For example, let's say you want to fetch webhooks. There's a helper function
`Get` specifically for fetching stuff so this will work:

```go
// Declare a model for the webhook
type Webhook struct {
    Id int         `json:"id"`
    Address string `json:"address"`
}

// Declare a model for the resource root.
type WebhooksResource struct {
    Webhooks []Webhook `json:"webhooks"`
}

func FetchWebhooks(ctx context.Context) ([]Webhook, error) {
    path := "webhooks.json"
    resource := new(WebhooksResource)
    client, err := goshopify.NewClient(app, "shopname", "token")
    if err != nil {
        return nil, err
    }

    // resource gets modified when calling Get
    err = client.Get(ctx, path, resource, nil)

    return resource.Webhooks, err
}
```

#### Webhooks verification

In order to be sure that a webhook is sent from ShopifyApi you could easily verify
it with the `VerifyWebhookRequest` method.

For example:

```go
func ValidateWebhook(httpRequest *http.Request) (bool) {
    shopifyApp := goshopify.App{ApiSecret: "ratz"}
    return shopifyApp.VerifyWebhookRequest(httpRequest)
}
```

## Develop and test

`docker` and `docker-compose` must be installed

### Mac/Linux/Windows with make

Using the make file is the easiest way to get started with the tests and wraps the manual steps below with easy to use
make commands.

```shell
make && make test
```

#### Makefile goals

- `make` or `make container`: default goal is to make the `go-shopify:latest` build container
- `make test`: run go test in the container
- `make clean`: deletes the `go-shopify:latest` image and coverage output
- `make coverage`: generates the coverage.html and opens it

### Manually

To run the tests you will need the `go-shopify:latest` image built to run your tests, to do this run

```
docker-compose build test
```

To run tests you can use run

```shell
docker-compose run --rm test
```

To create a coverage profile run the following to generate a coverage.html

```
docker-compose run --rm test sh -c 'go test -coverprofile=coverage.out ./... && go tool cover -html coverage.out -o coverage.html'
```

When done testing and you want to cleanup simply run

```
docker image rm go-shopify:latest
```

Read the docker-compose.yml and Dockerfile for further details.
