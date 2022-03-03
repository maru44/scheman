.PHONY: psql mysql

psql:
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@v4.8.6 && go run main.go psql \
		--services "NOTION" --services "File" --erd-outputs "NOTION" --erd-outputs "file" --erd-file ./testdata/postgres/erd.md --def-file ./testdata/postgres/def.tsv \
		--notion-token ${NOTION_TOKEN} --notion-page-id ${PSQL_NOTION_PAGE_ID} --notion-table-index ${PSQL_NOTION_TABLE_INDEX_ID}

mysql:
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.8.6 && go run main.go mysql --notion-token ${NOTION_TOKEN} --notion-page-id ${MYSQL_NOTION_PAGE_ID} --notion-table-index ${MYSQL_NOTION_TABLE_INDEX_ID} -c './sqlboiler.toml'

mssql:
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mssql@v4.8.6 && go run main.go mssql --notion-token ${NOTION_TOKEN} --notion-page-id ${MSSQL_NOTION_PAGE_ID}

test:
	@go test ./...
