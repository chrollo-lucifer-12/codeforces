package controllers

import (
	"strconv"
	"time"

	"github.com/chrollo-lucifer-12/backend/src/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetContests(ctx *gin.Context, contestService *services.ContestService, status string) {
	limitStr := ctx.Query("limit")
	offsetStr := ctx.Query("offset")

	limit := 10
	offset := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	contests, err := contestService.GetContests(status, limit, offset)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, contests)
}

func GetContest(ctx *gin.Context, contestService *services.ContestService) {
	contestIdStr := ctx.Param("contestId") // returns string

	// convert string to uuid.UUID
	contestId, err := uuid.Parse(contestIdStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid contest ID"})
		return
	}

	// now you can use contestId in your service
	contest, err := contestService.GetContest(contestId)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "contest not found"})
		return
	}

	ctx.JSON(200, contest)
}

func GetChallenge(ctx *gin.Context, contestService *services.ContestService) {
	contestIdStr := ctx.Param("contestId")
	contestId, err := uuid.Parse(contestIdStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid contest ID"})
		return
	}

	// Get challengeId from path
	challengeIdStr := ctx.Param("challengeId")
	challengeId, err := uuid.Parse(challengeIdStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid challenge ID"})
		return
	}

	// Fetch the challenge for this contest
	challenge, err := contestService.GetChallenge(challengeId, contestId)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "challenge not found in contest"})
		return
	}

	ctx.JSON(200, challenge)
}

type CreateContestRequest struct {
	Title     string    `json:"title" binding:"required"`
	StartTime time.Time `json:"startTime" binding:"required"`
}

func CreateContest(ctx *gin.Context, contestService *services.ContestService) {
	var req CreateContestRequest

	// Bind JSON body to struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call the service to create the contest
	contest, err := contestService.CreateContest(req.Title, req.StartTime)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create contest"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Contest created successfully",
		"contest": contest,
	})
}
