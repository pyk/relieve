# Relieve endpoint

## Development

1. setup postgres database with the following [schema][schema]

[schema]: link-to-schema.sql

2. export port and datatabase environment variable
    
    export PORT=8080
    export DATABASE_URL=postgres://user:password@[host]:[port]/[dbname]

3. Build & run it
    
    go install
    relieve

4. start testing endpoint

## Testing endpoint

1. users endpoint

```
curl -i -H "Accept: application/json" -X POST -d '{"user_email":"test2", "user_gender":"male", "user_age":"17", "user_profession":"student"}' http://localhost:8080/v0/users
```
