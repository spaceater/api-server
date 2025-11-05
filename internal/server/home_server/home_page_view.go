package home_server

import (
	"ismismcube-backend/internal/model"
)

func GetPageViewCount() (int, error) {
	return model.GetHomePageViewCount()
}

func AddPageView(visitorIP, userAgent string) (*model.HomePageView, error) {
	view := &model.HomePageView{
		VisitorIP: visitorIP,
		UserAgent: userAgent,
	}
	return model.AddHomePageView(view)
}

func GetPageViews(limit, offset int) ([]model.HomePageView, error) {
	return model.GetHomePageViews(limit, offset)
}

func DeletePageView(id int) error {
	return model.DeleteHomePageView(id)
}
