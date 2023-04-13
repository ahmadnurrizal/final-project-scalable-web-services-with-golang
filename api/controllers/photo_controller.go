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

// CreatePhoto godoc
// @Summary     Create Photo
// @Description Add a new Photo
// @Tags        Photo
// @Accept      json
// @Produce     json
// @Param       CreatePhoto body models.CreatePhoto true "Photo Data"
// @Security ApiKeyAuth
// @Success     200  {object} models.Photo
// @Router      /photos [post]
func (server *Server) CreatePhoto(c *gin.Context) {

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
	photo := models.Photo{}

	err = json.Unmarshal(body, &photo)
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

	photo.UserID = uid //the authenticated user is the one creating the photo

	photo.Prepare()
	errorMessages := photo.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	photoCreated, err := photo.SavePhoto(server.DB)
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
		"response": photoCreated,
	})
}

// GetPhotos godoc
// @Summary Get All Photos
// @Description Retrieve all photos
// @Tags Photo
// @Accept json
// @Produce json
// @Success 200 {array} models.Photo
// @Router /photos [get]
func (server *Server) GetPhotos(c *gin.Context) {

	photo := models.Photo{}

	photos, err := photo.FindAllPhotos(server.DB)
	if err != nil {
		errList["No_photo"] = "No Photo Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": photos,
	})
}

// GetPhotoByID godoc
// @Summary Get Photo by ID
// @Description Retrieve a photo by ID
// @Tags Photo
// @Accept json
// @Produce json
// @Param id path int true "Photo ID"
// @Success 200 {object} models.Photo
// @Router /photos/{id} [get]
func (server *Server) GetPhoto(c *gin.Context) {

	photoID := c.Param("id")
	pid, err := strconv.ParseUint(photoID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	photo := models.Photo{}

	photoReceived, err := photo.FindPhotoByID(server.DB, pid)
	if err != nil {
		errList["No_photo"] = "No Photo Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": photoReceived,
	})
}

// UpdatePhotoByID godoc
// @Summary Update Photo by ID
// @Description Update a photo by ID
// @Tags Photo
// @Accept json
// @Produce json
// @Param id path int true "Photo ID"
// @Param UpdatePhoto body models.UpdatePhoto true "Photo Data"
// @Security ApiKeyAuth
// @Success 200 {object} models.Photo
// @Router /photos/{id} [put]
func (server *Server) UpdatePhoto(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	photoID := c.Param("id")
	// Check if the photo id is valid
	pid, err := strconv.ParseUint(photoID, 10, 64)
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
	//Check if the photo exist
	origPhoto := models.Photo{}
	err = server.DB.Debug().Model(models.Photo{}).Where("id = ?", pid).Take(&origPhoto).Error
	if err != nil {
		errList["No_photo"] = "No Photo Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	if uid != origPhoto.UserID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Read the data photoed
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
	photo := models.Photo{}
	err = json.Unmarshal(body, &photo)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	photo.ID = origPhoto.ID //this is important to tell the model the photo id to update, the other update field are set above
	photo.UserID = origPhoto.UserID

	photo.Prepare()
	errorMessages := photo.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	photoUpdated, err := photo.UpdateAPhoto(server.DB)
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
		"response": photoUpdated,
	})
}

// DeletePhotoByID godoc
// @Summary Delete Photo by ID
// @Description Delete a photo by ID
// @Tags Photo
// @Accept json
// @Produce json
// @Param id path int true "Photo ID"
// @Security ApiKeyAuth
// @Success 200 {string} string "Photo deleted"
// @Router /photos/{id} [delete]
func (server *Server) DeletePhoto(c *gin.Context) {

	photoID := c.Param("id")
	// Is a valid photo id given to us?
	pid, err := strconv.ParseUint(photoID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	fmt.Println("this is delete photo sir")

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
	// Check if the photo exist
	photo := models.Photo{}
	err = server.DB.Debug().Model(models.Photo{}).Where("id = ?", pid).Take(&photo).Error
	if err != nil {
		errList["No_photo"] = "No Photo Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Is the authenticated user, the owner of this photo?
	if uid != photo.UserID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	_, err = photo.DeleteAPhoto(server.DB)
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
		"response": "Photo deleted",
	})
}

func (server *Server) GetUserPhotos(c *gin.Context) {

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
	photo := models.Photo{}
	photos, err := photo.FindUserPhotos(server.DB, uint32(uid))
	if err != nil {
		errList["No_photo"] = "No Photo Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": photos,
	})
}
