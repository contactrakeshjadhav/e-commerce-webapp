// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package database

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const deleteObjectProductById = `-- name: DeleteObjectProductById :exec
Delete
FROM object_has_tags
where tag_id = $1::TEXT
`

func (q *Queries) DeleteObjectProductById(ctx context.Context, tagID string) error {
	_, err := q.db.ExecContext(ctx, deleteObjectProductById, tagID)
	return err
}

const deleteObjectProductByObjectId = `-- name: DeleteObjectProductByObjectId :exec
DELETE
FROM object_has_tags
WHERE object_id = $1::TEXT
`

func (q *Queries) DeleteObjectProductByObjectId(ctx context.Context, objectID string) error {
	_, err := q.db.ExecContext(ctx, deleteObjectProductByObjectId, objectID)
	return err
}

const deleteObjectProductMapping = `-- name: DeleteObjectProductMapping :exec
Delete
FROM object_has_tags
where tag_id = $1::TEXT and object_id = $2::TEXT
`

type DeleteObjectProductMappingParams struct {
	ProductID    string
	ObjectID string
}

func (q *Queries) DeleteObjectProductMapping(ctx context.Context, arg DeleteObjectProductMappingParams) error {
	_, err := q.db.ExecContext(ctx, deleteObjectProductMapping, arg.ProductID, arg.ObjectID)
	return err
}

const deleteProductById = `-- name: DeleteProductById :exec
Delete
FROM tags
where id = $1::TEXT
`

func (q *Queries) DeleteProductById(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteProductById, id)
	return err
}

const findAll = `-- name: FindAll :many
SELECT id, name, description, color
FROM tags
`

func (q *Queries) FindAll(ctx context.Context) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, findAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Color,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findManyProducts = `-- name: FindManyProducts :many
    SELECT tags.id, tags.name, tags.description, tags.color, COUNT(*) OVER()
    FROM tags
    WHERE
        (CASE WHEN $1::text = 'name' THEN tags.name ILIKE $2::TEXT ELSE TRUE END)
    AND (CASE WHEN array_length($3::text[], 1) > 0 THEN tags.id = ANY ($3::TEXT[]) ELSE TRUE END)

    ORDER BY
        CASE WHEN $4::boolean AND $5::text = 'name' THEN tags.name END ASC,
        CASE WHEN $4::boolean AND $5::text = 'description' THEN tags.description END ASC,
        CASE WHEN $4::boolean AND $5::text = 'color' THEN tags.color END ASC,
        CASE WHEN NOT $4::boolean AND $5::text = 'name' THEN tags.name END DESC,
        CASE WHEN NOT $4::boolean AND $5::text = 'description' THEN tags.description END DESC,
        CASE WHEN NOT $4::boolean AND $5::text = 'color' THEN tags.color END DESC

    LIMIT $7::INT OFFSET $6::INT
`

type FindManyProductsParams struct {
	FilterBy   string
	Pattern    string
	Ids        []string
	Ascending  bool
	OrderBy    string
	PageOffset int32
	PageLimit  int32
}

type FindManyProductsRow struct {
	ID          string
	Name        sql.NullString
	Description sql.NullString
	Color       sql.NullString
	Count       int64
}

func (q *Queries) FindManyProducts(ctx context.Context, arg FindManyProductsParams) ([]FindManyProductsRow, error) {
	rows, err := q.db.QueryContext(ctx, findManyProducts,
		arg.FilterBy,
		arg.Pattern,
		pq.Array(arg.Ids),
		arg.Ascending,
		arg.OrderBy,
		arg.PageOffset,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindManyProductsRow
	for rows.Next() {
		var i FindManyProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Color,
			&i.Count,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findObjectProductMapping = `-- name: FindObjectProductMapping :one
Select id, object_id, object_type, tag_id from object_has_tags
Where tag_id = $1::TEXT and object_id = $2::TEXT
`

type FindObjectProductMappingParams struct {
	ProductID    string
	ObjectID string
}


const findObjectsByProductIdsAndObjectTypeIntersection = `-- name: FindObjectsByProductIdsAndObjectTypeIntersection :many
SELECT result.object_id from (SELECT oht.object_id, count(object_id) FROM object_has_tags oht
WHERE tag_id = ANY ($1::TEXT[]) and object_type = $2::TEXT
GROUP BY object_id) as result where count = $3::INT
`

type FindObjectsByProductIdsAndObjectTypeIntersectionParams struct {
	ProductIds       []string
	ObjectType   string
	ProductIdsLenght int32
}

func (q *Queries) FindObjectsByProductIdsAndObjectTypeIntersection(ctx context.Context, arg FindObjectsByProductIdsAndObjectTypeIntersectionParams) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, findObjectsByProductIdsAndObjectTypeIntersection, pq.Array(arg.ProductIds), arg.ObjectType, arg.ProductIdsLenght)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var object_id sql.NullString
		if err := rows.Scan(&object_id); err != nil {
			return nil, err
		}
		items = append(items, object_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findObjectsByProductIdsAndObjectTypeUnion = `-- name: FindObjectsByProductIdsAndObjectTypeUnion :many
SELECT object_id FROM object_has_tags
WHERE tag_id = ANY ($1::TEXT[]) and object_type = $2::TEXT
GROUP by object_id
`

type FindObjectsByProductIdsAndObjectTypeUnionParams struct {
	ProductIds     []string
	ObjectType string
}

func (q *Queries) FindObjectsByProductIdsAndObjectTypeUnion(ctx context.Context, arg FindObjectsByProductIdsAndObjectTypeUnionParams) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, findObjectsByProductIdsAndObjectTypeUnion, pq.Array(arg.ProductIds), arg.ObjectType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var object_id sql.NullString
		if err := rows.Scan(&object_id); err != nil {
			return nil, err
		}
		items = append(items, object_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findProductById = `-- name: FindProductById :one
SELECT id, name, description, color
FROM tags
WHERE id = $1::TEXT
`

func (q *Queries) FindProductById(ctx context.Context, id string) (Product, error) {
	row := q.db.QueryRowContext(ctx, findProductById, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Color,
	)
	return i, err
}

const findProductByName = `-- name: FindProductByName :one
SELECT id, name, description, color
FROM tags
WHERE name = $1::TEXT
`

func (q *Queries) FindProductByName(ctx context.Context, name string) (Product, error) {
	row := q.db.QueryRowContext(ctx, findProductByName, name)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Color,
	)
	return i, err
}

const findProductsByObjectId = `-- name: FindProductsByObjectId :many
Select tag_id from object_has_tags
Where object_id = $1::TEXT
`

func (q *Queries) FindProductsByObjectId(ctx context.Context, objectID string) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, findProductsByObjectId, objectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var tag_id sql.NullString
		if err := rows.Scan(&tag_id); err != nil {
			return nil, err
		}
		items = append(items, tag_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findProductsByPartialName = `-- name: FindProductsByPartialName :many
SELECT id, name, description, color
From tags
WHERE
        (CASE WHEN $1::text = 'name' THEN tags.name ILIKE $2::TEXT ELSE TRUE END)
ORDER BY
        CASE WHEN $3::boolean AND $4::text = 'name' THEN tags.name END ASC

LIMIT $5::INT
`

type FindProductsByPartialNameParams struct {
	FilterBy  string
	Pattern   string
	Ascending bool
	OrderBy   string
	PageLimit int32
}

func (q *Queries) FindProductsByPartialName(ctx context.Context, arg FindProductsByPartialNameParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, findProductsByPartialName,
		arg.FilterBy,
		arg.Pattern,
		arg.Ascending,
		arg.OrderBy,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Color,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}


const getProductsByIDs = `-- name: GetProductsByIDs :many
SELECT id, name, description, color
FROM tags
WHERE id = ANY ($1::TEXT[])
`

func (q *Queries) GetProductsByIDs(ctx context.Context, tagIds []string) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByIDs, pq.Array(tagIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Color,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductsByObjectIDs = `-- name: GetProductsByObjectIDs :many
SELECT t.id, t.name, t.description, t.color, oht.object_id
FROM tags t
JOIN object_has_tags oht on t.id = oht.tag_id
WHERE oht.object_id = ANY ($1::TEXT[])
`

type GetProductsByObjectIDsRow struct {
	ID          string
	Name        sql.NullString
	Description sql.NullString
	Color       sql.NullString
	ObjectID    sql.NullString
}

func (q *Queries) GetProductsByObjectIDs(ctx context.Context, objectIds []string) ([]GetProductsByObjectIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByObjectIDs, pq.Array(objectIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsByObjectIDsRow
	for rows.Next() {
		var i GetProductsByObjectIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Color,
			&i.ObjectID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertObject = `-- name: InsertObject :exec
INSERT INTO object_has_tags(id, object_id, object_type, tag_id)
VALUES ($1::TEXT, $2::TEXT, $3::TEXT, $4::TEXT)
`

type InsertObjectParams struct {
	ID         string
	ObjectID   string
	ObjectType string
	ProductID      string
}

func (q *Queries) InsertObject(ctx context.Context, arg InsertObjectParams) error {
	_, err := q.db.ExecContext(ctx, insertObject,
		arg.ID,
		arg.ObjectID,
		arg.ObjectType,
		arg.ProductID,
	)
	return err
}

const insertProduct = `-- name: InsertProduct :exec
INSERT INTO tags(id, name, description, color)
VALUES ($1::TEXT, $2::TEXT, $3::TEXT, $4::TEXT)
`

type InsertProductParams struct {
	ID          string
	Name        string
	Description string
	Color       string
}

func (q *Queries) InsertProduct(ctx context.Context, arg InsertProductParams) error {
	_, err := q.db.ExecContext(ctx, insertProduct,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Color,
	)
	return err
}

const updateProduct = `-- name: UpdateProduct :exec
UPDATE tags
SET name             = $1::TEXT,
    description      = $2::TEXT,
    color            = $3::TEXT

WHERE id = $4::text
`

type UpdateProductParams struct {
	Name        string
	Description string
	Color       string
	ID          string
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.ExecContext(ctx, updateProduct,
		arg.Name,
		arg.Description,
		arg.Color,
		arg.ID,
	)
	return err
}
