Banister Notes
==============

Banister is a tool for:

- Generating type-safe database client
- Generating type-safe migrations
- Running migrations

Architecture
------------

**`banister`**: CLI entry-point to do everything
**`generate`**: Tool to generate database client code
**`migrations`**: Tool to generate database migration code
**`common`**: A library used by generated client and generated migrations
**`backends/<name>`**: Backends for each supported database (postgres, sqlite, etc)

General Notes
-------------

- Generated client code is independent of generated migrations
- `common` might be able to be split into two
    - One to be used by generated client
    - One to be used by generated migrations
- 

Generated Types
---------------

Top-level helpers:

- **Store**: Stored connection to DB and provides access to model managers (from `banister.Connect()`)
- **modelManager**: Provides entrypoint to DB helpers for specific model (from `store.Users`)
- **modelQueryset**: Holds state for current filter operation (from `store.Users.Filter(...)`)

Static accessor types:
`${Model}${Operation}Options`

- **modelWhereOptions**: Holds references to "Where" helpers for each field (eg. `userWhereOptions { ... }`)
- **modelSortOptions**: Holds references to "OrderBy" helpers for each field (eg. `userOrderByOptions { ... }`)
- **modelSetterOptions**: Holds references to "Setter" helpers for each field (eg. `userSetterOptions { ... }`)

Argument types:
`${Model}${Operation}Arg`

- **modelFilterArg**: Arg used to pass into Queryset.Filter (from `WhereUser.Email.Eq("...")`)
- **modelOrderByArg**: Arg used to pass into Queryset.OrderBy (from `OrderByUser.EmailAscending`)
- **modelSetterArg**: Arg used to pass into Queryset.Insert or Queryset.Update (from `SetUser.Email("...")`)

Values:

- **WhereModel**: Instance of `modelWhereOptions`
- **OrderByModel**: Instance of `modelOrderByOptions`
- **SetModel**: Instance of `modelSetterOptions`
