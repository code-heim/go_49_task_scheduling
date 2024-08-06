package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func main() {
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		panic(err)
	}

	// add a simple job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			30*time.Second,
		),
		gocron.NewTask(
			func(a string) {
				log.Println(a)
			},
			"Every 30 seconds",
		),
		gocron.WithName("job: Every 30 seconds"),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(
				func(jobID uuid.UUID, jobName string) {
					log.Printf("Job starting: %s, %s \n", jobID.String(), jobName)
				},
			),
			gocron.AfterJobRuns(
				func(jobID uuid.UUID, jobName string) {
					log.Printf("Job completed: %s, %s \n", jobID.String(), jobName)
				},
			),
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					log.Printf("Job had an error: %s, %s \n", jobID.String(), jobName)
				},
			),
		),
	)
	if err != nil {
		// handle error
	}

	log.Println(j.ID())

	// add a Cron job to the scheduler
	s.NewJob(
		gocron.CronJob(
			"*/10 * * * *",
			false,
		),
		gocron.NewTask(
			func(a string) {
				log.Println(a)
			},
			"Cronjob: Every 10 mins",
		),
		gocron.WithName("CronJob: Every 10 mins"),
	)

	// add a Daily job to the scheduler
	s.NewJob(
		gocron.DailyJob(
			1, // Runs every day
			gocron.NewAtTimes(
				gocron.NewAtTime(23, 10, 00),
				gocron.NewAtTime(05, 30, 00),
			),
		),
		gocron.NewTask(
			func(a string, b string) {
				log.Println(a, b)
			},
			"Dailyjob", "Runs everyday",
		),
		gocron.WithName("Dailyjob"),
	)

	// add a One time job to the scheduler
	s.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(time.Now().Add(10*time.Second)),
		),
		gocron.NewTask(
			func(a string) {
				log.Println(a)
			},
			"One time job",
		),
	)

	// add a Random duration job to the scheduler
	s.NewJob(
		gocron.DurationRandomJob(2*time.Minute, 4*time.Minute),
		gocron.NewTask(
			func(a string) {
				log.Println(a)
			},
			"Random job",
		),
	)

	// Start the scheduler
	s.Start()

	// Set up a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\nInterrupt signal received. Exiting...")
		_ = s.Shutdown()
		os.Exit(0)
	}()

	for {

	}
}
