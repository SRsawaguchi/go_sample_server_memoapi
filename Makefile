ENV_LOCAL=\
	MEMO_APP_PORT=8080 \
	MEMO_APP_HOST=127.0.0.1 \
	MEMO_DB_HOST=localhost \
	MEMO_DB_PORT=5432 \
	MEMO_DB_NAME=memo_app \
	MEMO_DB_USER=postgres \
	MEMO_DB_PASSWORD=postgres

ENV_INTEGRATION=\
	MEMO_APP_PORT=8080 \
	MEMO_APP_HOST=127.0.0.1 \
	MEMO_DB_HOST=localhost \
	MEMO_DB_PORT=5433 \
	MEMO_DB_NAME=memo_app_test \
	MEMO_DB_USER=postgres \
	MEMO_DB_PASSWORD=postgres

app.start:
	$(ENV_LOCAL) \
	go run main.go

test.unit:
	go test ./memo/handler -v

test.integration:
	$(ENV_INTEGRATION) \
	go test -tags=integration ./it -v -count 1
