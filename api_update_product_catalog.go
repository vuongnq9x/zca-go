package zca

import (
	"context"
	"net/http"
)

type UpdateProductCatalogPayload struct {
	CatalogID   string
	ProductID   string
	ProductName string
	Price       string
	Description string
	CreateTime  int64
	// Files: up to 5 media files, uploaded with UploadProductPhoto and appended to ProductPhotos.
	Files []AttachmentSource
	// ProductPhotos: product photo URLs (up to 5 in total), e.g. from UploadProductPhoto.
	ProductPhotos []string
}

type UpdateProductCatalogResponse struct {
	Item             ProductCatalogItem `json:"item"`
	VersionLsCatalog int64              `json:"version_ls_catalog"`
	VersionCatalog   int64              `json:"version_catalog"`
}

// UpdateProductCatalog updates a catalog product (zBusiness, max 5 media).
func (a *API) UpdateProductCatalog(ctx context.Context, payload UpdateProductCatalogPayload) (*UpdateProductCatalogResponse, error) {
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
		"product_id": payload.ProductID, "product_name": payload.ProductName, "price": payload.Price,
		"description": payload.Description, "product_photos": photos, "catalog_id": payload.CatalogID,
		"currency_unit": "₫", "create_time": payload.CreateTime,
	}
	return call[*UpdateProductCatalogResponse](ctx, a.Session, http.MethodPost,
		a.MakeURL(a.svc("catalog")+"/api/prodcatalog/product/update", nil, true), params)
}
