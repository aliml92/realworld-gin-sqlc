package api

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/aliml92/realworld-gin-sqlc/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
)

type listQuery struct {
	Tag       string `form:"tag" binding:"omitempty"`
	Author    string `form:"author" binding:"omitempty"`
	Favorited string `form:"favorited" binding:"omitempty"`
	Offset    int    `form:"offset" binding:"omitempty"`
	Limit     int    `form:"limit" binding:"omitempty"`
}

// ListArticles godoc
// @Summary 	List articles
// @Description List articles
// @Tags 		articles
// @Accept  	json
// @Produce  	json
// @Param 		tag 	  query 	string		false 	"Tag"
// @Param 		author    query 	string		false 	"Author"
// @Param 		favorited query		string		false 	"Favorited"
// @Param 		limit	  query 	int   		false 	"Limit"
// @Param 		offset    query 	int    		false 	"Offset"
// @Success 	200 	  {object} 	map[string]interface{}
// @Failure 	500 	  {object} 	Error
// @Router /articles [get]
func (s *Server) ListArticles(c *gin.Context) { // TODO:✅ GET /articles - ListArticles
	var (
		query listQuery
		count int64
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if query.Limit == 0 {
		query.Limit = 20
	}
	if query.Tag != "" {
		p := db.ListArticlesByTagParams{
			Name:   query.Tag,
			Limit:  int32(query.Limit),
			Offset: int32(query.Offset),
		}
		articles, err := NullableList(s.store.ListArticlesByTag(c, p))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if articles == nil {
			c.JSON(http.StatusOK, gin.H{"articles": []string{}, "articlesCount": 0})
			return
		}
		count, err = s.store.CountArticlesByTag(c, query.Tag)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"articles": articles, "articlesCount": count})
		return
	} else if query.Author != "" {
		p := db.ListArticlesByAuthorParams{
			Username: query.Author,
			Limit:    int32(query.Limit),
			Offset:   int32(query.Offset),
		}
		articles, err := NullableList(s.store.ListArticlesByAuthor(c, p))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if articles == nil {
			c.JSON(http.StatusOK, gin.H{"articles": []string{}, "articlesCount": 0})
			return
		}
		count, err = s.store.CountArticlesByAuthor(c, query.Author)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"articles": articles, "articlesCount": count})
		return
	} else if query.Favorited != "" {
		p := db.ListArticlesByFavoritedParams{
			Username: query.Favorited,
			Limit:    int32(query.Limit),
			Offset:   int32(query.Offset),
		}
		articles, err := NullableList(s.store.ListArticlesByFavorited(c, p))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if articles == nil {
			c.JSON(http.StatusOK, gin.H{"articles": []string{}, "articlesCount": 0})
			return
		}
		count, err = s.store.CountArticlesByFavorited(c, query.Favorited)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"articles": articles, "articlesCount": count})
		return
	}
	p := db.ListArticlesParams{
		Limit:  int32(query.Limit),
		Offset: int32(query.Offset),
	}
	articles, err := NullableList(s.store.ListArticles(c, p))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if articles == nil {
		c.JSON(http.StatusOK, gin.H{"articles": []string{}, "articlesCount": 0})
		return
	}
	count, err = s.store.CountArticles(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"articles": articles, "articlesCount": count})
}

// FeedArticles godoc
// @Summary Feed articles
// @Description Feed articles
// @Tags articles
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} Error
// @Security Bearer
// @Router /articles/feed [get]
func (s *Server) FeedArticles(c *gin.Context) { // TODO:✅ GET /articles/feed - FeedArticles
	var (
		query listQuery
		count int64
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if query.Limit == 0 {
		query.Limit = 20
	}
	p := db.ListArticlesByFollowingParams{
		Username: c.GetString("username"),
		Limit:    int32(query.Limit),
		Offset:   int32(query.Offset),
	}
	articles, err := NullableList(s.store.ListArticlesByFollowing(c, p))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if articles == nil {
		c.JSON(http.StatusOK, gin.H{"articles": []string{}, "articlesCount": 0})
		return
	}
	count, err = s.store.CountArticlesByFollowing(c, c.GetString("username"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"articles": articles, "articlesCount": count})
}

type articleResponse struct {
	Article struct {
		Slug           string    `json:"slug"`
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		Body           string    `json:"body"`
		TagList        []string  `json:"tagList"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
		Favorited      bool      `json:"favorited"`
		FavoritesCount int64     `json:"favoritesCount"`
		Author         struct {
			Username  *string `json:"username"`
			Bio       *string `json:"bio"`
			Image     *string `json:"image"`
			Following bool    `json:"following"`
		} `json:"author"`
	} `json:"article"`
}

func newArticleResponse(article *db.GetArticleBySlugRow, favorited, following bool) *articleResponse {
	resp := new(articleResponse)
	resp.Article.Slug = article.Slug
	resp.Article.Title = article.Title
	resp.Article.Description = article.Description
	resp.Article.Body = article.Body
	var tags []string
	if article.TagList != nil {
		a, _ := article.TagList.(pgtype.TextArray)
		a.AssignTo(&tags)
		if tags[0] != "" {
			resp.Article.TagList = tags
		}		
	}
	resp.Article.CreatedAt = article.CreatedAt
	resp.Article.UpdatedAt = article.UpdatedAt
	resp.Article.Favorited = favorited
	resp.Article.FavoritesCount = article.FavoritesCount
	resp.Article.Author.Username = article.Username
	resp.Article.Author.Bio = article.Bio
	resp.Article.Author.Image = article.Image
	resp.Article.Author.Following = following
	return resp
}

// GetArticle godoc
// @Summary Get an article
// @Description Get an article by slug
// @Tags articles
// @Accept  json
// @Produce  json
// @Param slug path string true "Article slug"
// @Success 200 {object} articleResponse
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /articles/{slug} [get]
func (s *Server) GetArticle(c *gin.Context) { // TODO:✅ GET /articles/:slug - GetArticle
	// get token from header if exists        // TODO: GET /articles/:slug - GetArticle - Following
	slug := c.Param("slug")
	article, err := Nullable(s.store.GetArticleBySlug(c, slug))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if article == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}
	c.JSON(http.StatusOK, newArticleResponse(article, false, false))
}

type createArticleReq struct {
	Article struct {
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description" binding:"required"`
		Body        string   `json:"body" binding:"required"`
		TagList     []string `json:"tagList" binding:"omitempty"`
	} `json:"article" binding:"required"`
}

func (req *createArticleReq) bind(c *gin.Context, p *db.CreateArticleTxParams) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}
	p.CreateArticleParams.ID = generateID()
	p.CreateArticleParams.Title = req.Article.Title
	p.CreateArticleParams.Description = req.Article.Description
	p.CreateArticleParams.Body = req.Article.Body
	p.Tags = req.Article.TagList
	return nil
}

// CreateArticle godoc
// @Summary Create article
// @Description Create Article
// @Tags articles
// @Accept json
// @Produce json
// @Param article body createArticleReq true "Article"
// @Success 201 {object} articleResponse
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Security Bearer
// @Router /articles [post]
func (s *Server) CreateArticle(c *gin.Context) { // TODO:✅ POST /articles - CreateArticle
	id := GetIDFromHeader(c)
	var (
		req createArticleReq
		p   db.CreateArticleTxParams
	)
	if err := req.bind(c, &p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, NewValidationError(err))
	}
	p.AuthorID = id
	uniqueSlug, err := s.findUniqueSlug(c, p.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	p.Slug = uniqueSlug
	articleTx, err := s.store.CreateArticleTx(c, p)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.JSON(http.StatusCreated, newArticleTxResponse(articleTx))
}

func newArticleTxResponse(articleTx *db.CreateArticleTxResult) *articleResponse {
	resp := new(articleResponse)
	resp.Article.Slug = articleTx.Article.Slug
	resp.Article.Title = articleTx.Article.Title
	resp.Article.Description = articleTx.Article.Description
	resp.Article.Body = articleTx.Article.Body
	resp.Article.TagList = articleTx.Tags
	resp.Article.CreatedAt = articleTx.Article.CreatedAt
	resp.Article.UpdatedAt = articleTx.Article.UpdatedAt
	resp.Article.Favorited = false
	resp.Article.FavoritesCount = 0
	resp.Article.Author.Username = &articleTx.User.Username
	resp.Article.Author.Bio = articleTx.User.Bio
	resp.Article.Author.Image = articleTx.User.Image
	return resp
}

type updateArticleReq struct {
	Article struct {
		Title       *string `json:"title" binding:"omitempty"`
		Description *string `json:"description" binding:"omitempty"`
		Body        *string `json:"body" binding:"omitempty"`
	} `json:"article" binding:"required"`
}

func (req *updateArticleReq) bind(c *gin.Context, p *db.UpdateArticleParams) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}
	p.Title = req.Article.Title
	p.Description = req.Article.Description
	p.Body = req.Article.Body
	return nil
}

// UpdateArticle godoc
// @Summary Update article
// @Description Update Article
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Param article body updateArticleReq true "Article"
// @Success 200 {object} articleResponse
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Security Bearer
// @Router /articles/{slug} [put]
func (s *Server) UpdateArticle(c *gin.Context) { // TODO:✅ PUT /articles/:slug - UpdateArticle
	authorID := GetIDFromHeader(c)
	slug := c.Param("slug")
	var (
		req updateArticleReq
		p   db.UpdateArticleParams
	)
	if err := req.bind(c, &p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, NewValidationError(err))
	}
	p.AuthorID = authorID
	articleID, err := NullableID(s.store.GetArticleIDBySlug(c, slug))
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	if articleID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}
	p.ID = articleID
	if p.Slug != nil {
		uniqueSlug, err := s.findUniqueSlug(c, *p.Slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewValidationError(err))
			return
		}
		p.Slug = &uniqueSlug
	}
	updatedArticle, err := s.store.UpdateArticle(c, p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	article, err := s.store.GetArticleBySlug(c, updatedArticle.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	c.JSON(http.StatusOK, newArticleResponse(article, false, false))
}

// DeleteArticle godoc
// @Summary Delete article
// @Description Delete Article by slug
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Success 204
// @Failure 401 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Security Bearer
// @Router /articles/{slug} [delete]
func (s *Server) DeleteArticle(c *gin.Context) { // TODO:✅ DELETE /articles/:slug - DeleteArticle
	authorID := GetIDFromHeader(c)
	slug := c.Param("slug")
	p := db.DeleteArticleParams{
		Slug:     slug,
		AuthorID: authorID,
	}
	if err := s.store.DeleteArticle(c, p); err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

type addCommentReq struct {
	Comment struct {
		Body string `json:"body" binding:"required"`
	} `json:"comment" binding:"required"`
}

func (req *addCommentReq) bind(c *gin.Context, p *db.AddCommentParams) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}
	p.ID = generateID()
	p.Body = req.Comment.Body
	return nil
}

// AddComment godoc
// @Summary Add comment
// @Description Add comment to article
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Param comment body addCommentReq true "Comment"
// @Success 200 {object} commentResponse
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Security Bearer
// @Router /articles/{slug}/comments [post]
func (s *Server) AddComment(c *gin.Context) { // TODO:✅ POST /articles/:slug/comments - AddComment
	authorID := GetIDFromHeader(c)
	slug := c.Param("slug")
	var (
		req addCommentReq
		p   db.AddCommentParams
	)
	if err := req.bind(c, &p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, NewValidationError(err))
	}
	p.AuthorID = authorID
	articleID, err := NullableID(s.store.GetArticleIDBySlug(c, slug))
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	if articleID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}
	p.ArticleID = articleID
	comment, err := s.store.AddComment(c, p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	user, err := Nullable(s.store.GetUser(c, comment.AuthorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusCreated, newCommentResponse(comment, user))
}

type commentResponse struct {
	Comment struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Body      string    `json:"body"`
		Author    struct {
			Username  string  `json:"username"`
			Bio       *string `json:"bio"`
			Image     *string `json:"image"`
			Following bool    `json:"following"`
		} `json:"author"`
	} `json:"comment"`
}

func newCommentResponse(comment *db.Comment, user *db.User) *commentResponse {
	res := &commentResponse{}
	res.Comment.ID = comment.ID
	res.Comment.CreatedAt = comment.CreatedAt
	res.Comment.UpdatedAt = comment.UpdatedAt
	res.Comment.Body = comment.Body
	res.Comment.Author.Username = user.Username
	res.Comment.Author.Bio = user.Bio
	res.Comment.Author.Image = user.Image
	return res
}

// GetComments godoc
// @Summary Get comments
// @Description Get comments by article slug
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} commentsResponse
// @Failure 401 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Security Bearer
// @Router /articles/{slug}/comments [get]
func (s *Server) GetComments(c *gin.Context) { // TODO:✅ GET /articles/:slug/comments - GetComments
	var followerID string
	token := GetJWTFromHeader(c)
	if token != "" {
		followerID = GetIDFromToken(token)
	}
	slug := c.Param("slug")
	comments, err := s.store.GetCommentsBySlug(c, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.JSON(http.StatusOK, newCommentsResponse(s, c, followerID, comments))   // UGLY HACK:
}

type commentsResponse struct {
	Comments []struct {
		ID        *string    `json:"id"`
		CreatedAt *time.Time `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
		Body      *string    `json:"body"`
		Author    struct {
			Username  *string `json:"username"`
			Bio       *string `json:"bio"`
			Image     *string `json:"image"`
			Following bool    `json:"following"`
		} `json:"author"`
	} `json:"comments"`
}

func newCommentsResponse(
	s *Server,
	c *gin.Context, 
	followerID string,
	comments []*db.GetCommentsBySlugRow,
	) *commentsResponse {
	res := &commentsResponse{
		Comments: make([]struct {
			ID        *string    `json:"id"`
			CreatedAt *time.Time `json:"createdAt"`
			UpdatedAt *time.Time `json:"updatedAt"`
			Body      *string    `json:"body"`
			Author    struct {
				Username  *string `json:"username"`
				Bio       *string `json:"bio"`
				Image     *string `json:"image"`
				Following bool    `json:"following"`
			} `json:"author"`
		}, len(comments)),
	}
	for i, comment := range comments {
		res.Comments[i].ID = comment.ID
		res.Comments[i].CreatedAt = comment.CreatedAt
		res.Comments[i].UpdatedAt = comment.UpdatedAt
		res.Comments[i].Body = comment.Body
		res.Comments[i].Author.Username = comment.Username
		res.Comments[i].Author.Bio = comment.Bio
		res.Comments[i].Author.Image = comment.Image
		p := db.IsFollowingParams{
			FollowerID: followerID,
			FollowingID: *comment.AuthorID,			
		}
		isFollowing, err := s.store.IsFollowing(c, p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewError(err))
			return nil
		}
		res.Comments[i].Author.Following = isFollowing // TODO:✅ GET /articles/:slug/comments - GetComments, add following
	}

	return res
}

// DeleteComment godoc
// @Summary Delete comment
// @Description Delete comment by id
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Param id path string true "Comment id"
// @Success 200 {object} commentResponse
// @Failure 401 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Security Bearer
// @Router /articles/{slug}/comments/{id} [delete]
func (s *Server) DeleteComment(c *gin.Context) { // TODO:✅ DELETE /articles/:slug/comments/:id - DeleteComment
	// slug := c.Param("slug")
	id := c.Param("id")
	err := s.store.DeleteComment(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}

// FavoriteArticle godoc
// @Summary Favorite article
// @Description Favorite article by slug
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} articleResponse
// @Failure 401 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Security Bearer
// @Router /articles/{slug}/favorite [post]
func (s *Server) FavoriteArticle(c *gin.Context) { // TODO:✅ POST /articles/:slug/favorite - FavoriteArticle
	userID := GetIDFromHeader(c)
	slug := c.Param("slug")
	p := db.FavoriteArticleTxParams{
		UserID: userID,
		Slug:   slug,
	}
	a, err := s.store.FavoriteArticleTx(c, p)
	if err != nil {
		if err == db.ErrArticleNotFound {
			c.JSON(http.StatusNotFound, NewError(err))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.JSON(http.StatusOK, newArticleResponse(a.Article, a.Favorited, a.Following))
}

// UnfavoriteArticle godoc
// @Summary Unfavorite article
// @Description Unfavorite article by slug
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} articleResponse
// @Failure 401 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Security Bearer
// @Router /articles/{slug}/favorite [delete]
func (s *Server) UnfavoriteArticle(c *gin.Context) { // TODO:✅ DELETE /articles/:slug/favorite - UnfavoriteArticle
	userID := GetIDFromHeader(c)
	slug := c.Param("slug")
	p := db.UnfavoriteArticleTxParams{
		UserID: userID,
		Slug:   slug,
	}
	a, err := s.store.UnfavoriteArticleTx(c, p)
	if err != nil {
		if err == db.ErrArticleNotFound {
			c.JSON(http.StatusNotFound, NewError(err))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.JSON(http.StatusOK, newArticleResponse(a.Article, a.Favorited, a.Following))
}

// GetTags godoc
// @Summary Get tags
// @Description Get tags
// @Tags articles
// @Accept json
// @Produce json
// @Success 200 {object} tagsResponse
// @Failure 500 {object} Error
// @Router /tags [get]
func (s *Server) GetTags(c *gin.Context) { // TODO:✅ GET /tags - GetTags
	tags, err := s.store.GetTags(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	c.JSON(http.StatusOK, newTagsResponse(tags))
}

type tagsResponse struct {
	Tags []string `json:"tags"`
}

func newTagsResponse(tags []string) *tagsResponse {
	return &tagsResponse{
		Tags: tags,
	}
}