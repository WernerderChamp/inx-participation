package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	participationEvents      *prometheus.GaugeVec
	ballotAnswersCurrent     *prometheus.GaugeVec
	ballotAnswersAccumulated *prometheus.GaugeVec
	stakingCurrent           *prometheus.GaugeVec
	stakingPerMilestone      *prometheus.GaugeVec
	stakingAccumulated       *prometheus.GaugeVec
)

func configureParticipationMetrics(registry *prometheus.Registry) {
	participationEvents = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "iota",
			Subsystem: "participation",
			Name:      "participation_events",
			Help:      "Number of participation events added to the node.",
		},
		[]string{"status"},
	)

	ballotAnswersCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "iota",
			Subsystem: "participation",
			Name:      "ballot_answers_current",
			Help:      "Current amount of tokens voting for this answer",
		},
		[]string{"eventID", "eventName", "questionID", "answerID"},
	)
	ballotAnswersAccumulated = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "iota",
			Subsystem: "participation",
			Name:      "ballot_answers_total",
			Help:      "Accumulated amount of tokens voting for this answer",
		},
		[]string{"eventID", "eventName", "questionID", "answerID"},
	)

	stakingCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "iota",
			Subsystem: "participation",
			Name:      "staking_current",
			Help:      "Current amount of tokens staked for this event",
		},
		[]string{"eventID", "eventName"},
	)
	stakingPerMilestone = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "iota",
			Subsystem: "participation",
			Name:      "staking_per_milestone",
			Help:      "Current amount of new tokens rewarded per milestone",
		},
		[]string{"eventID", "eventName"},
	)
	stakingAccumulated = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "iota",
			Subsystem: "participation",
			Name:      "staking_total",
			Help:      "Accumulated amount of new tokens rewarded",
		},
		[]string{"eventID", "eventName"},
	)

	registry.MustRegister(participationEvents)
	registry.MustRegister(ballotAnswersCurrent)
	registry.MustRegister(ballotAnswersAccumulated)
	registry.MustRegister(stakingCurrent)
	registry.MustRegister(stakingPerMilestone)
	registry.MustRegister(stakingAccumulated)
}

func collectParticipationMetrics() error {
	metrics, err := deps.ParticipationManager.ProvidePrometheusMetrics()
	if err != nil {
		return err
	}
	participationEvents.WithLabelValues("upcomming").Set(metrics.ParticipationEvents.EventStatusUpcomming)
	participationEvents.WithLabelValues("commencing").Set(metrics.ParticipationEvents.EventStatusCommencing)
	participationEvents.WithLabelValues("holding").Set(metrics.ParticipationEvents.EventStatusHolding)
	participationEvents.WithLabelValues("ended").Set(metrics.ParticipationEvents.EventStatusEnded)
	participationEvents.WithLabelValues("total").Set(metrics.ParticipationEvents.EventsTotal)

	for eventID, stakingMetric := range metrics.StakingMetrics {
		stakingCurrent.WithLabelValues(eventID.ToHex(), stakingMetric.EventName).Set(stakingMetric.StakingCurrent)
		stakingPerMilestone.WithLabelValues(eventID.ToHex(), stakingMetric.EventName).Set(stakingMetric.StakingPerMilestone)
		stakingAccumulated.WithLabelValues(eventID.ToHex(), stakingMetric.EventName).Set(stakingMetric.StakingAccumulated)
	}

	for eventID, ballotMetric := range metrics.BallotMetrics {
		for questionIndex, questionMetrics := range ballotMetric.QuestionMetrics {
			for answerIndex, answerMetric := range questionMetrics {
				ballotAnswersCurrent.WithLabelValues(eventID.ToHex(), ballotMetric.EventName,
					fmt.Sprint(questionIndex), fmt.Sprint(answerIndex)).Set(answerMetric.StakingCurrent)
				ballotAnswersAccumulated.WithLabelValues(eventID.ToHex(), ballotMetric.EventName,
					fmt.Sprint(questionIndex), fmt.Sprint(answerIndex)).Set(answerMetric.StakingAccumulated)
			}
		}
	}
	return nil
}
