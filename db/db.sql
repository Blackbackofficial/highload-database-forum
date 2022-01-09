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
    Id        SERIAL      PRIMARY KEY,
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

create index if not exists user_nickname ON users using hash(nickname);
create index if not exists user_email ON users using hash(email);
create index if not exists forum_slug ON forum using hash(slug);
create index if not exists thr_date ON threads (created);
create index if not exists thr_forum_date ON threads(forum, created);
create index if not exists thr_forum ON threads using hash(forum);
create index if not exists thr_slug ON threads using hash(slug);
create index if not exists post_id_path on posts(id, (path[1]));
create index if not exists post_thread_path_id on posts(thread, path, id);
create index if not exists post_thread_id_path1_parent on posts(thread, id, (path[1]), parent);
create index if not exists post_path1 on posts((path[1]));
create index if not exists post_thr_id ON posts(thread);
create index if not exists post_thread_id on posts(thread, id);
create unique index if not exists forum_users_unique on users_forum(slug, nickname);
create unique index if not exists vote_unique on votes(Author, Thread);

CREATE OR REPLACE FUNCTION insertVotes() RETURNS TRIGGER AS
$update_vote$
BEGIN
    UPDATE thread SET votes=(votes+NEW.voice) WHERE id=NEW.thread;
    return NEW;
end
$update_vote$ LANGUAGE plpgsql;

CREATE TRIGGER add_voice
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insertVotes();


CREATE OR REPLACE FUNCTION updatePostUserForum() RETURNS TRIGGER AS
$update_forum_posts$
DECLARE
    m_fullname CITEXT;
    m_about    CITEXT;
    m_email CITEXT;
BEGIN
    SELECT fullname, about, email FROM users WHERE nickname = NEW.author INTO m_fullname, m_about, m_email;
    INSERT INTO users_forum (nickname, fullname, about, email, Slug)
    VALUES (New.Author,m_fullname, m_about, m_email, NEW.forum) on conflict do nothing;
    return NEW;
end
$update_forum_posts$ LANGUAGE plpgsql;

CREATE TRIGGER post_insert_user_forum
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE updatePostUserForum();

CREATE OR REPLACE FUNCTION updateThreadUserForum() RETURNS TRIGGER AS
$update_forum_threads$
DECLARE
    author_nick CITEXT;
    m_fullname CITEXT;
    m_about    CITEXT;
    m_email CITEXT;
BEGIN
    SELECT Nickname, fullname, about, email
    FROM users WHERE Nickname = new.Author INTO author_nick, m_fullname, m_about, m_email;
    INSERT INTO users_forum (nickname, fullname, about, email, Slug)
    VALUES (author_nick,m_fullname, m_about, m_email, NEW.forum) on conflict do nothing;
    return NEW;
end
$update_forum_threads$ LANGUAGE plpgsql;

CREATE TRIGGER thread_insert_user_forum
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE updateThreadUserForum();


CREATE OR REPLACE FUNCTION updateVotes() RETURNS TRIGGER AS
$update_votes$
BEGIN
    IF OLD.Voice <> NEW.Voice THEN
        UPDATE threads SET votes=(votes+NEW.Voice*2) WHERE id=NEW.Thread;
    END IF;
    return NEW;
end
$update_votes$ LANGUAGE plpgsql;

CREATE TRIGGER edit_voice
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE updateVotes();


CREATE OR REPLACE FUNCTION updateCountOfThreads() RETURNS TRIGGER AS
$update_forums$
BEGIN
    UPDATE forum SET Threads=(Threads+1) WHERE slug=NEW.forum;
    return NEW;
end
$update_forums$ LANGUAGE plpgsql;

CREATE TRIGGER addThreadInForum
    BEFORE INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE updateCountOfThreads();


CREATE OR REPLACE FUNCTION updatePath() RETURNS TRIGGER AS
$update_paths$
DECLARE
    parentPath         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent IS NULL) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT path FROM posts WHERE id = new.parent INTO parentPath;
        SELECT thread FROM posts WHERE id = parentPath[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread <> NEW.thread THEN
            RAISE EXCEPTION 'parent is from different thread' USING ERRCODE = '00409';
        end if;

        NEW.path := NEW.path || parentPath || new.id;
    end if;
    UPDATE forum SET Posts=Posts + 1 WHERE forum.slug = new.forum;
    RETURN new;
end
$update_paths$ LANGUAGE plpgsql;

CREATE TRIGGER update_path_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE updatePath();

cluster users_forum using forum_users_unique;

analyse users_forum;
analyse threads;
analyse posts;