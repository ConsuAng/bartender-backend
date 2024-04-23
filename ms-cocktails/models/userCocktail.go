package models

type UserCocktail struct {
	UserCocktailID int    `gorm:"primaryKey"`
	UserID         int    `json:"user_id"`
	CocktailID     string `json:"cocktail_id"`
	CocktailName   string `json:"cocktail_name"`
	CocktailImage  string `json:"cocktail_image"`
}

func (UserCocktail) TableName() string {
	return "user_cocktail"
}
