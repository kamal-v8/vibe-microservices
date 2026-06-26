# Graph Report - .  (2026-06-23)

## Corpus Check
- Corpus is ~10,050 words - fits in a single context window. You may not need a graph.

## Summary
- 253 nodes · 341 edges · 36 communities (18 shown, 18 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 56 edges (avg confidence: 0.67)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Post Service Core|Post Service Core]]
- [[_COMMUNITY_Post Service API|Post Service API]]
- [[_COMMUNITY_Notification Service Core|Notification Service Core]]
- [[_COMMUNITY_Pulse Kubernetes Manifests|Pulse Kubernetes Manifests]]
- [[_COMMUNITY_User Service Handlers|User Service Handlers]]
- [[_COMMUNITY_User Service Data Layer|User Service Data Layer]]
- [[_COMMUNITY_User Service Entrypoint|User Service Entrypoint]]
- [[_COMMUNITY_Notification NPM Package|Notification NPM Package]]
- [[_COMMUNITY_Shared Infrastructure Config|Shared Infrastructure Config]]
- [[_COMMUNITY_Postgres Execution Functions|Postgres Execution Functions]]
- [[_COMMUNITY_Application Configuration|Application Configuration]]
- [[_COMMUNITY_Event Subscriber Logic|Event Subscriber Logic]]
- [[_COMMUNITY_Postgres Kubernetes Stack|Postgres Kubernetes Stack]]
- [[_COMMUNITY_Postgres SQL Tables|Postgres SQL Tables]]
- [[_COMMUNITY_User Service Models|User Service Models]]
- [[_COMMUNITY_Postgres Init DB|Postgres Init DB]]
- [[_COMMUNITY_Notification CI Workflow|Notification CI Workflow]]
- [[_COMMUNITY_Post Service CI Workflow|Post Service CI Workflow]]
- [[_COMMUNITY_User Service CI Workflow|User Service CI Workflow]]
- [[_COMMUNITY_Kubernetes Kustomization|Kubernetes Kustomization]]
- [[_COMMUNITY_Notification Service HPA|Notification Service HPA]]
- [[_COMMUNITY_Notification Service API Route|Notification Service API Route]]
- [[_COMMUNITY_Go Module Path|Go Module Path]]
- [[_COMMUNITY_Post Service HPA|Post Service HPA]]
- [[_COMMUNITY_Python Async Dependencies|Python Async Dependencies]]
- [[_COMMUNITY_Python FastAPI Dependencies|Python FastAPI Dependencies]]
- [[_COMMUNITY_Python HTTPX Dependency|Python HTTPX Dependency]]
- [[_COMMUNITY_Python Pydantic Dependency|Python Pydantic Dependency]]
- [[_COMMUNITY_Python Uvicorn Dependency|Python Uvicorn Dependency]]
- [[_COMMUNITY_Post Service API Route|Post Service API Route]]
- [[_COMMUNITY_Gemini Chatbot Rules|Gemini Chatbot Rules]]
- [[_COMMUNITY_Redis Kubernetes Manifest|Redis Kubernetes Manifest]]
- [[_COMMUNITY_Kind Cluster Configuration|Kind Cluster Configuration]]

## God Nodes (most connected - your core abstractions)
1. `Database` - 15 edges
2. `APIResponse` - 14 edges
3. `EventPublisher` - 13 edges
4. `UserServiceClient` - 12 edges
5. `Request` - 10 edges
6. `create_post()` - 10 edges
7. `UserHandler` - 10 edges
8. `JSONFormatter` - 9 edges
9. `HealthResponse` - 9 edges
10. `Context` - 9 edges

## Surprising Connections (you probably didn't know these)
- `Post Service` --semantically_similar_to--> `Post Service Deployment`  [INFERRED] [semantically similar]
  docker-compose.yaml → k8s-manifest.yaml
- `User Service` --semantically_similar_to--> `User Service`  [INFERRED] [semantically similar]
  README.md → docker-compose.yaml
- `Post Service` --semantically_similar_to--> `Post Service`  [INFERRED] [semantically similar]
  README.md → docker-compose.yaml
- `Notification Service` --semantically_similar_to--> `Notification Service`  [INFERRED] [semantically similar]
  README.md → docker-compose.yaml
- `PostgreSQL` --semantically_similar_to--> `Postgres Service`  [INFERRED] [semantically similar]
  README.md → docker-compose.yaml

## Import Cycles
- 1-file cycle: `services/post-service/app/main.py -> services/post-service/app/main.py`

## Hyperedges (group relationships)
- **Pulse Microservices Platform** — pulse_readme_user_service, pulse_readme_post_service, pulse_readme_notification_service, pulse_readme_postgres, pulse_readme_redis [INFERRED 0.95]
- **Notification Service K8s Stack** — notification_service_deployment_main, notification_service_hpa_main, notification_service_service_main [INFERRED 0.95]
- **Post Service K8s Stack** — post_service_deployment_main, post_service_hpa_main, post_service_service_main [INFERRED 0.95]
- **User Service Components** — user_service_deployment_user_service, user_service_hpa_user_service_hpa, user_service_service_user_service [EXTRACTED 1.00]
- **Redis Components** — redis_deployment_redis, redis_pvc_redis_pvc, redis_service_redis [EXTRACTED 1.00]
- **Postgres Components** — postgres_deployment_postgres, postgres_pvc_postgres_pvc, postgres_service_postgres [EXTRACTED 1.00]

## Communities (36 total, 18 thin omitted)

### Community 0 - "Post Service Core"
Cohesion: 0.05
Nodes (40): Database, Database module — thin wrapper around psycopg2 for direct PostgreSQL access.  Th, Simple PostgreSQL connection wrapper with retry and dict cursors., Args:             dsn: PostgreSQL connection string (DATABASE_URL)., Establish a database connection with retry logic.          Args:             ret, Close the database connection gracefully., _configure_logging(), _handle_signal() (+32 more)

### Community 1 - "Post Service API"
Cohesion: 0.14
Nodes (28): APIResponse, Config, CreatePostRequest, PostResponse, Pydantic schemas (request / response models) for the Post Service.  These models, Payload for creating a new post., Serialized representation of a single post., Consistent envelope for all API responses.      Every endpoint returns this shap (+20 more)

### Community 2 - "Notification Service Core"
Cohesion: 0.09
Nodes (15): config, logger, config, config, EventSubscriber, Fastify, logger, notificationRoutes (+7 more)

### Community 3 - "Pulse Kubernetes Manifests"
Cohesion: 0.16
Nodes (17): Notification Service Deployment, Post Service Deployment, Notification Service, Post Service, Postgres Service, Redis Service, User Service, Notification Service Deployment (+9 more)

### Community 4 - "User Service Handlers"
Cohesion: 0.31
Nodes (7): jsonResponse, errorResponse(), NewUserHandler(), successResponse(), UserHandler, Context, UserRepository

### Community 5 - "User Service Data Layer"
Cohesion: 0.26
Nodes (7): CreateUserRequest, NewUserRepository(), UserRepository, Context, DB, UpdateUserRequest, User

### Community 6 - "User Service Entrypoint"
Cohesion: 0.24
Nodes (10): connectWithRetry(), logJSON(), main(), Config, getEnv(), Load(), Duration, HandlerFunc (+2 more)

### Community 7 - "Notification NPM Package"
Cohesion: 0.17
Nodes (11): dependencies, fastify, ioredis, uuid, description, main, name, scripts (+3 more)

### Community 8 - "Shared Infrastructure Config"
Cohesion: 0.20
Nodes (11): redis Requirement, redis Deployment, redis-pvc PersistentVolumeClaim, redis Service, pulse-shared-config ConfigMap, pulse-gateway Gateway, eg GatewayClass, pulse-route HTTPRoute (+3 more)

### Community 9 - "Postgres Execution Functions"
Cohesion: 0.22
Nodes (5): Any, Execute a query and return a single row as a dict, or None.          Args:, Execute a query and return all rows as a list of dicts.          Args:, Return True if the database is reachable., Execute a query and return the cursor (useful for INSERT/UPDATE/DELETE).

### Community 10 - "Application Configuration"
Cohesion: 0.25
Nodes (6): Config, Configuration module for the Post Service.  Uses pydantic-settings to load confi, Application settings loaded from environment variables., Settings, BaseSettings, HTTP client for the User Service.  The Post Service needs to verify that a user

### Community 12 - "Postgres Kubernetes Stack"
Cohesion: 0.50
Nodes (4): psycopg2-binary Requirement, postgres Deployment, postgres-pvc PersistentVolumeClaim, postgres Service

### Community 13 - "Postgres SQL Tables"
Cohesion: 0.67
Nodes (4): Posts Table, Users Table, Posts Table, Users Table

### Community 14 - "User Service Models"
Cohesion: 1.00
Nodes (3): CreateUserRequest, UpdateUserRequest, User

### Community 15 - "Postgres Init DB"
Cohesion: 0.67
Nodes (3): Postgres Init DB ConfigMap, Init DB ConfigMap, Postgres Deployment

## Knowledge Gaps
- **65 isolated node(s):** `name`, `version`, `description`, `main`, `start` (+60 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Database` connect `Post Service Core` to `Postgres Execution Functions`?**
  _High betweenness centrality (0.051) - this node is a cross-community bridge._
- **Why does `FastAPI` connect `Post Service Core` to `Post Service API`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `Database` (e.g. with `JSONFormatter` and `FastAPI`) actually correct?**
  _`Database` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 7 inferred relationships involving `APIResponse` (e.g. with `JSONFormatter` and `CreatePostRequest`) actually correct?**
  _`APIResponse` has 7 INFERRED edges - model-reasoned connections that need verification._
- **Are the 6 inferred relationships involving `EventPublisher` (e.g. with `JSONFormatter` and `lifespan()`) actually correct?**
  _`EventPublisher` has 6 INFERRED edges - model-reasoned connections that need verification._
- **Are the 6 inferred relationships involving `UserServiceClient` (e.g. with `JSONFormatter` and `lifespan()`) actually correct?**
  _`UserServiceClient` has 6 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `Request` (e.g. with `APIResponse` and `CreatePostRequest`) actually correct?**
  _`Request` has 3 INFERRED edges - model-reasoned connections that need verification._