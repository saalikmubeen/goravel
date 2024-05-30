-- ** This an example migration file. Write your up migrations here




-- Function to update the updated_at column with the current timestamp
CREATE TABLE some_table (
    id serial PRIMARY KEY,
    some_field VARCHAR ( 255 ) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);



-- add auto update of updated_at. If you already have this trigger
-- you can delete the next 7 lines
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Trigger to call the function before an update on some_table
CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON some_table
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();