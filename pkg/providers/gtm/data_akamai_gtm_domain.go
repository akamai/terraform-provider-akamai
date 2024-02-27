package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &domainDataSource{}

var _ datasource.DataSourceWithConfigure = &domainDataSource{}

// NewGTMDomainDataSource returns a new GTM domain data source
func NewGTMDomainDataSource() datasource.DataSource {
	return &domainDataSource{}
}

var (
	domainBlock = map[string]schema.Block{
		"status": schema.SingleNestedBlock{
			Description: "Status information for the configuration.",
			Attributes: map[string]schema.Attribute{
				"message": schema.StringAttribute{
					Computed:    true,
					Description: "A notification generated when a change occurs to the domain.",
				},
				"change_id": schema.StringAttribute{
					Computed:    true,
					Description: "A unique identifier generated when a change occurs to the domain.",
				},
				"propagation_status": schema.StringAttribute{
					Computed:    true,
					Description: "Tracks the status of the domain's propagation state.",
				},
				"propagation_status_date": schema.StringAttribute{
					Computed:    true,
					Description: "An ISO 8601 timestamp indicating when a change occurs to the domain.",
				},
				"passing_validation": schema.BoolAttribute{
					Computed:    true,
					Description: "Indicates if the domain validates.",
				},
			},
			Blocks: map[string]schema.Block{
				"links": &schema.SetNestedBlock{
					Description: "Specifies the URL path that allows direct navigation to the domain.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"rel": schema.StringAttribute{
								Computed:    true,
								Description: "Indicates the link relationship of the object.",
							},
							"href": schema.StringAttribute{
								Computed:    true,
								Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
							},
						},
					},
				},
			},
		},
		"resources": schema.SetNestedBlock{
			Description: "List of resources associated with the domain.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"aggregation_type": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies how GTM handles different load numbers when multiple load servers are used for a data center or property.",
					},
					"constrained_property": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies the name of the property that this resource constraints.",
					},
					"decay_rate": schema.Float64Attribute{
						Computed:    true,
						Description: "For internal use only.",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive note to help you track what the resource constraints.",
					},
					"host_header": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies the host header used when fetching the load object.",
					},
					"leader_string": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies the text that comes before the loadObject.",
					},
					"least_squares_decay": schema.Float64Attribute{
						Computed:    true,
						Description: "For internal use only.",
					},
					"load_imbalance_percentage": schema.Float64Attribute{
						Computed:    true,
						Description: "Indicates the percent of load imbalance factor for the domain.",
					},
					"max_u_multiplicative_increment": schema.Float64Attribute{
						Computed:    true,
						Description: "For internal use only.",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive label for the resource.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "Indicates the kind of loadObject format used to determine the load on the resource.",
					},
					"upper_bound": schema.Int64Attribute{
						Computed:    true,
						Description: "An optional sanity check that specifies the maximum allowed value for any component of the load object.",
					},
				},
				Blocks: map[string]schema.Block{
					"resource_instances": &schema.SetNestedBlock{
						Description: "List of resource instances.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"load_object": schema.StringAttribute{
									Computed:    true,
									Description: "Identifies the load object file used to report real-time information about the current load, maximum allowable load and target load on each resource.",
								},
								"load_object_port": schema.Int64Attribute{
									Computed:    true,
									Description: "Specifies the TCP port of the loadObject.",
								},
								"load_servers": schema.ListAttribute{
									Description: "Specifies the list of servers to requests the load object from.",
									Computed:    true,
									ElementType: types.StringType,
								},
								"datacenter_id": schema.Int64Attribute{
									Computed:    true,
									Description: "A unique identifier for an existing data center in the domain.",
								},
								"use_default_load_object": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to use default loadObject.",
								},
							},
						},
					},
					"links": &schema.SetNestedBlock{
						Description: "Specifies the URL path that allows direct navigation to the resource.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Computed:    true,
									Description: "Indicates the link relationship of the object.",
								},
								"href": schema.StringAttribute{
									Computed:    true,
									Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
								},
							},
						},
					},
				},
			},
		},
		"properties": schema.SetNestedBlock{
			Description: "List of properties associated with the domain.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"backup_cname": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies a backup CNAME.",
					},
					"backup_ip": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies a backup IP.",
					},
					"balance_by_download_score": schema.BoolAttribute{
						Computed:    true,
						Description: "Indicates whether download score based load balancing is enabled.",
					},
					"cname": schema.StringAttribute{
						Computed:    true,
						Description: "Indicates the fully qualified name aliased to a particular property.",
					},
					"comments": schema.StringAttribute{
						Computed:    true,
						Description: "Descriptive comments for the property.",
					},
					"dynamic_ttl": schema.Int64Attribute{
						Computed:    true,
						Description: "Indicates the TTL in seconds for records that might change dynamically based on liveness and load balancing.",
					},
					"failover_delay": schema.Int64Attribute{
						Computed:    true,
						Description: "Specifies the failover delay in seconds.",
					},
					"failback_delay": schema.Int64Attribute{
						Computed:    true,
						Description: "Specifies the failback delay in seconds.",
					},
					"ghost_demand_reporting": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether an alternate way to collect load feedback from a GTM Performance domain is enabled.",
					},
					"handout_mode": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies how IPs are returned when more than one IP is alive and available.",
					},
					"handout_limit": schema.Int64Attribute{
						Computed:    true,
						Description: "Indicates the limit for the number of live IPs handed out to a DNS request.",
					},
					"health_max": schema.Float64Attribute{
						Computed:    true,
						Description: "Defines the absolute limit beyond which IPs are declared unhealthy.",
					},
					"health_multiplier": schema.Float64Attribute{
						Computed:    true,
						Description: "Configures a cutoff value that is computed from the median scores.",
					},
					"health_threshold": schema.Float64Attribute{
						Computed:    true,
						Description: "Configures a cutoff value that is computed from the median scores.",
					},
					"last_modified": schema.StringAttribute{
						Computed:    true,
						Description: "An ISO 8601 timestamp that indicates when the property was last changed.",
					},
					"load_imbalance_percentage": schema.Float64Attribute{
						Computed:    true,
						Description: "Indicates the percent of load imbalance factor for the domain.",
					},
					"map_name": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive label for a geographic or a CIDR map that's required if the property is either geographic or cidrmapping.",
					},
					"max_unreachable_penalty": schema.Int64Attribute{
						Computed:    true,
						Description: "For performance domains, this specifies a penalty value that's added to liveness test scores when data centers show an aggregated loss fraction higher than the penalty value.",
					},
					"min_live_fraction": schema.Float64Attribute{
						Computed:    true,
						Description: "Specifies what fraction of the servers need to respond to requests so GTM considers the data center up and able to receive traffic.",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive label for the property.",
					},
					"score_aggregation_type": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies how GTM aggregates liveness test scores across different tests, when multiple tests are configured.",
					},
					"stickness_bonus_constant": schema.Int64Attribute{
						Computed:    true,
						Description: "Specifies a percentage used to configure data center affinity.",
					},
					"stickness_bonus_percentage": schema.Int64Attribute{
						Computed:    true,
						Description: "Specifies a percentage used to configure data center affinity.",
					},
					"static_ttl": schema.Int64Attribute{
						Computed:    true,
						Description: "Specifies the TTL in seconds for static resource records that don't change based on the requesting name server IP.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies the load balancing behvior for the property.",
					},
					"unreachable_threshold": schema.Float64Attribute{
						Computed:    true,
						Description: "For performance domains, this specifies a penalty value that's added to liveness test scores when data centers have an aggregated loss fraction higher than this value.",
					},
					"use_computed_targets": schema.BoolAttribute{
						Computed:    true,
						Description: "For load-feedback domains only, this specifies that you want GTM to automatically compute target load.",
					},
					"ipv6": schema.BoolAttribute{
						Computed:    true,
						Description: "Indicates the type of IP address handed out by a property.",
					},
					"weighted_hash_bits_for_ipv4": schema.Int64Attribute{
						Computed:    true,
						Description: "For weighted hashed properties, how many leading bits of the client nameserver IP address to include when computing a hash for picking a datacenter for a client nameserver using IPv4; the default value is 32 (the entire address).",
					},
					"weighted_hash_bits_for_ipv6": schema.Int64Attribute{
						Computed:    true,
						Description: "For weighted hashed properties, how many leading bits of the client nameserver IP address to include when computing a hash for picking a datacenter for a client nameserver using IPv6; the default value is 128 (the entire address).",
					},
				},
				Blocks: map[string]schema.Block{
					"links": schema.SetNestedBlock{
						Description: "Provides a URL path that allows direct navigation to the property.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Computed:    true,
									Description: "Indicates the link relationship of the object.",
								},
								"href": schema.StringAttribute{
									Computed:    true,
									Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
								},
							},
						},
					},
					"static_rr_sets": schema.SetNestedBlock{
						Description: "Contains static recordsets.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "The record type.",
								},
								"ttl": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of seconds that this record should live in a resolver's cache before being refetched.",
								},
								"rdata": schema.ListAttribute{
									Description: "An array of data strings, representing multiple records within a set.",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
					"traffic_targets": schema.SetNestedBlock{
						Description: "Traffic targets for the property.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"datacenter_id": schema.Int64Attribute{
									Computed:    true,
									Description: "A unique identifier for an existing data center in the domain.",
								},
								"enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Indicates whether the traffic target is used.",
								},
								"weight": schema.Float64Attribute{
									Computed:    true,
									Description: "Specifies the traffic target weight for the target.",
								},
								"handout_cname": schema.StringAttribute{
									Computed:    true,
									Description: "Specifies an optional data center for the property.",
								},
								"name": schema.StringAttribute{
									Computed:    true,
									Description: "An alternative label for the traffic target.",
								},
								"servers": schema.ListAttribute{
									Description: "Identifies the IP address or the hostnames of the servers.",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
					"liveness_tests": schema.SetNestedBlock{
						Description: "Contains information about liveness tests.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"answers_required": schema.BoolAttribute{
									Computed:    true,
									Description: "If testObjectProtocol is DNS, DOH or DOT, requires an answer to the DNS query to be considered a success.",
								},
								"disabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Disables the liveness test.",
								},
								"disable_nonstandard_port_warning": schema.BoolAttribute{
									Computed:    true,
									Description: "Disables warnings when non-standard ports are used.",
								},
								"error_penalty": schema.Float64Attribute{
									Computed:    true,
									Description: "Specifies the score that's reported if the liveness test encounters an error other than timeout, such as connection refused, and 404.",
								},
								"http_error3xx": schema.BoolAttribute{
									Computed:    true,
									Description: "Treats a 3xx HTTP response as a failure if the testObjectProtocol is http, https or ftp.",
								},
								"http_error4xx": schema.BoolAttribute{
									Computed:    true,
									Description: "Treats a 4xx HTTP response as a failure if the testObjectProtocol is http, https or ftp.",
								},
								"http_error5xx": schema.BoolAttribute{
									Computed:    true,
									Description: "Treats a 5xx HTTP response as a failure if the testObjectProtocol is http, https or ftp.",
								},
								"name": schema.StringAttribute{
									Computed:    true,
									Description: "A descriptive name for the liveness test.",
								},
								"peer_certificate_verification": schema.BoolAttribute{
									Computed:    true,
									Description: "Validates the origin certificate. Applies only to tests with testObjectProtocol of https.",
								},
								"request_string": schema.StringAttribute{
									Computed:    true,
									Description: "Specifies a request string.",
								},
								"response_string": schema.StringAttribute{
									Computed:    true,
									Description: "Specifies a response string.",
								},
								"resource_type": schema.StringAttribute{
									Computed:    true,
									Description: "Specifies the query type, if testObjectProtocol is DNS.",
								},
								"recursion_requested": schema.BoolAttribute{
									Computed:    true,
									Description: "Indicates that if testObjectProtocol is DNS, DOH or DOT, the DNS query is recursive.",
								},
								"test_interval": schema.Int64Attribute{
									Computed:    true,
									Description: "Indicates the interval at which the liveness test is run, in seconds.",
								},
								"test_object": schema.StringAttribute{
									Computed:    true,
									Description: "Specifies the static text that acts as a stand-in for the data that you're sending on the network.",
								},
								"test_object_port": schema.Int64Attribute{
									Computed:    true,
									Description: "Specifies the port number for the testObject.",
								},
								"test_object_protocol": schema.StringAttribute{
									Computed:    true,
									Description: "Specifies the test protocol.",
								},
								"test_object_username": schema.StringAttribute{
									Computed:    true,
									Description: "A descriptive name for the testObject.",
								},
								"test_object_password": schema.StringAttribute{
									Computed:    true,
									Description: "Specifies the test object's password.",
								},
								"test_timeout": schema.Float64Attribute{
									Computed:    true,
									Description: "Specifies the duration of the liveness test before it fails.",
								},
								"timeout_penalty": schema.Float64Attribute{
									Computed:    true,
									Description: "Specifies the timeout penalty score.",
								},
								"ssl_client_certificate": schema.StringAttribute{
									Computed:    true,
									Description: "Indicates a base64-encoded certificate.",
								},
								"ssl_client_private_key": schema.StringAttribute{
									Computed:    true,
									Description: "Indicates a base64-encoded private key.",
								},
							},
							Blocks: map[string]schema.Block{
								"http_headers": &schema.SetNestedBlock{
									Description: "List of HTTP headers for the liveness test.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Computed:    true,
												Description: "Name of the HTTP header.",
											},
											"value": schema.StringAttribute{
												Computed:    true,
												Description: "Value of the HTTP header.",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"datacenters": schema.SetNestedBlock{
			Description: "List of data centers associated with the domain.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"datacenter_id": schema.Int64Attribute{
						Computed:    true,
						Description: "A unique identifier for an existing data center in the domain.",
					},
					"nickname": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive label for the datacenter.",
					},
					"score_penalty": schema.Int64Attribute{
						Computed:    true,
						Description: "Influences the score for a datacenter.",
					},
					"city": schema.StringAttribute{
						Computed:    true,
						Description: "The name of the city where the data center is located.",
					},
					"state_or_province": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies a two-letter ISO 3166 country code for the state of province, where the data center is located.",
					},
					"country": schema.StringAttribute{
						Computed:    true,
						Description: "A two-letter ISO 3166 country code that specifies the country where the data center is located.",
					},
					"latitude": schema.Float64Attribute{
						Computed:    true,
						Description: "Specifies the geographic latitude of the data center's position.",
					},
					"longitude": schema.Float64Attribute{
						Computed:    true,
						Description: "Specifies the geographic longitude of the data center's position.",
					},
					"clone_of": schema.Int64Attribute{
						Computed:    true,
						Description: "Identifies the data center's ID of which this data center is a clone.",
					},
					"virtual": schema.BoolAttribute{
						Computed:    true,
						Description: "Indicates whether or not the data center is virtual or physical.",
					},
					"continent": schema.StringAttribute{
						Computed:    true,
						Description: "A two-letter code that specifies the continent where the data center maps to.",
					},
					"server_monitor_pool": schema.StringAttribute{
						Computed:    true,
						Description: "The name of the pool from which servermonitors are drawn for liveness tests in this datacenter. If omitted (null), the domain-wide default is used. (If no domain-wide default is specified, the pool used is all servermonitors in the same continent as the datacenter.).",
					},
					"cloud_server_targeting": schema.BoolAttribute{
						Computed:    true,
						Description: "Balances load between two or more servers in a cloud environment.",
					},
					"cloud_server_host_header_override": schema.BoolAttribute{
						Computed:    true,
						Description: "Balances load between two or more servers in a cloud environment.",
					},
				},
				Blocks: map[string]schema.Block{
					"links": schema.SetNestedBlock{
						Description: "Provides a URL path that allows direct navigation to a data center.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Computed:    true,
									Description: "Indicates the link relationship of the object.",
								},
								"href": schema.StringAttribute{
									Computed:    true,
									Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
								},
							},
						},
					},
					"default_load_object": schema.SetNestedBlock{
						Description: "Specifies the load reporting interface between you and the GTM system.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"load_object": schema.StringAttribute{
									Computed:    true,
									Description: "Specifies the load object that GTM requests.",
								},
								"load_object_port": schema.Int64Attribute{
									Computed:    true,
									Description: "Specifies the TCP port to connect to when requesting the load object.",
								},
								"load_servers": schema.ListAttribute{
									Computed:    true,
									ElementType: types.StringType,
									Description: "Specifies the list of servers to requests the load object from.",
								},
							},
						},
					},
				},
			},
		},
		"geographic_maps": schema.SetNestedBlock{
			Description: "List of geographic maps associated with the domain.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive label for the geographic map.",
					},
				},
				Blocks: map[string]schema.Block{
					"assignments": schema.SetNestedBlock{
						Description: "Contains information about the geographic zone groupings of countries.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"countries": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
									Description: "Specifies an array of two-letter ISO 3166 `country` codes.",
								},
								"datacenter_id": schema.Int64Attribute{
									Computed:    true,
									Description: "A unique identifier for an existing data center in the domain.",
								},
								"nickname": schema.StringAttribute{
									Computed:    true,
									Description: "A descriptive label for all other AS zones.",
								},
							},
						},
					},
					"links": schema.SetNestedBlock{
						Description: "Specifies the URL path that allows direct navigation to the geographic map.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Computed:    true,
									Description: "Indicates the link relationship of the object.",
								},
								"href": schema.StringAttribute{
									Computed:    true,
									Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
								},
							},
						},
					},
					"default_datacenter": schema.SingleNestedBlock{
						Description: "A placeholder for all other geographic zones, countries not found in these geographic zones.",
						Attributes: map[string]schema.Attribute{
							"datacenter_id": schema.Int64Attribute{
								Computed:    true,
								Description: "An identifier for all other geographic zones' CNAME.",
							},
							"nickname": schema.StringAttribute{
								Computed:    true,
								Description: "A descriptive label for all other geographic zones.",
							},
						},
					},
				},
			},
		},
		"cidr_maps": schema.SetNestedBlock{
			Description: "Contains information about the set of CIDR maps assigned to this domain.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Unique name for the CIDR map.",
					},
				},
				Blocks: map[string]schema.Block{
					"assignments": schema.SetNestedBlock{
						Description: "Contains information about the CIDR zone groupings of CIDR blocks.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"blocks": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
									Description: "Specifies an array of CIDR blocks.",
								},
								"datacenter_id": schema.Int64Attribute{
									Computed:    true,
									Description: "A unique identifier for an existing data center in the domain.",
								},
								"nickname": schema.StringAttribute{
									Computed:    true,
									Description: "A descriptive label for all other AS zones.",
								},
							},
						},
					},
					"default_datacenter": schema.SingleNestedBlock{
						Description: "A placeholder for all other CIDR zones, CIDR blocks not found in these CIDR zones.",
						Attributes: map[string]schema.Attribute{
							"datacenter_id": schema.Int64Attribute{
								Computed:    true,
								Description: "For each property, an identifier for all other CIDR zones' CNAME.",
							},
							"nickname": schema.StringAttribute{
								Computed:    true,
								Description: "A descriptive label for all other CIDR blocks.",
							},
						},
					},
					"links": schema.SetNestedBlock{
						Description: "Specifies the URL path that allows direct navigation to the CIDR map.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Computed:    true,
									Description: "Indicates the link relationship of the object.",
								},
								"href": schema.StringAttribute{
									Computed:    true,
									Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
								},
							},
						},
					},
				},
			},
		},
		"as_maps": schema.SetNestedBlock{
			Description: "Contains information about the set of AS maps assigned to this domain.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive label for the AS map.",
					},
				},
				Blocks: map[string]schema.Block{
					"assignments": schema.SetNestedBlock{
						Description: "Contains information about the AS zone groupings of AS IDs.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"as_numbers": schema.SetAttribute{
									Computed:    true,
									ElementType: types.Int64Type,
									Description: "Specifies an array of AS numbers.",
								},
								"datacenter_id": schema.Int64Attribute{
									Computed:    true,
									Description: "A unique identifier for an existing data center in the domain.",
								},
								"nickname": schema.StringAttribute{
									Computed:    true,
									Description: "A descriptive label for all other AS zones.",
								},
							},
						},
					},
					"default_datacenter": schema.SingleNestedBlock{
						Description: "A placeholder for all other AS zones, AS IDs not found in these AS zones.",
						Attributes: map[string]schema.Attribute{
							"datacenter_id": schema.Int64Attribute{
								Computed:    true,
								Description: "For each property, an identifier for all other AS zones' CNAME.",
							},
							"nickname": schema.StringAttribute{
								Computed:    true,
								Description: "A descriptive label for all other AS zones.",
							},
						},
					},
					"links": schema.SetNestedBlock{
						Description: "Specifies the URL path that allows direct navigation to the As map.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Computed:    true,
									Description: "Indicates the link relationship of the object.",
								},
								"href": schema.StringAttribute{
									Computed:    true,
									Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
								},
							},
						},
					},
				},
			},
		},
		"links": schema.SetNestedBlock{
			Description: "Provides a URL path that allows direct navigation to the domain.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Computed:    true,
						Description: "Indicates the link relationship of the object.",
					},
					"href": schema.StringAttribute{
						Computed:    true,
						Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
					},
				},
			},
		},
	}
)

type domainDataSource struct {
	meta meta.Meta
}

type (
	domainDataSourceModel struct {
		ID                           types.String     `tfsdk:"id"`
		Name                         types.String     `tfsdk:"name"`
		CNameCoalescingEnabled       types.Bool       `tfsdk:"cname_coalescing_enabled"`
		DefaultErrorPenalty          types.Int64      `tfsdk:"default_error_penalty"`
		DefaultHealthMax             types.Float64    `tfsdk:"default_health_max"`
		DefaultHealthMultiplier      types.Float64    `tfsdk:"default_health_multiplier"`
		DefaultHealthThreshold       types.Float64    `tfsdk:"default_health_threshold"`
		DefaultMaxUnreachablePenalty types.Int64      `tfsdk:"default_max_unreachable_penalty"`
		DefaultSSLClientCertificate  types.String     `tfsdk:"default_ssl_client_certificate"`
		DefaultSSLClientPrivateKey   types.String     `tfsdk:"default_ssl_client_private_key"`
		DefaultTimeoutPenalty        types.Int64      `tfsdk:"default_timeout_penalty"`
		DefaultUnreachableThreshold  types.Float64    `tfsdk:"default_unreachable_threshold"`
		EmailNotificationList        types.List       `tfsdk:"email_notification_list"`
		EndUserMappingEnabled        types.Bool       `tfsdk:"end_user_mapping_enabled"`
		LastModified                 types.String     `tfsdk:"last_modified"`
		LastModifiedBy               types.String     `tfsdk:"last_modified_by"`
		LoadFeedback                 types.Bool       `tfsdk:"load_feedback"`
		MapUpdateInterval            types.Int64      `tfsdk:"map_update_interval"`
		MaxProperties                types.Int64      `tfsdk:"max_properties"`
		MaxResources                 types.Int64      `tfsdk:"max_resources"`
		MaxTestTimeout               types.Float64    `tfsdk:"max_test_timeout"`
		MaxTTL                       types.Int64      `tfsdk:"max_ttl"`
		MinPingableRegionFraction    types.Float64    `tfsdk:"min_pingable_region_fraction"`
		MinTestInterval              types.Int64      `tfsdk:"min_test_interval"`
		MinTTL                       types.Int64      `tfsdk:"min_ttl"`
		ModificationComments         types.String     `tfsdk:"modification_comments"`
		RoundRobinPrefix             types.String     `tfsdk:"round_robin_prefix"`
		ServerMonitorPool            types.String     `tfsdk:"server_monitor_pool"`
		Type                         types.String     `tfsdk:"type"`
		Status                       *status          `tfsdk:"status"`
		LoadImbalancePercentage      types.Float64    `tfsdk:"load_imbalance_percentage"`
		Resources                    []domainResource `tfsdk:"resources"`
		Properties                   []property       `tfsdk:"properties"`
		Datacenters                  []datacenter     `tfsdk:"datacenters"`
		GeographicMaps               []geographicMap  `tfsdk:"geographic_maps"`
		CIDRMaps                     []cidrMap        `tfsdk:"cidr_maps"`
		ASMaps                       []asMap          `tfsdk:"as_maps"`
		Links                        []link           `tfsdk:"links"`
	}

	status struct {
		Message               types.String `tfsdk:"message"`
		ChangeID              types.String `tfsdk:"change_id"`
		PropagationStatus     types.String `tfsdk:"propagation_status"`
		PropagationStatusDate types.String `tfsdk:"propagation_status_date"`
		PassingValidation     types.Bool   `tfsdk:"passing_validation"`
		Links                 []link       `tfsdk:"links"`
	}

	domainResource struct {
		AggregationType             types.String       `tfsdk:"aggregation_type"`
		ConstrainedProperty         types.String       `tfsdk:"constrained_property"`
		DecayRate                   types.Float64      `tfsdk:"decay_rate"`
		Description                 types.String       `tfsdk:"description"`
		HostHeader                  types.String       `tfsdk:"host_header"`
		LeaderString                types.String       `tfsdk:"leader_string"`
		LeastSquaresDecay           types.Float64      `tfsdk:"least_squares_decay"`
		LoadImbalancePercentage     types.Float64      `tfsdk:"load_imbalance_percentage"`
		MaxUMultiplicativeIncrement types.Float64      `tfsdk:"max_u_multiplicative_increment"`
		Name                        types.String       `tfsdk:"name"`
		ResourceInstances           []resourceInstance `tfsdk:"resource_instances"`
		Type                        types.String       `tfsdk:"type"`
		UpperBound                  types.Int64        `tfsdk:"upper_bound"`
		Links                       []link             `tfsdk:"links"`
	}

	livenessTest struct {
		AnswersRequired               types.Bool    `tfsdk:"answers_required"`
		Disabled                      types.Bool    `tfsdk:"disabled"`
		DisableNonstandardPortWarning types.Bool    `tfsdk:"disable_nonstandard_port_warning"`
		ErrorPenalty                  types.Float64 `tfsdk:"error_penalty"`
		HTTPHeaders                   []httpHeader  `tfsdk:"http_headers"`
		HTTPError3xx                  types.Bool    `tfsdk:"http_error3xx"`
		HTTPError4xx                  types.Bool    `tfsdk:"http_error4xx"`
		HTTPError5xx                  types.Bool    `tfsdk:"http_error5xx"`
		Name                          types.String  `tfsdk:"name"`
		PeerCertificateVerification   types.Bool    `tfsdk:"peer_certificate_verification"`
		RequestString                 types.String  `tfsdk:"request_string"`
		ResponseString                types.String  `tfsdk:"response_string"`
		ResourceType                  types.String  `tfsdk:"resource_type"`
		RecursionRequested            types.Bool    `tfsdk:"recursion_requested"`
		TestInterval                  types.Int64   `tfsdk:"test_interval"`
		TestObject                    types.String  `tfsdk:"test_object"`
		TestObjectPort                types.Int64   `tfsdk:"test_object_port"`
		TestObjectProtocol            types.String  `tfsdk:"test_object_protocol"`
		TestObjectUsername            types.String  `tfsdk:"test_object_username"`
		TestObjectPassword            types.String  `tfsdk:"test_object_password"`
		TestTimeout                   types.Float64 `tfsdk:"test_timeout"`
		TimeoutPenalty                types.Float64 `tfsdk:"timeout_penalty"`
		SSLClientCertificate          types.String  `tfsdk:"ssl_client_certificate"`
		SSLClientPrivateKey           types.String  `tfsdk:"ssl_client_private_key"`
	}

	httpHeader struct {
		Name  types.String `tfsdk:"name"`
		Value types.String `tfsdk:"value"`
	}

	staticRRSet struct {
		Type  types.String `tfsdk:"type"`
		TTL   types.Int64  `tfsdk:"ttl"`
		RData types.List   `tfsdk:"rdata"`
	}

	trafficTarget struct {
		DatacenterID types.Int64   `tfsdk:"datacenter_id"`
		Enabled      types.Bool    `tfsdk:"enabled"`
		Weight       types.Float64 `tfsdk:"weight"`
		HandoutCNAME types.String  `tfsdk:"handout_cname"`
		Name         types.String  `tfsdk:"name"`
		Servers      types.List    `tfsdk:"servers"`
	}

	property struct {
		BackupCNAME               types.String    `tfsdk:"backup_cname"`
		BackupIP                  types.String    `tfsdk:"backup_ip"`
		BalanceByDownloadScore    types.Bool      `tfsdk:"balance_by_download_score"`
		CName                     types.String    `tfsdk:"cname"`
		Comments                  types.String    `tfsdk:"comments"`
		DynamicTTL                types.Int64     `tfsdk:"dynamic_ttl"`
		FailoverDelay             types.Int64     `tfsdk:"failover_delay"`
		FailbackDelay             types.Int64     `tfsdk:"failback_delay"`
		GhostDemandReporting      types.Bool      `tfsdk:"ghost_demand_reporting"`
		HandoutMode               types.String    `tfsdk:"handout_mode"`
		HandoutLimit              types.Int64     `tfsdk:"handout_limit"`
		HealthMax                 types.Float64   `tfsdk:"health_max"`
		HealthMultiplier          types.Float64   `tfsdk:"health_multiplier"`
		HealthThreshold           types.Float64   `tfsdk:"health_threshold"`
		LastModified              types.String    `tfsdk:"last_modified"`
		LivenessTests             []livenessTest  `tfsdk:"liveness_tests"`
		LoadImbalancePercentage   types.Float64   `tfsdk:"load_imbalance_percentage"`
		MapName                   types.String    `tfsdk:"map_name"`
		MaxUnreachablePenalty     types.Int64     `tfsdk:"max_unreachable_penalty"`
		MinLiveFraction           types.Float64   `tfsdk:"min_live_fraction"`
		Name                      types.String    `tfsdk:"name"`
		ScoreAggregationType      types.String    `tfsdk:"score_aggregation_type"`
		StickinessBonusConstant   types.Int64     `tfsdk:"stickness_bonus_constant"`
		StickinessBonusPercentage types.Int64     `tfsdk:"stickness_bonus_percentage"`
		StaticTTL                 types.Int64     `tfsdk:"static_ttl"`
		StaticRRSets              []staticRRSet   `tfsdk:"static_rr_sets"`
		TrafficTargets            []trafficTarget `tfsdk:"traffic_targets"`
		Type                      types.String    `tfsdk:"type"`
		UnreachableThreshold      types.Float64   `tfsdk:"unreachable_threshold"`
		UseComputedTargets        types.Bool      `tfsdk:"use_computed_targets"`
		IPv6                      types.Bool      `tfsdk:"ipv6"`
		WeightedHashBitsForIPv4   types.Int64     `tfsdk:"weighted_hash_bits_for_ipv4"`
		WeightedHashBitsForIPv6   types.Int64     `tfsdk:"weighted_hash_bits_for_ipv6"`
		Links                     []link          `tfsdk:"links"`
	}

	loadObject struct {
		LoadObject     types.String `tfsdk:"load_object"`
		LoadObjectPort types.Int64  `tfsdk:"load_object_port"`
		LoadServers    types.List   `tfsdk:"load_servers"`
	}

	datacenter struct {
		DatacenterID                  types.Int64   `tfsdk:"datacenter_id"`
		Nickname                      types.String  `tfsdk:"nickname"`
		ScorePenalty                  types.Int64   `tfsdk:"score_penalty"`
		City                          types.String  `tfsdk:"city"`
		StateOrProvince               types.String  `tfsdk:"state_or_province"`
		Country                       types.String  `tfsdk:"country"`
		Latitude                      types.Float64 `tfsdk:"latitude"`
		Longitude                     types.Float64 `tfsdk:"longitude"`
		CloneOf                       types.Int64   `tfsdk:"clone_of"`
		Virtual                       types.Bool    `tfsdk:"virtual"`
		DefaultLoadObject             []loadObject  `tfsdk:"default_load_object"`
		Continent                     types.String  `tfsdk:"continent"`
		ServerMonitorPool             types.String  `tfsdk:"server_monitor_pool"`
		CloudServerTargeting          types.Bool    `tfsdk:"cloud_server_targeting"`
		CloudServerHostHeaderOverride types.Bool    `tfsdk:"cloud_server_host_header_override"`
		Links                         []link        `tfsdk:"links"`
	}

	geographicMap struct {
		Name              types.String              `tfsdk:"name"`
		Assignments       []geographicMapAssignment `tfsdk:"assignments"`
		DefaultDatacenter defaultDatacenter         `tfsdk:"default_datacenter"`
		Links             []link                    `tfsdk:"links"`
	}

	geographicMapAssignment struct {
		Countries    types.Set    `tfsdk:"countries"`
		DatacenterID types.Int64  `tfsdk:"datacenter_id"`
		Nickname     types.String `tfsdk:"nickname"`
	}

	cidrMap struct {
		Name              types.String        `tfsdk:"name"`
		Assignments       []cidrMapAssignment `tfsdk:"assignments"`
		DefaultDatacenter defaultDatacenter   `tfsdk:"default_datacenter"`
		Links             []link              `tfsdk:"links"`
	}

	cidrMapAssignment struct {
		DatacenterID types.Int64  `tfsdk:"datacenter_id"`
		Nickname     types.String `tfsdk:"nickname"`
		Blocks       types.Set    `tfsdk:"blocks"`
	}

	asMap struct {
		Name              types.String      `tfsdk:"name"`
		Assignments       []asMapAssignment `tfsdk:"assignments"`
		DefaultDatacenter defaultDatacenter `tfsdk:"default_datacenter"`
		Links             []link            `tfsdk:"links"`
	}

	asMapAssignment struct {
		DatacenterID types.Int64   `tfsdk:"datacenter_id"`
		Nickname     types.String  `tfsdk:"nickname"`
		ASNumbers    []types.Int64 `tfsdk:"as_numbers"`
	}

	defaultDatacenter struct {
		DatacenterID types.Int64  `tfsdk:"datacenter_id"`
		Nickname     types.String `tfsdk:"nickname"`
	}
)

// Schema is used to define data source's terraform schema
func (d *domainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "GTM Domain data source",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The full GTM domain name.",
				Required:            true,
			},
			"cname_coalescing_enabled": schema.BoolAttribute{
				MarkdownDescription: "If enabled, GTM collapses CNAME redirections in DNS answers when it knows the target of the CNAME.",
				Computed:            true,
			},
			"default_error_penalty": schema.Int64Attribute{
				MarkdownDescription: "Specifies the download penalty score.",
				Computed:            true,
			},
			"default_health_max": schema.Float64Attribute{
				MarkdownDescription: "Default value for healthMax if none specified at the property level.",
				Computed:            true,
			},
			"default_health_multiplier": schema.Float64Attribute{
				MarkdownDescription: "Default value for healthMultiplier if none specified at the property level.",
				Computed:            true,
			},
			"default_health_threshold": schema.Float64Attribute{
				MarkdownDescription: "Default value for healthThreshold if none specified at the property level.",
				Computed:            true,
			},
			"default_max_unreachable_penalty": schema.Int64Attribute{
				MarkdownDescription: "Applicable only to Performance Plus domains.",
				Computed:            true,
			},
			"default_ssl_client_certificate": schema.StringAttribute{
				MarkdownDescription: "Specifies an optional Base64-encoded certificate that corresponds with the private key for TLS-based liveness tests (HTTPS, SMTPS, POPS, and TCPS).",
				Computed:            true,
			},
			"default_ssl_client_private_key": schema.StringAttribute{
				MarkdownDescription: "Specifies an optional Base64-encoded private key that corresponds with the TLS certificate for TLS-based liveness tests (HTTPS, SMTPS, POPS, and TCPS).",
				Computed:            true,
			},
			"default_timeout_penalty": schema.Int64Attribute{
				MarkdownDescription: "Specifies the timeout penalty score.",
				Computed:            true,
			},
			"default_unreachable_threshold": schema.Float64Attribute{
				MarkdownDescription: "Applicable only to Performance Plus domains. If the fraction of core points that cannot reach the datacenter (they have 100% packet loss to it) exceeds this threshold, a score penalty is added to liveness scores for servers in the datacenter; the penalty is equal to maxUnreachablePenalty * (fractionUnreachable - unreachableThreshold) / (1 - unreachableThreshold). ",
				Computed:            true,
			},
			"email_notification_list": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Email addresses where notifications will be sent.",
			},
			"end_user_mapping_enabled": schema.BoolAttribute{
				MarkdownDescription: "A boolean indicating whether whether the GTM Domain is using end user client subnet mapping.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "An ISO 8601 timestamp that indicates the time of the last domain change.",
				Computed:            true,
			},
			"last_modified_by": schema.StringAttribute{
				MarkdownDescription: "The email address of the administrator who made the last change to the domain.",
				Computed:            true,
			},
			"load_feedback": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether you're using resources to control load balancing.",
				Computed:            true,
			},
			"map_update_interval": schema.Int64Attribute{
				MarkdownDescription: "How often new maps are generated for performance domains. Not applicable to non-performance domains.",
				Computed:            true,
			},
			"max_properties": schema.Int64Attribute{
				MarkdownDescription: "Maximum amount of properties that may be associated with the domain.",
				Computed:            true,
			},
			"max_resources": schema.Int64Attribute{
				MarkdownDescription: "Maximum amount of resources that may be associated with the domain.",
				Computed:            true,
			},
			"max_test_timeout": schema.Float64Attribute{
				MarkdownDescription: "Maximum timeout for a test.",
				Computed:            true,
			},
			"max_ttl": schema.Int64Attribute{
				MarkdownDescription: "The largest TTL allowed. Configurations specifying TTLs greater than this will fail validation.",
				Computed:            true,
			},
			"min_pingable_region_fraction": schema.Float64Attribute{
				MarkdownDescription: "Applicable only to Performance Plus domains. If set (nonzero), any core point that cannot ping more than this fraction of datacenters is rejected and will not be mapped by ping scores.",
				Computed:            true,
			},
			"min_test_interval": schema.Int64Attribute{
				MarkdownDescription: "The smallest allowed liveness test interval. Configurations specifying liveness test intervals smaller than this will fail validation.",
				Computed:            true,
			},
			"min_ttl": schema.Int64Attribute{
				MarkdownDescription: "The smallest TTL allowed. Configurations specifying TTLs smaller than this will fail validation.",
				Computed:            true,
			},
			"modification_comments": schema.StringAttribute{
				MarkdownDescription: "A descriptive note about changes to the domain.",
				Computed:            true,
			},
			"round_robin_prefix": schema.StringAttribute{
				MarkdownDescription: "A string that when configured automatically creates a shadow property for each normal property.",
				Computed:            true,
			},
			"server_monitor_pool": schema.StringAttribute{
				MarkdownDescription: "The name of the pool from which servermonitors are drawn for liveness tests in this datacenter. If omitted (null), the domain-wide default is used. (If no domain-wide default is specified, the pool used is all servermonitors in the same continent as the datacenter.)",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Specifies the load balancing behavior for the property. ",
				Computed:            true,
			},
			"load_imbalance_percentage": schema.Float64Attribute{
				MarkdownDescription: "Indicates the percent of load imbalance factor (LIF) for the domain.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data source.",
				Computed:            true,
			},
		},
		Blocks: domainBlock,
	}
}

// Configure  configures data source at the beginning of the lifecycle
func (d *domainDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta, ok := request.ProviderData.(meta.Meta)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
	}
	d.meta = meta
}

// Metadata configures data source's meta information
func (d *domainDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_gtm_domain"
}

// Read is called when the provider must read data source values in order to update state
func (d *domainDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "GTM Domain DataSource Read")
	var data *domainDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	domain, err := client.GetDomain(ctx, data.Name.ValueString())
	if err != nil {
		response.Diagnostics.AddError("fetching domain failed", err.Error())
		return
	}

	data, diags := populateDomain(ctx, domain)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func populateDomain(ctx context.Context, domain *gtm.Domain) (*domainDataSourceModel, diag.Diagnostics) {
	emailNotificationList, diags := types.ListValueFrom(ctx, types.StringType, domain.EmailNotificationList)
	if diags.HasError() {
		return nil, diags
	}
	datacenters, diags := getDatacenters(ctx, domain.Datacenters)
	if diags.HasError() {
		return nil, diags
	}
	properties, diags := getProperties(ctx, domain.Properties)
	if diags.HasError() {
		return nil, diags
	}
	cidrMaps, diags := getCIDRMaps(ctx, domain.CidrMaps)
	if diags.HasError() {
		return nil, diags
	}
	geoMaps, diags := getGeographicMaps(ctx, domain.GeographicMaps)
	if diags.HasError() {
		return nil, diags
	}
	return &domainDataSourceModel{
		Name:                         types.StringValue(domain.Name),
		CNameCoalescingEnabled:       types.BoolValue(domain.CnameCoalescingEnabled),
		DefaultErrorPenalty:          types.Int64Value(int64(domain.DefaultErrorPenalty)),
		DefaultHealthMax:             types.Float64Value(domain.DefaultHealthMax),
		DefaultHealthMultiplier:      types.Float64Value(domain.DefaultHealthMultiplier),
		DefaultHealthThreshold:       types.Float64Value(domain.DefaultHealthThreshold),
		DefaultMaxUnreachablePenalty: types.Int64Value(int64(domain.DefaultMaxUnreachablePenalty)),
		DefaultSSLClientCertificate:  types.StringValue(domain.DefaultSslClientCertificate),
		DefaultSSLClientPrivateKey:   types.StringValue(domain.DefaultSslClientPrivateKey),
		DefaultTimeoutPenalty:        types.Int64Value(int64(domain.DefaultTimeoutPenalty)),
		DefaultUnreachableThreshold:  types.Float64Value(float64(domain.DefaultUnreachableThreshold)),
		EmailNotificationList:        emailNotificationList,
		EndUserMappingEnabled:        types.BoolValue(domain.EndUserMappingEnabled),
		LastModified:                 types.StringValue(domain.LastModified),
		LastModifiedBy:               types.StringValue(domain.LastModifiedBy),
		LoadFeedback:                 types.BoolValue(domain.LoadFeedback),
		MapUpdateInterval:            types.Int64Value(int64(domain.MapUpdateInterval)),
		MaxProperties:                types.Int64Value(int64(domain.MaxProperties)),
		MaxResources:                 types.Int64Value(int64(domain.MaxResources)),
		MaxTestTimeout:               types.Float64Value(domain.MaxTestTimeout),
		MaxTTL:                       types.Int64Value(domain.MaxTTL),
		MinPingableRegionFraction:    types.Float64Value(float64(domain.MinPingableRegionFraction)),
		MinTestInterval:              types.Int64Value(int64(domain.MinTestInterval)),
		MinTTL:                       types.Int64Value(domain.MinTTL),
		ModificationComments:         types.StringValue(domain.ModificationComments),
		RoundRobinPrefix:             types.StringValue(domain.RoundRobinPrefix),
		ServerMonitorPool:            types.StringValue(domain.ServermonitorPool),
		Type:                         types.StringValue(domain.Type),
		LoadImbalancePercentage:      types.Float64Value(domain.LoadImbalancePercentage),
		ID:                           types.StringValue(domain.Name),
		Status:                       getStatus(domain.Status),
		Resources:                    getResources(domain.Resources),
		Properties:                   properties,
		Datacenters:                  datacenters,
		GeographicMaps:               geoMaps,
		CIDRMaps:                     cidrMaps,
		ASMaps:                       getASMaps(domain.AsMaps),
		Links:                        getLinks(domain.Links),
	}, nil
}

func getLinks(links []*gtm.Link) []link {
	var result []link
	if links != nil {
		result = make([]link, len(links))
		for i, l := range links {
			result[i] = link{
				Rel:  types.StringValue(l.Rel),
				Href: types.StringValue(l.Href),
			}
		}
	}
	return result
}

func getASMaps(maps []*gtm.AsMap) []asMap {
	var result []asMap
	for _, am := range maps {
		asMapInstance := asMap{
			Name: types.StringValue(am.Name),
		}

		if am.Links != nil {
			asMapInstance.Links = populateLinks(am.Links)
		}

		if am.DefaultDatacenter != nil {
			defaultDataCenter := defaultDatacenter{
				Nickname:     types.StringValue(am.DefaultDatacenter.Nickname),
				DatacenterID: types.Int64Value(int64(am.DefaultDatacenter.DatacenterId)),
			}
			asMapInstance.DefaultDatacenter = defaultDataCenter
		}
		if am.Assignments != nil {
			asMapInstance.Assignments = make([]asMapAssignment, len(am.Assignments))
			for i, asg := range am.Assignments {
				asMapInstance.Assignments[i] = populateASMapAssignment(asg)
			}
		}
		result = append(result, asMapInstance)
	}
	return result
}

func getCIDRMaps(ctx context.Context, maps []*gtm.CidrMap) ([]cidrMap, diag.Diagnostics) {
	var result []cidrMap
	for _, cm := range maps {
		cidrMapInstance := cidrMap{
			Name: types.StringValue(cm.Name),
		}

		if cm.Links != nil {
			cidrMapInstance.Links = populateLinks(cm.Links)
		}

		if cm.DefaultDatacenter != nil {
			defaultDataCenter := defaultDatacenter{
				Nickname:     types.StringValue(cm.DefaultDatacenter.Nickname),
				DatacenterID: types.Int64Value(int64(cm.DefaultDatacenter.DatacenterId)),
			}
			cidrMapInstance.DefaultDatacenter = defaultDataCenter
		}

		if cm.Assignments != nil {
			cidrMapInstance.Assignments = make([]cidrMapAssignment, len(cm.Assignments))
			for i, asg := range cm.Assignments {
				popCIDRMapAssignment, diags := populateCIDRMapAssignment(ctx, asg)
				if diags.HasError() {
					return nil, diags
				}
				cidrMapInstance.Assignments[i] = popCIDRMapAssignment
			}
		}
		result = append(result, cidrMapInstance)
	}
	return result, nil
}

func getGeographicMaps(ctx context.Context, maps []*gtm.GeoMap) ([]geographicMap, diag.Diagnostics) {
	var result []geographicMap
	for _, gm := range maps {
		geoMapInstance := geographicMap{
			Name: types.StringValue(gm.Name),
		}

		if gm.Links != nil {
			geoMapInstance.Links = populateLinks(gm.Links)
		}

		if gm.DefaultDatacenter != nil {
			defaultDataCenter := defaultDatacenter{
				Nickname:     types.StringValue(gm.DefaultDatacenter.Nickname),
				DatacenterID: types.Int64Value(int64(gm.DefaultDatacenter.DatacenterId)),
			}
			geoMapInstance.DefaultDatacenter = defaultDataCenter
		}

		if gm.Assignments != nil {
			geoMapInstance.Assignments = make([]geographicMapAssignment, len(gm.Assignments))
			for i, asg := range gm.Assignments {
				popGeoMap, diags := populateGeographicMapAssignment(ctx, asg)
				if diags.HasError() {
					return nil, diags
				}
				geoMapInstance.Assignments[i] = popGeoMap
			}
		}
		result = append(result, geoMapInstance)
	}
	return result, nil
}

func getDatacenters(ctx context.Context, datacenters []*gtm.Datacenter) ([]datacenter, diag.Diagnostics) {
	var result []datacenter
	for _, dc := range datacenters {
		dataCenterInstance := datacenter{
			DatacenterID:                  types.Int64Value(int64(dc.DatacenterId)),
			Nickname:                      types.StringValue(dc.Nickname),
			ScorePenalty:                  types.Int64Value(int64(dc.ScorePenalty)),
			City:                          types.StringValue(dc.City),
			StateOrProvince:               types.StringValue(dc.StateOrProvince),
			Country:                       types.StringValue(dc.Country),
			Latitude:                      types.Float64Value(dc.Latitude),
			Longitude:                     types.Float64Value(dc.Longitude),
			CloneOf:                       types.Int64Value(int64(dc.CloneOf)),
			Virtual:                       types.BoolValue(dc.Virtual),
			Continent:                     types.StringValue(dc.Continent),
			ServerMonitorPool:             types.StringValue(dc.ServermonitorPool),
			CloudServerTargeting:          types.BoolValue(dc.CloudServerTargeting),
			CloudServerHostHeaderOverride: types.BoolValue(dc.CloudServerHostHeaderOverride),
		}

		if dc.DefaultLoadObject != nil {
			loadObj, diags := populateLoadObject(ctx, dc.DefaultLoadObject)
			if diags.HasError() {
				return nil, diags
			}
			dataCenterInstance.DefaultLoadObject = []loadObject{loadObj}
		}

		if dc.Links != nil {
			dataCenterInstance.Links = populateLinks(dc.Links)
		}

		result = append(result, dataCenterInstance)
	}
	return result, nil
}

func getProperties(ctx context.Context, properties []*gtm.Property) ([]property, diag.Diagnostics) {
	var result []property
	for _, prop := range properties {
		propertyInstance := property{
			BackupCNAME:               types.StringValue(prop.BackupCName),
			BackupIP:                  types.StringValue(prop.BackupIp),
			BalanceByDownloadScore:    types.BoolValue(prop.BalanceByDownloadScore),
			CName:                     types.StringValue(prop.CName),
			Comments:                  types.StringValue(prop.Comments),
			DynamicTTL:                types.Int64Value(int64(prop.DynamicTTL)),
			FailoverDelay:             types.Int64Value(int64(prop.FailoverDelay)),
			FailbackDelay:             types.Int64Value(int64(prop.FailbackDelay)),
			GhostDemandReporting:      types.BoolValue(prop.GhostDemandReporting),
			HandoutMode:               types.StringValue(prop.HandoutMode),
			HandoutLimit:              types.Int64Value(int64(prop.HandoutLimit)),
			HealthMax:                 types.Float64Value(float64(prop.HandoutLimit)),
			HealthMultiplier:          types.Float64Value(prop.HealthMultiplier),
			HealthThreshold:           types.Float64Value(prop.HealthThreshold),
			LastModified:              types.StringValue(prop.LastModified),
			LoadImbalancePercentage:   types.Float64Value(prop.LoadImbalancePercentage),
			MapName:                   types.StringValue(prop.MapName),
			MaxUnreachablePenalty:     types.Int64Value(int64(prop.MaxUnreachablePenalty)),
			MinLiveFraction:           types.Float64Value(prop.MinLiveFraction),
			Name:                      types.StringValue(prop.Name),
			ScoreAggregationType:      types.StringValue(prop.ScoreAggregationType),
			StickinessBonusConstant:   types.Int64Value(int64(prop.StickinessBonusConstant)),
			StickinessBonusPercentage: types.Int64Value(int64(prop.StickinessBonusPercentage)),
			StaticTTL:                 types.Int64Value(int64(prop.StaticTTL)),
			Type:                      types.StringValue(prop.Type),
			UnreachableThreshold:      types.Float64Value(prop.UnreachableThreshold),
			UseComputedTargets:        types.BoolValue(prop.UseComputedTargets),
			IPv6:                      types.BoolValue(prop.Ipv6),
			WeightedHashBitsForIPv4:   types.Int64Value(int64(prop.WeightedHashBitsForIPv4)),
			WeightedHashBitsForIPv6:   types.Int64Value(int64(prop.WeightedHashBitsForIPv6)),
		}

		if prop.LivenessTests != nil {
			propertyInstance.LivenessTests = make([]livenessTest, len(prop.LivenessTests))
			for i, lt := range prop.LivenessTests {
				propertyInstance.LivenessTests[i] = populateLivenessTest(lt)
			}
		}

		if prop.StaticRRSets != nil {
			propertyInstance.StaticRRSets = make([]staticRRSet, len(prop.StaticRRSets))
			for i, s := range prop.StaticRRSets {
				popStaticRRSet, diags := populateStaticRRSet(ctx, s)
				if diags.HasError() {
					return nil, diags
				}
				propertyInstance.StaticRRSets[i] = popStaticRRSet
			}
		}

		if prop.TrafficTargets != nil {
			propertyInstance.TrafficTargets = make([]trafficTarget, len(prop.TrafficTargets))
			for i, t := range prop.TrafficTargets {
				popTrafficTarget, diags := populateTrafficTarget(ctx, t)
				if diags.HasError() {
					return nil, diags
				}
				propertyInstance.TrafficTargets[i] = popTrafficTarget
			}
		}

		if prop.Links != nil {
			propertyInstance.Links = populateLinks(prop.Links)
		}

		result = append(result, propertyInstance)
	}
	return result, nil
}

func getStatus(st *gtm.ResponseStatus) *status {
	if st == nil {
		return nil
	}
	statusInstance := status{
		Message:               types.StringValue(st.Message),
		ChangeID:              types.StringValue(st.ChangeId),
		PropagationStatus:     types.StringValue(st.PropagationStatus),
		PropagationStatusDate: types.StringValue(st.PropagationStatusDate),
		PassingValidation:     types.BoolValue(st.PassingValidation),
	}
	if st.Links != nil {
		statusInstance.Links = make([]link, len(*st.Links))
		for i, l := range *st.Links {
			statusInstance.Links[i] = link{
				Rel:  types.StringValue(l.Rel),
				Href: types.StringValue(l.Href),
			}
		}
	}
	return &statusInstance
}

func getResources(resources []*gtm.Resource) []domainResource {
	var result []domainResource
	for _, res := range resources {
		resource := domainResource{
			AggregationType:             types.StringValue(res.AggregationType),
			ConstrainedProperty:         types.StringValue(res.ConstrainedProperty),
			DecayRate:                   types.Float64Value(res.DecayRate),
			Description:                 types.StringValue(res.Description),
			HostHeader:                  types.StringValue(res.HostHeader),
			LeaderString:                types.StringValue(res.LeaderString),
			LeastSquaresDecay:           types.Float64Value(res.LeastSquaresDecay),
			LoadImbalancePercentage:     types.Float64Value(res.LoadImbalancePercentage),
			MaxUMultiplicativeIncrement: types.Float64Value(res.MaxUMultiplicativeIncrement),
			Name:                        types.StringValue(res.Name),
			Type:                        types.StringValue(res.Type),
			UpperBound:                  types.Int64Value(int64(res.UpperBound)),
		}

		if res.ResourceInstances != nil {
			resInstances := make([]resourceInstance, len(res.ResourceInstances))
			for i, ri := range res.ResourceInstances {
				resInstances[i] = resourceInstance{
					LoadObject:           types.StringValue(ri.LoadObject.LoadObject),
					LoadObjectPort:       types.Int64Value(int64(ri.LoadObject.LoadObjectPort)),
					DataCenterID:         types.Int64Value(int64(ri.DatacenterId)),
					UseDefaultLoadObject: types.BoolValue(ri.UseDefaultLoadObject),
				}
				if ri.LoadObject.LoadServers != nil {
					loadServers := make([]types.String, len(ri.LoadObject.LoadServers))
					for i, s := range ri.LoadObject.LoadServers {
						loadServers[i] = types.StringValue(s)
					}
					resInstances[i].LoadServers = loadServers
				}
			}
			resource.ResourceInstances = resInstances
		}

		if res.Links != nil {
			resource.Links = populateLinks(res.Links)
		}

		result = append(result, resource)
	}
	return result
}

func populateLivenessTest(lt *gtm.LivenessTest) livenessTest {
	return livenessTest{
		AnswersRequired:               types.BoolValue(lt.AnswersRequired),
		Disabled:                      types.BoolValue(lt.Disabled),
		DisableNonstandardPortWarning: types.BoolValue(lt.DisableNonstandardPortWarning),
		ErrorPenalty:                  types.Float64Value(lt.ErrorPenalty),
		HTTPError3xx:                  types.BoolValue(lt.HttpError3xx),
		HTTPError4xx:                  types.BoolValue(lt.HttpError4xx),
		HTTPError5xx:                  types.BoolValue(lt.HttpError5xx),
		Name:                          types.StringValue(lt.Name),
		PeerCertificateVerification:   types.BoolValue(lt.PeerCertificateVerification),
		RequestString:                 types.StringValue(lt.RequestString),
		ResponseString:                types.StringValue(lt.ResponseString),
		ResourceType:                  types.StringValue(lt.ResourceType),
		RecursionRequested:            types.BoolValue(lt.RecursionRequested),
		TestInterval:                  types.Int64Value(int64(lt.TestInterval)),
		TestObject:                    types.StringValue(lt.TestObject),
		TestObjectPort:                types.Int64Value(int64(lt.TestObjectPort)),
		TestObjectProtocol:            types.StringValue(lt.TestObjectProtocol),
		TestObjectUsername:            types.StringValue(lt.TestObjectUsername),
		TestObjectPassword:            types.StringValue(lt.TestObjectPassword),
		TestTimeout:                   types.Float64Value(float64(lt.TestTimeout)),
		TimeoutPenalty:                types.Float64Value(lt.TimeoutPenalty),
		SSLClientCertificate:          types.StringValue(lt.SslClientCertificate),
		SSLClientPrivateKey:           types.StringValue(lt.SslClientPrivateKey),
		HTTPHeaders:                   populateHTTPHeaders(lt.HttpHeaders),
	}
}

func populateHTTPHeaders(headers []*gtm.HttpHeader) []httpHeader {
	result := make([]httpHeader, len(headers))
	for i, h := range headers {
		result[i] = httpHeader{
			Name:  types.StringValue(h.Name),
			Value: types.StringValue(h.Value),
		}
	}
	return result
}

func populateStaticRRSet(ctx context.Context, s *gtm.StaticRRSet) (staticRRSet, diag.Diagnostics) {
	Rdata, diags := types.ListValueFrom(ctx, types.StringType, s.Rdata)
	if diags.HasError() {
		return staticRRSet{}, diags
	}
	return staticRRSet{
		Type:  types.StringValue(s.Type),
		TTL:   types.Int64Value(int64(s.TTL)),
		RData: Rdata,
	}, diags
}

func populateLinks(links []*gtm.Link) []link {
	result := make([]link, len(links))
	for i, l := range links {
		result[i] = link{
			Rel:  types.StringValue(l.Rel),
			Href: types.StringValue(l.Href),
		}
	}
	return result
}

func populateTrafficTarget(ctx context.Context, t *gtm.TrafficTarget) (trafficTarget, diag.Diagnostics) {
	servers, diags := types.ListValueFrom(ctx, types.StringType, t.Servers)
	if diags.HasError() {
		return trafficTarget{}, diags
	}
	return trafficTarget{
		DatacenterID: types.Int64Value(int64(t.DatacenterId)),
		Enabled:      types.BoolValue(t.Enabled),
		Weight:       types.Float64Value(t.Weight),
		HandoutCNAME: types.StringValue(t.HandoutCName),
		Name:         types.StringValue(t.Name),
		Servers:      servers,
	}, nil
}

func populateLoadObject(ctx context.Context, lo *gtm.LoadObject) (loadObject, diag.Diagnostics) {
	loadObj := loadObject{
		LoadObject:     types.StringValue(lo.LoadObject),
		LoadObjectPort: types.Int64Value(int64(lo.LoadObjectPort)),
	}
	if lo.LoadServers != nil {
		loadServers, diags := types.ListValueFrom(ctx, types.StringType, lo.LoadServers)
		if diags.HasError() {
			return loadObj, diags
		}
		loadObj.LoadServers = loadServers
	}
	return loadObj, nil
}

func populateGeographicMapAssignment(ctx context.Context, asg *gtm.GeoAssignment) (geographicMapAssignment, diag.Diagnostics) {
	result := geographicMapAssignment{}

	if asg.Countries != nil {
		countries, diags := types.SetValueFrom(ctx, types.StringType, asg.Countries)
		if diags.HasError() {
			return result, diags
		}
		result.Countries = countries
	}

	result.Nickname = types.StringValue(asg.DatacenterBase.Nickname)
	result.DatacenterID = types.Int64Value(int64(asg.DatacenterBase.DatacenterId))

	return result, nil
}

func populateCIDRMapAssignment(ctx context.Context, asg *gtm.CidrAssignment) (cidrMapAssignment, diag.Diagnostics) {
	result := cidrMapAssignment{}

	if asg.Blocks != nil {
		lstStr, diags := types.SetValueFrom(ctx, types.StringType, asg.Blocks)
		if diags.HasError() {
			return result, diags
		}
		result.Blocks = lstStr
	}

	result.Nickname = types.StringValue(asg.DatacenterBase.Nickname)
	result.DatacenterID = types.Int64Value(int64(asg.DatacenterBase.DatacenterId))

	return result, nil
}

func populateASMapAssignment(asg *gtm.AsAssignment) asMapAssignment {
	result := asMapAssignment{}

	if asg.AsNumbers != nil {
		asNumbers := make([]types.Int64, len(asg.AsNumbers))
		for i, c := range asg.AsNumbers {
			asNumbers[i] = types.Int64Value(c)
		}
		result.ASNumbers = asNumbers
	}

	result.Nickname = types.StringValue(asg.DatacenterBase.Nickname)
	result.DatacenterID = types.Int64Value(int64(asg.DatacenterBase.DatacenterId))

	return result
}
