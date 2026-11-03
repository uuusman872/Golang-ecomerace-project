package server

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/uuusman872/learning-go-shop/graph"
	"github.com/uuusman872/learning-go-shop/graph/resolver"
	"github.com/uuusman872/learning-go-shop/internal/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

func (s *Server) createGraphQLHandler() *handler.Server {
	rvr := resolver.NewResolver(s.authServices, s.userService, s.productService, s.cartService, s.orderService)
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: rvr})
	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})
	//http.Handle("/graphql", srv)
	return srv
}

func (s *Server) graphqlHandler() gin.HandlerFunc {
	h := s.createGraphQLHandler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground (Public)", "/graphql/public")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// GraphQL playground handler for protected endpoint
func (s *Server) playgroundProtectedHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground (Protected)", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// optionalAuthMiddleware parses a Bearer token when present but does not reject unauthenticated requests.
func (s *Server) optionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], s.config.JWT.Secret)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func (s *Server) graphqlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if userID, ok := c.Get("user_id"); ok {
			ctx = context.WithValue(ctx, utils.UserIDKey, userID)
		}
		if userEmail, ok := c.Get("user_email"); ok {
			ctx = context.WithValue(ctx, utils.UserEmailKey, userEmail)
		}
		if userRole, ok := c.Get("user_role"); ok {
			ctx = context.WithValue(ctx, utils.UserRoleKey, userRole)
		}
		ctx = context.WithValue(ctx, utils.GinContextKey, c)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
