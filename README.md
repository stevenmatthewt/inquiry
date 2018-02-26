# inquiry

`inquiry` is a simple package desiged to facilitate with handling Query String parameters in Golang HTTP Servers. It allows Query String parameters to be automatically Unmarshalled into a struct, thus promoting type safety throughout the rest of your code.

## How to use

`inquiry` aims to maintain a very simple, intuitive interface that matches other Go packages such as `json` and `xml`.

```go
type queryFormat struct {
    AccountID      int      `query:"account_id"`
    FirstName      string   `query:"first_name"`
    LastName       string   `query:"last_name"`
    FavoriteColors []string `query:"fav_colors"`
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var queryParameters queryFormat
    err := inquiry.UnmarshalMap(r.URL.Query(), &queryParameters)
    if err != nil {
        panic(err)
    }

    // You can now access you query parameters in a type-safe manner
    fmt.Printf("%s %s belongs to account %d",
        queryParameters.FirstName,
        queryParameters.LastName,
        queryParameters.AccountID,
    )
}
```
