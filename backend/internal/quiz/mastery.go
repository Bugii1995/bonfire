package quiz

import "time"

//
// -------- Constants (tunable pedagogy knobs) --------
//

const (
	BaseCorrectDelta = 5.0
	BaseWrongDelta   = -7.0

	EasyMultiplier   = 0.6
	MediumMultiplier = 1.0
	HardMultiplier   = 1.4

	CorrectStreakBonusEvery = 3
	CorrectStreakBonus      = 2.0

	WrongStreakPenaltyEvery = 2
	WrongStreakPenalty      = 3.0

	DecayAfterDays      = 7
	DecayMediumRate     = 0.3 // per day (8–30 days)
	DecayHeavyRate      = 0.7 // per day (>30 days)
	DecayHeavyThreshold = 30

	MasteredThreshold     = 100.0
	UnmasteredHysteresis  = 90.0
	MinimumPassiveMastery = 20.0
)

//
// -------- Input snapshot --------
//

// MasteryUpdateInput describes a single answered question.
type MasteryUpdateInput struct {
	WasCorrect  bool
	Difficulty  int       // 1, 2, 3
	AnsweredAt  time.Time
	CurrentTime time.Time
}

//
// -------- Output snapshot --------
//

// MasteryUpdateResult is the updated topic state.
type MasteryUpdateResult struct {
	Mastery       float64
	CorrectStreak int
	WrongStreak   int
	IsMastered    bool
	LastSeen      time.Time
}

//
// -------- Public API --------
//

// UpdateMastery applies correctness, difficulty, streaks, and decay
// to produce a new mastery state.
func UpdateMastery(
	current TopicProgress,
	input MasteryUpdateInput,
) MasteryUpdateResult {

	mastery := applyDecay(current, input.CurrentTime)

	// Base delta
	delta := baseDelta(input.WasCorrect)

	// Difficulty scaling
	delta *= difficultyMultiplier(input.Difficulty)

	// Update streaks
	correctStreak, wrongStreak := updateStreaks(
		current.CorrectStreak,
		current.WrongStreak,
		input.WasCorrect,
	)

	// Apply streak effects
	delta += streakBonus(correctStreak)
	delta -= streakPenalty(wrongStreak)

	// Apply delta
	mastery = clamp(mastery+delta, 0, 100)

	// Mastery flags with hysteresis
	isMastered := mastery >= MasteredThreshold
	if mastery < UnmasteredHysteresis {
		isMastered = false
	}

	return MasteryUpdateResult{
		Mastery:       mastery,
		CorrectStreak: correctStreak,
		WrongStreak:   wrongStreak,
		IsMastered:    isMastered,
		LastSeen:      input.AnsweredAt,
	}
}

//
// -------- Internal helpers --------
//

func baseDelta(correct bool) float64 {
	if correct {
		return BaseCorrectDelta
	}
	return BaseWrongDelta
}

func difficultyMultiplier(difficulty int) float64 {
	switch difficulty {
	case 1:
		return EasyMultiplier
	case 3:
		return HardMultiplier
	default:
		return MediumMultiplier
	}
}

func updateStreaks(correctStreak, wrongStreak int, correct bool) (int, int) {
	if correct {
		return correctStreak + 1, 0
	}
	return 0, wrongStreak + 1
}

func streakBonus(correctStreak int) float64 {
	if correctStreak > 0 && correctStreak%CorrectStreakBonusEvery == 0 {
		return CorrectStreakBonus
	}
	return 0
}

func streakPenalty(wrongStreak int) float64 {
	if wrongStreak > 0 && wrongStreak%WrongStreakPenaltyEvery == 0 {
		return WrongStreakPenalty
	}
	return 0
}

func applyDecay(current TopicProgress, now time.Time) float64 {
	// If topic has never been seen, do not apply decay
	if current.LastSeen.IsZero() {
		return current.Mastery
	}

	days := int(now.Sub(current.LastSeen).Hours() / 24)
	mastery := current.Mastery

	if days <= DecayAfterDays {
		return mastery
	}

	if days <= DecayHeavyThreshold {
		mastery -= float64(days-DecayAfterDays) * DecayMediumRate
	} else {
		mastery -= float64(DecayHeavyThreshold-DecayAfterDays)*DecayMediumRate +
			float64(days-DecayHeavyThreshold)*DecayHeavyRate
	}

	if mastery < MinimumPassiveMastery {
		return MinimumPassiveMastery
	}
	return mastery
}


func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
