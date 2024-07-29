# Mock Bot
---

## Init

- create a SQLite file named `mockbot.db` from `mockbot.db.sql`
- `insert into bot_settings (bot_id,key,value) values ('<bot_id>','token','<bot_token>');`
- `insert into staff (bot_id,role,user_id) values ('<bot_id>','owner','<user_id>');`
- `go run main.go --db_path=mockbot.db`
- ...

## WordCloud

use `MiSans-Medium.ttf`

## TODO

- [ ] // TODO
- [ ] i18n, Chinese simplify only now
- [ ] ...