-- User
-- INSERT INTO users(user_email, user_gender, user_age, user_profession) VALUES ('test1', 'test L', 18, 'test job');
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    user_email text NOT NULL UNIQUE,
    user_gender text,
    user_age integer,
    user_profession text
);

-- Psikolog
-- INSERT INTO psikologs(psikolog_email, psikolog_name, psikolog_image_url, psikolog_bio) VALUES ('test1', 'test name', 'test image', 'tst bio');
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
    post_image_url text DEFAULT '',
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

-- Wisdom
-- cek the sum of psikolog wisdom points
-- SELECT SUM(wisdom_point) FROM wisdom_points WHERE wisdom_psikolog_id=1;
-- cek if record exists
-- SELECT EXISTS(SELECT 1 FROM wisdom_points WHERE wisdom_user_id=1 AND wisdom_psikolog_id=5);
-- insert into table
-- INSERT INTO wisdom_points(wisdom_user_id,wisdom_psikolog_id) VALUES (1,1);   
CREATE TABLE IF NOT EXISTS wisdom_points (
    wisdom_point integer DEFAULT 10,
    wisdom_user_id integer REFERENCES users(user_id) ON DELETE CASCADE,
    wisdom_psikolog_id integer REFERENCES psikologs(psikolog_id) ON DELETE CASCADE,
    UNIQUE(wisdom_user_id, wisdom_psikolog_id)
);