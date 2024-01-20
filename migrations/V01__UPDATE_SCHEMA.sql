\c market_list;

ALTER TABLE PURCHASE
    ADD COLUMN NAME VARCHAR(500) NOT NULL default '';

ALTER TABLE PURCHASE
    ADD COLUMN IS_FAVORITE BOOLEAN DEFAULT FALSE;

CREATE TABLE TAG
(
    ID         BIGSERIAL PRIMARY KEY,
    NAME       VARCHAR(500)                       NOT NULL,
    USER_ID    BIGINT REFERENCES MARKET_USER (ID) NOT NULL,
    CREATED_AT TIMESTAMP DEFAULT NOW()
);


CREATE TABLE TAG_PURCHASE
(
    PURCHASE_ID BIGINT REFERENCES PURCHASE (ID),
    TAG_ID      BIGINT REFERENCES TAG (ID),
    CONSTRAINT TAG_PURCHASE_PK PRIMARY KEY (PURCHASE_ID, TAG_ID)
);

CREATE TABLE PURCHASE_USER (
  PURCHASE_ID BIGINT REFERENCES PURCHASE(ID),
  USER_ID BIGINT REFERENCES MARKET_USER(ID),
  CONSTRAINT PURCHASE_USER_PK  PRIMARY KEY (USER_ID, PURCHASE_ID)
);


INSERT INTO PURCHASE_USER (PURCHASE_ID, USER_ID)
    SELECT id as PURCHASE_ID, USER_ID FROM PURCHASE;


ALTER TABLE PURCHASE DROP COLUMN USER_ID;