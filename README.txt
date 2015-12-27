To write to table call users in a postgres database called chitchat

First create the database using command

CREATE DATABASE chitchat
  WITH OWNER = adrianjackson
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       LC_COLLATE = 'en_US.UTF-8'
       LC_CTYPE = 'en_US.UTF-8'
       CONNECTION LIMIT = -1;

Then create the table using

CREATE TABLE users
(
  id serial NOT NULL,
  uuid character varying(64) NOT NULL,
  name character varying(255),
  email character varying(255) NOT NULL,
  password character varying(255) NOT NULL,
  created_at timestamp without time zone NOT NULL,
  CONSTRAINT users_pkey PRIMARY KEY (id),
  CONSTRAINT users_email_key UNIQUE (email),
  CONSTRAINT users_uuid_key UNIQUE (uuid)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE users
  OWNER TO adrianjackson;



Then you can add users like so

curl -i -v -X POST -d"name=Jim Defries&email=jdf@hollowspace.com&password=darkgrounds" http://localhost:3000/users

