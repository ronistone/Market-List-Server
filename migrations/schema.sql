CREATE database market_list;

CREATE TABLE MARKET_USER
(
    ID         BIGSERIAL PRIMARY KEY,
    EMAIL      VARCHAR(300),
    NAME       VARCHAR(300),
    PASSWORD   VARCHAR(100),
    CREATED_AT TIMESTAMP DEFAULT now(),
    UPDATED_AT TIMESTAMP DEFAULT now()
);

CREATE TABLE PRODUCT
(
    ID         BIGSERIAL PRIMARY KEY,
    EAN        VARCHAR(20) UNIQUE,
    NAME       VARCHAR(300),
    UNIT       VARCHAR(30),
    SIZE       INT,
    CREATED_AT TIMESTAMP DEFAULT now(),
    UPDATED_AT TIMESTAMP DEFAULT now()
);

CREATE TABLE MARKET
(
    ID         BIGSERIAL PRIMARY KEY,
    NAME       VARCHAR(500),
    CREATED_AT TIMESTAMP DEFAULT now(),
    UPDATED_AT TIMESTAMP DEFAULT now(),
    ENABLED    BOOLEAN   DEFAULT TRUE
);

CREATE TABLE PRODUCT_INSTANCE
(
    ID         BIGSERIAL PRIMARY KEY,
    PRODUCT_ID BIGINT REFERENCES PRODUCT (ID),
    MARKET_ID  BIGINT REFERENCES MARKET (ID),
    PRICE      INT,
    PRECISION  INT,
    CREATED_AT TIMESTAMP DEFAULT now()
);

CREATE TABLE PURCHASE
(
    ID         BIGSERIAL PRIMARY KEY,
    CREATED_AT TIMESTAMP DEFAULT now(),
    USER_ID    BIGINT REFERENCES MARKET_USER (ID),
    MARKET_ID  BIGINT REFERENCES MARKET (ID)
);

CREATE TABLE PURCHASE_ITEM
(
    ID                  BIGSERIAL PRIMARY KEY,
    PRODUCT_INSTANCE_ID BIGINT REFERENCES PRODUCT_INSTANCE (ID),
    PURCHASE_ID         BIGINT REFERENCES PURCHASE (ID),
    PURCHASED           BOOLEAN default FALSE,
    QUANTITY            INT     default 1
);