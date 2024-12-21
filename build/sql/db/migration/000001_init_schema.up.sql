CREATE TABLE IF NOT EXISTS profile (
                                       id SERIAL PRIMARY KEY,
                                       firstname text NOT NULL,
                                       lastname text NOT NULL,
                                       gender text NOT NULL,
                                       birthday_date text NOT NULL,
                                       target text NOT NULL,
                                       about text NOT NULL,
                                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username text,
                                     password text,
                                     profile INT NOT NULL,

                                     CONSTRAINT fk_profile FOREIGN KEY (profile)
    REFERENCES profile (id)
    ON DELETE SET NULL
    ON UPDATE CASCADE,

    CONSTRAINT unique_username UNIQUE (username)
    );
CREATE TABLE IF NOT EXISTS photo (
                                     id SERIAL PRIMARY KEY,
                                     user_id bigint NOT NULL,
                                     link text NOT NULL UNIQUE,
                                     number bigint NOT NULL,

                                     CONSTRAINT fk_user FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
    );


CREATE TABLE IF NOT EXISTS reaction (
                                        id SERIAL PRIMARY KEY ,
                                        author bigint NOT NULL ,
                                        receiver bigint NOT NULL,
                                        type boolean,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                        CONSTRAINT fk_author FOREIGN KEY (author)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    CONSTRAINT fk_receiver FOREIGN KEY (receiver)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    CONSTRAINT unique_pair UNIQUE (author, receiver)
    );

CREATE TABLE IF NOT EXISTS report (
    id SERIAL PRIMARY KEY,
    author INT NOT NULL,
    receiver INT NOT NULL,
    reason text NOT NULL CHECK (reason in ('about', 'photo', 'fake', 'abuse', 'another')),
    body text NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_author FOREIGN KEY (author)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    CONSTRAINT fk_receiver FOREIGN KEY (receiver)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    author INT NOT NULL,
    receiver INT NOT NULL,
    body text NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_author FOREIGN KEY (author)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    CONSTRAINT fk_receiver FOREIGN KEY (receiver)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS survey (
    id SERIAL PRIMARY KEY,
    author INT NOT NULL,
    question text NOT NULL,
    comment text NOT NULL,
    rating INT NOT NULL,
    grade int NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_author FOREIGN KEY (author)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS question (
    id SERIAL PRIMARY KEY,
    content text NOT NULL UNIQUE,
    grade INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS daily_likes (
    id SERIAL PRIMARY KEY,
    userID INT NOT NULL,
    likes_count INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY (userID)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS purchased_likes (
    id SERIAL PRIMARY KEY,
    userID INT NOT NULL,
    likes_count INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY (userID)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS balance (
    id SERIAL PRIMARY KEY,
    userID INT NOT NULL,
    balance INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY (userID)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS product (
    id SERIAL PRIMARY KEY,
    title text NOT NULL UNIQUE,
    description text NOT NULL,
    imagelink text NOT NULL,
    price INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);