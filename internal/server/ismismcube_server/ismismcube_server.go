package ismismcube_server

import (
	"ismismcube-backend/internal/model"
)

func GetPageViewCount() (int, error) {
	return model.GetPageViewCount()
}

func AddPageView(visitorIP, userAgent string) (*model.PageView, error) {
	view := &model.PageView{
		VisitorIP: visitorIP,
		UserAgent: userAgent,
	}
	return model.AddPageView(view)
}

func GetPageViews(limit, offset int) ([]model.PageView, error) {
	return model.GetPageViews(limit, offset)
}

func DeletePageView(id int) error {
	return model.DeletePageView(id)
}
