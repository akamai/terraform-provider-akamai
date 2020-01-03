package apiendpoints

type Category struct {
	APICategoryID          int    `json:"apiCategoryId,omitempty"`
	APICategoryName        string `json:"apiCategoryName"`
	APICategoryDescription string `json:"apiCategoryDescription"`
	Link                   string `json:"link"`
	LockVersion            int    `json:"lockVersion"`
	CreatedBy              string `json:"createdBy,omitempty"`
	CreateDate             string `json:"createDate,omitempty"`
	UpdatedBy              string `json:"updatedBy,omitempty"`
	UpdateDate             string `json:"updateDate,omitempty"`
}
