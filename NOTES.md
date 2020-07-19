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
