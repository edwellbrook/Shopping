# Database Layout

```sql
-- login details for a card and shopping list
CREATE TABLE IF NOT EXISTS cards (
  id          varchar(16)                    NOT NULL,
  password    char(60)                       NOT NULL,
  list        varchar(22)[]   DEFAULT '{}'   NOT NULL,

  CONSTRAINT cards_id_key PRIMARY KEY (id)
);
```
