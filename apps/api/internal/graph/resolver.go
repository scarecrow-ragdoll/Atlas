// FILE: apps/api/internal/graph/resolver.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Hold GraphQL resolver dependencies.
//   SCOPE: User and admin-auth service dependencies for generated gqlgen resolvers; excludes resolver method behavior.
//   DEPENDS: apps/api/internal/service.
//   LINKS: M-API / V-M-API / M-GRAPHQL-SCHEMA / V-M-GRAPHQL-SCHEMA.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Resolver - Dependency container for GraphQL resolver implementations.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added AdminAuthService dependency for backend admin auth.
// END_CHANGE_SUMMARY

package graph

import "monorepo-template/apps/api/internal/service"

// Resolver holds dependencies for GraphQL resolvers.
type Resolver struct {
	UserService      *service.UserService
	AdminAuthService *service.AdminAuthService
}
