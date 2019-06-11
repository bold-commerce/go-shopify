# go-shopify

The new home of Conversio's Shopify Go library.

**Note**: The library does not have implementations of all Shopify resources, but it is being used in production and should be stable for usage. PRs for new resources and endpoints are welcome, or you can simply implement some yourself as-you-go. See the section "Using your own models" for more info.

[![Build Status](https://travis-ci.org/bold-commerce/go-shopify.svg?branch=master)](https://travis-ci.org/bold-commerce/go-shopify)
[![codecov](https://codecov.io/gh/bold-commerce/go-shopify/branch/master/graph/badge.svg)](https://codecov.io/gh/bold-commerce/go-shopify) [![Join the chat at https://gitter.im/bold-commerce/go-shopify](https://badges.gitter.im/bold-commerce/go-shopify.svg)](https://gitter.im/bold-commerce/go-shopify?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Install

```console
$ go get github.com/bold-commerce/go-shopify
```

## Use

```go
import "github.com/bold-commerce/go-shopify"
```

This gives you access to the `goshopify` package.

#### Oauth

If you don't have an access token yet, you can obtain one with the oauth flow.
Something like this will work:

```go
// Create an app somewhere.
app := goshopify.App{
    ApiKey: "abcd",
    ApiSecret: "efgh",
    RedirectUrl: "https://example.com/shopify/callback",
    Scope: "read_products,read_orders",
}

// Create an oauth-authorize url for the app and redirect to it.
// In some request handler, you probably want something like this:
func MyHandler(w http.ResponseWriter, r *http.Request) {
    shopName := r.URL.Query().Get("shop")
    authUrl := app.AuthorizeURL(shopName)
    http.Redirect(w, r, authUrl, http.StatusFound)
}

// Fetch a permanent access token in the callback
func MyCallbackHandler(w http.ResponseWriter, r *http.Request) {
    // Check that the callback signature is valid
    if ok, _ := app.VerifyAuthorizationURL(r.URL); !ok {
        http.Error(w, "Invalid Signature", http.StatusUnauthorized)
        return
    }

    query := r.URL.Query()
    shopName := query.Get("shop")
    code := query.Get("code")
    token, err := app.GetAccessToken(shopName, code)

    // Do something with the token, like store it in a DB.
}
```

#### Api calls with a token

With a permanent access token, you can make API calls like this:

```go
// Create an app somewhere.
app := goshopify.App{
    ApiKey: "abcd",
    ApiSecret: "efgh",
    RedirectUrl: "https://example.com/shopify/callback",
    Scope: "read_products",
}

// Create a new API client
client := goshopify.NewClient(app, "shopname", "token")

// Fetch the number of products.
numProducts, err := client.Product.Count(nil)
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
client := goshopify.NewClient(app, "shopname", "")

// Fetch the number of products.
numProducts, err := client.Product.Count(nil)
```

#### Query options

Most API functions take an options `interface{}` as parameter. You can use one
from the library or create your own. For example, to fetch the number of
products created after January 1, 2016, you can do:

```go
// Create standard CountOptions
date := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
options := goshopify.CountOptions{createdAtMin: date}

// Use the options when calling the API.
numProducts, err := client.Product.Count(options)
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
orderCount, err := client.Order.Count(options)
```

#### Using your own models

Not all endpoints are implemented right now. In those case, feel free to
implement them and make a PR, or you can create your own struct for the data
and use `NewRequest` with the API client. This is how the existing endpoints
are implemented.

For example, let's say you want to fetch webhooks. There's a helper function
`Get` specifically for fetching stuff so this will work:

```go
// Declare a model for the webhook
type Webhook struct {
    ID int         `json:"id"`
    Address string `json:"address"`
}

// Declare a model for the resource root.
type WebhooksResource struct {
    Webhooks []Webhook `json:"webhooks"`
}

func FetchWebhooks() ([]Webhook, error) {
    path := "admin/webhooks.json"
    resource := new(WebhooksResource)
    client := goshopify.NewClient(app, "shopname", "token")

    // resource gets modified when calling Get
    err := client.Get(path, resource, nil)

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

There's nothing special to note about the tests except that if you have Docker
and Compose installed, you can test like this:

    $ docker-compose build dev
    $ docker-compose run --rm dev

Testing the package is the default command for the dev container. To create a
coverage profile:

    $ docker-compose run --rm dev bash -c 'go test -coverprofile=coverage.out ./... && go tool cover -html coverage.out -o coverage.html'
