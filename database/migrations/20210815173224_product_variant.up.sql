create table if not exists "public"."product_variant"(
    -- "id"           uuid        not null default gen_random_uuid(),
    "id"           int           not null generated by default as identity(start with 1 increment by 1),
    "product_id"   int           not null,
    "name"         citext        not null,
    "price"        decimal(19,4) not null,
    "stock"        int           not null,
    "created_at"   timestamptz   not null,
    "updated_at"   timestamptz   null,
    "deleted_at"   timestamptz   null,
    foreign key ("product_id")   references "product"("id")   on delete restrict deferrable initially deferred,
    constraint "product_variant_id_pkey" primary key("id"),
    constraint "product_stock_check" check("stock"::int >= 0),
    constraint "product_name_check"  check((length(("name")::text) >= 2) and (length(("name")::text) <= 32)) -- and ("name" ~ '^[[:alnum:]_]+$'::citext))
);

create index if not exists "product_variant_product_id"
on "public"."product_variant"(
	"product_id"
);

create index if not exists "product_variant_name"
on "public"."product_variant"(
	"name"
);

create trigger "_timestamps" before insert or update or delete
on "public"."product_variant" for each row
    execute procedure "public"."tg__timestamps"();