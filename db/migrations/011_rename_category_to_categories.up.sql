ALTER TABLE category RENAME TO categories;

ALTER INDEX idx_category_is_active RENAME TO idx_categories_is_active;
ALTER INDEX idx_category_deleted_at RENAME TO idx_categories_deleted_at;
