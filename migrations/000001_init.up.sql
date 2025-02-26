/*
WALLET
*/

CREATE TABLE IF NOT EXISTS wallets
(
    uuid       uuid PRIMARY KEY DEFAULT gen_random_uuid() UNIQUE NOT NULL,
    amount     BIGINT           DEFAULT 0                        NOT NULL,
    version    BIGINT           DEFAULT 0                        NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE                          NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE                          NOT NULL
);

/*
TEST WALLET
*/

INSERT INTO wallets
    (uuid, amount, version, created_at, updated_at)
VALUES ('00000000-0000-0000-0000-000000000001', 0, 0, now(), now());

/*
TRANSACTIONS
*/

CREATE TABLE IF NOT EXISTS transactions
(
    id              BIGSERIAL                NOT NULL,
    wallet_uuid     uuid                     NOT NULL,
    idempotency_key uuid                     NOT NULL,
    operation       VARCHAR(50)              NOT NULL,
    amount          BIGINT                   NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL
)

/*
     CONSTRAINT fk_transactions_wallet_uuid FOREIGN KEY (wallet_uuid)
         REFERENCES wallets (uuid)
*/