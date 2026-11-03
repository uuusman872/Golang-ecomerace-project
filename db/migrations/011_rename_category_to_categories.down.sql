ALTER INDEX idx_categories_is_active RENAME TO idx_category_is_active;
ALTER INDEX idx_categories_deleted_at RENAME TO idx_category_deleted_at;

ALTER TABLE categories RENAME TO category;
