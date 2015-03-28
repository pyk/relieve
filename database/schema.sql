-- User
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    user_email text NOT NULL UNIQUE,
    user_gender text,
    user_age integer,
    user_profession text
);

-- Psikolog
CREATE TABLE IF NOT EXISTS psikologs (
    psikolog_id SERIAL PRIMARY KEY,
    psikolog_email text NOT NULL UNIQUE,
    psikolog_name text,
    psikolog_image_url text,
    psikolog_wisdom integer DEFAULT 0,
    psikolog_bio text
);

-- Post
CREATE TABLE IF NOT EXISTS posts (
    post_id SERIAL PRIMARY KEY,
    post_user_id integer REFERENCES users(user_id) ON DELETE CASCADE,
    post_psikolog_id integer REFERENCES psikologs(psikolog_id),
    post_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    post_title text,
    post_category text,
    post_content text,
    post_image_url text,
    post_report_count integer DEFAULT 0
);

-- Comment
CREATE TABLE IF NOT EXISTS comments (
    comment_id SERIAL PRIMARY KEY,
    comment_user_id integer REFERENCES users(user_id) ON DELETE CASCADE,
    comment_psikolog_id integer REFERENCES psikologs(psikolog_id) ON DELETE CASCADE,
    comment_post_id integer REFERENCES posts(post_id),
    comment_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    comment_text text
);

-- Report
CREATE TABLE IF NOT EXISTS reports (
    report_id SERIAL PRIMARY KEY,
    report_user_id integer REFERENCES users(user_id) ON DELETE CASCADE,
    report_post_id integer REFERENCES posts(post_id)  ON DELETE CASCADE
);
