DO $$ BEGIN
  CREATE TYPE room_type AS ENUM ('text','voice');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE friend_policy AS ENUM ('everyone','friends_of_friends');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE app_provider AS ENUM ('spotify','reddit','steam','twitter');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;
