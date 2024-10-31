CREATE TYPE reaction_type AS ENUM ('like', 'dislike');
CREATE TYPE message_status AS ENUM ('sent', 'delivered', 'read');
CREATE TYPE subscription_status AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE "profile" (
   "id" bigint PRIMARY KEY ,
   "firstname" text NOT NULL,
   "lastname" text NOT NULL,
   "age" bigint NOT NULL CHECK (age >= 18),
   "gender" gender NOT NULL,
   "target" text NOT NULL,
   "about" text NOT NULL,
   "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "subscription" (
    "id" bigint PRIMARY KEY ,
    "status" subscription_status NOT NULL,
    "start_date" TIMESTAMP,
    "end_date" TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "user" (
    "id" bigint PRIMARY KEY ,
    "username" text NOT NULL UNIQUE,
    "email" text NOT NULL UNIQUE,
    "password_hash" text NOT NULL,
    "profile" bigint NOT NULL UNIQUE,
    "subscribe" bigint NOT NULL,
    "balance" bigint NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_profile FOREIGN KEY (profile)
    REFERENCES profile (id)
    ON DELETE SET NULL
    ON UPDATE CASCADE,

    CONSTRAINT fk_subscribe FOREIGN KEY (subscribe)
    REFERENCES subscription (id)
    ON DELETE SET NULL
    ON UPDATE CASCADE
);



CREATE TABLE "reaction" (
    "id" bigint PRIMARY KEY ,
    "author" bigint NOT NULL ,
    "receiver" bigint NOT NULL,
    "type" reaction_type,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_author FOREIGN KEY (author)
    REFERENCES "user" (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    CONSTRAINT fk_receiver FOREIGN KEY (receiver)
    REFERENCES "user" (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE

    CONSTRAINT unique_pair UNIQUE (author, receiver)
);

CREATE TABLE "purchase" (
    "id" bigint PRIMARY KEY ,
    "author" bigint NOT NULL,
    "price" bigint NOT NULL,
    "about" text NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_author FOREIGN KEY (author)
    REFERENCES "user" (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE "complaint" (
    "id" bigint PRIMARY KEY ,
    "author" bigint NOT NULL,
    "category" text NOT NULL,
    "body" text NOT NULL,
    "status" text NOT NULL,
    "answer" text NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_author FOREIGN KEY (author)
    REFERENCES "user" (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE "message" (
   "id" bigint PRIMARY KEY,
   "sender" bigint NOT NULL,
   "receiver" bigint NOT NULL,
   "body" text NOT NULL,
   "status" message_status,
   "sendDate" TIMESTAMP,
   "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

   CONSTRAINT fk_sender FOREIGN KEY (sender)
   REFERENCES "user" (id)
   ON DELETE CASCADE
   ON UPDATE CASCADE,

   CONSTRAINT fk_receiver FOREIGN KEY (receiver)
   REFERENCES "user" (id)
   ON DELETE CASCADE
   ON UPDATE CASCADE
);

CREATE TABLE "photo" (
    "id" bigint PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "link" text NOT NULL UNIQUE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY (user_id)
    REFERENCES "user" (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

