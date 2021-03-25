package main

import (
	"log"
	"os"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

func main() {
	var (
		pactDir             = os.Getenv("PACT_DIR")
		pactBrokerURL       = os.Getenv("PACT_BROKER_URL")
		pactBrokerUsername  = os.Getenv("PACT_BROKER_USERNAME")
		pactBrokerPassword  = os.Getenv("PACT_BROKER_PASSWORD")
		pactConsumerVersion = os.Getenv("PACT_CONSUMER_VERSION")
		pactTag             = os.Getenv("PACT_TAG")
	)

	log.Println("Publishing Pact files to broker", pactDir, pactBrokerURL)

	p := &dsl.Publisher{}
	r := types.PublishRequest{
		PactURLs: []string{
			pactDir + "/sirius-workflow-sirius.json",
		},
		PactBroker:      pactBrokerURL,
		ConsumerVersion: pactConsumerVersion,
		Tags:            []string{pactTag},
		BrokerUsername:  pactBrokerUsername,
		BrokerPassword:  pactBrokerPassword,
	}

	if err := p.Publish(r); err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
}
