CREATE TABLE account_type (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE account (
    id SERIAL PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    account_type_id INTEGER NOT NULL,
    CONSTRAINT fk_account_type FOREIGN KEY(account_type_id) REFERENCES account_type(id)
);

CREATE TABLE session (
    id SERIAL PRIMARY KEY,
    jwt TEXT,
    account_id INTEGER NOT NULL,
    CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES account(id)
);

