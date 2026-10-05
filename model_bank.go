package zca

type BankInfo struct {
	Bin           int64  `json:"bin"`
	Logo          string `json:"logo"`
	Name          string `json:"name"`
	NameEng       string `json:"name_eng"`
	ShortName     string `json:"short_name"`
	SearchKeyWord string `json:"search_key_word"`
}

type BankAccount struct {
	ID         string `json:"id"`
	Bin        int64  `json:"bin"`
	Default    bool   `json:"default"`
	BankNumber string `json:"bank_number"`
	// Not returned by getBankAccounts, createBankAccount.
	BankLogo   string `json:"bank_logo,omitempty"`
	HolderName string `json:"holder_name"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
	AccountID  int64  `json:"account_id"`
	// Not returned by getBankAccounts, createBankAccount.
	BankName  string `json:"bank_name,omitempty"`
	IsDefault bool   `json:"is_default"`
}

// BinBankCard is the bank BIN code list (after MITM on Mobile, WEB and banks supported by Zalo).
// Docs: https://developers.zalo.me/docs/zalo-notification-service/phu-luc/danh-sach-bin-code
type BinBankCard int64

const (
	// NH TMCP An Bình
	BinBankCardABBank BinBankCard = 970425
	// NH TMCP Á Châu
	BinBankCardACB BinBankCard = 970416
	// NH Nông nghiệp và Phát triển Nông thôn Việt Nam
	BinBankCardAgribank BinBankCard = 970405
	// NH TMCP Đầu tư và Phát triển Việt Nam
	BinBankCardBIDV BinBankCard = 970418
	// Ngân hàng BNP Paribas - CN TP. Hồ Chí Minh
	BinBankCardBNP_Paribas_HCM BinBankCard = 963666
	// Ngân hàng BNP Paribas - CN Hà Nội
	BinBankCardBNP_Paribas_HN BinBankCard = 963668
	// NH TMCP Bản Việt
	BinBankCardBVBank BinBankCard = 970454
	// NH TMCP Bắc Á
	BinBankCardBacA_Bank BinBankCard = 970409
	// NH TMCP Bảo Việt
	BinBankCardBaoViet_Bank BinBankCard = 970438
	// NH số CAKE by VPBank - TMCP Việt Nam Thịnh Vượng
	BinBankCardCAKE BinBankCard = 546034
	// Ngân hàng Cathay United - CN TP. Hồ Chí Minh
	BinBankCardCathay_United_HCM BinBankCard = 168999
	// NH Thương mại TNHH MTV Xây dựng Việt Nam
	// @note Also observed as "VCBNeo" in supported banks list (same bin) - Old is bank: CB_Bank
	BinBankCardVCBNeo BinBankCard = 970444
	// NH TNHH MTV CIMB Việt Nam
	BinBankCardCIMB_Bank BinBankCard = 422589
	// NH Hợp tác xã Việt Nam
	BinBankCardCoop_Bank BinBankCard = 970446
	// NH TNHH MTV Phát triển Singapore - CN TP. Hồ Chí Minh
	BinBankCardDBS_Bank BinBankCard = 796500
	// NH TMCP Đông Á
	BinBankCardDongA_Bank BinBankCard = 970406
	// NH TMCP Xuất Nhập khẩu Việt Nam
	BinBankCardEximbank BinBankCard = 970431
	// Ngân hàng Citibank Việt Nam
	BinBankCardCitibank BinBankCard = 533948
	// NH TMCP Dầu khí Toàn cầu
	BinBankCardGPBank BinBankCard = 970408
	// NH TMCP Phát triển TP. Hồ Chí Minh
	BinBankCardHDBank BinBankCard = 970437
	// NH TNHH MTV HSBC (Việt Nam)
	BinBankCardHSBC BinBankCard = 458761
	// NH TNHH MTV Hong Leong Việt Nam
	BinBankCardHongLeong_Bank BinBankCard = 970442
	// NH Công nghiệp Hàn Quốc - CN TP. Hồ Chí Minh
	BinBankCardIBK_HCM BinBankCard = 970456
	// NH Công nghiệp Hàn Quốc - CN Hà Nội
	BinBankCardIBK_HN BinBankCard = 970455
	// NH TNHH Indovina
	BinBankCardIndovina_Bank BinBankCard = 970434
	// NH Đại chúng TNHH Kasikornbank - CN TP. Hồ Chí Minh
	BinBankCardKBank BinBankCard = 668888
	// NH TMCP Kiên Long
	BinBankCardKienlongBank BinBankCard = 970452
	// NH Kookmin - CN TP. Hồ Chí Minh
	BinBankCardKookmin_Bank_HCM BinBankCard = 970463
	// NH Kookmin - CN Hà Nội
	BinBankCardKookmin_Bank_HN BinBankCard = 970462
	// Liobank by OCB
	BinBankCardLiobank BinBankCard = 963369
	// NH TMCP Lộc Phát Việt Nam
	BinBankCardLPBank BinBankCard = 970449
	// NH TMCP Quân đội
	BinBankCardMB_Bank BinBankCard = 970422
	// NH TMCP Hàng Hải
	BinBankCardMSB BinBankCard = 970426
	// MoMo
	BinBankCardMoMo BinBankCard = 971025
	// NH TMCP Quốc Dân
	BinBankCardNCB BinBankCard = 970419
	// NH TMCP Nam Á
	BinBankCardNam_A_Bank BinBankCard = 970428
	// NH Nonghyup - CN Hà Nội
	BinBankCardNongHyup_Bank BinBankCard = 801011
	// NH TMCP Phương Đông
	BinBankCardOCB BinBankCard = 970448
	// NH Thương mại TNHH MTV Đại Dương
	BinBankCardOcean_Bank BinBankCard = 970414
	// NH TMCP Thịnh vượng và Phát triển
	BinBankCardPGBank BinBankCard = 970430
	// NH TMCP Đại Chúng Việt Nam
	BinBankCardPVcomBank BinBankCard = 970412
	// NH TNHH MTV Public Việt Nam
	BinBankCardPublic_Bank_Vietnam BinBankCard = 970439
	// NH TMCP Sài Gòn
	BinBankCardSCB BinBankCard = 970429
	// NH TMCP Sài Gòn - Hà Nội
	BinBankCardSHB BinBankCard = 970443
	// NH TMCP Sài Gòn Thương Tín
	BinBankCardSacombank BinBankCard = 970403
	// NH TMCP Sài Gòn Công Thương
	BinBankCardSaigon_Bank BinBankCard = 970400
	// NH TMCP Đông Nam Á
	BinBankCardSeABank BinBankCard = 970440
	// NH TNHH MTV Shinhan Việt Nam
	BinBankCardShinhan_Bank BinBankCard = 970424
	// NH TNHH MTV Standard Chartered Bank Việt Nam
	BinBankCardStandard_Chartered_Vietnam BinBankCard = 970410
	// NH số TNEX
	BinBankCardTNEX BinBankCard = 963326
	// NH TMCP Tiên Phong
	BinBankCardTPBank BinBankCard = 970423
	// NH TMCP Kỹ thương Việt Nam
	BinBankCardTechcombank BinBankCard = 970407
	// NH số Timo by Bản Việt Bank
	BinBankCardTimo BinBankCard = 963388
	// NH số UBank by VPBank
	BinBankCardUBank BinBankCard = 546035
	// NH United Overseas Bank Việt Nam
	BinBankCardUnited_Overseas_Bank_Vietnam BinBankCard = 970458
	// NH TMCP Quốc tế Việt Nam
	BinBankCardVIB BinBankCard = 970441
	// NH TMCP Việt Nam Thịnh Vượng
	BinBankCardVPBank BinBankCard = 970432
	// NH Liên doanh Việt - Nga
	BinBankCardVRB BinBankCard = 970421
	// NH TMCP Việt Á
	BinBankCardVietABank BinBankCard = 970427
	// NH TMCP Việt Nam Thương Tín
	BinBankCardVietBank BinBankCard = 970433
	// NH TMCP Ngoại Thương Việt Nam
	BinBankCardVietcombank BinBankCard = 970436
	// NH TMCP Công thương Việt Nam
	BinBankCardVietinBank BinBankCard = 970415
	// NH TNHH MTV Woori Việt Nam
	BinBankCardWoori_Bank BinBankCard = 970457
)
