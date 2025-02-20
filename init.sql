DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'cabri') THEN
      PERFORM dblink_exec('dbname=postgres user=' || current_user, 'CREATE DATABASE cabri');
   END IF;
END
$$;
