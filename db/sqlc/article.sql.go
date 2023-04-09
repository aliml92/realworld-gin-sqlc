// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: article.sql

package db

import (
	"context"
	"time"
)

const addComment = `-- name: AddComment :one
INSERT INTO comments (
    id,
    author_id,
    article_id,
    body
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, author_id, article_id, body, created_at, updated_at
`

type AddCommentParams struct {
	ID        string `json:"id"`
	AuthorID  string `json:"author_id"`
	ArticleID string `json:"article_id"`
	Body      string `json:"body"`
}

func (q *Queries) AddComment(ctx context.Context, arg AddCommentParams) (*Comment, error) {
	row := q.db.QueryRow(ctx, addComment,
		arg.ID,
		arg.AuthorID,
		arg.ArticleID,
		arg.Body,
	)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.ArticleID,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const countArticles = `-- name: CountArticles :one
SELECT count(*)
FROM articles
`

func (q *Queries) CountArticles(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countArticles)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countArticlesByAuthor = `-- name: CountArticlesByAuthor :one
SELECT count(*)
FROM articles a
LEFT JOIN users u ON a.author_id = u.id
WHERE u.username = $1
`

func (q *Queries) CountArticlesByAuthor(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRow(ctx, countArticlesByAuthor, username)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countArticlesByFavorited = `-- name: CountArticlesByFavorited :one
SELECT count(*)
FROM articles a
LEFT JOIN favorites fav ON a.id = fav.article_id
LEFT JOIN users u2 ON fav.user_id = u2.id
WHERE u2.username = $1
`

func (q *Queries) CountArticlesByFavorited(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRow(ctx, countArticlesByFavorited, username)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countArticlesByFollowing = `-- name: CountArticlesByFollowing :one
SELECT count(*)
FROM articles a
LEFT JOIN users u ON a.author_id = u.id
LEFT JOIN follows f2 ON u.id = f2.followee_id
LEFT JOIN users u2 ON f2.follower_id = u2.id
WHERE u2.username = $1
`

func (q *Queries) CountArticlesByFollowing(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRow(ctx, countArticlesByFollowing, username)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countArticlesByTag = `-- name: CountArticlesByTag :one
SELECT count(*)
FROM articles a
LEFT JOIN article_tags art ON a.id = art.article_id
LEFT JOIN tags t ON art.tag_id = t.id
WHERE t.name = $1
`

func (q *Queries) CountArticlesByTag(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRow(ctx, countArticlesByTag, name)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createArticle = `-- name: CreateArticle :one
INSERT INTO articles (
    id,
    author_id,
    slug,
    title,
    description,
    body
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, author_id, slug, title, description, body, created_at, updated_at
`

type CreateArticleParams struct {
	ID          string `json:"id"`
	AuthorID    string `json:"author_id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (*Article, error) {
	row := q.db.QueryRow(ctx, createArticle,
		arg.ID,
		arg.AuthorID,
		arg.Slug,
		arg.Title,
		arg.Description,
		arg.Body,
	)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const createArticleTag = `-- name: CreateArticleTag :one
INSERT INTO article_tags (
    article_id,
    tag_id
) VALUES (
    $1,
    $2
)
RETURNING article_id, tag_id
`

type CreateArticleTagParams struct {
	ArticleID string `json:"article_id"`
	TagID     string `json:"tag_id"`
}

func (q *Queries) CreateArticleTag(ctx context.Context, arg CreateArticleTagParams) (*ArticleTag, error) {
	row := q.db.QueryRow(ctx, createArticleTag, arg.ArticleID, arg.TagID)
	var i ArticleTag
	err := row.Scan(&i.ArticleID, &i.TagID)
	return &i, err
}

const createTag = `-- name: CreateTag :one
INSERT INTO tags (
    id,
    name
) VALUES (
    $1,
    $2
)
ON CONFLICT 
    ON CONSTRAINT tags_name_key
DO NOTHING
RETURNING id
`

type CreateTagParams struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) CreateTag(ctx context.Context, arg CreateTagParams) (string, error) {
	row := q.db.QueryRow(ctx, createTag, arg.ID, arg.Name)
	var id string
	err := row.Scan(&id)
	return id, err
}

const deleteArticle = `-- name: DeleteArticle :exec
DELETE FROM articles
WHERE slug = $1 and author_id = $2
`

type DeleteArticleParams struct {
	Slug     string `json:"slug"`
	AuthorID string `json:"author_id"`
}

func (q *Queries) DeleteArticle(ctx context.Context, arg DeleteArticleParams) error {
	_, err := q.db.Exec(ctx, deleteArticle, arg.Slug, arg.AuthorID)
	return err
}

const deleteComment = `-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1
`

func (q *Queries) DeleteComment(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteComment, id)
	return err
}

const doesFavoriteExist = `-- name: DoesFavoriteExist :one
SELECT EXISTS (
  SELECT 1
  FROM favorites
  WHERE user_id = $1 AND article_id = $2
)
`

type DoesFavoriteExistParams struct {
	UserID    string `json:"user_id"`
	ArticleID string `json:"article_id"`
}

func (q *Queries) DoesFavoriteExist(ctx context.Context, arg DoesFavoriteExistParams) (bool, error) {
	row := q.db.QueryRow(ctx, doesFavoriteExist, arg.UserID, arg.ArticleID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const favoriteArticle = `-- name: FavoriteArticle :exec
INSERT INTO favorites (user_id, article_id)
VALUES ($1, $2)
`

type FavoriteArticleParams struct {
	UserID    string `json:"user_id"`
	ArticleID string `json:"article_id"`
}

func (q *Queries) FavoriteArticle(ctx context.Context, arg FavoriteArticleParams) error {
	_, err := q.db.Exec(ctx, favoriteArticle, arg.UserID, arg.ArticleID)
	return err
}

const getArticleAuthorIDBySlug = `-- name: GetArticleAuthorIDBySlug :one
SELECT id, author_id
FROM articles
WHERE slug = $1
`

type GetArticleAuthorIDBySlugRow struct {
	ID       string `json:"id"`
	AuthorID string `json:"author_id"`
}

func (q *Queries) GetArticleAuthorIDBySlug(ctx context.Context, slug string) (*GetArticleAuthorIDBySlugRow, error) {
	row := q.db.QueryRow(ctx, getArticleAuthorIDBySlug, slug)
	var i GetArticleAuthorIDBySlugRow
	err := row.Scan(&i.ID, &i.AuthorID)
	return &i, err
}

const getArticleBySlug = `-- name: GetArticleBySlug :one
SELECT a.id,
       a.slug,
       a.title,
       a.description,
       a.body,
       array_agg(t.name) filter (where t.name is not null) AS tag_list,
       a.created_at,
       a.updated_at,
       count(distinct f.user_id) as favorites_count,
       u.id as author_id,
       u.username,
       u.bio,
       u.image
FROM articles a
LEFT  JOIN article_tags art ON a.id = art.article_id
LEFT JOIN tags t ON art.tag_id = t.id
LEFT  JOIN favorites f ON a.id = f.article_id
LEFT  JOIN users u ON a.author_id = u.id
WHERE a.slug = $1
GROUP BY  a.id, a.slug, a.title, a.description, a.body, 
          a.created_at, a.updated_at, u.id
`

type GetArticleBySlugRow struct {
	ID             string      `json:"id"`
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        interface{} `json:"tag_list"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	FavoritesCount int64       `json:"favorites_count"`
	AuthorID       *string     `json:"author_id"`
	Username       *string     `json:"username"`
	Bio            *string     `json:"bio"`
	Image          *string     `json:"image"`
}

func (q *Queries) GetArticleBySlug(ctx context.Context, slug string) (*GetArticleBySlugRow, error) {
	row := q.db.QueryRow(ctx, getArticleBySlug, slug)
	var i GetArticleBySlugRow
	err := row.Scan(
		&i.ID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.Body,
		&i.TagList,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FavoritesCount,
		&i.AuthorID,
		&i.Username,
		&i.Bio,
		&i.Image,
	)
	return &i, err
}

const getArticleIDBySlug = `-- name: GetArticleIDBySlug :one
SELECT id
FROM articles
WHERE slug = $1
`

func (q *Queries) GetArticleIDBySlug(ctx context.Context, slug string) (string, error) {
	row := q.db.QueryRow(ctx, getArticleIDBySlug, slug)
	var id string
	err := row.Scan(&id)
	return id, err
}

const getCommentAuthorID = `-- name: GetCommentAuthorID :one
SELECT c.author_id
FROM comments c
WHERE c.id = $1
LIMIT 1
`

func (q *Queries) GetCommentAuthorID(ctx context.Context, id string) (string, error) {
	row := q.db.QueryRow(ctx, getCommentAuthorID, id)
	var author_id string
	err := row.Scan(&author_id)
	return author_id, err
}

const getCommentsBySlug = `-- name: GetCommentsBySlug :many
SELECT 
    c.id,
    c.body,
    c.created_at,
    c.updated_at,
    u.id as author_id,
    u.username,
    u.bio,
    u.image     
FROM (
    SELECT id
    FROM articles
    WHERE slug = $1
) a
LEFT JOIN comments c ON a.id = c.article_id 
LEFT JOIN users u ON c.author_id = u.id
`

type GetCommentsBySlugRow struct {
	ID        *string    `json:"id"`
	Body      *string    `json:"body"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	AuthorID  *string    `json:"author_id"`
	Username  *string    `json:"username"`
	Bio       *string    `json:"bio"`
	Image     *string    `json:"image"`
}

func (q *Queries) GetCommentsBySlug(ctx context.Context, slug string) ([]*GetCommentsBySlugRow, error) {
	rows, err := q.db.Query(ctx, getCommentsBySlug, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetCommentsBySlugRow
	for rows.Next() {
		var i GetCommentsBySlugRow
		if err := rows.Scan(
			&i.ID,
			&i.Body,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AuthorID,
			&i.Username,
			&i.Bio,
			&i.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTags = `-- name: GetTags :many
SELECT name
FROM tags
`

func (q *Queries) GetTags(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const isFollowingList = `-- name: IsFollowingList :one
SELECT ARRAY(
  SELECT EXISTS (
    SELECT 1 FROM follows
    WHERE follower_id = $1 AND following_id = id
  )
  FROM unnest($2::text[]) AS id
)::bool[]
`

type IsFollowingListParams struct {
	FollowerID  string   `json:"follower_id"`
	FollowingID []string `json:"following_id"`
}

func (q *Queries) IsFollowingList(ctx context.Context, arg IsFollowingListParams) ([]bool, error) {
	row := q.db.QueryRow(ctx, isFollowingList, arg.FollowerID, arg.FollowingID)
	var column_1 []bool
	err := row.Scan(&column_1)
	return column_1, err
}

const listArticles = `-- name: ListArticles :many
SELECT a.id,
       a.slug,
       a.title,
       a.description,
       a.body,
       array_agg(t.name) AS tag_list,
       a.created_at,
       a.updated_at,
       coalesce(count(f.article_id), 0)::int as favorites_count, 
       u.username,
       u.bio,
       u.image
FROM articles a
LEFT JOIN article_tags art ON a.id = art.article_id
LEFT JOIN tags t ON art.tag_id = t.id
LEFT JOIN (
  SELECT   article_id
  FROM     favorites
  GROUP BY article_id
) f ON a.id = f.article_id
LEFT JOIN users u ON a.author_id = u.id
GROUP BY  a.id, a.slug, a.title, a.description, a.body, 
          a.created_at, a.updated_at, u.id
ORDER BY a.created_at DESC
LIMIT $1 OFFSET $2
`

type ListArticlesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListArticlesRow struct {
	ID             string      `json:"id"`
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        interface{} `json:"tag_list"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	FavoritesCount int32       `json:"favorites_count"`
	Username       *string     `json:"username"`
	Bio            *string     `json:"bio"`
	Image          *string     `json:"image"`
}

func (q *Queries) ListArticles(ctx context.Context, arg ListArticlesParams) ([]*ListArticlesRow, error) {
	rows, err := q.db.Query(ctx, listArticles, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListArticlesRow
	for rows.Next() {
		var i ListArticlesRow
		if err := rows.Scan(
			&i.ID,
			&i.Slug,
			&i.Title,
			&i.Description,
			&i.Body,
			&i.TagList,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FavoritesCount,
			&i.Username,
			&i.Bio,
			&i.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listArticlesByAuthor = `-- name: ListArticlesByAuthor :many
SELECT a.id,
       a.slug,
       a.title,
       a.description,
       a.body,
       array_agg(t.name) AS tag_list,
       a.created_at,
       a.updated_at,
       coalesce(count(f.article_id), 0)::int as favorites_count, 
       u.username,
       u.bio,
       u.image
FROM articles a
LEFT JOIN article_tags art ON a.id = art.article_id
LEFT JOIN tags t ON art.tag_id = t.id
LEFT JOIN (
  SELECT   article_id
  FROM     favorites
  GROUP BY article_id
) f ON a.id = f.article_id
LEFT JOIN users u ON a.author_id = u.id
WHERE u.username = $1
GROUP BY  a.id, a.slug, a.title, a.description, a.body, 
          a.created_at, a.updated_at, u.id
ORDER BY a.created_at DESC
LIMIT $2 OFFSET $3
`

type ListArticlesByAuthorParams struct {
	Username string `json:"username"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type ListArticlesByAuthorRow struct {
	ID             string      `json:"id"`
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        interface{} `json:"tag_list"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	FavoritesCount int32       `json:"favorites_count"`
	Username       *string     `json:"username"`
	Bio            *string     `json:"bio"`
	Image          *string     `json:"image"`
}

func (q *Queries) ListArticlesByAuthor(ctx context.Context, arg ListArticlesByAuthorParams) ([]*ListArticlesByAuthorRow, error) {
	rows, err := q.db.Query(ctx, listArticlesByAuthor, arg.Username, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListArticlesByAuthorRow
	for rows.Next() {
		var i ListArticlesByAuthorRow
		if err := rows.Scan(
			&i.ID,
			&i.Slug,
			&i.Title,
			&i.Description,
			&i.Body,
			&i.TagList,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FavoritesCount,
			&i.Username,
			&i.Bio,
			&i.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listArticlesByFavorited = `-- name: ListArticlesByFavorited :many
SELECT a.id,
       a.slug,
       a.title,
       a.description,
       a.body,
       array_agg(t.name) AS tag_list,
       a.created_at,
       a.updated_at,
       coalesce(count(f.article_id), 0)::int as favorites_count, 
       u.username,
       u.bio,
       u.image
FROM articles a
LEFT JOIN article_tags art ON a.id = art.article_id
LEFT JOIN tags t ON art.tag_id = t.id
LEFT JOIN (
  SELECT   article_id
  FROM     favorites
  GROUP BY article_id
) f ON a.id = f.article_id
LEFT JOIN users u ON a.author_id = u.id
LEFT JOIN favorites fav ON a.id = fav.article_id
LEFT JOIN users u2 ON fav.user_id = u2.id
WHERE u2.username = $1
GROUP BY  a.id, a.slug, a.title, a.description, a.body, 
          a.created_at, a.updated_at, u.id
ORDER BY a.created_at DESC
LIMIT $2 OFFSET $3
`

type ListArticlesByFavoritedParams struct {
	Username string `json:"username"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type ListArticlesByFavoritedRow struct {
	ID             string      `json:"id"`
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        interface{} `json:"tag_list"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	FavoritesCount int32       `json:"favorites_count"`
	Username       *string     `json:"username"`
	Bio            *string     `json:"bio"`
	Image          *string     `json:"image"`
}

func (q *Queries) ListArticlesByFavorited(ctx context.Context, arg ListArticlesByFavoritedParams) ([]*ListArticlesByFavoritedRow, error) {
	rows, err := q.db.Query(ctx, listArticlesByFavorited, arg.Username, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListArticlesByFavoritedRow
	for rows.Next() {
		var i ListArticlesByFavoritedRow
		if err := rows.Scan(
			&i.ID,
			&i.Slug,
			&i.Title,
			&i.Description,
			&i.Body,
			&i.TagList,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FavoritesCount,
			&i.Username,
			&i.Bio,
			&i.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listArticlesByFollowing = `-- name: ListArticlesByFollowing :many
SELECT a.id,
       a.slug,
       a.title,
       a.description,
       a.body,
       array_agg(t.name) AS tag_list,
       a.created_at,
       a.updated_at,
       coalesce(count(f.article_id), 0)::int as favorites_count, 
       u.username,
       u.bio,
       u.image
FROM articles a
LEFT JOIN article_tags art ON a.id = art.article_id
LEFT JOIN tags t ON art.tag_id = t.id
LEFT JOIN (
  SELECT   article_id
  FROM     favorites
  GROUP BY article_id
) f ON a.id = f.article_id
LEFT JOIN users u ON a.author_id = u.id
LEFT JOIN follows f2 ON u.id = f2.followee_id
LEFT JOIN users u2 ON f2.follower_id = u2.id
WHERE u2.username = $1
GROUP BY  a.id, a.slug, a.title, a.description, a.body, 
          a.created_at, a.updated_at, u.id
ORDER BY a.created_at DESC
LIMIT $2 OFFSET $3
`

type ListArticlesByFollowingParams struct {
	Username string `json:"username"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type ListArticlesByFollowingRow struct {
	ID             string      `json:"id"`
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        interface{} `json:"tag_list"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	FavoritesCount int32       `json:"favorites_count"`
	Username       *string     `json:"username"`
	Bio            *string     `json:"bio"`
	Image          *string     `json:"image"`
}

func (q *Queries) ListArticlesByFollowing(ctx context.Context, arg ListArticlesByFollowingParams) ([]*ListArticlesByFollowingRow, error) {
	rows, err := q.db.Query(ctx, listArticlesByFollowing, arg.Username, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListArticlesByFollowingRow
	for rows.Next() {
		var i ListArticlesByFollowingRow
		if err := rows.Scan(
			&i.ID,
			&i.Slug,
			&i.Title,
			&i.Description,
			&i.Body,
			&i.TagList,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FavoritesCount,
			&i.Username,
			&i.Bio,
			&i.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listArticlesByTag = `-- name: ListArticlesByTag :many
SELECT a.id,
       a.slug,
       a.title,
       a.description,
       a.body,
       array_agg(t.name) AS tag_list,
       a.created_at,
       a.updated_at,
       coalesce(count(f.article_id), 0)::int as favorites_count, 
       u.username,
       u.bio,
       u.image
FROM articles a
LEFT JOIN article_tags art ON a.id = art.article_id
LEFT JOIN tags t ON art.tag_id = t.id
LEFT JOIN (
  SELECT   article_id
  FROM     favorites
  GROUP BY article_id
) f ON a.id = f.article_id
LEFT JOIN users u ON a.author_id = u.id
WHERE t.name = $1
GROUP BY  a.id, a.slug, a.title, a.description, a.body, 
          a.created_at, a.updated_at, u.id
ORDER BY a.created_at DESC
LIMIT $2 OFFSET $3
`

type ListArticlesByTagParams struct {
	Name   string `json:"name"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ListArticlesByTagRow struct {
	ID             string      `json:"id"`
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        interface{} `json:"tag_list"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	FavoritesCount int32       `json:"favorites_count"`
	Username       *string     `json:"username"`
	Bio            *string     `json:"bio"`
	Image          *string     `json:"image"`
}

func (q *Queries) ListArticlesByTag(ctx context.Context, arg ListArticlesByTagParams) ([]*ListArticlesByTagRow, error) {
	rows, err := q.db.Query(ctx, listArticlesByTag, arg.Name, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListArticlesByTagRow
	for rows.Next() {
		var i ListArticlesByTagRow
		if err := rows.Scan(
			&i.ID,
			&i.Slug,
			&i.Title,
			&i.Description,
			&i.Body,
			&i.TagList,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FavoritesCount,
			&i.Username,
			&i.Bio,
			&i.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unfavoriteArticle = `-- name: UnfavoriteArticle :exec
DELETE FROM favorites
WHERE user_id = $1 AND article_id = $2
`

type UnfavoriteArticleParams struct {
	UserID    string `json:"user_id"`
	ArticleID string `json:"article_id"`
}

func (q *Queries) UnfavoriteArticle(ctx context.Context, arg UnfavoriteArticleParams) error {
	_, err := q.db.Exec(ctx, unfavoriteArticle, arg.UserID, arg.ArticleID)
	return err
}

const updateArticle = `-- name: UpdateArticle :one
UPDATE articles
SET slug = coalesce($1, slug),
    title = coalesce($2, title),
    description = coalesce($3, description),
    body = coalesce($4, body),
    updated_at = now()
WHERE id = $5 and author_id = $6
RETURNING id, author_id, slug, title, description, body, created_at, updated_at
`

type UpdateArticleParams struct {
	Slug        *string `json:"slug"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Body        *string `json:"body"`
	ID          string  `json:"id"`
	AuthorID    string  `json:"author_id"`
}

func (q *Queries) UpdateArticle(ctx context.Context, arg UpdateArticleParams) (*Article, error) {
	row := q.db.QueryRow(ctx, updateArticle,
		arg.Slug,
		arg.Title,
		arg.Description,
		arg.Body,
		arg.ID,
		arg.AuthorID,
	)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
