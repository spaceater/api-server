package ismismcube_server

import (
	"ismismcube-backend/internal/model"
  "ismismcube-backend/internal/websocket"
)

func GetPageViewCount() (int, error) {
	return model.GetIsmismcubePageViewCount()
}

func AddPageView(visitorIP, userAgent string) (*model.IsmismcubePageView, error) {
	view := &model.IsmismcubePageView{
		VisitorIP: visitorIP,
		UserAgent: userAgent,
	}
	return model.AddIsmismcubePageView(view)
}

func GetPageViews(limit, offset int) ([]model.IsmismcubePageView, error) {
	return model.GetIsmismcubePageViews(limit, offset)
}

func DeletePageView(id int) error {
	return model.DeleteIsmismcubePageView(id)
}

func SendDanmu(clientIP, content string) error {
	err := model.CheckAndSetRateLimit(clientIP)
	if err != nil {
		return err
	}
	ws.BroadcastIsmismcubeDanmu(content)
	return nil
}
