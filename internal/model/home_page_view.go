package model

import (
	"ismismcube-backend/internal/config"
	"log"
	"time"
)

type HomePageView struct {
	ID        int       `json:"id"`
	Time      time.Time `json:"time"`
	VisitorIP string    `json:"visitor_ip"`
	UserAgent string    `json:"user_agent"`
}

func GetHomePageViewCount() (int, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM home_page_views").Scan(&count)
	if err != nil {
		log.Printf("Error getting home page view count: %v", err)
		return -1, err
	}
	return count, nil
}

func AddHomePageView(view *HomePageView) (*HomePageView, error) {
	result, err := config.DB.Exec(
		"INSERT INTO home_page_views (visitor_ip, user_agent) VALUES (?, ?)",
		view.VisitorIP, view.UserAgent,
	)
	if err != nil {
		log.Printf("Error adding home page view: %v", err)
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

func GetHomePageViews(limit, offset int) ([]HomePageView, error) {
	rows, err := config.DB.Query(
		"SELECT id, time, visitor_ip, user_agent FROM home_page_views ORDER BY time DESC LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Printf("Error getting home page views: %v", err)
		return nil, err
	}
	defer rows.Close()

	var views []HomePageView
	for rows.Next() {
		var view HomePageView
		err := rows.Scan(&view.ID, &view.Time, &view.VisitorIP, &view.UserAgent)
		if err != nil {
			log.Printf("Error scanning home page view: %v", err)
			return nil, err
		}
		views = append(views, view)
	}

	return views, nil
}

func DeleteHomePageView(id int) error {
	_, err := config.DB.Exec("DELETE FROM home_page_views WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting home page view: %v", err)
		return err
	}

	return nil
}
