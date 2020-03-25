// Package cloudwatch cleans up old LogGroups (empty and olde than 90 days and old
// LogStreams (last event timestamp is over 30 days old or if the logstream
// is empty and has been created over 30 days ago) from AWS Cloudwatch
package cloudwatch

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/rs/zerolog/log"
)

// CwProcessor holds the cloudwatch-related actions
type CwProcessor struct {
	svc    *cloudwatchlogs.CloudWatchLogs
	lsChan chan map[*string]*string
	lgChan chan *string
	daysLs int
	daysLg int
	lgName string
}

const timeDiv = 1000
const yearsLStream = 0
const monthsLStream = 0
const yearsLGroup = 0
const monthLGroup = 0

// NewCwProcessor creates a new instance of CwProcessor containing an already
// initialized cloudwatchlogs client
func NewCwProcessor(daysLs int, daysLg int, lgName string) *CwProcessor {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &CwProcessor{svc: cloudwatchlogs.New(sess),
		lsChan: make(chan map[*string]*string),
		lgChan: make(chan *string),
		daysLs: daysLs, daysLg: daysLg, lgName: lgName}
}

// ProcessLogStreams add the logstreams of the given loggroup to the lsChan
// channel if their last event timestamp is over 30 days old or if the logstream
// is empty and has been created over 30 days ago
func (p *CwProcessor) ProcessLogStreams(groupName *string) (int, error) {
	daysLs := -(p.daysLs)
	nbrLogStreams := 0
	err := p.svc.DescribeLogStreamsPages(&cloudwatchlogs.DescribeLogStreamsInput{LogGroupName: groupName},
		func(ls *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
			nbrLogStreams = len(ls.LogStreams)
			for _, stream := range ls.LogStreams {
				lastTs := *stream.CreationTime / timeDiv
				if stream.LastEventTimestamp != nil {
					lastTs = *stream.LastEventTimestamp / timeDiv
				}
				if lastTs < time.Now().AddDate(yearsLStream, monthsLStream, daysLs).Unix() {
					p.lsChan <- map[*string]*string{groupName: stream.LogStreamName}
				}
			}
			return !lastPage
		})

	return nbrLogStreams, err
}

// ProcessLogGroups does a cleanup of all LogGroups using the ProcessLogStreams
// function and add the Log Groups that have no logStream and are older than
// 90 days to the lgChan channel for deletion
func (p *CwProcessor) ProcessLogGroups() {
	daysLg := -(p.daysLg)
	err := p.svc.DescribeLogGroupsPages(&cloudwatchlogs.DescribeLogGroupsInput{},
		func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
			if p.lgName == "all" {
				for _, lg := range page.LogGroups {
					nbrLs, err := p.ProcessLogStreams(lg.LogGroupName)
					if err != nil {
						log.Error().Msgf("Getting Log group %s streams: %v", *lg.Arn, err)
					} else if nbrLs == 0 && (*lg.CreationTime/timeDiv) < time.Now().AddDate(yearsLGroup, monthLGroup, daysLg).Unix() {
						p.lgChan <- lg.LogGroupName
					}
				}
				return !lastPage
			}

			nbrLs, err := p.ProcessLogStreams(&p.lgName)
			if err != nil {
				log.Error().Msgf("Getting Log group %s streams: %v", p.lgName, err)
			} else if nbrLs == 0 {
				p.lgChan <- &p.lgName
			}
			return !lastPage
		})

	if err != nil {
		log.Fatal().Msgf("DescribeLogGroups returned: %v", err)
	}

	close(p.lsChan)
	close(p.lgChan)
}

// CleanupLogStreams delete the LogGroups that exist in the channel lsChan
func (p *CwProcessor) CleanupLogStreams() {
	for elt := range p.lsChan {
		for k, v := range elt {
			log.Info().Msgf("Deleting LogStream: %s from LogGroup: %s", *v, *k)

			if _, err := p.svc.DeleteLogStream(&cloudwatchlogs.DeleteLogStreamInput{LogGroupName: k,
				LogStreamName: v}); err != nil {
				log.Error().Msgf("Deleting LogGroup: %s, LogStream: %s --> %s", *k, *v, err.Error())
			}
		}
	}
}

// CleanupLogGroups delete the LogGroups that exist in the channel lgChan
func (p *CwProcessor) CleanupLogGroups() {
	for elt := range p.lgChan {
		log.Info().Msgf("Deleting LogGroup: %s", *elt)

		if _, err := p.svc.DeleteLogGroup(&cloudwatchlogs.DeleteLogGroupInput{LogGroupName: elt}); err != nil {
			log.Error().Msgf("Deleting LogGroup: %s --> %s", *elt, err.Error())
		}
	}
}
