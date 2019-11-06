provider "akamai" {
    edgerc = "/home/demo/.edgerc"
    gtm_section = "default"
}

locals {
  contract = "9-CONTRACT"
  group = "12345"
}


resource "akamai_gtm_domain" "tfexample_domain" {
    // 
    // Domain auto generated id format [name], e.g. tfexample.akadns.net
    //
    // Required
    name = "tfexample.akadns.net"
    type = "weighted"
    //
    // Computed -- DO NOT CONFIGURE
    // default_unreachable_threshold
    // min_pingable_region_fraction
    // servermonitor_liveness_count
    // round_robin_prefix
    // servermonitor_load_count
    // ping_interval
    // max_ttl
    // default_health_max
    // map_update_interval
    // max_properties
    // max_resources
    // max_test_timeout
    // default_health_multiplier
    // servermonitor_pool
    // min_ttl
    // default_max_unreachable_penalty
    // default_health_threshold
    // min_test_interval
    // ping_packet_size
    //
    //Optional - list all?
    contract = "${locals.contract}"
    group = "${locals.group}"
    email_notification_list = []
    load_imbalance_percentage = 20
    wait_on_complete = false
}

resource "akamai_gtm_datacenter" "tfexample_dc_1" {
    //
    // Datacenter auto generated id format [domain:datacenter_id], e.g. tfexample.akadns.net:3131
    //
    // Required
    domain = "${akamai_gtm_domain.tfexample_domain.name}"
    //
    // Computed - DO NOT CONFIGURE
    datacenter_id = 3131
    // ping_interval
    // ping_packet_size
    // score_penalty
    // servermonitor_liveness_count
    // servermonitor_load_count
    // servermonitor_pool
    //
    // Optional
    nickname = "tfexample_dc_1"
    wait_on_complete = false
    default_load_object = [{
        load_object = "test"
	load_object_port = 80 
	load_servers = ["1.2.3.4", "1.2.3.5"]
    }]
    depends_on = [
         "akamai_gtm_domain.tfexample_domain"
    ]
}

resource "akamai_gtm_datacenter" "tfexample_dc_2" {
    domain = "${akamai_gtm_domain.tfexample_domain.name}"
    nickname = "tfexample_dc_2"
    datacenter_id = 3132
    wait_on_complete = false
    //
    // Datacenters need strict dependencies for multiple creation since dcids are auto generated
    //
    depends_on = [
        "akamai_gtm_domain.tfexample_domain",
        "akamai_gtm_datacenter.tfexample_dc_1"
    ]
}

resource "akamai_gtm_property" "tfexample_prop_1" {
    //
    // Property auto generated id format [domain:name], e.g. tfexample.akadns.net:tfexample_prop_1
    //
    // Required
    domain = "${akamai_gtm_domain.tfexample_domain.name}"
    name = "tfexample_prop_1"
    type = "weighted-round-robin"
    score_aggregation_type = "median"
    handout_limit = 5
    handout_mode = "normal"
    //
    // Computed - DO NOT CONFIGURE
    // weighted_hash_bits_for_ipv4
    // weighted_hash_bits_for_ipv6
    //
    // Optional 
    failover_delay = 0
    failback_delay = 0
    static_ttl = 45
    traffic_targets = [{
	datacenter_id = "${akamai_gtm_datacenter.tfexample_dc_1.datacenter_id}"
	enabled = true 
	weight = 100
	servers = ["1.2.3.4"]
	// optional
    	name = ""
    	handout_cname = ""
	}]
    liveness_tests = [{
	name = "lt1"
	test_interval = 10
	test_object_protocol = "HTTP"
	test_timeout = 20
	// optional
	answer_required = false
	disable_nonstandard_port_warning = false
	error_penalty = 0
	host_header = ""
	http_error3xx = false
	http_error4xx = false
	http_error5xx = false
	peer_certificate_verification = false
	recursion_requested = false
	request_string = ""
	resource_type = ""
	response_string = ""
	ssl_client_certificate = ""
	ssl_client_private_key = ""
	test_object = "junk"
	test_object_password = ""
	test_object_port = 1
	test_object_username = ""
	timeout_penalty = 0
	}]
    /*
    mx_records = [{
	exchange = "test_e"
	preference = 100
	}]
    */
    wait_on_complete = false
    depends_on = [
         "akamai_gtm_domain.tfexample_domain",
	 "akamai_gtm_datacenter.tfexample_dc_1"
    ]
}

resource "akamai_gtm_resource" "tfexample_resource_1" {
    //
    // Resource auto generated id format [domain:name], e.g. tfexample.akadns.net:tfexample_resource_1
    //
    // Required
    domain = "${akamai_gtm_domain.tfexample_domain.name}"
    name = "tfexample_resource_1"
    aggregation_type = "latest"
    type = "XML load object via HTTP"
    //
    // Optional
    resource_instances = [{
        datacenter_id = "${akamai_gtm_datacenter.tfexample_dc_1.datacenter_id}"
        use_default_load_object = false
        load_object = "/test1"
        load_servers = ["1.2.3.4"]
        load_object_port = 80
        }]
    wait_on_complete = false
    depends_on = [
         "akamai_gtm_domain.tfexample_domain", "akamai_gtm_datacenter.tfexample_dc_1"
    ]
}

resource "akamai_gtm_resource" "tfexample_resource_2" {
    domain = "${akamai_gtm_domain.tfexample_domain.name}"
    name = "tfexample_resource_2"
    aggregation_type = "median"
    type = "XML load object via HTTP"
    resource_instances = [{
        datacenter_id = "${akamai_gtm_datacenter.tfexample_dc_2.datacenter_id}"
        use_default_load_object = false
        load_object = "/test"
        load_servers = ["1.2.3.4"]
        load_object_port = 80
        }]
    wait_on_complete = false
    depends_on = [
         "akamai_gtm_domain.tfexample_domain", "akamai_gtm_datacenter.tfexample_dc_2"
    ]
}

resource "akamai_gtm_cidrmap" "tfexample_cidr_1" {
    //
    // CIDRmap auto generated id format [domain:name], e.g tfexample.akadns.net:tfexample_cidr_1
    //
    // Required
    domain = "${akamai_gtm_domain.tfexample_domain.name}"
    name = "tfexample_cidr_1"
    default_datacenter = [{
        datacenter_id = 5400
        nickname = "All Other CIDR Blocks"
        }]
    //
    // Optional 
    assignments = [{
        datacenter_id = "${akamai_gtm_datacenter.tfexample_dc_1.datacenter_id}"
        nickname = "${akamai_gtm_datacenter.tfexample_dc_1.nickname}"
        // Optional
        blocks = ["1.2.3.9/24"]
        }]
    wait_on_complete = true
    depends_on = [
        "akamai_gtm_domain.tfexample_domain", 
	"akamai_gtm_datacenter.tfexample_dc_1"
    ]
}

resource "akamai_gtm_cidrmap" "tfexample_cidr_2" {
    domain = "${akamai_gtm_domain.tfexample_domain.name}"
    name = "tfexample_cidr_2"
    default_datacenter = [{
        datacenter_id = 5400
        nickname = "All Other CIDR Blocks"
        }]
    wait_on_complete = true
    depends_on = [
        "akamai_gtm_domain.tfexample_domain",
    ]
}

resource "akamai_gtm_asmap" "tfexample_as_1" {
    //
    // ASmap auto generated id format [domain:name], e.g. tfexample.akadns.net:tfexample_as_1
    //
    // Required
    domain = "${akamai_gtm_domain.tfexample_domain.name}"
    name = "tfexample_as_1"
    default_datacenter = [{
        datacenter_id = 5400
        nickname = "All Other AS numbers"
        }]
    //
    // Optional
    assignments = [{
        datacenter_id = "${akamai_gtm_datacenter.tfexample_dc_1.datacenter_id}"
        nickname = "${akamai_gtm_datacenter.tfexample_dc_1.nickname}"
        as_numbers = [12222, 16702,17334]
        },
	{
        datacenter_id = "${akamai_gtm_datacenter.tfexample_dc_2.datacenter_id}"
        nickname = "${akamai_gtm_datacenter.tfexample_dc_2.nickname}"
        as_numbers = [12229, 16703,17335]
        }]
    wait_on_complete = true
    depends_on = [
        "akamai_gtm_domain.tfexample_domain",
        "akamai_gtm_datacenter.tfexample_dc_1",
	"akamai_gtm_datacenter.tfexample_dc_2"
    ]
}

resource "akamai_gtm_geomap" "tfexample_geo_2" {
    //
    // Geomap auto generated id format [domain:name], e.g. tfexample.akadns.net:tfexample_geo_2
    //
    // Required
    domain = "${akamai_gtm_domain.tfexample_domain.name}"
    name = "tfexample_geo_2"
    default_datacenter = [{
        datacenter_id = 5400
        nickname = "All Others"
        }]
    //
    // Optional
    assignments = [{
        datacenter_id = "${akamai_gtm_datacenter.tfexample_dc_2.datacenter_id}"
        nickname = "${akamai_gtm_datacenter.tfexample_dc_2.nickname}"
        // Optional
        countries = ["GB"]
        }]
    wait_on_complete = true
    depends_on = [
        "akamai_gtm_domain.tfexample_domain",
        "akamai_gtm_datacenter.tfexample_dc_2"
    ]
}


