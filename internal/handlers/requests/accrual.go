package requests

type AccrualOrderRequest struct {
	Order string                    `json:"order" validate:"required"`
	Goods []AccrualOrderGoodRequest `json:"goods,omitempty"`
}
type AccrualOrderGoodRequest struct {
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}
type RewardRuleRequest struct {
	Match      string  `json:"match" validate:"required,min=1" message:"Match is required"`
	Reward     float64 `json:"reward" validate:"required,gt=0" message:"Reward is required"`
	RewardType string  `json:"reward_type" validate:"required,oneof=% pt" message:"RewardType is required and one of % or pt"` // "%" или "pt"
}
