DROP TABLE IF EXISTS lists CASCADE;
DROP TABLE IF EXISTS texts;

CREATE TABLE lists (
  id serial UNIQUE NOT NULL,
  title text UNIQUE NOT NULL
);

INSERT INTO lists (
  title
)
VALUES
  ('Primary lit'),
  ('Secondary lit');

CREATE TABLE texts (
  id serial UNIQUE NOT NULL,
  title text NOT NULL,
  author_last text NOT NULL,
  author_first text NOT NULL,
  publication_date integer CHECK (publication_date > -1000 AND publication_date < 3000),
  list integer REFERENCES lists (id) ON DELETE CASCADE
);

INSERT INTO
  texts (
    title,
    author_last,
    author_first,
    publication_date,
    list
  )
VALUES
  ('The Greats: Duns Scotus', 'Ward', 'Thomas', 2020, 2),
  ('The Greats: Shakespeare', 'Smith', 'Matthew', 2019, 2),
  ('The Greats: Avicenna', 'Fatigati', 'Michael', 2021, 2),
  ('Logic', 'Ibn Sina', 'Avicenna', 1000, 1),
  ('Metaphysics', 'Ibn Sina', 'Avicenna', 1010, 1),
  ('Psychology', 'Ibn Sina', 'Avicenna', 1020, 1);