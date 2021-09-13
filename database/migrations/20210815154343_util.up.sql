create extension if not exists plpgsql     with schema pg_catalog;
create extension if not exists "uuid-ossp" with schema public;
create extension if not exists citext      with schema public;
create extension if not exists pgcrypto    with schema public;
create extension if not exists unaccent    with schema public;

create function "public"."tg__timestamps"() returns trigger as $$
begin
    NEW."created_at" = (
        case
            when TG_OP = 'INSERT'
                then now()
                else OLD."created_at"
        end
    );
    NEW."updated_at" = (
        case
            when TG_OP = 'UPDATE' and OLD."updated_at" >= now()
                then OLD."updated_at" + interval '1 millisecond'
                else now()
        end
    );
    NEW."deleted_at" = (
        case when TG_OP = 'DELETE' and OLD."deleted_at" >= now()
            then OLD."deleted_at" + interval '1 millisecond'
            else now()
        end
    );
    return NEW;
end;
$$ language plpgsql volatile set search_path to pg_catalog, public, pg_temp;
