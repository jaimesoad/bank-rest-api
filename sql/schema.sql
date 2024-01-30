CREATE TABLE IF NOT EXISTS account
(
    account_number INTEGER      PRIMARY KEY auto_increment,
    account_state  BOOLEAN      NOT NULL,
    balance        FLOAT        NOT NULL,
    client_name    varchar(255) NOT NULL
);