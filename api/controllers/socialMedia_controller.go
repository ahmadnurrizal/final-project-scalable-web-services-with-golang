package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/auth"
	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/models"
	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/utils/formaterror"
	"github.com/gin-gonic/gin"
)

// CreateSocialMedia godoc
// @Summary     Create Social Media
// @Description Add a new Social Media
// @Tags        Social Media
// @Accept      json
// @Produce     json
// @Param       CreateSocialMedia body models.CreateSocialMedia true "SocialMedia Data"
// @Success     200  {object} models.SocialMedia
// @Router      /social-media [post]
func (server *Server) CreateSocialMedia(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	socialMedia := models.SocialMedia{}

	err = json.Unmarshal(body, &socialMedia)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// check if the user exist:
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	socialMedia.UserID = uid //the authenticated user is the one creating the socialMedia

	socialMedia.Prepare()
	errorMessages := socialMedia.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	socialMediaCreated, err := socialMedia.SaveSocialMedia(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": socialMediaCreated,
	})
}

// GetSocialMediaAll godoc
// @Summary Get All Social Media
// @Description Retrieve all social media
// @Tags Social Media
// @Accept json
// @Produce json
// @Success 200 {array} models.SocialMedia
// @Router /social-media-all [get]
func (server *Server) GetSocialMediaAll(c *gin.Context) {

	socialMedia := models.SocialMedia{}

	socialMedias, err := socialMedia.FindAllSocialMedia(server.DB)
	if err != nil {
		errList["No_socialMedia"] = "No Social Media Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": socialMedias,
	})
}

// GetSocialMediaByID godoc
// @Summary Get Social Media by ID
// @Description Retrieve a social media by ID
// @Tags Social Media
// @Accept json
// @Produce json
// @Success 200 {object} models.SocialMedia
// @Router /social-media/{id} [get]
func (server *Server) GetSocialMedia(c *gin.Context) {

	socialMediaID := c.Param("id")
	pid, err := strconv.ParseUint(socialMediaID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	socialMedia := models.SocialMedia{}

	socialMediaReceived, err := socialMedia.FindSocialMediaByID(server.DB, pid)
	if err != nil {
		errList["No_socialMedia"] = "No Social Media Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": socialMediaReceived,
	})
}

// UpdateSocialMediaByID godoc
// @Summary Update Social Media by ID
// @Description Update a social Media by ID
// @Tags Social Media
// @Accept json
// @Produce json
// @Param id path int true "SocialMedia ID"
// @Param UpdateSocialMedia body models.UpdateSocialMedia true "SocialMedia Data"
// @Success 200 {object} models.SocialMedia
// @Router /social-media/{id} [put]
func (server *Server) UpdateSocialMedia(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	socialMediaID := c.Param("id")
	// Check if the socialMedia id is valid
	pid, err := strconv.ParseUint(socialMediaID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	//Check if the socialMedia exist
	origSocialMedia := models.SocialMedia{}
	err = server.DB.Debug().Model(models.SocialMedia{}).Where("id = ?", pid).Take(&origSocialMedia).Error
	if err != nil {
		errList["No_socialMedia"] = "No Social Media Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	if uid != origSocialMedia.UserID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Read the data socialMediaed
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Start processing the request data
	socialMedia := models.SocialMedia{}
	err = json.Unmarshal(body, &socialMedia)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	socialMedia.ID = origSocialMedia.ID //this is important to tell the model the socialMedia id to update, the other update field are set above
	socialMedia.UserID = origSocialMedia.UserID

	socialMedia.Prepare()
	errorMessages := socialMedia.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	socialMediaUpdated, err := socialMedia.UpdateASocialMedia(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": socialMediaUpdated,
	})
}

// DeleteSocialMediaByID godoc
// @Summary Delete social media by ID
// @Description Delete a social media by ID
// @Tags Social Media
// @Accept json
// @Produce json
// @Param id path int true "SocialMedia ID"
// @Success 200 {string} string "Social Media deleted"
// @Router /social-media/{id} [delete]
func (server *Server) DeleteSocialMedia(c *gin.Context) {

	socialMediaID := c.Param("id")
	// Is a valid socialMedia id given to us?
	pid, err := strconv.ParseUint(socialMediaID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	fmt.Println("this is delete socialMedia sir")

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Check if the socialMedia exist
	socialMedia := models.SocialMedia{}
	err = server.DB.Debug().Model(models.SocialMedia{}).Where("id = ?", pid).Take(&socialMedia).Error
	if err != nil {
		errList["No_socialMedia"] = "No Social Media Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Is the authenticated user, the owner of this socialMedia?
	if uid != socialMedia.UserID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	_, err = socialMedia.DeleteASocialMedia(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Social Media deleted",
	})
}

func (server *Server) GetUserSocialMedias(c *gin.Context) {

	userID := c.Param("id")
	// Is a valid user id given to us?
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	socialMedia := models.SocialMedia{}
	socialMedias, err := socialMedia.FindUserSocialMedias(server.DB, uint32(uid))
	if err != nil {
		errList["No_socialMedia"] = "No Social Media Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": socialMedias,
	})
}
