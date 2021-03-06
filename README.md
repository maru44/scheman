# scheman

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/maru44/scheman/blob/master/LICENSE)
![ActionsCI](https://github.com/maru44/scheman/workflows/Test%20Lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/maru44/scheman)](https://goreportcard.com/report/github.com/maru44/scheman)

Scheman is a tool to visualize tables schema or ERD of connected RDB.

Main purpose of this is accelerating working collaboration between engineers and non-engineers.

You can choose output destination from `Notion`, `File` or both of them.

### Sample Images

Here is examples for output.

#### Notion

![](https://user-images.githubusercontent.com/46714011/155822065-f0f9f785-b2b1-4abd-b98b-052496dff169.png)

![](https://user-images.githubusercontent.com/46714011/155862202-77e81b99-681a-44fb-bf1c-669dae7f1f5a.png)

#### File

##### definition

https://github.com/maru44/scheman/blob/master/testdata/postgres/def.csv

##### ERD

https://github.com/maru44/scheman/blob/master/testdata/postgres/erd.md

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
| File               | 👌             |

## How to Use

Install this package and write settings for connection.

**_installation_**

```shell: installation
go install github.com/maru44/scheman@v1.3.0
```

**_example for PostgreSQL - Notion)_**

```shell: Notion - PostgreSQL
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@v4.8.6 \
 && scheman psql --services notion --erd-outputs notion \
  --notion-token ${NOTION_TOKEN} --notion-page-id ${PSQL_NOTION_PAGE_ID}
```

If you want to overwrite your schema-definition tables or ERD, you have to set `notion-table-index` after this command done. This value is oututted in your command line.

![](https://user-images.githubusercontent.com/46714011/156856299-67bed77d-d744-458f-9967-61d50983eade.png)

**_example for MySQL - File with sqlboiler.toml)_**

```shell: File - MySQL
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.8.6 \
 && scheman mysql -c sqlboiler.toml --def-file def.csv --erd-file erd.md
```

#### Generic config options

| Name               | Defaults  |                                                                                                   |
| ------------------ | --------- | ------------------------------------------------------------------------------------------------- |
| config             | "scheman" | if you use sqlboiler, you can use sqlboiler.toml(.yaml) with setting this "sqlboiler.toml(.yaml)" |
| services           | [ ]       | notion                                                                                            |
| erd-outputs        | [ ]       | notion                                                                                            |
| notion-page-id     | ""        | required if output destinations contain "notion"                                                  |
| notion-page-token  | ""        | required if output destinations contain "notion"                                                  |
| notion-table-index | ""        | if you want to overwrite definition table, please fill this                                       |
| def-file           | ""        | The file name. required if output destinations if you want to output tables schema to file        |
| erd-file           | ""        | The file name. required if output destinations if you want to output ERD to file                  |
| disable-views      | false     |                                                                                                   |
| attr-ignore        | [ ]       |                                                                                                   |

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

## thx

https://github.com/volatiletech/sqlboiler

https://github.com/dstotijn/go-notion
