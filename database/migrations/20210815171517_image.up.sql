create table if not exists "public"."image"(
    -- "id"            uuid        not null default gen_random_uuid(),
    "id"            int         not null generated by default as identity(start with 1 increment by 1),
    "name"          citext      null,
    "image_url"     text        not null,
    "thumbnail_url" text        not null,
    "created_at"    timestamptz not null,
    "updated_at"    timestamptz null,
    "deleted_at"    timestamptz null,
    constraint "image_id_pkey"       primary key("id"),
    constraint "image_name_check"    check(("name" = '') is true or (length(("name")::text) >= 2) and (length(("name")::text) <= 64)) -- and ("name" ~ '^[[:alnum:]_]+$'::citext)),
    -- constraint "store_image_url"     check("image_url" ~ 'https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,255}\.[a-z]{2,9}\y([-a-zA-Z0-9@:%_\+.,~#?!&>//=]*)$'::text),
    -- constraint "store_thumbnail_url" check("thumbnail_url" ~ 'https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,255}\.[a-z]{2,9}\y([-a-zA-Z0-9@:%_\+.,~#?!&>//=]*)$'::text)
);

create index if not exists "image_name"
on "public"."image"(
	"name"
);

create trigger "_timestamps" before insert or update or delete
on "public"."image" for each row
    execute procedure "public"."tg__timestamps"();
