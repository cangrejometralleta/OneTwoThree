# Roster Service

The same Service, Served by three Frameworks.
Pick one at Startup.

```
go run .                 # net/http
SERVER=chi go run .      # chi
SERVER=gin go run .      # gin
```

Every Framework Returns byte-identical Answers,
because none of them Reaches the Handlers.

## Where the Framework Stops

```
main.go        Casts the Players and Picks an Adapter
adapters.go    THE ONLY FILE THAT IMPORTS gin OR chi
transport.go   Request, Response, Handler, Route
handlers.go    the Script. Framework-free
dto.go         the Wire Shapes
store.go       GORM. THE ONLY FILE THAT IMPORTS gorm
person.go      the Business Truth
```

Two Files Hold every Dependency.
Swapping a Framework Edits one of them, and nothing else.

## The Contract

```go
type Handler func(Request) Response
```

A Handler Receives a Request and Returns a Response.
It never Sees a ResponseWriter, a gin.Context or a chi.Router.

An Adapter Owes one Method:

```go
type Server interface {
	ServeRoutes(routes []Route, address string) error
}
```

Adding Fiber or Echo Means Writing that Method once.
No Handler Changes.

## What the Interfaces Buy

`RosterAPI` Depends on `PeopleStore`, never on GORM.
So a Test Hands it a Map and Runs in six Milliseconds,
with no Database and no Port.

```
go test ./...
```

## Three Shapes of one Person

| Shape        | Trusted | Speaks   |
|--------------|---------|----------|
| `PersonBody` | no      | the Wire |
| `Person`     | yes     | Business |
| `PersonRow`  | yes     | Storage  |
