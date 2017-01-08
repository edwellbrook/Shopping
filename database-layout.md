# Database Layout

```sql
-- login details for a card and shopping list
CREATE TABLE IF NOT EXISTS cards (
  card_id     varchar(16)                  PRIMARY KEY NOT NULL,
  password    char(60)                                 NOT NULL,
  list        varchar(22)[]   DEFAULT {}               NOT NULL
);

-- available items in the store
CREATE TABLE IF NOT EXISTS items (
  name        varchar(22)                  PRIMARY KEY NOT NULL,
  location    varchar(16)                              NOT NULL
);

-- the different locations within a store
CREATE TABLE IF NOT EXISTS locations (
  beacon_id   varchar(16)                  PRIMARY KEY NOT NULL
);
```
