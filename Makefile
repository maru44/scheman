.PHONY: psql mysql

# for tests

psql:
	@go run main.go psql --notion-page-id ${PSQL_NOTION_PAGE_ID} --notion-token ${NOTION_TOKEN} --notion-table-index ${PSQL_NOTION_TABLE_LIST_ID}

mysql:
	@go run main.go mysql --notion-page-id ${MYSQL_NOTION_PAGE_ID} --notion-token ${NOTION_TOKEN} --notion-table-index ${MYSQL_NOTION_TABLE_LIST_ID}
