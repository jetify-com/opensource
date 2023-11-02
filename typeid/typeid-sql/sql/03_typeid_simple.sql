-- Implementation of typeids in SQL (Postgres).
-- This file:
-- + Defines functions to generate and validate typeids in SQL.

-- Create a `text` type.
-- Note that the "uuid" field should be a UUID v7.

-- Function that generates a random typeid of the given prefix.
-- This depends on the `uuid_generate_v7` function defined in `uuid_v7.sql`.
create or replace function typeid_generate(prefix text)
    returns text
as $$
begin
    if (prefix is null) or not (prefix ~ '^[a-z]{0,63}$') then
        raise exception 'typeid prefix must match the regular expression [a-z]{0,63}';
    end if;
    return prefix || '_' || base32_encode(uuid_generate_v7());
end
$$
    language plpgsql
    volatile;

-- Function that checks if a typeid is valid, for the given type prefix.
-- It also enforces that the UUID is a v7 UUID.
-- NOTE: we might want to make the version check optional.
create or replace function typeid_check(type_id text, expected_type text)
    returns boolean
as $$
declare
    prefix text;
    bytes bytea;
    ver int;
begin
    prefix = split_part(type_id, '_', 1);
    bytes = uuid_send(base32_decode(split_part(type_id, '_', 2)));
    ver = (get_byte(bytes, 6) >> 4)::bit(4)::int;
    return prefix = expected_type AND (ver = 7 OR split_part(type_id, '_', 2) = '00000000-0000-0000-0000-000000000000');
end
$$
    language plpgsql
    immutable;


-- Function that parses a string into a typeid.
create or replace function typeid_parse(typeid text)
    returns text
as $$
declare
    prefix text;
    suffix text;
begin
    if (typeid is null) then
        return null;
    end if;
    if position('_' in typeid) = 0 then
        return ('', base32_decode(typeid))::text;
    end if;
    prefix = split_part(typeid, '_', 1);
    suffix = split_part(typeid, '_', 2);
    if prefix is null or prefix = '' then
        raise exception 'typeid prefix cannot be empty with a delimiter';
    end if;
    -- prefix must match the regular expression [a-z]{0,63}
    if not prefix ~ '^[a-z]{0,63}$' then
        raise exception 'typeid prefix must match the regular expression [a-z]{0,63}';
    end if;

    return (prefix, base32_decode(suffix))::text;
end
$$
    language plpgsql
    immutable;