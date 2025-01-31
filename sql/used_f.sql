
-- триггер для автоматического обновления path каждой категории в таблице при вставке или обновлении
CREATE OR REPLACE FUNCTION update_category_path() RETURNS TRIGGER AS $$
BEGIN
  IF NEW.parent_id IS NULL THEN
    NEW.path := NEW.id::TEXT;  -- Если это корневая категория, путь будет просто её ID
  ELSE
    SELECT path INTO NEW.path FROM categories WHERE id = NEW.parent_id;
    NEW.path := NEW.path || '-' || NEW.id::TEXT;  -- Добавляем ID категории к пути родителя
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_path_trigger
BEFORE INSERT OR UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION update_category_path();