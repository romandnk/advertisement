CREATE TABLE adverts (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE images (
    id VARCHAR(36) PRIMARY KEY,
    advert_id VARCHAR(36) REFERENCES adverts(id),
    url VARCHAR(255)
);