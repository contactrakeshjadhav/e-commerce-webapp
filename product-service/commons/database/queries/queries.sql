-- name: InsertTag :exec
INSERT INTO products(id, name, description, color)
VALUES (@id::TEXT, @name::TEXT, @description::TEXT, @color::TEXT);

-- name: UpdateTag :exec
UPDATE products
SET name             = @name::TEXT,
    description      = @description::TEXT,
    color            = @color::TEXT

WHERE id = @id::text;

-- name: FindTagById :one
SELECT *
FROM products
WHERE id = @id::TEXT;

-- name: FindTagByName :one
SELECT *
FROM products
WHERE name = @name::TEXT;

-- name: FindManyTags :many
    SELECT products.*, COUNT(*) OVER()
    FROM products
    WHERE
        (CASE WHEN @filter_by::text = 'name' THEN products.name ILIKE @pattern::TEXT ELSE TRUE END)
    AND (CASE WHEN array_length(@ids::text[], 1) > 0 THEN products.id = ANY (@ids::TEXT[]) ELSE TRUE END)

    ORDER BY
        CASE WHEN @ascending::boolean AND @order_by::text = 'name' THEN products.name END ASC,
        CASE WHEN @ascending::boolean AND @order_by::text = 'description' THEN products.description END ASC,
        CASE WHEN @ascending::boolean AND @order_by::text = 'color' THEN products.color END ASC,
        CASE WHEN NOT @ascending::boolean AND @order_by::text = 'name' THEN products.name END DESC,
        CASE WHEN NOT @ascending::boolean AND @order_by::text = 'description' THEN products.description END DESC,
        CASE WHEN NOT @ascending::boolean AND @order_by::text = 'color' THEN products.color END DESC

    LIMIT @page_limit::INT OFFSET @page_offset::INT;
-- name: FindAll :many
SELECT *
FROM products;

-- name: FindTagsByObjectId :many
Select product_id from object_has_products
Where object_id = @object_id::TEXT;

-- name: InsertObject :exec
INSERT INTO object_has_products(id, object_id, object_type, product_id)
VALUES (@id::TEXT, @object_id::TEXT, @object_type::TEXT, @product_id::TEXT);

-- name: FindTagsByPartialName :many
SELECT *
From products
WHERE
        (CASE WHEN @filter_by::text = 'name' THEN products.name ILIKE @pattern::TEXT ELSE TRUE END)
ORDER BY
        CASE WHEN @ascending::boolean AND @order_by::text = 'name' THEN products.name END ASC

LIMIT @page_limit::INT;

-- name: GetObjectTagMappingByTagID :many
SELECT * FROM object_has_products
WHERE product_id = @product_id::TEXT;

-- name: DeleteTagById :exec
Delete
FROM products
where id = @id::TEXT;

-- name: DeleteObjectTagById :exec
Delete
FROM object_has_products
where product_id = @product_id::TEXT;

-- name: DeleteObjectTagByObjectId :exec
DELETE
FROM object_has_products
WHERE object_id = @object_id::TEXT;

-- name: FindObjectTagMapping :one
Select * from object_has_products
Where product_id = @product_id::TEXT and object_id = @object_id::TEXT;

-- name: DeleteObjectTagMapping :exec
Delete
FROM object_has_products
where product_id = @product_id::TEXT and object_id = @object_id::TEXT;

-- name: FindObjectsByTagIdsAndObjectTypeUnion :many
SELECT object_id FROM object_has_products
WHERE product_id = ANY (@product_ids::TEXT[]) and object_type = @object_type::TEXT
GROUP by object_id;

-- name: FindObjectsByTagIdsAndObjectTypeIntersection :many
SELECT result.object_id from (SELECT oht.object_id, count(object_id) FROM object_has_products oht
WHERE product_id = ANY (@product_ids::TEXT[]) and object_type = @object_type::TEXT
GROUP BY object_id) as result where count = @product_ids_lenght::INT;

-- name: GetTagsByIDs :many
SELECT *
FROM products
WHERE id = ANY (@product_ids::TEXT[]);

-- name: GetTagsByObjectIDs :many
SELECT t.*, oht.object_id
FROM products t
JOIN object_has_products oht on t.id = oht.product_id
WHERE oht.object_id = ANY (@object_ids::TEXT[]);