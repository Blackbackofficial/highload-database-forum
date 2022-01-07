CREATE EXTENSION IF NOT EXISTS CITEXT; -- eliminate calls to lower

CREATE UNLOGGED TABLE users
(
    Nickname   CITEXT PRIMARY KEY,
    FullName   TEXT NOT NULL,
    About      TEXT NOT NULL DEFAULT '',
    Email      CITEXT UNIQUE
);

CREATE UNLOGGED TABLE forum
(
    Title    TEXT   NOT NULL,
    "user"   CITEXT REFERENCES "users"(Nickname),
    Slug     CITEXT PRIMARY KEY,
    Posts    INT    DEFAULT 0,
    Threads  INT    DEFAULT 0
);

CREATE UNLOGGED TABLE threads
(
    Id      SERIAL    PRIMARY KEY,
    Title   TEXT      NOT NULL,
    Author  CITEXT    REFERENCES "users"(Nickname),
    Forum   CITEXT    REFERENCES "forum"(Slug),
    Message TEXT      NOT NULL,
    Votes   INT       DEFAULT 0,
    Slug    CITEXT    UNIQUE,
    Created TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNLOGGED TABLE posts
(
    Parent    INT         DEFAULT 0,
    Author    CITEXT      REFERENCES "users"(Nickname),
    Message   TEXT        NOT NULL,
    IsEdited  BOOLEAN     DEFAULT FALSE,
    Forum     CITEXT      DEFAULT '',
    Thread    INT         DEFAULT 0 REFERENCES "threads"(Id),
    Created   TIMESTAMP   WITH TIME ZONE DEFAULT now(),
    Path      INT[]       DEFAULT ARRAY []::INTEGER[]
);

CREATE UNLOGGED TABLE votes
(
    Author  CITEXT   UNIQUE REFERENCES "users"(Nickname),
    Voice   INT      NOT NULL DEFAULT 0,
    Thread  INT      UNIQUE REFERENCES "threads"(Id)
);


CREATE UNLOGGED TABLE users_forum
(
    Nickname  CITEXT  REFERENCES "users"(Nickname) UNIQUE NOT NULL,
    FullName  TEXT    NOT NULL,
    -- ?????
    About     TEXT    NOT NULL DEFAULT '',
    Email     CITEXT  UNIQUE,
    Slug      CITEXT  REFERENCES "forum"(Slug) UNIQUE NOT NULL
);