package request

type UpdateWikiRequest struct {
	Language *int      `json:"language"`
	Code     *string   `json:"code"`
	Public   *int      `json:"public"`
	Title    *string   `json:"title"`
	Icon     *string   `json:"icon"`
	Keywords *string   `json:"keywords"`
	Level    *int      `json:"level"`
	Unit     *string   `json:"unit"`
	Elements []Element `json:"elements"`
}
