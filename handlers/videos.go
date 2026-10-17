package handlers

import (
	"encoding/json"
	"net/http"
	"robotics-project/database"
	"robotics-project/models"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, youtube_id, description FROM videos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		rows.Scan(&v.ID, &v.Title, &v.YouTubeID, &v.Description)
		videos = append(videos, v)
	}
	json.NewEncoder(w).Encode(videos)
}
