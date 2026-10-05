package zca

type ZBusinessPackage struct {
	Label map[string]string `json:"label,omitempty"`
	PkgID int64             `json:"pkgId"`
}

type BusinessCategory int

const (
	BusinessCategoryOther                  BusinessCategory = 0
	BusinessCategoryRealEstate             BusinessCategory = 1
	BusinessCategoryTechnologyAndDevices   BusinessCategory = 2
	BusinessCategoryTravelAndHospitality   BusinessCategory = 3
	BusinessCategoryEducationAndTraining   BusinessCategory = 4
	BusinessCategoryShoppingAndRetail      BusinessCategory = 5
	BusinessCategoryCosmeticsAndBeauty     BusinessCategory = 6
	BusinessCategoryRestaurantAndCafe      BusinessCategory = 7
	BusinessCategoryAutoAndMotorbike       BusinessCategory = 8
	BusinessCategoryFashionAndApparel      BusinessCategory = 9
	BusinessCategoryFoodAndBeverage        BusinessCategory = 10
	BusinessCategoryMediaAndEntertainment  BusinessCategory = 11
	BusinessCategoryInternalCommunications BusinessCategory = 12
	BusinessCategoryTransportation         BusinessCategory = 13
	BusinessCategoryTelecommunications     BusinessCategory = 14
)

var BusinessCategoryName = map[BusinessCategory]string{
	BusinessCategoryOther:                  "Dịch vụ khác (Không hiển thị)",
	BusinessCategoryRealEstate:             "Bất động sản",
	BusinessCategoryTechnologyAndDevices:   "Công nghệ & Thiết bị",
	BusinessCategoryTravelAndHospitality:   "Du lịch & Lưu trú",
	BusinessCategoryEducationAndTraining:   "Giáo dục & Đào tạo",
	BusinessCategoryShoppingAndRetail:      "Mua sắm & Bán lẻ",
	BusinessCategoryCosmeticsAndBeauty:     "Mỹ phẩm & Làm đẹp",
	BusinessCategoryRestaurantAndCafe:      "Nhà hàng & Quán",
	BusinessCategoryAutoAndMotorbike:       "Ô tô & Xe máy",
	BusinessCategoryFashionAndApparel:      "Thời trang & May mặc",
	BusinessCategoryFoodAndBeverage:        "Thực phẩm & Đồ uống",
	BusinessCategoryMediaAndEntertainment:  "Truyền thông & Giải trí",
	BusinessCategoryInternalCommunications: "Truyền thông nội bộ",
	BusinessCategoryTransportation:         "Vận tải",
	BusinessCategoryTelecommunications:     "Viễn thông",
}
