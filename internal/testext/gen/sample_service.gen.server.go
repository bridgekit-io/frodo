// Code generated by Frodo - DO NOT EDIT.
//
//	Timestamp: Fri, 15 Nov 2024 12:49:04 EST
//	Source:    sample_service.go
//	Generator: https://github.com/bridgekitio/frodo
package testext

import (
	"context"

	"github.com/bridgekitio/frodo/fail"
	"github.com/bridgekitio/frodo/internal/testext"
	"github.com/bridgekitio/frodo/services"
)

// SampleServiceServer accepts your "real" SampleService instance (the thing that really does
// the work), and returns a set of endpoint routes which allow this service to be consumed
// via the gateways/listeners you configure in main().
//
//	// Example
//	serviceHandler := testext.SampleServiceHandler{ /* set up to your liking */ }
//	server := services.New(
//		services.Listen(apis.NewGateway()),
//		services.Register(testextgen.SampleServiceServer(serviceHandler)),
//	)
//	server.Listen()
//
// From there, you can add middleware, event sourcing support and more. Look at the frodo
// documentation for more details/examples on how to make your service production ready.
func SampleServiceServer(handler testext.SampleService, middleware ...services.MiddlewareFunc) *services.Service {
	middlewareFuncs := services.MiddlewareFuncs(middleware)

	return &services.Service{
		Name:    "SampleService",
		Version: "0.0.1",
		Handler: handler,
		Endpoints: []services.Endpoint{

			{
				ServiceName: "SampleService",
				Name:        "Authorization",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Authorization(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Authorization",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Chain1",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Chain1(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Chain1",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Chain1GroupFooBar",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Chain1GroupFooBar(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Chain1GroupFooBar",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.Chain1",
						PathParams:  []string{},
						Group:       "FooBar",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Chain1GroupStar",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Chain1GroupStar(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Chain1GroupStar",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.Chain1",
						PathParams:  []string{},
						Group:       "*",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Chain2",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Chain2(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Chain2",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.Chain1",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Chain2OnError",
				NewInput:    func() services.StructPointer { return &testext.FailAlwaysErrorRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.FailAlwaysErrorRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Chain2OnError(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.Chain2:Error",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Chain2OnSuccess",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Chain2OnSuccess(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.Chain2",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "ComplexValues",
				NewInput:    func() services.StructPointer { return &testext.SampleComplexRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleComplexRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.ComplexValues(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.ComplexValues",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "ComplexValuesPath",
				NewInput:    func() services.StructPointer { return &testext.SampleComplexRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleComplexRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.ComplexValuesPath(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/complex/values/{InUser.ID}/{InUser.Name}/woot",
						PathParams: []string{
							"InUser.ID",

							"InUser.Name",
						},
						Group:  "",
						Status: 200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "CustomRoute",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.CustomRoute(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/custom/route/1/{ID}/{Text}",
						PathParams: []string{
							"ID",

							"Text",
						},
						Group:  "",
						Status: 202,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "CustomRouteBody",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.CustomRouteBody(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "PUT",
						Path:        "/v2/custom/route/3/{ID}",
						PathParams: []string{
							"ID",
						},
						Group:  "",
						Status: 201,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "CustomRouteQuery",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.CustomRouteQuery(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/custom/route/2/{ID}",
						PathParams: []string{
							"ID",
						},
						Group:  "",
						Status: 202,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Defaults",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Defaults(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Defaults",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Download",
				NewInput:    func() services.StructPointer { return &testext.SampleDownloadRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleDownloadRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Download(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/download",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "DownloadResumable",
				NewInput:    func() services.StructPointer { return &testext.SampleDownloadRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleDownloadRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.DownloadResumable(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/download/resumable",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Fail4XX",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Fail4XX(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Fail4XX",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Fail5XX",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Fail5XX(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Fail5XX",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "FailAlways",
				NewInput:    func() services.StructPointer { return &testext.FailAlwaysRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.FailAlwaysRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.FailAlways(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.FailAlways",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "ListenerA",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.ListenerA(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/ListenerA/Woot",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.TriggerUpperCase",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "ListenerB",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.ListenerB(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.TriggerUpperCase",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.TriggerLowerCase",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.TriggerFailure",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.ListenerA",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "OtherService.SpaceOut",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "OmitMe",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.OmitMe(ctx, typedReq)
				}),
				Roles:  []string{},
				Routes: []services.EndpointRoute{},
			},

			{
				ServiceName: "SampleService",
				Name:        "OnFailAlways",
				NewInput:    func() services.StructPointer { return &testext.FailAlwaysErrorRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.FailAlwaysErrorRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.OnFailAlways(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.OnFailAlways",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.FailAlways:Error",
						PathParams:  []string{},
						Group:       "",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Panic",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Panic(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Panic",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Redirect",
				NewInput:    func() services.StructPointer { return &testext.SampleRedirectRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRedirectRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Redirect(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/redirect",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "SecureWithRoles",
				NewInput:    func() services.StructPointer { return &testext.SampleSecurityRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleSecurityRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.SecureWithRoles(ctx, typedReq)
				}),
				Roles: []string{
					"admin.write",
					"user.{ID}.write",
					"user.{User.ID}.admin",
					"junk.{NotReal}.crap",
				},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.SecureWithRoles",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "SecureWithRolesAliased",
				NewInput:    func() services.StructPointer { return &testext.SampleSecurityRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleSecurityRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.SecureWithRolesAliased(ctx, typedReq)
				}),
				Roles: []string{
					"admin.write",
					"user.{FancyID}.write",
					"user.{User.FancyID}.admin",
					"junk.{NotReal}.crap",
				},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.SecureWithRolesAliased",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Sleep",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Sleep(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Sleep",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "TriggerFailure",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.TriggerFailure(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.TriggerFailure",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "TriggerLowerCase",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.TriggerLowerCase(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.TriggerLowerCase",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "TriggerUpperCase",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.TriggerUpperCase(ctx, typedReq)
				}),
				Roles: []string{},
				Routes: []services.EndpointRoute{
					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/Upper/Case/WootyAndTheBlowfish",
						PathParams:  []string{},
						Group:       "",
						Status:      200,
					},
				},
			},
		},
	}
}
