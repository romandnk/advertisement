FROM migrate/migrate:latest

COPY ./migrations /migrations

ENTRYPOINT ["migrate", "-path", "./migrations", "-database", "postgres://postgres:1234@postgres:5432/adverts_db?sslmode=disable", "up"]
