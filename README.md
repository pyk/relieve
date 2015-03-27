# Relieve endpoint

## Development

1. setup postgres database with the following [schema][schema]

2. export port and datatabase environment variable
    
    ```
    export PORT=8080
    export DATABASE_URL=postgres://user:password@[host]:[port]/[dbname]
    ```

3. Build & run it
    
    ```
    go install
    relieve
    ``` 

4. start testing endpoint

## Testing endpoint

1. users endpoint
    
    ```
    curl -i -H "Accept: application/json" -X POST -d '{"user_email":"test2", "user_gender":"male", "user_age":"17", "user_profession":"student"}' http://localhost:8080/v0/users
    ```

## Deploy to heroku

Create app with custom buildpack

    heroku apps:create relieve-endpoint -b https://github.com/kr/heroku-buildpack-go.git
    echo 'web: relieve' > Procfile
    heroku addons:add heroku-postgresql

Push code to heroku

    godep save
    git add . && git commit -m "deploy to heroku"
    git push heroku master

Make sure [godep][godep] installed

    go get github.com/tools/godep

[godep]: https://github.com/tools/godep

Update heroku database

    heroku pg:psql

and create the following [schema][schema]

[schema]: https://github.com/pyk/relieve/blob/master/database/schema.sql