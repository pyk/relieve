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
    curl -i -H "Accept: application/json" -X POST -d '{"user_email":"test", "user_gender":"male", "user_age":17, "user_profession":"student"}' http://localhost:8080/v0/users
    ```

2. psikolog endpoint
    
    ```
    curl -i -H "Accept: application/json" -X POST -d '{"psikolog_email":"test", "psikolog_name":"Prof. huru hara", "psikolog_image_url":"", "psikolog_bio":"biography"}' http://localhost:8080/v0/psikologs
    ```

3. post endpoint
    
    ```
    curl -i -H "Accept: application/json" -X POST -d '{"post_user_id":1, "post_psikolog_id":1, "post_title":"curhat nih", "post_category":"curhat", "post_content": "sample content", "post_image_url":""}' http://localhost:8080/v0/posts
    ```

4. comment endpoint
    
    ```
    curl -i -H "Accept: application/json" -X POST -d '{"comment_user_id":1, "comment_psikolog_id":1, "comment_post_id":1, "comment_text":"sample comment"}' http://localhost:8080/v0/comments
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

## Database

Change column type `text` to `integer`

    ALTER TABLE psikologs
        ALTER COLUMN psikolog_wisdom TYPE integer USING cast(psikolog_wisdom as int);

    ALTER TABLE posts
        ALTER COLUMN post_date SET DEFAULT CURRENT_TIMESTAMP;

    ALTER TABLE psikologs
        ALTER COLUMN psikolog_wisdom SET DEFAULT 0;
    
    ALTER TABLE posts
        ALTER COLUMN post_report_count SET DEFAULT 0;
    
    ALTER TABLE comments
        ALTER COLUMN comment_date SET DEFAULT CURRENT_TIMESTAMP;