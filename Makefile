.PHONY: psql mysql

psql:
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@v4.8.6 && go run main.go psql --notion-token ${NOTION_TOKEN} --notion-page-id ${PSQL_NOTION_PAGE_ID} --notion-table-index ${PSQL_NOTION_TABLE_INDEX_ID} --notion-mermaid-id ${PSQL_NOTION_MERMAID_ID}

mysql:
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.8.6 && go run main.go mysql --notion-token ${NOTION_TOKEN} --notion-page-id ${MYSQL_NOTION_PAGE_ID} --notion-table-index ${MYSQL_NOTION_TABLE_INDEX_ID} -c './sqlboiler.toml'

mssql:
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mssql@v4.8.6 && go run main.go mssql --notion-token ${NOTION_TOKEN} --notion-page-id ${MSSQL_NOTION_PAGE_ID}
