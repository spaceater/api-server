package model

import (
	"ismismcube-backend/internal/config"
	"log"
	"time"
)

type PageView struct {
	ID        int       `json:"id"`
	Time      time.Time `json:"time"`
	VisitorIP string    `json:"visitor_ip"`
	UserAgent string    `json:"user_agent"`
}

func GetPageViewCount() (int, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM ismismcube_page_views").Scan(&count)
	if err != nil {
		log.Printf("Error getting page view count: %v", err)
		return -1, err
	}
	return count, nil
}

func AddPageView(view *PageView) (*PageView, error) {
	result, err := config.DB.Exec(
		"INSERT INTO ismismcube_page_views (visitor_ip, user_agent) VALUES (?, ?)",
		view.VisitorIP, view.UserAgent,
	)
	if err != nil {
		log.Printf("Error adding page view: %v", err)
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

func GetPageViews(limit, offset int) ([]PageView, error) {
	rows, err := config.DB.Query(
		"SELECT id, time, visitor_ip, user_agent FROM ismismcube_page_views ORDER BY time DESC LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Printf("Error getting page views: %v", err)
		return nil, err
	}
	defer rows.Close()

	var views []PageView
	for rows.Next() {
		var view PageView
		err := rows.Scan(&view.ID, &view.Time, &view.VisitorIP, &view.UserAgent)
		if err != nil {
			log.Printf("Error scanning page view: %v", err)
			return nil, err
		}
		views = append(views, view)
	}

	return views, nil
}

func DeletePageView(id int) error {
	_, err := config.DB.Exec("DELETE FROM ismismcube_page_views WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting page view: %v", err)
		return err
	}

	return nil
}
