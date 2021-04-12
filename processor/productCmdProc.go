//
// Command Processors accept commands, evaluate their validity, and determine if an event
// should be created. The way to think about it - a command is a request to do something
// that has not yet been done.
//
package processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/efvincent/archex5/commands"
	"github.com/efvincent/archex5/eventStore"
	"github.com/efvincent/archex5/eventStore/MemoryEventStore"
	"github.com/efvincent/archex5/eventStore/esErrors.go"
	"github.com/efvincent/archex5/events"
	"github.com/efvincent/archex5/models"
	validation "github.com/go-ozzo/ozzo-validation"
)

type CmdProc struct {
	es eventStore.EventStore
}

func MakeCmdProc() *CmdProc {
	return &CmdProc{MemoryEventStore.SingletonMemoryEventStore}
}

// Private utility function that gets a product aggregate from the event store given the namespace
// and sku. It first gets all the events, then it runs them through the product reducer to get
// the up to date aggregate. THIS IS WHERE YOU'D use redis or memstore, or depending on the situation
// even use a two stage local/remote cache to store aggregates so you don't have to fold over
// all the events every time to get an aggregate
func (cp CmdProc) GetProduct(ns string, sku string) (*models.ProductModel, error) {
	es, err := cp.es.GetEventRange(ns, sku, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(es) == 0 {
		return nil, errors.New(fmt.Sprintf("No such SKU %s on %s", sku, ns))
	}
	return ProductReducer(&models.ProductModel{}, es)
}

// Dispatches the product command to the appropriate command handler
func (cp CmdProc) ProcessProductCommand(cmd interface{}) error {
	switch c := cmd.(type) {
	case *commands.CreateProductCmd:
		return cp.processCreateProduct(c)
	case *commands.HeadCheckCmd:
		return cp.performHeadCheck(c)
	case *commands.SetActiveCmd:
		return cp.setProductActiveState(c)
	case *commands.UpdatePriceCmd:
		return cp.updatePrice(c)
	case *commands.UpdateProductAttributesCmd:
		return cp.updateAttribs(c)
	case *commands.UpdateProductImagesCmd:
		return cp.updateImages(c)
	default:
		return errors.New(fmt.Sprintf("Unknown command type: %v", c))
	}
}

func (cp CmdProc) updateImages(cmd *commands.UpdateProductImagesCmd) error {
	log.Printf("Updating images for %s in %s", cmd.SKU, cmd.Namespace)
	return nil
}

func (cp CmdProc) updateAttribs(cmd *commands.UpdateProductAttributesCmd) error {
	log.Printf("Updating product attributes for %s in %s", cmd.SKU, cmd.Namespace)
	return nil
}

func (cp CmdProc) updatePrice(cmd *commands.UpdatePriceCmd) error {
	product, err := cp.GetProduct(cmd.Namespace, cmd.SKU)
	if err != nil {
		return err
	}

	// There are different things we may want to check for price changes. First, is the
	// price change valid? Lets just say that negative prices are invalid and the command
	// should be rejected, and never become an event.
	// But what if the price change is technically valid but does not make sense?
	// A pack of gum who's price goes to $245,000; or a high end digital SLR for $4.99?
	// One approach would be to apply some logic here, and record an event that makes sense
	// for the new price, even if it's different from the command. But then we've
	// *lost information* - there was a request for a new price and we didn't capture that.
	//
	// Another approach might be to record the request for new price, examine it in the context
	// of past price change requests, and determine whether or not to apply that new price.
	// There's choice here too, we could split update price command into two events, one that
	// records the request, and if the analysis of past events indicates we chould change the
	// price, another event that actually changes the current active price. This way we record
	// the request, and separately record the decision to set the active price.

	if cmd.Price <= 0 {
		return errors.New(fmt.Sprintf("Invalid price %v for sku %s in %s", cmd.Price, cmd.SKU, cmd.Namespace))
	}

	// create the event
	e := events.PriceUpdated{
		Namespace: cmd.Namespace,
		SKU:       cmd.SKU,
		Price:     cmd.Price,
	}

	// serialize event
	data, err := json.Marshal(&e)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not marshal event"))
	}

	// wrap in envelope
	env := eventStore.EventEnvelope{
		EventType: events.PriceUpdatedT,
		Timestamp: time.Now().Local().UnixNano(),
		Data:      data,
	}

	// write the event see the head check event for deeper notes on checking consistency errors
	newId, err := cp.es.WriteEvent(cmd.Namespace, cmd.SKU, eventStore.EXPECTING_SEQ_NUM, product.SequenceNum, &env)
	if err != nil {
		return err
	}
	log.Printf("processor: Wrote price with sequence %v on stream %s in namespace %s ",
		newId, cmd.SKU, cmd.Namespace)
	return nil
}

func (cp CmdProc) setProductActiveState(cmd *commands.SetActiveCmd) error {
	product, err := cp.GetProduct(cmd.Namespace, cmd.SKU)
	if err != nil {
		return err
	}

	// In this design, we've decided to record the event even if it's redundenat, it
	// may still be useful information that a command was sent to activate a product
	// that was alreay active. We'll check tho and log a message if a redundant command
	// is being sent, just to show that we can take action if needed
	if product.IsActive == cmd.Active {
		var a string
		if cmd.Active {
			a = "Active"
		} else {
			a = "Inactive"
		}
		log.Printf("Redundant SetActive command. SKU %s in %s was already %s", cmd.SKU, cmd.Namespace, a)
	}

	// create the event
	e := events.ActiveStateSet{
		Namespace: cmd.Namespace,
		SKU:       cmd.SKU,
		Active:    cmd.Active,
	}

	// serialize event
	data, err := json.Marshal(&e)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not marshal set active event"))
	}

	// wrap in envelope
	env := eventStore.EventEnvelope{
		EventType: events.ActiveStateSetT,
		Timestamp: time.Now().Local().UnixNano(),
		Data:      data,
	}

	// write the event see the head check event for deeper notes on checking consistency errors
	newId, err := cp.es.WriteEvent(cmd.Namespace, cmd.SKU, eventStore.EXPECTING_SEQ_NUM, product.SequenceNum, &env)
	if err != nil {
		return err
	}
	log.Printf("processor: Wrote ActiveStateSet with sequence %v on stream %s in namespace %s ",
		newId, cmd.SKU, cmd.Namespace)
	return nil
}

// Validates that a head check can be performed. If it cannot the command fails,
// if it can per performed we simulate the headcheck and record the result as an event
func (cp CmdProc) performHeadCheck(cmd *commands.HeadCheckCmd) error {
	product, err := cp.GetProduct(cmd.Namespace, cmd.SKU)
	if err != nil {
		return err
	}

	// HEADCHECK SIMULATED!
	// perform the headcheck
	headCheckBad := (time.Now().Nanosecond())/1000%4 == 0

	// Build an event that records the headcheck and serialize it
	hce := events.HeadCheckPerformed{
		Namespace: cmd.Namespace,
		SKU:       cmd.SKU,
		Reason:    "command",
		Success:   !headCheckBad,
	}
	data, err := json.Marshal(&hce)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not marshal head check event"))
	}

	// Wrap the event in the event envelope
	e := eventStore.EventEnvelope{
		EventType: events.HeadCheckPerformedT,
		Timestamp: time.Now().Local().UnixNano(),
		Data:      data,
	}

	// Write the event into the event store, using the consistency mode that expects a specific sequence number.
	// We want to only write this event if no one "snuck in" and wrote another event after we got our product
	// aggregate from the event store, but before we were able to write our new event. If that did happen, in this
	// case, we're failing the call. Another option would be to start at the top of this function again, pull
	// the latest version of the aggregate, and try again. With some commands this makes sense, with other
	// commands, it does not. You may want to examine the last head check timestamp and see another one should
	// be done in the elapsed time. Behavior during a consistency failure is up to the domain, command, current
	// state, and business rules
	newId, err := cp.es.WriteEvent(cmd.Namespace, cmd.SKU, eventStore.EXPECTING_SEQ_NUM, product.SequenceNum, &e)
	if err != nil {
		if e, ok := err.(*esErrors.ESError); ok {
			// we specifically have an event store error, which will have information about the
			// consistency failure
			if e.ErrCode == esErrors.SEQ_NUM_EXPECTATION_FAILED {
				// we might handle the consistency failure different here
				log.Printf(fmt.Sprintf(
					"Consistency Failure writing HeadCheck event. Expected %v, actual %v", e.Expected, e.Actual))
				return e
			}
		}
		return err
	}
	log.Printf("processor: Wrote HeadCheckPerformed with sequence %v on stream %s in namespace %s ",
		newId, cmd.SKU, cmd.Namespace)
	return nil
}

func (cp CmdProc) processCreateProduct(cmd *commands.CreateProductCmd) error {
	// perform simple validation using the ozzo-validation library
	p := &cmd.Product
	if err := validation.ValidateStruct(p,
		validation.Field(&p.Namespace, validation.Required),
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.Price, validation.Required),
	); err != nil {
		return err
	}

	// If valid, make a product created event and attempt to save it to the
	// event store with the expectation that the stream does not yet exist,
	// since we're creating a new product
	pe := events.ProductCreated{
		Namespace: p.Namespace,
		SKU:       p.SKU,
		Source:    "Test",
		Product:   p,
	}
	data, err := json.Marshal(&pe)
	if err != nil {
		log.Fatal("Failed to encode JSON")
	}
	e := eventStore.EventEnvelope{
		SeqNum:    0,
		Timestamp: time.Now().UnixNano(),
		EventType: events.ProductCreatedT,
		Data:      data,
	}

	newId, err := cp.es.WriteEvent(p.Namespace, p.SKU, eventStore.NEW_STREAM, 0, &e)
	if err != nil {
		// how an error is handled depends on the command being processed, the type of
		// error, and whether or not the command processor is being run synchronously.
		// For this example, at the time this is being written, the pipeline is synchronous
		// (the API call is awaiting this result before returning to the caller), so we'll
		// return the error to the caller. In the asynchronous mode (reading from a topic
		// for example), you need to decide what to do with a failed event of this type.
		return err
	}
	log.Printf("processor: Wrote ProductCreated with sequence %v on stream %s in namespace %s ",
		newId, p.SKU, p.Namespace)

	return nil
}

// The reducer's job is to assemble the aggregate model (ProductModel) from a starting point
// and a series of events. It should be a pure function, requiring nothing that's not passed
// into the function as a formal parameter. This way, the same startingModel and set of events
// always produces the same aggregate model.
func ProductReducer(startingModel *models.ProductModel, es []eventStore.EventEnvelope) (*models.ProductModel, error) {
	cur := *startingModel
	for _, e := range es {
		switch e.EventType {
		case events.ProductCreatedT:
			// a product created event produces the product on the event, ignoring
			// the startingModel.
			var pc events.ProductCreated
			if err := json.Unmarshal(e.Data, &pc); err != nil {
				return nil, errors.New(fmt.Sprintf("Could not unmarshal ProductCreated event"))
			}
			cur = *pc.Product
			cur.SequenceNum = e.SeqNum

		case events.AttribsUpdatedT:
			var au events.AttribsUpdated
			if err := json.Unmarshal(e.Data, &au); err != nil {
				return nil, errors.New(fmt.Sprintf("Could not unmarshal AttrubutesUpdated event"))
			}
			cur.Title = au.Title
			cur.Description = au.Description
			cur.Url = au.Url
			cur.SequenceNum = e.SeqNum

		case events.ImagesUpdatedT:
			var iu events.ImagesUpdated
			if err := json.Unmarshal(e.Data, &iu); err != nil {
				return nil, errors.New(fmt.Sprintf("Could not unmarshal ImagesUpdated event"))
			}
			cur.Images = iu.Images
			cur.PrimaryImgIdx = iu.PrimaryImgIdx
			cur.SequenceNum = e.SeqNum

		case events.PriceUpdatedT:
			var pu events.PriceUpdated
			if err := json.Unmarshal(e.Data, &pu); err != nil {
				return nil, errors.New(fmt.Sprintf("Could not unmarshal ImagesUpdated event"))
			}
			cur.PriceChangeRequests = append(cur.PriceChangeRequests, models.PriceChange{
				RequestedPrice: pu.Price,
				Timestamp:      e.Timestamp,
			})
			cur.Price = pu.Price
			cur.SequenceNum = e.SeqNum

		case events.HeadCheckPerformedT:
			var hcp events.HeadCheckPerformed
			if err := json.Unmarshal(e.Data, &hcp); err != nil {
				return nil, errors.New(fmt.Sprintf("Could not unmarshal HeadCheckPerformed event"))
			}
			cur.HeadCheckOk = hcp.Success
			cur.LastHeadCheck = e.Timestamp
			cur.SequenceNum = e.SeqNum

		default:
			return nil, errors.New(fmt.Sprintf("Invalid event type in ProductReducer: %s", e.EventType))
		}
	}
	return &cur, nil
}
