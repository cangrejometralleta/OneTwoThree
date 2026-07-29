# People Service

A small Microservice that Follows the Rules on Purpose.
Run it, then Read it.

```
go run .
```

## Three Shapes of one Person

The same Person Wears three Types, and never two at once.

| Shape        | Lives in    | Trusted | Speaks   |
|--------------|-------------|---------|----------|
| `PersonBody` | `dto.go`    | no      | the Wire |
| `Person`     | `person.go` | yes     | Business |
| `PersonRow`  | `store.go`  | yes     | Storage  |

Three Shapes Generate three Crossings, and no Shape Holds the Center.
`BuildPersonRecord` Promotes.
`EncodePersonRow` Stores.
`RenderPersonView` Answers.

## Typing as Documentation

`FullName`, `EmailAddress` and `PhoneNumber` are not Decoration.
They are `string` underneath, and still the Compiler Refuses
to Pass an Email where a Name Belongs.

The Business Rule Stops being a Comment.
It Becomes a Build Failure.

## Five Endpoints, five Stories

| Method   | Path           | Handler             |
|----------|----------------|---------------------|
| `GET`    | `/people`      | `ListPeopleRecords` |
| `POST`   | `/people`      | `AddPersonRecord`   |
| `GET`    | `/people/{id}` | `ShowPersonRecord`  |
| `PUT`    | `/people/{id}` | `SavePersonRecord`  |
| `DELETE` | `/people/{id}` | `DropPersonRecord`  |

Every Handler Reads in three Sections, and Names no Query.
Open one and you Find the Story, never the Mechanism.
