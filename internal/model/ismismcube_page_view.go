package model

import (
	"ismismcube-backend/internal/config"
	"log"
	"time"
)

type IsmismcubePageView struct {
	ID        int       `json:"id"`
	Time      time.Time `json:"time"`
	VisitorIP string    `json:"visitor_ip"`
	UserAgent string    `json:"user_agent"`
}

func GetIsmismcubePageViewCount() (int, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM ismismcube_page_views").Scan(&count)
	if err != nil {
		log.Printf("Error getting ismismcube page view count: %v", err)
		return -1, err
	}
	return count, nil
}

func AddIsmismcubePageView(view *IsmismcubePageView) (*IsmismcubePageView, error) {
	result, err := config.DB.Exec(
		"INSERT INTO ismismcube_page_views (visitor_ip, user_agent) VALUES (?, ?)",
		view.VisitorIP, view.UserAgent,
	)
	if err != nil {
		log.Printf("Error adding ismismcube page view: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return nil, err
	}

	view.ID = int(id)
	return view, nil
}

func GetIsmismcubePageViews(limit, offset int) ([]IsmismcubePageView, error) {
	rows, err := config.DB.Query(
		"SELECT id, time, visitor_ip, user_agent FROM ismismcube_page_views ORDER BY time DESC LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Printf("Error getting ismismcube page views: %v", err)
		return nil, err
	}
	defer rows.Close()

	var views []IsmismcubePageView
	for rows.Next() {
		var view IsmismcubePageView
		err := rows.Scan(&view.ID, &view.Time, &view.VisitorIP, &view.UserAgent)
		if err != nil {
			log.Printf("Error scanning ismismcube page view: %v", err)
			return nil, err
		}
		views = append(views, view)
	}

	return views, nil
}

func DeleteIsmismcubePageView(id int) error {
	_, err := config.DB.Exec("DELETE FROM ismismcube_page_views WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting ismismcube page view: %v", err)
		return err
	}

	return nil
}
