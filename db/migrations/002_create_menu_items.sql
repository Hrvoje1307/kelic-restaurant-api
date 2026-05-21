CREATE TABLE IF NOT EXISTS menu_items (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id  UUID         REFERENCES menu_categories(id) ON DELETE SET NULL,
    name         TEXT         NOT NULL,
    description  TEXT         NOT NULL DEFAULT '',
    price        NUMERIC(10,2) NOT NULL,
    image_url    TEXT         NOT NULL DEFAULT '',
    is_available BOOLEAN      NOT NULL DEFAULT TRUE,
    allergens    TEXT[]       NOT NULL DEFAULT '{}',
    tags         TEXT[]       NOT NULL DEFAULT '{}',
    sort_order   INTEGER      NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_menu_items_category_id  ON menu_items(category_id);
CREATE INDEX IF NOT EXISTS idx_menu_items_is_available ON menu_items(is_available);
