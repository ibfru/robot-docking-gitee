package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"

	"community-robot-lib/mq"
	"community-robot-lib/utils"
)

// groupConsumer represents a Sarama consumer group consumer
type groupConsumer struct {
	kOpts   mq.Options
	handler mq.Handler
	subOpts mq.SubscribeOptions

	notifyReady func()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (gc *groupConsumer) Setup(sarama.ConsumerGroupSession) error {
	gc.notifyReady()

	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (gc *groupConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a groupConsumer loop of ConsumerGroupClaim's Messages().
func (gc *groupConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	handle := gc.genHanler(session)

	for {
		select {
		case message := <-claim.Messages():
			handle(message)

			if gc.subOpts.AutoAck {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

func (gc *groupConsumer) genHanler(session sarama.ConsumerGroupSession) func(*sarama.ConsumerMessage) {
	handler := gc.handler
	if handler == nil {
		handler = func(event mq.Event) error {
			return nil
		}
	}

	log := gc.kOpts.Log

	eh := gc.kOpts.ErrorHandler
	if eh == nil {
		eh = func(e mq.Event) error {
			log.Error(e.Error())

			return nil
		}
	}

	unmarshal := gc.kOpts.Codec.Unmarshal

	return func(msg *sarama.ConsumerMessage) {
		ke := &event{
			km:   msg,
			m:    new(mq.Message),
			sess: session,
		}

		if err := unmarshal(msg.Value, ke.m); err != nil {
			ke.err = fmt.Errorf("unmarshal msg failed, err: %v", err)
			ke.m.Body = msg.Value

			if err := eh(ke); err != nil {
				log.Error(err)
			}

			return
		}

		if err := handler(ke); err != nil {
			ke.err = fmt.Errorf("handle event, err: %v", err)

			if err := eh(ke); err != nil {
				log.Error(err)
			}
		}
	}
}

type subscriber struct {
	cli sarama.Client
	cg  sarama.ConsumerGroup

	t  string
	gc groupConsumer

	once  sync.Once
	ready chan struct{}
	stop  chan struct{}
	done  chan struct{}

	cancel context.CancelFunc
}

func newSubscriber(
	topic string,
	cli sarama.Client, cg sarama.ConsumerGroup,
	gc groupConsumer,

) (s *subscriber) {
	s = &subscriber{
		t:   topic,
		cli: cli,
		cg:  cg,
		gc:  gc,

		ready: make(chan struct{}),
		stop:  make(chan struct{}),
		done:  make(chan struct{}),
	}

	s.gc.notifyReady = s.notifyReady

	return
}

func (s *subscriber) Options() mq.SubscribeOptions {
	return s.gc.subOpts
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	mErr := utils.MultiError{}

	s.once.Do(func() {
		s.cancel()

		// wait
		<-s.done

		mErr.AddError(s.cg.Close())

		mErr.AddError(s.cli.Close())
	})

	return mErr.Err()
}

func (s *subscriber) start() error {
	f := func(ctx context.Context) {
		defer close(s.done)

		log := s.gc.kOpts.Log
		topic := []string{s.t}

		// doesn't support re-consuming because of server-side rebalance
		if err := s.cg.Consume(ctx, topic, &s.gc); err != nil {
			log.Errorf("Consume err: %s", err.Error())
			close(s.stop)

			return
		}

		if err := ctx.Err(); err != nil {
			log.Infof("exit by unsubscribing, err:%s", err.Error())
		} else {
			log.Fatal("maybe, server-side rebalance happens. Restart to fix it!")
		}
	}

	ctx, cancel := context.WithCancel(s.gc.subOpts.Context)

	go f(ctx)

	select {
	case <-s.ready:
		s.cancel = cancel

		return nil

	case <-s.stop:
		cancel()

		return fmt.Errorf("start failed")
	}
}

func (s *subscriber) notifyReady() {
	close(s.ready)
}
