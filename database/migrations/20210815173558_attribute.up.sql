create table if not exists "public"."attribute"(
    -- "id"         uuid        not null default gen_random_uuid(),
    "id"         int         not null generated by default as identity(start with 1 increment by 1),
    "name"       citext      not null,
    "type"       citext      not null,
    "created_at" timestamptz not null,
    "updated_at" timestamptz null,
    "deleted_at" timestamptz null,
    constraint "attribute_id_pkey"    primary key("id"),
    constraint "attribute_name_check" check((length(("name")::text) >= 1) and (length(("name")::text) <= 32)), -- and ("name" ~ '^[[:alnum:]_]+$'::citext))
    constraint "attribute_type_check" check((length(("type")::text) >= 1) and (length(("type")::text) <= 32)) -- and ("type" ~ '^[[:alnum:]_]+$'::citext))
);

create index if not exists "attribute_name"
on "public"."attribute"(
	"name"
);

create trigger "_timestamps" before insert or update or delete
on "public"."attribute" for each row
    execute procedure "public"."tg__timestamps"();
