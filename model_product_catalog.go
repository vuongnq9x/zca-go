package zca

type ProductCatalogItem struct {
	Price       string `json:"price"`
	Description string `json:"description"`
	// Relative path used to build the product URL: https://catalog.zalo.me/${path}
	Path          string   `json:"path"`
	ProductID     string   `json:"product_id"`
	ProductName   string   `json:"product_name"`
	CurrencyUnit  string   `json:"currency_unit"`
	ProductPhotos []string `json:"product_photos"`
	CreateTime    int64    `json:"create_time"`
	CatalogID     string   `json:"catalog_id"`
	OwnerID       string   `json:"owner_id"`
}
