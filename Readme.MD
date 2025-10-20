## How to use

1. Copy files to a folder (module path used above `https://github.com/bakwanx/xyz_test`). Update the module path in `go.mod` if needed.
2. Create `.env` from `.env.example` and adjust values.
3. Start MySQL (either `docker-compose up -d` or your own DB) and run the SQL in `migrations/migrations.sql`.
4. `go mod tidy`
5. `make run` or `go run ./cmd/server`
6. Register admin: `POST /api/v1/register` with `{ "name":"Admin", "email":"a@a.com", "password":"secret" }`.
7. Login: `POST /api/v1/login` with `{ "email":"a@a.com", "password":"secret" }` — you'll get JWT.
8. Use the JWT as `Authorization: Bearer <token>` to call protected endpoints.
