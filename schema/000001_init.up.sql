CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(256) NOT NULL,
    role VARCHAR(256) CHECK (role IN ('employee', 'moderator')) NOT NULL,
    password VARCHAR(256) NOT NULL
);

CREATE TABLE pvz (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    registrationDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    city VARCHAR(256) CHECK (city IN ('Москва', 'Казань', 'Санкт-Петербург')) NOT NULL
);

CREATE TABLE reception (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    pvzId UUID REFERENCES pvz(id) ON DELETE SET NULL,
    status VARCHAR(256) CHECK (status IN ('in_progress', 'close')) NOT NULL
);

CREATE TABLE product (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type VARCHAR(256) CHECK (typr IN ('электроника', 'одежда', 'обувь')) NOT NULL,
    receptionId UUID REFERENCES reception(id) ON DELETE SET NULL
);