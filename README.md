# scheman

Scheman is a tool to write database schema from connected database.

#### Supported Databases

| Database     | Test Confirmed |
| ------------ | -------------- |
| PostgreSQL   | 👌             |
| MySQL        | 👌             |
| MSSQL Server |                |
| SQLite3      |                |
| CockroachDB  |                |

#### Supported Output Destination

Only Notion is supported as output destination now. But I am going to add output destination like spread-sheat.

| Output Destination | Test Confirmed |
| ------------------ | -------------- |
| Notion             | 👌             |

## How to Use

Install this package and write settings for connection.

**_installation_**

```
go install github.com/maru44/scheman@v1.0.0
```

**_example for PostgreSQL - Notion)_**

```
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@v4.8.6 && scheman psql --notion-token ${NOTION_TOKEN} --notion-page-id ${PSQL_NOTION_PAGE_ID}
```

#### Generic config options

| Name               | Defaults   |                                                  |
| ------------------ | ---------- | ------------------------------------------------ |
| config             | "scheman"  |                                                  |
| services           | ["NOTION"] |                                                  |
| notion-page-id     | ""         | required if output destinations contain "NOTION" |
| notion-page-token  | ""         | required if output destinations contain "NOTION" |
| notion-table-index | ""         |                                                  |
| disable-views      | false      |                                                  |
| attr-ignore        | [ ]        |                                                  |

#### Database Driver Configuration

Settings for database you want to connect.

| Name      | Required | Postgres Default | MySQL Default | MSSQL Default |
| --------- | -------- | ---------------- | ------------- | ------------- |
| schema    | no       | "public"         | none          | "dbo"         |
| dbname    | yes      | none             | none          | none          |
| host      | yes      | none             | none          | none          |
| port      | no       | 5432             | 3306          | 1433          |
| user      | yes      | none             | none          | none          |
| pass      | no       | none             | none          | none          |
| sslmode   | no       | "require"        | "true"        | "true"        |
| whitelist | no       | []               | []            | []            |
| blacklist | no       | []               | []            | []            |

ref: https://github.com/volatiletech/sqlboiler#database-driver-configuration
