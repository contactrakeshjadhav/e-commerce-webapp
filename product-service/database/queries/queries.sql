-- name: InsertTag :exec
INSERT INTO tags(id, name, description, color)
VALUES (@id::TEXT, @name::TEXT, @description::TEXT, @color::TEXT);

-- name: UpdateTag :exec
UPDATE tags
SET name             = @name::TEXT,
    description      = @description::TEXT,
    color            = @color::TEXT

WHERE id = @id::text;

-- name: FindTagById :one
SELECT *
FROM tags
WHERE id = @id::TEXT;

-- name: FindTagByName :one
SELECT *
FROM tags
WHERE name = @name::TEXT;

-- name: FindManyTags :many
    SELECT tags.*, COUNT(*) OVER()
    FROM tags
    WHERE
        (CASE WHEN @filter_by::text = 'name' THEN tags.name ILIKE @pattern::TEXT ELSE TRUE END)
    AND (CASE WHEN array_length(@ids::text[], 1) > 0 THEN tags.id = ANY (@ids::TEXT[]) ELSE TRUE END)

    ORDER BY
        CASE WHEN @ascending::boolean AND @order_by::text = 'name' THEN tags.name END ASC,
        CASE WHEN @ascending::boolean AND @order_by::text = 'description' THEN tags.description END ASC,
        CASE WHEN @ascending::boolean AND @order_by::text = 'color' THEN tags.color END ASC,
        CASE WHEN NOT @ascending::boolean AND @order_by::text = 'name' THEN tags.name END DESC,
        CASE WHEN NOT @ascending::boolean AND @order_by::text = 'description' THEN tags.description END DESC,
        CASE WHEN NOT @ascending::boolean AND @order_by::text = 'color' THEN tags.color END DESC

    LIMIT @page_limit::INT OFFSET @page_offset::INT;
-- name: FindAll :many
SELECT *
FROM tags;

-- name: FindTagsByObjectId :many
Select tag_id from object_has_tags
Where object_id = @object_id::TEXT;

-- name: InsertObject :exec
INSERT INTO object_has_tags(id, object_id, object_type, tag_id)
VALUES (@id::TEXT, @object_id::TEXT, @object_type::TEXT, @tag_id::TEXT);

-- name: FindTagsByPartialName :many
SELECT *
From tags
WHERE
        (CASE WHEN @filter_by::text = 'name' THEN tags.name ILIKE @pattern::TEXT ELSE TRUE END)
ORDER BY
        CASE WHEN @ascending::boolean AND @order_by::text = 'name' THEN tags.name END ASC

LIMIT @page_limit::INT;

-- name: GetObjectTagMappingByTagID :many
SELECT * FROM object_has_tags
WHERE tag_id = @tag_id::TEXT;

-- name: DeleteTagById :exec
Delete
FROM tags
where id = @id::TEXT;

-- name: DeleteObjectTagById :exec
Delete
FROM object_has_tags
where tag_id = @tag_id::TEXT;

-- name: DeleteObjectTagByObjectId :exec
DELETE
FROM object_has_tags
WHERE object_id = @object_id::TEXT;

-- name: FindObjectTagMapping :one
Select * from object_has_tags
Where tag_id = @tag_id::TEXT and object_id = @object_id::TEXT;

-- name: DeleteObjectTagMapping :exec
Delete
FROM object_has_tags
where tag_id = @tag_id::TEXT and object_id = @object_id::TEXT;

-- name: FindObjectsByTagIdsAndObjectTypeUnion :many
SELECT object_id FROM object_has_tags
WHERE tag_id = ANY (@tag_ids::TEXT[]) and object_type = @object_type::TEXT
GROUP by object_id;

-- name: FindObjectsByTagIdsAndObjectTypeIntersection :many
SELECT result.object_id from (SELECT oht.object_id, count(object_id) FROM object_has_tags oht
WHERE tag_id = ANY (@tag_ids::TEXT[]) and object_type = @object_type::TEXT
GROUP BY object_id) as result where count = @tag_ids_lenght::INT;

-- name: GetTagsByIDs :many
SELECT *
FROM tags
WHERE id = ANY (@tag_ids::TEXT[]);

-- name: GetTagsByObjectIDs :many
SELECT t.*, oht.object_id
FROM tags t
JOIN object_has_tags oht on t.id = oht.tag_id
WHERE oht.object_id = ANY (@object_ids::TEXT[]);