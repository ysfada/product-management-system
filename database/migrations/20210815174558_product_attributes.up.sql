create table if not exists "public"."product_attributes"(
    -- "id"                 uuid        not null default gen_random_uuid(),
    "id"                 int         not null generated by default as identity(start with 1 increment by 1),
    "product_variant_id" int         not null,
    "attribute_id"       int         not null,
    "created_at"         timestamptz not null,
    "updated_at"         timestamptz null,
    "deleted_at"         timestamptz null,
    foreign key("product_variant_id")   references "product_variant"("id")   on delete restrict deferrable initially deferred,
    foreign key("attribute_id")         references "attribute"("id") on delete restrict deferrable initially deferred,
    constraint "product_attributes_id_pkey" primary key("id")
);

create unique index if not exists "product_attributes_product_variant_id_attribute_id_un_iq"
on "public"."product_attributes"(
	"product_variant_id",
	"attribute_id"
);

create index if not exists "product_attributes_product_variant_id"
on "public"."product_attributes"(
	"product_variant_id"
);

create index if not exists "product_attributes_attribute_id"
on "public"."product_attributes"(
	"attribute_id"
);

create trigger "_timestamps" before insert or update or delete
on "public"."product_attributes" for each row
    execute procedure "public"."tg__timestamps"();
