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

        CONSTRAINT fk_user FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    );

	CREATE TABLE IF NOT EXISTS profile (
		id SERIAL PRIMARY KEY,
       firstname text NOT NULL,
       lastname text NOT NULL,
       age bigint NOT NULL,
       gender text NOT NULL,
       target text NOT NULL,
       about text NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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