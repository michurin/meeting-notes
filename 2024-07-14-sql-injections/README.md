# SQL-injection example

## run

```sh
sudo docker compose -f postgresql.yaml up
```

```
http://localhost:8080/?pgsql=db&username=user&db=user&ns=public
```

## queries

Formatted example

```sql
select string_agg(
  concat(
    application_name, e'\n',
    user, e'\n',
    client_addr, e'\n',
    query),
  e'\n\n')
from pg_stat_activity
where query <> ''
on conflict (login) do
update set email=excluded.email
```

How to avoid `where x %s y`

```sql
prepare Q (text, bool) as select login from accounts where (login = $1 and $2) or (login <> $1 and not $2);
execute Q('a', true);
execute Q('a', false);
```
