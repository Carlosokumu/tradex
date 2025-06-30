package controllers

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/carlosokumu/dubbedapi/config"
	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/dtos"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCommunity will create a new trading community
func CreateCommunity(ctx *gin.Context) {
	var communityDto dtos.CommunityDto

	if err := ctx.ShouldBindJSON(&communityDto); err != nil {
		log.Printf("Invalid request payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request payload",
			"details": err.Error(),
		})
		return
	}

	if strings.TrimSpace(communityDto.Name) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "community name cannot be empty",
		})
		return
	}

	if exists, err := checkCommunityExists(communityDto.Name); err != nil {
		log.Printf("Database error checking community existence: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	} else if exists {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "community with this name already exists",
		})
		return
	}

	tx := database.Instance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var userModel models.UserModel
	if err := tx.Where("user_name = ?", communityDto.Creator).First(&userModel).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "creator user not found",
			})
			return
		}
		log.Printf("Database error fetching user %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	community := &models.Community{
		Name:        communityDto.Name,
		Description: communityDto.Description,
		AdminID:     userModel.ID,
		Admin:       userModel,
		Members:     []models.UserModel{userModel},
	}

	if err := tx.Create(community).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error creating community: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create community",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Database transaction commit failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create community",
		})
		return
	}

	type CommunityResponse struct {
		ID          uint             `json:"id"`
		Name        string           `json:"name"`
		Description string           `json:"description"`
		AdminID     uint             `json:"admin_id"`
		Admin       models.UserModel `json:"admin"`
		MemberCount int              `json:"member_count"`
		CreatedAt   time.Time        `json:"created_at"`
	}

	response := CommunityResponse{
		ID:          community.ID,
		Name:        community.Name,
		Description: community.Description,
		AdminID:     community.AdminID,
		Admin:       community.Admin,
		MemberCount: len(community.Members),
		CreatedAt:   community.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": response,
	})
}

// GetAllCommunities will retrieve all tradeshare communities with each page size set at 10.
func GetAllCommunities(c *gin.Context) {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		maxPageSize     = 100
	)

	page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(defaultPage)))
	if err != nil || page < 1 {
		page = defaultPage
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(defaultPageSize)))
	if err != nil || pageSize < 1 || pageSize > maxPageSize {
		pageSize = defaultPageSize
	}

	offset := (page - 1) * pageSize

	var totalRows int64
	if err := database.Instance.Model(&models.Community{}).Count(&totalRows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count communities",
		})
		return
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))

	// Fetch paginated communities
	var communities []models.Community
	if err := database.Instance.
		Preload("Admin", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_name", "avatar")
		}).
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&communities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch communities",
		})
		return
	}

	type CommunityResponse struct {
		ID          uint      `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		Admin       struct {
			ID       uint   `json:"id"`
			UserName string `json:"user_name"`
			Avatar   string `json:"avatar"`
		} `json:"admin"`
		MemberCount int `json:"member_count"`
	}

	response := make([]CommunityResponse, len(communities))
	for i, comm := range communities {
		response[i] = CommunityResponse{
			ID:          comm.ID,
			Name:        comm.Name,
			Description: comm.Description,
			CreatedAt:   comm.CreatedAt,
			Admin: struct {
				ID       uint   `json:"id"`
				UserName string `json:"user_name"`
				Avatar   string `json:"avatar"`
			}{
				ID:       comm.Admin.ID,
				UserName: comm.Admin.UserName,
				Avatar:   comm.Admin.Avatar,
			},
			MemberCount: len(comm.Members),
		}
	}

	pagination := gin.H{
		"current_page":  page,
		"page_size":     pageSize,
		"total_pages":   totalPages,
		"total_records": totalRows,
		"has_previous":  page > 1,
		"has_next":      page < totalPages,
	}

	if page > 1 {
		pagination["previous_page"] = page - 1
	}

	if page < totalPages {
		pagination["next_page"] = page + 1
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       response,
		"pagination": pagination,
	})
}

// PostToCommunity  allow sharing  posts to communities
func PostToCommunity(ctx *gin.Context) {
	var pstToCommunityDto dtos.PostToCommunityDto

	err := ctx.ShouldBind(&pstToCommunityDto)
	if err != nil {
		log.Printf("Invalid request payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request payload",
			"details": err.Error(),
		})
		return
	}

	cld, err := cloudinary.NewFromParams(
		config.CloudinaryCloudName,
		config.CloudinaryAPIKey,
		config.CloudinaryAPISecret,
	)

	if err != nil {
		log.Printf("Failed to initialize cloudinary: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to initialize Cloudinary"})
		return
	}

	tx := database.Instance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.UserModel
	if err := tx.Where("user_name = ?", pstToCommunityDto.Username).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check user",
		})
		return
	}

	var community models.Community
	if err := tx.Where("name = ?", pstToCommunityDto.CommunityName).First(&community).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "community not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check community",
		})
		return
	}

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 2<<20) // 2MB

	file, _, err := ctx.Request.FormFile("image")
	if err != nil {
		log.Printf("Error processing image: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file or size > 2MB"})
		return
	}
	defer file.Close()

	uploadResult, err := cld.Upload.Upload(
		ctx.Request.Context(),
		file,
		uploader.UploadParams{
			PublicID: fmt.Sprintf("%s_%d", community.Name, time.Now().Unix()),
			Folder:   "community_uploads",
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Upload failed: %v", err)})
		return
	}

	post := models.Post{
		Content:     pstToCommunityDto.Content,
		ImageURL:    uploadResult.SecureURL,
		PosterID:    user.ID,
		CommunityID: community.ID,
	}
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post: " + err.Error()})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Post created successfully",
		"post_id": post.ID,
		"image":   uploadResult.SecureURL,
	})
}

// AddNewMemberToCommunity adds a new member to a  community
func AddNewMemberToCommunity(ctx *gin.Context) {
	var joinCommunityDto dtos.JoinCommunityDto
	if err := ctx.ShouldBindJSON(&joinCommunityDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request payload",
			"details": err.Error(),
		})
		return
	}

	tx := database.Instance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var community models.Community
	if err := tx.Where("name = ?", joinCommunityDto.CommunityName).First(&community).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "community not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check community",
		})
		return
	}

	var user models.UserModel
	if err := tx.Where("user_name = ?", joinCommunityDto.Username).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check user",
		})
		return
	}

	existingMembership := tx.Model(&community).
		Where("user_communities.user_model_id = ?", user.ID).
		Association("Members").
		Count()

	if existingMembership > 0 {
		tx.Rollback()
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "user is already a member of this community",
		})
		return
	}

	if err := tx.Model(&community).Association("Members").Append(&user); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add member to community",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to complete operation",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user successfully added to community",
		"data": gin.H{
			"community_id": community.ID,
			"user_id":      user.ID,
		},
	})
}

// GetCommunityByName will retrieve a specific community by name
func GetCommunityByName(c *gin.Context) {
	communityName := strings.TrimSpace(c.Query("community_name"))
	if communityName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "community_name parameter is required",
		})
		return
	}

	var community models.Community
	err := database.Instance.
		Select("id", "name", "description", "admin_id", "created_at").
		Where("name = ?", communityName).
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_name", "avatar")
		}).
		Preload("Posts", func(db *gorm.DB) *gorm.DB {
			return db.Select(
				"id",
				"content",
				"poster_id",
				"community_id",
				"created_at",
				"image_url",
			).Order("created_at DESC").Limit(20)
		}).First(&community).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Community not found with name: %s", communityName)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Community not found",
			})
			return
		}

		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	type CommunityResponse struct {
		ID          uint               `json:"id"`
		Name        string             `json:"name"`
		Description string             `json:"description"`
		MemberCount int                `json:"member_count"`
		Members     []models.UserModel `json:"members"`
		Posts       []models.Post      `json:"posts"`
	}

	response := CommunityResponse{
		ID:          community.ID,
		Name:        community.Name,
		Description: community.Description,
		MemberCount: len(community.Members),
		Members:     community.Members,
		Posts:       community.Posts,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func checkCommunityExists(name string) (bool, error) {
	var count int64
	err := database.Instance.Model(&models.Community{}).
		Where("LOWER(name) = LOWER(?)", name).
		Count(&count).
		Error
	return count > 0, err
}
