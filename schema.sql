CREATE SCHEMA IF NOT EXISTS content;


##Table for News
CREATE TABLE IF NOT EXISTS content.`News` (
                                `Id` bigint NOT NULL AUTO_INCREMENT PRIMARY KEY ,
                                `Title` tinytext NOT NULL,
                                `Content` longtext NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



##Table for connection news and categories
CREATE TABLE IF NOT EXISTS content.NewsCategories (
                                        NewsId bigint,
                                        CategoryId bigint
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


## Ex. procedure for getting all the rows from News
CREATE PROCEDURE content.GetNews()
    BEGIN
        SELECT n.Id, n.Title, n.Content FROM content.News n;
    END;


## Ex. procedure for getting all the rows from Categories for ID
CREATE PROCEDURE content.GetCategories(IN NewsId bigint)
    BEGIN
        SELECT CategoryId from content.NewsCategories nc where nc.NewsId = NewsId;
    end;


## Ex. procedure for modifying news
CREATE PROCEDURE content.UpdateNews(IN newsId INT, IN newTitle TEXT, IN newContent TEXT)
BEGIN
    -- Удаление строки по Id
DELETE FROM content.News WHERE Id = newsId;

-- Добавление новой строки
INSERT INTO content.News (Id, Title, Content) VALUES (newsId, newTitle, newContent);

END;

## Ex. procedure for modifying Categories
CREATE PROCEDURE content.UpdateCategory(IN p_newsID INT, IN p_CategoryIDs  varchar(255))
    BEGIN
        DECLARE v_CategoryId INT;

        START TRANSACTION;

        -- Удаление всех строк из таблицы NewsCategories для заданного NewsId
        DELETE FROM content.NewsCategories WHERE NewsId = p_NewsId;

        -- Разделение списка категорий
        SET @categories := p_CategoryIDs;
        WHILE CHAR_LENGTH(@categories) > 0 DO
                SET @category := TRIM(SUBSTRING_INDEX(@categories, ',', 1));
                SET @categories := TRIM(SUBSTRING(@categories, CHAR_LENGTH(@category) + 2));

                -- Преобразование строки в число
                SET v_CategoryId = CAST(@category AS UNSIGNED);

                -- Вставка новых данных в таблицу NewsCategories
                INSERT INTO content.NewsCategories (NewsId, CategoryId) VALUES (p_NewsId, v_CategoryId);
            END WHILE;

        COMMIT;
    END;