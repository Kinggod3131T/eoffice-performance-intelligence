package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGoal(c *gin.Context) {
	var goal Goal
	c.BindJSON(&goal)

	err := DB.QueryRow(
		"INSERT INTO goals (title, status) VALUES ($1, $2) RETURNING id",
		goal.Title, goal.Status,
	).Scan(&goal.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

func GetGoals(c *gin.Context) {
	rows, err := DB.Query("SELECT id, title, status FROM goals")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var goals []Goal

	for rows.Next() {
		var g Goal
		err := rows.Scan(&g.ID, &g.Title, &g.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		goals = append(goals, g)
	}

	c.JSON(http.StatusOK, goals)
}
