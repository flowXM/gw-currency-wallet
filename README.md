Запуск сервиса
```bash
go run cmd/main.go -c config.env
```

Миграция базы данных
```bash
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE TABLE IF NOT EXISTS users (
        user_id uuid primary key,
        username varchar(255) not null unique,
        email varchar(255) not null unique,
        password_hash varchar(255) not null,
        salt BYTEA not null
    );

    CREATE TABLE IF NOT EXISTS wallets (
    	wallet_id uuid primary key,
        user_id uuid,
    	rub_amount numeric(32, 2) not null default 0.00 CHECK (rub_amount >= 0),
        usd_amount numeric(32, 2) not null default 0.00 CHECK (usd_amount >= 0),
        eur_amount numeric(32, 2) not null default 0.00 CHECK (eur_amount >= 0),
        CONSTRAINT fk_user
            FOREIGN KEY(user_id)
                REFERENCES users(user_id)
    );
EOSQL
```