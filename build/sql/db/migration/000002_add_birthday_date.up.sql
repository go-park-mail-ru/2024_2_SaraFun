ALTER TABLE profile ADD birthday_date text NOT NULL;
UPDATE profile SET birthday_date = '2000-01-01';