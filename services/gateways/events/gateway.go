package events

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/bridgekit-io/frodo/codec"
	"github.com/bridgekit-io/frodo/eventsource"
	"github.com/bridgekit-io/frodo/eventsource/local"
	"github.com/bridgekit-io/frodo/fail"
	"github.com/bridgekit-io/frodo/internal/wait"
	"github.com/bridgekit-io/frodo/metadata"
	"github.com/bridgekit-io/frodo/services"
)

// NewGateway creates an event-sourced gateway that executes service methods based on event subscriptions.
// By default, events are sourced using local.Broker(). This means that events will only be available to
// services running in the services.Server that this gateway is added to. You can provide a broker to a
// different event source like NATS JetStream using the WithBroker() option.
func NewGateway(options ...GatewayOption) *Gateway {
	jsonEncoder := codec.JSONEncoder{}
	jsonDecoder := codec.JSONDecoder{Loose: true}
	gw := Gateway{
		encoder:        jsonEncoder,
		decoder:        jsonDecoder,
		valueEncoder:   jsonEncoder,
		valueDecoder:   jsonDecoder,
		broker:         local.Broker(),
		listening:      &sync.WaitGroup{},
		activeRequests: &sync.WaitGroup{},
		errorListener: func(route metadata.EndpointRoute, err error) {
			// I'm leaving this commented out rather than deleting it to document that
			// it was an intentional decision to make the default listener do nothing. It's
			// better to have the user opt into whatever logging/handling they want rather
			// than forcing a random log line that doesn't conform to their standard log/slog
			// format. If they're following remotely sane practices, they'll have logging
			// middleware on the service that prints the failure anyway.
			//
			// log.Printf("[events error] [%s] %v\n", route.QualifiedName(), err)
		},
	}
	for _, option := range options {
		option(&gw)
	}
	return &gw
}

// Gateway encapsulates the logic to invoke service operations based on event sourcing. You
// should not create one of these yourself - use the NewGateway() constructor instead.
type Gateway struct {
	encoder        codec.Encoder
	decoder        codec.Decoder
	valueEncoder   codec.ValueEncoder
	valueDecoder   codec.ValueDecoder
	broker         eventsource.Broker
	errorListener  ErrorListener
	routes         []*route
	listening      *sync.WaitGroup
	activeRequests *sync.WaitGroup
}

// Type returns "EVENTS" to indicate the tagging value for this gateway.
func (gw *Gateway) Type() services.GatewayType {
	return services.GatewayTypeEvents
}

// Register adds the given service endpoint to the routing rules for this gateway. You will
// not invoke this yourself! The services.Server will utilize this as necessary.
func (gw *Gateway) Register(endpoint services.Endpoint, endpointRoute services.EndpointRoute) {
	if endpointRoute.GatewayType != services.GatewayTypeEvents {
		return
	}

	// We use the fully qualified endpoint name as the group to create a "consumer group".
	// This prevents more than one instance of the service method handling the same event.
	// Here's an example - let's say we have an operation OrderService.PlaceOrder(). When
	// that successfully occurs, we want 3 things to happen:
	//
	// * Send a confirmation email
	// * Send a coupon if it's your first purchase ever
	// * Charge the hold on the credit card
	//
	// To do this, you probably set up some doc options for "ON OrderService.PlaceOrder"
	// on each of those 3 methods. Two of them are on the EmailService and one is on the
	// TransactionService.
	//
	// Even if we have 4 email service instances and 8 transaction service instances, we
	// only want each of those items to occur one time for a given order. As a result,
	// when we register routes for these handlers, we'll use the fully qualified endpoint
	// name:
	//
	// * EmailService.SendOrderConfirmation
	// * EmailService.SendFirstOrderCoupon
	// * TransactionService.ChargeHold
	//
	// We can't just use the service name because in this case we have 2 handlers that need
	// to fire from the email service. If we just used the service name, only one of them
	// would ever trigger. No good. We can't just use the method name because you might also
	// have a SendOrderConfirmation method on the SMSService as well, so they could be stealing
	// each other's events.
	//
	// By using the fully qualified name, we ensure that no matter how much redundancy we have in
	// our system, every handler fires at most once.
	//
	// Lastly, we're not going to actually perform these subscriptions to NATS/Redis/etc yet. The
	// broker might not have been started up yet, so we just want to construct and capture the
	// handler information for what we *will* subscribe to once Listen() is fired on this gateway.
	//
	//
	var consumerGroup string
	switch endpointRoute.Group {
	case "":
		// The parser saw something like "ON FooService.Bar" where no group was specified. Default to the
		// group described in the giant fuck-off comment above so that only one instance of the endpoint handles it.
		consumerGroup = endpoint.QualifiedName()
	case "*":
		// The parser saw something like "ON FooService.Bar GROUP *" which indicates that you want ALL instances of the
		// servers running this endpoint to get a copy of this message. You want a more traditional Pub/Sub, so we just
		// don't specify a group. On Listen(), this will result in calling broker.Subscribe() instead of broker.SubscribeGroup().
		consumerGroup = ""
	default:
		// The parser saw something like "ON FooService.Bar GROUP Hug", so we're being asked to register this endpoint
		// handler as a consumer in the "Hug" group/pool.
		consumerGroup = endpointRoute.Group
	}

	gw.routes = append(gw.routes, &route{
		key:     endpointRoute.Path,
		group:   consumerGroup,
		handler: gw.toStreamHandler(endpoint, endpointRoute),
	})
}

func (gw *Gateway) toStreamHandler(endpoint services.Endpoint, route services.EndpointRoute) eventsource.EventHandlerFunc {
	// decodeRequestBody uses the (likely JSON) decoder to overlay the response values of the source method
	// onto the request of the handling method.
	decodeRequestBody := func(event message, serviceRequest services.StructPointer) error {
		valueReader := strings.NewReader(string(event.Values))
		if err := gw.decoder.Decode(valueReader, &serviceRequest); err != nil {
			return fmt.Errorf("error decoding source event request values: %w", err)
		}
		return nil
	}

	// decodeErrorFields only concerns itself with events for source methods that returned a
	// non-nil error (i.e. they failed). It tries to decode the error message/status onto the
	// request structure using some fairly common naming conventions.
	decodeErrorFields := func(event message, serviceRequest services.StructPointer) error {
		if !event.ErrorPresent() {
			return nil
		}

		// Allow you to define your error as simply the field 'Error string' on your request or an error struct
		// that has either an Error or Message string attribute. The status follows the naming conventions of
		// status codes in the 'fail' package - any will work. We don't know what names/shape they'll use in their
		// request struct, but we'll try to accommodate these fairly common field layouts.
		errorStatus := strconv.FormatInt(int64(event.ErrorStatus), 10)
		errorValues := url.Values{
			// Flat attributes directly on the request struct.
			"Error":               []string{event.ErrorMessage},
			"ErrorCode":           []string{errorStatus},
			"ErrorStatus":         []string{errorStatus},
			"ErrorStatusCode":     []string{errorStatus},
			"ErrorHTTPStatusCode": []string{errorStatus},

			// You have an 'Error' field which is a struct that contains these values.
			"Error.Error":      []string{event.ErrorMessage},
			"Error.Message":    []string{event.ErrorMessage},
			"Error.Code":       []string{errorStatus},
			"Error.Status":     []string{errorStatus},
			"Error.StatusCode": []string{errorStatus},
		}

		if err := gw.valueDecoder.DecodeValues(errorValues, &serviceRequest); err != nil {
			return fmt.Errorf("error decoding source event error: %w", err)
		}
		return nil
	}

	return func(ctx context.Context, msg *eventsource.EventMessage) error {
		gw.activeRequests.Add(1)
		defer gw.activeRequests.Done()

		event := message{}
		serviceRequest := endpoint.NewInput()

		// Take the broker's message and read in the service event 'message' data from it.
		if err := gw.decoder.Decode(bytes.NewBuffer(msg.Payload), &event); err != nil {
			gw.errorListener(event.Route, fmt.Errorf("event decode error: %w", err))
			return nil
		}

		// If the source method was a failure/error rather than a success, also decode the
		// error details onto the request struct if they exist. For success events, this will simply nop.
		if err := decodeErrorFields(event, serviceRequest); err != nil {
			gw.errorListener(event.Route, err)
			return nil
		}

		// The message contains the raw encoded bytes for the response of the service
		// method that triggered the event. Overlay that data on this handler's input.
		if err := decodeRequestBody(event, serviceRequest); err != nil {
			gw.errorListener(event.Route, err)
			return nil
		}

		// We want to make sure that the metadata context is restored from the invocation
		// that triggered this originally. For example, we want to make sure that this
		// event handler uses the same request id as the HTTP/API request that originally
		// triggered this. It should also have the same authorization info and values, etc.
		ctx = metadata.Decode(ctx, event.Metadata)

		// This is a new invocation so the route should indicate THIS function, not the
		// thing that triggered us to execute.
		ctx = metadata.WithRoute(ctx, metadata.EndpointRoute{
			ServiceName: endpoint.ServiceName,
			Name:        endpoint.Name,
			Type:        gw.Type().String(),
			Method:      route.Method,
			Path:        route.Path,
			Group:       route.Group,
			Status:      200, // we don't have a doc option for setting this on event routes, so use sane default.
		})

		if _, err := endpoint.Handler(ctx, serviceRequest); err != nil {
			gw.errorListener(event.Route, err)
			return err
		}
		return nil
	}
}

// Middleware returns the middleware functions that ALL server routes should include in order
// to make sure that this gateway actually works. For instance, one of the middleware functions
// publishes the service operation's success/failure to the event source/stream. This happens
// regardless of whether the operation was invoked through the API gateway or the event one.
// Basically, these are agnostic of the gateway type and will be added to ALL gateway routes, not
// just the event gateway.
func (gw *Gateway) Middleware() services.MiddlewareFuncs {
	return services.MiddlewareFuncs{
		publishMiddleware(gw.broker, gw.encoder, gw.valueEncoder, gw.errorListener),
	}
}

// Listen causes the gateway to start subscribing/listening for events from the broker. This
// will block until we're told to stop by calling Shutdown().
func (gw *Gateway) Listen(ctx context.Context) error {
	errs, _ := fail.NewGroup(context.Background())

	for _, gatewayRoute := range gw.routes {
		// Make a separate variable from the loop variable to avoid this:
		// https://github.com/golang/go/discussions/56010
		r := gatewayRoute

		errs.Go(func() (err error) {
			switch r.group {
			case "":
				// The interface had "ON FooService.Bar GROUP *"
				r.subs, err = gw.broker.Subscribe(ctx, r.key, r.handler)
				return err
			default:
				// The interface had "ON FooService.Bar" without specifying a group to get the default grouping behavior.
				r.subs, err = gw.broker.SubscribeGroup(ctx, r.key, r.group, r.handler)
				return err
			}
		})
	}

	if err := errs.Wait(); err != nil {
		return fmt.Errorf("event gateway error: listen: %w", err)
	}

	gw.listening.Add(1)
	gw.listening.Wait()
	return nil
}

// Shutdown gracefully stops the event gateway. It will allow all of the in-progress requests
// to finish up before doing so. You can provide a deadline to the context parameter to limit
// how much time you're willing to give them before shutting down anyway.
func (gw *Gateway) Shutdown(ctx context.Context) error {
	errs, _ := fail.NewGroup(ctx)
	for _, r := range gw.routes {
		if r.subs != nil {
			errs.Go(r.subs.Close)
		}
	}

	// Make sure that we have stopped listening for all of our registered events.
	if err := errs.Wait(); err != nil {
		gw.listening.Done()
		return fmt.Errorf("event gateway error: shutdown: %w", err)
	}
	gw.listening.Done()

	// Any in-progress requests should get an opportunity to finish before
	// we consider shutdown 100% complete. They have until either the
	// context's deadline/cancellation is reached or the process receives
	// another SIGINT/SIGTERM signal. We'll exit once one of those 3 things happens.
	wait.ContextOrGroupOrInterrupt(ctx, gw.activeRequests)
	return nil
}

type route struct {
	key     string
	group   string
	handler eventsource.EventHandlerFunc
	subs    eventsource.Subscription
}

// GatewayOption defines a functional parameter that you can use to set up an event gateway.
type GatewayOption func(gw *Gateway)

// WithBroker defines the broker that the gateway will use to publish and listen for events. By
// default, the gateway will use a local broker that can only broadcast events to other services
// running inside the same services.Server instance.
func WithBroker(broker eventsource.Broker) GatewayOption {
	return func(gw *Gateway) {
		gw.broker = broker
	}
}

// WithEncoding allows you to customize how to marshal events to/from the broker. By default, the
// gateway will use standard library JSON encoding.
func WithEncoding(encoder codec.Encoder, decoder codec.Decoder) GatewayOption {
	return func(gw *Gateway) {
		gw.encoder = encoder
		gw.decoder = decoder
	}
}

// WithErrorListener sets a custom callback function that is invoked any time we encounter an error
// publishing an event, receiving an event, or executing a service handler. These are all invoked
// asynchronously, so this is the only way you can perform any custom error handling in those cases.
//
// This is fired IN ADDITION to any "On Service.Method:Error" routes you've set up. It's a way for you
// to provide some general purpose recovery/logging/whatever on every failed async call.
func WithErrorListener(listener ErrorListener) GatewayOption {
	return func(gw *Gateway) {
		gw.errorListener = listener
	}
}

// ErrorListener defines a custom callback that you can use to listen for any error generated by an async event
// gateway invocation failure.
type ErrorListener func(route metadata.EndpointRoute, err error)
