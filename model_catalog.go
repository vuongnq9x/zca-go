package zca

type CatalogItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Version   int64  `json:"version"`
	OwnerID   string `json:"ownerId"`
	IsDefault bool   `json:"isDefault"`
	// Relative path used to build the catalog URL: https://catalog.zalo.me/${path}
	Path         string  `json:"path"`
	CatalogPhoto *string `json:"catalogPhoto"`
	TotalProduct int64   `json:"totalProduct"`
	CreatedTime  int64   `json:"created_time"`
}
