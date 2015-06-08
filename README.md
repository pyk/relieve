# Relieve endpoint

## Development

1. setup postgres database with the following [schema][schema]
    
    ```
    sudo su - postgres
    createuser -P relieve
    createdb -O relieve relieve

    note:
    update pg_hba.conf from
    local  all      all          peer
    to
    local  all      all          md5
    
    psql -U relieve -d relieve
    ```

2. export port and datatabase environment variable
    
    ```
    export PORT=4000
    export DATABASE_URL=postgres://relieve:relieve@localhost/relieve
    ```

3. Build & run it
    
    ```
    go install
    relieve
    ```

4. start testing endpoint

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

## Database

Change column type `text` to `integer`

    ALTER TABLE psikologs
        ALTER COLUMN psikolog_wisdom TYPE integer USING cast(psikolog_wisdom as int);

Add default to existing columns

    ALTER TABLE posts
        ALTER COLUMN post_date SET DEFAULT CURRENT_TIMESTAMP;

    ALTER TABLE psikologs
        ALTER COLUMN psikolog_wisdom SET DEFAULT 0;
    
    ALTER TABLE posts
        ALTER COLUMN post_report_count SET DEFAULT 0;
    
    ALTER TABLE comments
        ALTER COLUMN comment_date SET DEFAULT CURRENT_TIMESTAMP;