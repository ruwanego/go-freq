package actions

import "gofreq/pkg/types"

type Builder struct {
    actions []Action
}

func NewBuilder() *Builder {
    return &Builder{actions: []Action{}}
}

func (b *Builder) BuyLimit(pair string, price, amount float64, tag string) {
    b.actions = append(b.actions, Action{
        Type:   ActionBuy,
        Pair:   pair,
        Side:   types.SideBuy,
        Price:  price,
        Amount: amount,
        Tag:    tag,
    })
}

func (b *Builder) SellLimit(pair string, price, amount float64, tag string) {
    b.actions = append(b.actions, Action{
        Type:   ActionSell,
        Pair:   pair,
        Side:   types.SideSell,
        Price:  price,
        Amount: amount,
        Tag:    tag,
    })
}

func (b *Builder) Build() []Action {
    return b.actions
}
