package models

// CardRank 扑克牌点数
type CardRank string

// CardSuit 扑克牌花色
type CardSuit string

// HandRank 牌型
type HandRank string

const (
	// 扑克牌点数
	RankAce   CardRank = "ace"
	RankTwo   CardRank = "2"
	RankThree CardRank = "3"
	RankFour  CardRank = "4"
	RankFive  CardRank = "5"
	RankSix   CardRank = "6"
	RankSeven CardRank = "7"
	RankEight CardRank = "8"
	RankNine  CardRank = "9"
	RankTen   CardRank = "10"
	RankJack  CardRank = "jack"
	RankQueen CardRank = "queen"
	RankKing  CardRank = "king"

	// 扑克牌花色
	SuitHearts   CardSuit = "hearts"
	SuitDiamonds CardSuit = "diamonds"
	SuitClubs    CardSuit = "clubs"
	SuitSpades   CardSuit = "spades"

	// 牌型
	HighCard      HandRank = "high_card"
	OnePair       HandRank = "one_pair"
	TwoPair       HandRank = "two_pair"
	ThreeOfAKind  HandRank = "three_of_a_kind"
	Straight      HandRank = "straight"
	Flush         HandRank = "flush"
	FullHouse     HandRank = "full_house"
	FourOfAKind   HandRank = "four_of_a_kind"
	StraightFlush HandRank = "straight_flush"
	RoyalFlush    HandRank = "royal_flush"

	// 德州扑克游戏回合
	RoundPreflop  = "preflop"
	RoundFlop     = "flop"
	RoundTurn     = "turn"
	RoundRiver    = "river"
	RoundShowdown = "showdown"
)

// Card 扑克牌
type Card struct {
	Rank CardRank `json:"rank"`
	Suit CardSuit `json:"suit"`
}

// PokerPlayerState 德州扑克游戏中的玩家状态
type PokerPlayerState struct {
	ID        uint     `json:"id"`
	Username  string   `json:"username"`
	Chips     int      `json:"chips"`                // 玩家筹码
	Bet       int      `json:"bet"`                  // 当前下注金额
	Folded    bool     `json:"folded"`               // 是否弃牌
	IsAllIn   bool     `json:"is_all_in"`            // 是否全押
	HandCards []Card   `json:"hand_cards,omitempty"` // 手牌，只对自己可见
	HandSize  int      `json:"hand_size,omitempty"`  // 手牌数量，对其他玩家可见
	Hand      HandRank `json:"hand,omitempty"`       // 最终牌型
}

// PokerGameState 德州扑克游戏状态
type PokerGameState struct {
	Status         string             `json:"status"` // waiting, playing, ended
	Players        []PokerPlayerState `json:"players"`
	CurrentPlayer  *Player            `json:"current_player,omitempty"`
	CommunityCards []Card             `json:"community_cards"`
	Pot            int                `json:"pot"`             // 底池金额
	CurrentBet     int                `json:"current_bet"`     // 当前回合最高下注额
	DealerPosition int                `json:"dealer_position"` // 庄家位置
	Round          string             `json:"round"`           // 当前回合
}

// Player 玩家信息（用于当前行动的玩家）
type Player struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// PokerGameResult 德州扑克游戏结果
type PokerGameResult struct {
	Message     string         `json:"message"`
	Winner      *Player        `json:"winner,omitempty"`
	WinningHand *WinningHand   `json:"winning_hand,omitempty"`
	Players     []PlayerResult `json:"players"`
}

// WinningHand 获胜牌型
type WinningHand struct {
	Rank  HandRank `json:"rank"`
	Cards []Card   `json:"cards"`
}

// PlayerResult 玩家结果
type PlayerResult struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Hand     []Card   `json:"hand"`
	HandRank HandRank `json:"hand_rank"`
	Winnings int      `json:"winnings"`
}
