ALTER TABLE
    posts
ADD constraint fk_user_id FOREIGN KEY (user_id)
REFERENCES  users (id);