package rest

func (s *echoServer) setupRouter() {
	apiV1 := s.echo.Group("/api/v1")
	{
		orderGroup := apiV1.Group("/orders")
		{
			orderGroup.POST("", s.handler.Order().Create)
			orderGroup.GET("", s.handler.Order().List)
			orderGroup.GET("/:id", s.handler.Order().Get)
			orderGroup.POST("/:id/cancel", s.handler.Order().Cancel)
		}
	}
}
