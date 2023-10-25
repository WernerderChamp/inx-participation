package participation

type EventListMetric struct {
	EventStatusUpcomming  float64
	EventStatusCommencing float64
	EventStatusHolding    float64
	EventStatusEnded      float64
	EventsTotal           float64
}

type StakingMetric struct {
	EventName           string
	StakingCurrent      float64
	StakingPerMilestone float64
	StakingAccumulated  float64
}

type AnswerMetric struct {
	StakingCurrent     float64
	StakingAccumulated float64
}

type BallotMetric struct {
	EventName       string
	QuestionMetrics []map[uint8]AnswerMetric
}

type PrometheusMetrics struct {
	ParticipationEvents EventListMetric
	StakingMetrics      map[EventID]StakingMetric
	BallotMetrics       map[EventID]BallotMetric
}

func (pm *Manager) ProvidePrometheusMetrics() (*PrometheusMetrics, error) {
	pm.Lock()
	defer pm.Unlock()

	currentMilestone, err := pm.readLedgerIndex()
	if err != nil {
		return nil, err
	}
	eventList := pm.eventsWithoutLocking()
	// eventCount := len(eventList)

	var metrics PrometheusMetrics
	metrics.StakingMetrics = map[EventID]StakingMetric{}
	metrics.BallotMetrics = map[EventID]BallotMetric{}

	for eventID, event := range eventList {
		switch event.Status(currentMilestone) {
		case "upcoming":
			metrics.ParticipationEvents.EventStatusUpcomming++
		case "commencing":
			metrics.ParticipationEvents.EventStatusCommencing++
		case "holding":
			metrics.ParticipationEvents.EventStatusHolding++
		case "ended":
			metrics.ParticipationEvents.EventStatusEnded++
		}
		metrics.ParticipationEvents.EventsTotal++

		switch event.Payload.(type) {
		case *Ballot:
			questions := event.BallotQuestions()
			questionMetrics := make([]map[uint8]AnswerMetric, len(questions))
			for questionIdx, question := range questions {
				questionMetrics[questionIdx] = map[uint8]AnswerMetric{}
				for _, answer := range question.Answers {
					currentBalance, err := pm.currentBallotVoteBalanceForQuestionAndAnswer(eventID, currentMilestone, uint8(questionIdx), answer.Value)
					if err != nil {
						return nil, err
					}
					accumaletedBalance, err := pm.accumulatedBallotVoteBalanceForQuestionAndAnswer(eventID, currentMilestone, uint8(questionIdx), answer.Value)
					if err != nil {
						return nil, err
					}
					questionMetrics[questionIdx][answer.Value] = AnswerMetric{
						StakingCurrent:     float64(currentBalance),
						StakingAccumulated: float64(accumaletedBalance),
					}
				}
				metrics.BallotMetrics[eventID] = BallotMetric{
					EventName:       event.Name,
					QuestionMetrics: questionMetrics,
				}
			}
		case *Staking:

			currentStaking, err := pm.currentRewardsPerMilestoneForStakingEvent(eventID, currentMilestone)
			if err != nil {
				return nil, err
			}
			stakingParticipation, err := pm.totalStakingParticipationForEvent(eventID, currentMilestone)
			if err != nil {
				return nil, err
			}
			stakingMetrics := StakingMetric{
				StakingCurrent:      float64(currentStaking),
				StakingPerMilestone: float64(stakingParticipation.staked),
				StakingAccumulated:  float64(stakingParticipation.rewarded),
			}
			metrics.StakingMetrics[eventID] = stakingMetrics
		}
	}

	return &metrics, nil
}
