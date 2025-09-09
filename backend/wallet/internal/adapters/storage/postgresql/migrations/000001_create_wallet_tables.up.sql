-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create asset_type ENUM for data integrity
CREATE TYPE asset_type AS ENUM ('CRYPTO', 'FIAT');

-- Create tables
CREATE TABLE assets
(
    id         UUID PRIMARY KEY            DEFAULT gen_random_uuid(),
    symbol     VARCHAR(10) UNIQUE NOT NULL,
    name       VARCHAR(50)        NOT NULL,
    asset_type asset_type         NOT NULL,
    precision  SMALLINT           NOT NULL,
    is_active  BOOLEAN            NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

CREATE TABLE networks
(
    id         UUID PRIMARY KEY            DEFAULT gen_random_uuid(),
    name       VARCHAR(50) UNIQUE NOT NULL,
    chain_id   BIGINT,
    rpc_url    VARCHAR(255)       NOT NULL,
    is_active  BOOLEAN            NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

CREATE TABLE asset_networks
(
    asset_id         UUID            NOT NULL REFERENCES assets (id),
    network_id       UUID            NOT NULL REFERENCES networks (id),
    contract_address VARCHAR(255),
    min_withdrawal   DECIMAL(36, 18) NOT NULL,
    withdrawal_fee   DECIMAL(36, 18) NOT NULL,
    is_active        BOOLEAN         NOT NULL DEFAULT TRUE,
    PRIMARY KEY (asset_id, network_id)
);

CREATE TABLE accounts
(
    id             UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    user_id        UUID            NOT NULL,
    asset_id       UUID            NOT NULL REFERENCES assets (id),
    balance        DECIMAL(36, 18) NOT NULL DEFAULT 0,
    locked_balance DECIMAL(36, 18) NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, asset_id)
);

CREATE TABLE deposit_addresses
(
    id              UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    account_id      UUID                NOT NULL REFERENCES accounts (id),
    address         VARCHAR(255) UNIQUE NOT NULL,
    derivation_path VARCHAR(100),
    memo_tag        VARCHAR(100),
    created_at      TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);

CREATE TABLE deposits
(
    id            UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    account_id    UUID            NOT NULL REFERENCES accounts (id),
    asset_id      UUID            NOT NULL REFERENCES assets (id),
    status        SMALLINT        NOT NULL DEFAULT 0,
    amount        DECIMAL(36, 18) NOT NULL,
    onchain_tx_id VARCHAR(255)    NOT NULL,
    from_address  VARCHAR(255),
    confirmations INT             NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    confirmed_at  TIMESTAMPTZ,
    UNIQUE (asset_id, onchain_tx_id)
);

CREATE TABLE withdrawals
(
    id            UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    account_id    UUID            NOT NULL REFERENCES accounts (id),
    asset_id      UUID            NOT NULL REFERENCES assets (id),
    status        SMALLINT        NOT NULL DEFAULT 0,
    amount        DECIMAL(36, 18) NOT NULL,
    to_address    VARCHAR(255)    NOT NULL,
    onchain_tx_id VARCHAR(255),
    network_fee   DECIMAL(36, 18),
    approved_by   UUID,
    priority      SMALLINT        NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE TABLE internal_transfers
(
    id              UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    from_account_id UUID            NOT NULL REFERENCES accounts (id),
    to_account_id   UUID            NOT NULL REFERENCES accounts (id),
    asset_id        UUID            NOT NULL REFERENCES assets (id),
    amount          DECIMAL(36, 18) NOT NULL,
    reference_id    UUID,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    CHECK (from_account_id <> to_account_id)
);