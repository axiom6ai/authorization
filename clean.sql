CREATE FUNCTION cleanAll() RETURNS void AS $$
DECLARE
    tmp VARCHAR(512);
DECLARE names CURSOR FOR
    select tablename from pg_tables where schemaname='public';
BEGIN
    FOR stmt IN names LOOP
        tmp := 'DROP TABLE '||quote_ident(stmt.tablename)||' CASCADE;';
        RAISE NOTICE 'notice:%', tmp;
        EXECUTE 'DROP TABLE '||quote_ident(stmt.tablename)||' CASCADE;';
    END LOOP;
    RAISE NOTICE 'finished....';
END;

$$ LANGUAGE plpgsql;
