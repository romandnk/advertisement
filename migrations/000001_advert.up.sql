CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(254) NOT NULL,
    password VARCHAR(60) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted BOOLEAN
);

CREATE TABLE adverts (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    deleted BOOLEAN NOT NULL
);

CREATE TABLE images (
    id VARCHAR(36) PRIMARY KEY,
    advert_id VARCHAR(36) REFERENCES adverts(id),
    created_at TIMESTAMP NOT NULL,
    deleted BOOLEAN NOT NULL
);


