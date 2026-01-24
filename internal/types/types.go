package types

// Resource represents a STOA resource (API, Subscription, etc.)
type Resource struct {
	APIVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	Metadata   Metadata `yaml:"metadata" json:"metadata"`
	Spec       any      `yaml:"spec" json:"spec"`
}

// Metadata contains resource metadata
type Metadata struct {
	Name      string            `yaml:"name" json:"name"`
	Namespace string            `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Labels    map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

// APISpec represents an API resource specification
type APISpec struct {
	Version        string              `yaml:"version,omitempty" json:"version,omitempty"`
	Description    string              `yaml:"description,omitempty" json:"description,omitempty"`
	Upstream       UpstreamSpec        `yaml:"upstream" json:"upstream"`
	Routing        RoutingSpec         `yaml:"routing,omitempty" json:"routing,omitempty"`
	Authentication AuthenticationSpec  `yaml:"authentication,omitempty" json:"authentication,omitempty"`
	Policies       PoliciesSpec        `yaml:"policies,omitempty" json:"policies,omitempty"`
	Catalog        CatalogSpec         `yaml:"catalog,omitempty" json:"catalog,omitempty"`
}

// UpstreamSpec defines the backend service
type UpstreamSpec struct {
	URL     string `yaml:"url" json:"url"`
	Timeout string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Retries int    `yaml:"retries,omitempty" json:"retries,omitempty"`
}

// RoutingSpec defines routing configuration
type RoutingSpec struct {
	Path      string   `yaml:"path" json:"path"`
	StripPath bool     `yaml:"stripPath,omitempty" json:"stripPath,omitempty"`
	Methods   []string `yaml:"methods,omitempty" json:"methods,omitempty"`
}

// AuthenticationSpec defines authentication requirements
type AuthenticationSpec struct {
	Required  bool           `yaml:"required" json:"required"`
	Providers []AuthProvider `yaml:"providers,omitempty" json:"providers,omitempty"`
}

// AuthProvider defines an authentication provider
type AuthProvider struct {
	Type   string `yaml:"type" json:"type"`
	Header string `yaml:"header,omitempty" json:"header,omitempty"`
	Issuer string `yaml:"issuer,omitempty" json:"issuer,omitempty"`
}

// PoliciesSpec defines API policies
type PoliciesSpec struct {
	RateLimit *RateLimitPolicy `yaml:"rateLimit,omitempty" json:"rateLimit,omitempty"`
	CORS      *CORSPolicy      `yaml:"cors,omitempty" json:"cors,omitempty"`
	Cache     *CachePolicy     `yaml:"cache,omitempty" json:"cache,omitempty"`
}

// RateLimitPolicy defines rate limiting
type RateLimitPolicy struct {
	RequestsPerHour int `yaml:"requestsPerHour,omitempty" json:"requestsPerHour,omitempty"`
}

// CORSPolicy defines CORS settings
type CORSPolicy struct {
	Origins []string `yaml:"origins,omitempty" json:"origins,omitempty"`
}

// CachePolicy defines caching settings
type CachePolicy struct {
	Enabled    bool `yaml:"enabled" json:"enabled"`
	TTLSeconds int  `yaml:"ttlSeconds,omitempty" json:"ttlSeconds,omitempty"`
}

// CatalogSpec defines catalog visibility
type CatalogSpec struct {
	Visibility string   `yaml:"visibility,omitempty" json:"visibility,omitempty"`
	Categories []string `yaml:"categories,omitempty" json:"categories,omitempty"`
}

// API represents an API from the STOA API
type API struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version,omitempty"`
	Description string   `json:"description,omitempty"`
	Tenant      string   `json:"tenant,omitempty"`
	Status      string   `json:"status,omitempty"`
	Upstream    string   `json:"upstream,omitempty"`
	Path        string   `json:"path,omitempty"`
	CreatedAt   string   `json:"createdAt,omitempty"`
	UpdatedAt   string   `json:"updatedAt,omitempty"`
}

// APIListResponse represents the response from GET /v1/portal/apis
type APIListResponse struct {
	Items      []API `json:"items"`
	TotalCount int   `json:"totalCount"`
}
