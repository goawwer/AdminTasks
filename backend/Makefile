login_db:
	psql ${DATABASE_URL}

create_migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate_upgrade:
	migrate -database ${DATABASE_URL} -path migrations up

migrate_downgrade:
	migrate -database ${DATABASE_URL} -path migrations down

migrate_downgrade_force:
	migrate -database ${DATABASE_URL} -path migrations force 2