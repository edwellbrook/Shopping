# Database Layout

```
CREATE TABLE IF NOT EXISTS cards (
  card_id  varchar(16) UNIQUE NOT NULL,
  password char(60)           NOT NULL
);
```
