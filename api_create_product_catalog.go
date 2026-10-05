package zca

import (
	"context"
	"net/http"
)

type CreateProductCatalogPayload struct {
	CatalogID   string
	ProductName string
	Price       string
	Description string
	// Up to 5 media files, uploaded via UploadProductPhoto.
	Files []AttachmentSource
	// Product photo URLs (up to 5 in total with Files), e.g. from UploadProductPhoto.
	ProductPhotos []string
}

type CreateProductCatalogResponse struct {
	Item             ProductCatalogItem `json:"item"`
	VersionLsCatalog int64              `json:"version_ls_catalog"`
	VersionCatalog   int64              `json:"version_catalog"`
}

// CreateProductCatalog creates a product in a catalog (zBusiness, max 5 photos).
func (a *API) CreateProductCatalog(ctx context.Context, payload CreateProductCatalogPayload) (*CreateProductCatalogResponse, error) {
	if len(payload.Files) > 5 {
		return nil, newError("Maximum 5 media files are allowed")
	}
	photos := append([]string{}, payload.ProductPhotos...)
	for _, f := range payload.Files {
		up, err := a.UploadProductPhoto(ctx, UploadProductPhotoPayload{File: f})
		if err != nil {
			return nil, err
		}
		u := up.NormalURL
		if u == "" {
			u = up.HdURL
		}
		photos = append(photos, u)
	}
	if len(photos) > 5 {
		return nil, newError("Maximum 5 media files are allowed")
	}
	params := map[string]any{
		"product_name":   payload.ProductName,
		"price":          payload.Price,
		"description":    payload.Description,
		"product_photos": photos,
		"catalog_id":     payload.CatalogID,
		"currency_unit":  "₫",
		"create_time":    nowMs(),
	}
	return call[*CreateProductCatalogResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("catalog")+"/api/prodcatalog/product/create", nil, true), params)
}
