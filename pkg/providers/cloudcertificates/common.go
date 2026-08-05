package cloudcertificates

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	fwdate "github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model structs and mapping helpers in this file are shared by multiple data sources and resources
// in this package.

type (
	lineageSubjectModel struct {
		CommonName         types.String `tfsdk:"common_name"`
		Organization       types.String `tfsdk:"organization"`
		OrganizationalUnit types.String `tfsdk:"organizational_unit"`
		Country            types.String `tfsdk:"country"`
		State              types.String `tfsdk:"state"`
		Locality           types.String `tfsdk:"locality"`
	}

	// algorithmModel holds per key-type certificate details shared by generation-related data sources. The key
	// type itself is the map key of the containing algorithms map, not a field here.
	algorithmModel struct {
		AlgorithmInstanceID                 types.Int64  `tfsdk:"algorithm_instance_id"`
		AlgorithmInstanceCreatedBy          types.String `tfsdk:"algorithm_instance_created_by"`
		AlgorithmInstanceCreatedTime        types.String `tfsdk:"algorithm_instance_created_time"`
		AlgorithmInstanceModifiedBy         types.String `tfsdk:"algorithm_instance_modified_by"`
		AlgorithmInstanceModifiedTime       types.String `tfsdk:"algorithm_instance_modified_time"`
		CertificateStatus                   types.String `tfsdk:"certificate_status"`
		CSRExpirationDate                   types.String `tfsdk:"csr_expiration_date"`
		CSRPEM                              types.String `tfsdk:"csr_pem"`
		SignedCertificateIssuer             types.String `tfsdk:"signed_certificate_issuer"`
		SignedCertificateNotValidAfterDate  types.String `tfsdk:"signed_certificate_not_valid_after_date"`
		SignedCertificateNotValidBeforeDate types.String `tfsdk:"signed_certificate_not_valid_before_date"`
		SignedCertificatePEM                types.String `tfsdk:"signed_certificate_pem"`
		SignedCertificateSerialNumber       types.String `tfsdk:"signed_certificate_serial_number"`
		SignedCertificateSHA256Fingerprint  types.String `tfsdk:"signed_certificate_sha256_fingerprint"`
		TrustChainPEM                       types.String `tfsdk:"trust_chain_pem"`
	}

	// commonGenerationModel holds the generation fields embedded by the generation data source's model.
	commonGenerationModel struct {
		Algorithms                    map[string]algorithmModel `tfsdk:"algorithms"`
		FirstPromotedToProductionTime types.String              `tfsdk:"first_promoted_to_production_time"`
		GenerationCreatedBy           types.String              `tfsdk:"generation_created_by"`
		GenerationCreatedTime         types.String              `tfsdk:"generation_created_time"`
		GenerationModifiedBy          types.String              `tfsdk:"generation_modified_by"`
		GenerationModifiedTime        types.String              `tfsdk:"generation_modified_time"`
	}

	// generationPointerModel is shared by the head, current_production, previous_production, and current_staging attributes.
	generationPointerModel struct {
		GenerationID     types.Int64  `tfsdk:"generation_id"`
		GenerationStatus types.String `tfsdk:"generation_status"`
		commonGenerationModel
	}

	// commonLineageModel holds lineage fields shared by the lineage data sources.
	commonLineageModel struct {
		AccountID           types.String            `tfsdk:"account_id"`
		ContractID          types.String            `tfsdk:"contract_id"`
		GeoClass            types.String            `tfsdk:"geo_class"`
		GroupID             types.Int64             `tfsdk:"group_id"`
		LineageName         types.String            `tfsdk:"lineage_name"`
		LineageType         types.String            `tfsdk:"lineage_type"`
		SecureNetwork       types.String            `tfsdk:"secure_network"`
		StackMode           types.String            `tfsdk:"stack_mode"`
		SANs                []types.String          `tfsdk:"sans"`
		Subject             *lineageSubjectModel    `tfsdk:"subject"`
		KeySpecs            map[string]string       `tfsdk:"key_specs"`
		LineageCreatedBy    types.String            `tfsdk:"lineage_created_by"`
		LineageCreatedTime  types.String            `tfsdk:"lineage_created_time"`
		LineageModifiedBy   types.String            `tfsdk:"lineage_modified_by"`
		LineageModifiedTime types.String            `tfsdk:"lineage_modified_time"`
		Head                *generationPointerModel `tfsdk:"head"`
		CurrentProduction   *generationPointerModel `tfsdk:"current_production"`
		PreviousProduction  *generationPointerModel `tfsdk:"previous_production"`
		CurrentStaging      *generationPointerModel `tfsdk:"current_staging"`
		SigningTarget       *generationPointerModel `tfsdk:"signing_target"`
	}

	// commonActivationModel holds the activation fields shared by activationDetailModel (used by
	// `akamai_cloudcertificates_activation_status` and the list items of `akamai_cloudcertificates_activations`)
	// and networkActivationModel (used by `akamai_cloudcertificates_activation`).
	commonActivationModel struct {
		ActivationType          types.String `tfsdk:"activation_type"`
		ActivationStatus        types.String `tfsdk:"activation_status"`
		GenerationID            types.Int64  `tfsdk:"generation_id"`
		ActivationCreatedTime   types.String `tfsdk:"activation_created_time"`
		ActivationModifiedTime  types.String `tfsdk:"activation_modified_time"`
		CreatedBy               types.String `tfsdk:"created_by"`
		ModifiedBy              types.String `tfsdk:"modified_by"`
		TotalHostnameCount      types.Int64  `tfsdk:"total_hostname_count"`
		InProgressHostnameCount types.Int64  `tfsdk:"in_progress_hostname_count"`
		PreEmptedBy             types.Int64  `tfsdk:"pre_empted_by"`
		ErrorTypes              types.String `tfsdk:"error_types"`
	}

	// activationDetailModel holds the computed detail fields shared by `akamai_cloudcertificates_activation_status`
	// and the list items of `akamai_cloudcertificates_activations`.
	activationDetailModel struct {
		commonActivationModel
		TargetEnvironment types.String `tfsdk:"target_environment"`
	}

	// activationModel is used as a list item model for `akamai_cloudcertificates_activations`.
	// Both ID fields are computed outputs echoed from the API.
	activationModel struct {
		ActivationID types.Int64 `tfsdk:"activation_id"`
		LineageID    types.Int64 `tfsdk:"lineage_id"`
		activationDetailModel
	}
)

// attr.Type maps mirroring the model structs above, needed by the resource to build types.Object/types.Map values.
var (
	algorithmObjectAttrTypes = map[string]attr.Type{
		"algorithm_instance_id":                    types.Int64Type,
		"algorithm_instance_created_by":            types.StringType,
		"algorithm_instance_created_time":          types.StringType,
		"algorithm_instance_modified_by":           types.StringType,
		"algorithm_instance_modified_time":         types.StringType,
		"certificate_status":                       types.StringType,
		"csr_expiration_date":                      types.StringType,
		"csr_pem":                                  types.StringType,
		"signed_certificate_issuer":                types.StringType,
		"signed_certificate_not_valid_after_date":  types.StringType,
		"signed_certificate_not_valid_before_date": types.StringType,
		"signed_certificate_pem":                   types.StringType,
		"signed_certificate_serial_number":         types.StringType,
		"signed_certificate_sha256_fingerprint":    types.StringType,
		"trust_chain_pem":                          types.StringType,
	}

	generationPointerObjectAttrTypes = map[string]attr.Type{
		"generation_id":                     types.Int64Type,
		"generation_status":                 types.StringType,
		"algorithms":                        types.MapType{ElemType: types.ObjectType{AttrTypes: algorithmObjectAttrTypes}},
		"first_promoted_to_production_time": types.StringType,
		"generation_created_by":             types.StringType,
		"generation_created_time":           types.StringType,
		"generation_modified_by":            types.StringType,
		"generation_modified_time":          types.StringType,
	}

	subjectObjectAttrTypes = map[string]attr.Type{
		"common_name":         types.StringType,
		"organization":        types.StringType,
		"organizational_unit": types.StringType,
		"country":             types.StringType,
		"state":               types.StringType,
		"locality":            types.StringType,
	}
)

// mapCommonGeneration maps a generation's fields for the generation data source, which flattens them directly
// into its own model via struct embedding.
func mapCommonGeneration(g cloudcertificates.Generation) commonGenerationModel {
	return commonGenerationModel{
		Algorithms:                    mapLineageAlgorithms(g.Algorithms),
		FirstPromotedToProductionTime: fwdate.TimeRFC3339PointerValue(g.FirstPromotedToProductionTime),
		GenerationCreatedBy:           types.StringPointerValue(g.GenerationCreatedBy),
		GenerationCreatedTime:         fwdate.TimeRFC3339PointerValue(g.GenerationCreatedTime),
		GenerationModifiedBy:          types.StringPointerValue(g.GenerationModifiedBy),
		GenerationModifiedTime:        fwdate.TimeRFC3339PointerValue(g.GenerationModifiedTime),
	}
}

// mapGenerationPointer maps the fields shared by every generation pointer. The server only ever returns these
// fields together, as a group, when the pointer was requested via expandGenerations - each field is nil when
// not expanded (or, for GenerationModifiedTime/FirstPromotedToProductionTime, when the generation was expanded
// but has never been modified/promoted).
func mapGenerationPointer(g cloudcertificates.GenerationPointer) *generationPointerModel {
	return &generationPointerModel{
		GenerationID:          types.Int64Value(g.ID()),
		GenerationStatus:      types.StringValue(g.Status()),
		commonGenerationModel: mapCommonGeneration(g.Common()),
	}
}

// mapCommonLineage maps a lineage's fields shared by the singular and plural lineage data sources.
func mapCommonLineage(lineage cloudcertificates.Lineage) commonLineageModel {
	data := commonLineageModel{
		AccountID:           types.StringValue(lineage.AccountID),
		ContractID:          types.StringValue(lineage.ContractID),
		GeoClass:            types.StringValue(lineage.GeoClass),
		GroupID:             types.Int64Value(lineage.GroupID),
		LineageName:         types.StringValue(lineage.LineageName),
		LineageType:         types.StringValue(lineage.LineageType),
		SecureNetwork:       types.StringValue(lineage.SecureNetwork),
		StackMode:           types.StringValue(lineage.StackMode),
		LineageCreatedBy:    types.StringValue(lineage.LineageCreatedBy),
		LineageCreatedTime:  fwdate.TimeRFC3339Value(lineage.LineageCreatedTime),
		LineageModifiedBy:   types.StringValue(lineage.LineageModifiedBy),
		LineageModifiedTime: fwdate.TimeRFC3339Value(lineage.LineageModifiedTime),
		Subject: &lineageSubjectModel{
			CommonName:         tf.StringValueOrNullIfEmpty(lineage.Subject.CommonName),
			Organization:       tf.StringValueOrNullIfEmpty(lineage.Subject.Organization),
			OrganizationalUnit: tf.StringValueOrNullIfEmpty(lineage.Subject.OrganizationalUnit),
			Country:            tf.StringValueOrNullIfEmpty(lineage.Subject.Country),
			State:              tf.StringValueOrNullIfEmpty(lineage.Subject.State),
			Locality:           tf.StringValueOrNullIfEmpty(lineage.Subject.Locality),
		},
	}

	data.SANs = make([]types.String, len(lineage.SANs))
	for i, san := range lineage.SANs {
		data.SANs[i] = types.StringValue(san)
	}

	data.KeySpecs = make(map[string]string, len(lineage.KeySpecs))
	for _, keySpec := range lineage.KeySpecs {
		data.KeySpecs[keySpec.KeyType] = keySpec.KeySize
	}

	if lineage.Head != nil {
		data.Head = mapGenerationPointer(lineage.Head)
	}

	if lineage.CurrentProduction != nil {
		data.CurrentProduction = mapGenerationPointer(lineage.CurrentProduction)
	}

	if lineage.PreviousProduction != nil {
		data.PreviousProduction = mapGenerationPointer(lineage.PreviousProduction)
	}

	if lineage.CurrentStaging != nil {
		data.CurrentStaging = mapGenerationPointer(lineage.CurrentStaging)
	}

	switch signingTarget(lineage) {
	case lineage.Head:
		data.SigningTarget = data.Head
	case lineage.CurrentProduction:
		data.SigningTarget = data.CurrentProduction
	}

	return data
}

// signingTarget picks the generation whose CSR currently needs to be signed and uploaded: the lineage's head
// generation if one exists, otherwise its current production generation, otherwise neither. Not a distinct API
// concept - a single, shared source of truth for this derived convenience, used by every lineage-shaped
// resource/data source that surfaces a signing_target attribute, so the definition can't drift between them.
func signingTarget(lineage cloudcertificates.Lineage) cloudcertificates.GenerationPointer {
	switch {
	case lineage.Head != nil:
		return lineage.Head
	case lineage.CurrentProduction != nil:
		return lineage.CurrentProduction
	default:
		return nil
	}
}

func mapLineageAlgorithms(algorithms []cloudcertificates.Algorithm) map[string]algorithmModel {
	// Return nil, not an empty map from make(), when there are no algorithms. Callers then convert this to a
	// null map/object attribute instead of a known-empty one. The API omits algorithms entirely for a
	// generation pointer that wasn't requested via expandGenerations. That "not requested" state must stay
	// distinguishable from an actual empty map.
	if len(algorithms) == 0 {
		return nil
	}
	result := make(map[string]algorithmModel, len(algorithms))
	for _, a := range algorithms {
		result[a.KeyType] = algorithmModel{
			AlgorithmInstanceID:                 types.Int64Value(a.AlgorithmInstanceID),
			AlgorithmInstanceCreatedBy:          tf.StringValueOrNullIfEmpty(a.AlgorithmInstanceCreatedBy),
			AlgorithmInstanceCreatedTime:        fwdate.TimeRFC3339PointerValue(a.AlgorithmInstanceCreatedTime),
			AlgorithmInstanceModifiedBy:         types.StringPointerValue(a.AlgorithmInstanceModifiedBy),
			AlgorithmInstanceModifiedTime:       fwdate.TimeRFC3339PointerValue(a.AlgorithmInstanceModifiedTime),
			CertificateStatus:                   types.StringValue(a.CertificateStatus),
			CSRExpirationDate:                   fwdate.TimeRFC3339PointerValue(a.CSRExpirationDate),
			CSRPEM:                              tf.StringValueOrNullIfEmpty(a.CSRPEM),
			SignedCertificateIssuer:             types.StringPointerValue(a.SignedCertificateIssuer),
			SignedCertificateNotValidAfterDate:  fwdate.TimeRFC3339PointerValue(a.SignedCertificateNotValidAfterDate),
			SignedCertificateNotValidBeforeDate: fwdate.TimeRFC3339PointerValue(a.SignedCertificateNotValidBeforeDate),
			SignedCertificatePEM:                types.StringPointerValue(a.SignedCertificatePEM),
			SignedCertificateSerialNumber:       types.StringPointerValue(a.SignedCertificateSerialNumber),
			SignedCertificateSHA256Fingerprint:  types.StringPointerValue(a.SignedCertificateSHA256Fingerprint),
			TrustChainPEM:                       types.StringPointerValue(a.TrustChainPEM),
		}
	}
	return result
}

// algorithmNestedAttributes returns the schema attributes for a certificate algorithm instance. key_type is not
// included here: it's the map key of the containing algorithms attribute, not a nested field.
func algorithmNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"algorithm_instance_id": schema.Int64Attribute{
			Computed:    true,
			Description: "Unique identifier of the algorithm instance.",
		},
		"algorithm_instance_created_by": schema.StringAttribute{
			Computed:    true,
			Description: "Username of the person who created the algorithm instance.",
		},
		"algorithm_instance_created_time": schema.StringAttribute{
			Computed:    true,
			Description: "Time the algorithm instance was created, in RFC3339 format.",
		},
		"algorithm_instance_modified_by": schema.StringAttribute{
			Computed:    true,
			Description: "Username of the person who last modified the algorithm instance. Null if never modified.",
		},
		"algorithm_instance_modified_time": schema.StringAttribute{
			Computed:    true,
			Description: "Time the algorithm instance was last modified, in RFC3339 format. Null if never modified.",
		},
		"certificate_status": schema.StringAttribute{
			Computed:    true,
			Description: "Status of the certificate for this key type: `CSR_READY`, `CERT_UPLOAD_PROCESSING`, `READY_FOR_USE`, or `ABANDONED`.",
		},
		"csr_expiration_date": schema.StringAttribute{
			Computed:    true,
			Description: "Date when the CSR expires, in RFC3339 format.",
		},
		"csr_pem": schema.StringAttribute{
			Computed:    true,
			Description: "PEM-encoded certificate signing request.",
		},
		"signed_certificate_issuer": schema.StringAttribute{
			Computed:    true,
			Description: "Issuer field of the signed certificate. Null until a certificate is uploaded.",
		},
		"signed_certificate_not_valid_after_date": schema.StringAttribute{
			Computed:    true,
			Description: "Expiration date of the signed certificate, in RFC3339 format. Null until a certificate is uploaded.",
		},
		"signed_certificate_not_valid_before_date": schema.StringAttribute{
			Computed:    true,
			Description: "Start of validity of the signed certificate, in RFC3339 format. Null until a certificate is uploaded.",
		},
		"signed_certificate_pem": schema.StringAttribute{
			Computed:    true,
			Description: "PEM-encoded signed certificate. Null until a certificate is uploaded.",
		},
		"signed_certificate_serial_number": schema.StringAttribute{
			Computed:    true,
			Description: "Serial number of the signed certificate in hex format. Null until a certificate is uploaded.",
		},
		"signed_certificate_sha256_fingerprint": schema.StringAttribute{
			Computed:    true,
			Description: "SHA-256 fingerprint of the signed certificate. Null until a certificate is uploaded.",
		},
		"trust_chain_pem": schema.StringAttribute{
			Computed:    true,
			Description: "PEM-encoded trust chain uploaded alongside the signed certificate. Null if none was uploaded.",
		},
	}
}

// lineageCommonAttributes returns the computed schema attributes shared by lineage data sources.
func lineageCommonAttributes() map[string]schema.Attribute {
	genPointerAttrs := mergeAttributes(generationCommonAttributes(), map[string]schema.Attribute{
		"generation_id": schema.Int64Attribute{
			Computed:    true,
			Description: "Unique identifier of the generation.",
		},
		"generation_status": schema.StringAttribute{
			Computed:    true,
			Description: "Status of the generation.",
		},
	})

	return map[string]schema.Attribute{
		"account_id": schema.StringAttribute{
			Computed:    true,
			Description: "Account identifier associated with the contract.",
		},
		"contract_id": schema.StringAttribute{
			Computed:    true,
			Description: "Contract identifier under which the lineage was created.",
		},
		"geo_class": schema.StringAttribute{
			Computed:    true,
			Description: "Geographic class of the certificate.",
		},
		"group_id": schema.Int64Attribute{
			Computed:    true,
			Description: "Unique identifier of the group.",
		},
		"lineage_name": schema.StringAttribute{
			Computed:    true,
			Description: "Name of the lineage.",
		},
		"lineage_type": schema.StringAttribute{
			Computed:    true,
			Description: "Type of the lineage, e.g. MULTIPLE_GENERATION or SINGLE_GENERATION.",
		},
		"secure_network": schema.StringAttribute{
			Computed:    true,
			Description: "Secure network type, e.g. ENHANCED_TLS or STANDARD_TLS.",
		},
		"stack_mode": schema.StringAttribute{
			Computed:    true,
			Description: "Stack mode of the lineage, e.g. SINGLE_STACK or MULTIPLE_STACK.",
		},
		"sans": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Subject Alternative Names (SANs) for the certificate.",
		},
		"subject": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "X.509 subject fields of the certificate. All fields are null if no subject was provided when the lineage was created.",
			Attributes: map[string]schema.Attribute{
				"common_name": schema.StringAttribute{
					Computed:    true,
					Description: "Common name (CN). Null if not provided.",
				},
				"organization": schema.StringAttribute{
					Computed:    true,
					Description: "Organization (O). Null if not provided.",
				},
				"organizational_unit": schema.StringAttribute{
					Computed:    true,
					Description: "Organizational unit (OU). Null if not provided.",
				},
				"country": schema.StringAttribute{
					Computed:    true,
					Description: "Two-letter ISO 3166 country code (C). Null if not provided.",
				},
				"state": schema.StringAttribute{
					Computed:    true,
					Description: "State or province name (ST). Null if not provided.",
				},
				"locality": schema.StringAttribute{
					Computed:    true,
					Description: "Locality or city name (L). Null if not provided.",
				},
			},
		},
		"key_specs": schema.MapAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Key specifications declared for the lineage, mapping key_type (e.g. RSA or ECDSA) to key_size (e.g. 2048 or P-256).",
		},
		"lineage_created_by": schema.StringAttribute{
			Computed:    true,
			Description: "Username of the person who created the lineage.",
		},
		"lineage_created_time": schema.StringAttribute{
			Computed:    true,
			Description: "Time the lineage was created, in RFC3339 format.",
		},
		"lineage_modified_by": schema.StringAttribute{
			Computed:    true,
			Description: "Username of the person who last modified the lineage.",
		},
		"lineage_modified_time": schema.StringAttribute{
			Computed:    true,
			Description: "Time the lineage was last modified, in RFC3339 format.",
		},
		"head": schema.SingleNestedAttribute{
			Computed: true,
			Description: "Head generation of the lineage. Null when no head generation exists. " +
				"Populated only when expand_generations is true.",
			Attributes: genPointerAttrs,
		},
		"current_production": schema.SingleNestedAttribute{
			Computed: true,
			Description: "Generation currently deployed to the production network. Null when no current production generation exists. " +
				"Populated only when expand_generations is true.",
			Attributes: genPointerAttrs,
		},
		"previous_production": schema.SingleNestedAttribute{
			Computed: true,
			Description: "Generation previously deployed to production (rollback candidate). Null when no previous production generation exists. " +
				"Populated only when expand_generations is true.",
			Attributes: genPointerAttrs,
		},
		"current_staging": schema.SingleNestedAttribute{
			Computed: true,
			Description: "Generation currently deployed to the staging network. Null when no current staging generation exists. " +
				"Populated only when expand_generations is true.",
			Attributes: genPointerAttrs,
		},
		"signing_target": schema.SingleNestedAttribute{
			Computed: true,
			Description: "Derived convenience field, not a distinct API concept: the generation whose CSR currently needs to be " +
				"signed and uploaded - the lineage's head generation if one exists, otherwise its current production generation " +
				"(e.g. while completing a MULTIPLE_STACK lineage's second algorithm). It is null if neither exists or if expand_generations is false.",
			Attributes: genPointerAttrs,
		},
	}
}

// generationCommonAttributes returns the schema attributes shared by all generation pointer types.
func generationCommonAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"algorithms": schema.MapNestedAttribute{
			Computed: true,
			Description: "Per key-type (RSA or ECDSA) certificate details for this generation, keyed by key_type. " +
				"Populated only when expand_generations is true.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: algorithmNestedAttributes(),
			},
		},
		"first_promoted_to_production_time": schema.StringAttribute{
			Computed: true,
			Description: "Time the generation was first promoted to production, in RFC3339 format. " +
				"Null if never promoted, or if not requested via expand_generations.",
		},
		"generation_created_by": schema.StringAttribute{
			Computed: true,
			Description: "Username of the person who created this generation. " +
				"Populated only when expand_generations is true.",
		},
		"generation_created_time": schema.StringAttribute{
			Computed: true,
			Description: "Time the generation was created, in RFC3339 format. " +
				"Populated only when expand_generations is true.",
		},
		"generation_modified_by": schema.StringAttribute{
			Computed: true,
			Description: "Username of the person who last modified this generation. " +
				"Populated only when expand_generations is true.",
		},
		"generation_modified_time": schema.StringAttribute{
			Computed: true,
			Description: "Time the generation was last modified, in RFC3339 format. Null if not requested via " +
				"expand_generations, or if the generation has never been modified since creation.",
		},
	}
}

// mergeAttributes merges two attribute maps into a new map. Keys from extra override base on conflict.
func mergeAttributes(base, extra map[string]schema.Attribute) map[string]schema.Attribute {
	merged := make(map[string]schema.Attribute, len(base)+len(extra))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range extra {
		merged[k] = v
	}
	return merged
}

// activationDetailAttributes returns the schema attributes for the computed detail fields of an
// activation request. It is used by both `akamai_cloudcertificates_activation_status` and the nested
// list items of `akamai_cloudcertificates_activations`.
func activationDetailAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"activation_type": schema.StringAttribute{
			Computed:    true,
			Description: "The type of the activation operation: `PROMOTE`, `ROLLBACK`, or `REPLACE_STAGING`.",
		},
		"activation_status": schema.StringAttribute{
			Computed:    true,
			Description: "The status of the activation request: `INIT`, `PENDING`, `IN_PROGRESS`, `COMPLETE`, `PARTIAL_SUCCESS`, `FAILED`, or `ABORTED`.",
		},
		"generation_id": schema.Int64Attribute{
			Computed:    true,
			Description: "The unique identifier of the generation being activated.",
		},
		"activation_created_time": schema.StringAttribute{
			Computed:    true,
			Description: "The time the activation request was created.",
		},
		"activation_modified_time": schema.StringAttribute{
			Computed:    true,
			Description: "The time the activation request was last modified.",
		},
		"target_environment": schema.StringAttribute{
			Computed:    true,
			Description: "The target network for the generation activation: `STAGING` or `PRODUCTION`.",
		},
		"total_hostname_count": schema.Int64Attribute{
			Computed:    true,
			Description: "The total number of hostnames being deployed as part of this activation, or null if not yet known.",
		},
		"in_progress_hostname_count": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of hostnames still in progress for this activation, or null if not yet known.",
		},
		"pre_empted_by": schema.Int64Attribute{
			Computed:    true,
			Description: "The activation request that pre-empted (superseded) this one, or null if this activation was not pre-empted.",
		},
		"created_by": schema.StringAttribute{
			Computed:    true,
			Description: "The user who created the activation request.",
		},
		"modified_by": schema.StringAttribute{
			Computed:    true,
			Description: "The user who last modified the activation request.",
		},
		"error_types": schema.StringAttribute{
			Computed:    true,
			Description: "Error type information when the activation failed, or null otherwise.",
		},
	}
}

// algorithmResourceAttributes is algorithmNestedAttributes' resource-schema counterpart, shared by the lineage
// and upload resources. terraform-plugin-framework's resource and data source schemas are distinct Go types, so
// the data source version can't be reused here directly.
func algorithmResourceAttributes() map[string]resourceschema.Attribute {
	return map[string]resourceschema.Attribute{
		"algorithm_instance_id": resourceschema.Int64Attribute{
			Computed:    true,
			Description: "Unique identifier of the algorithm instance.",
		},
		"algorithm_instance_created_by": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Username of the person who created the algorithm instance.",
		},
		"algorithm_instance_created_time": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Time the algorithm instance was created, in RFC3339 format.",
		},
		"algorithm_instance_modified_by": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Username of the person who last modified the algorithm instance. Null if never modified.",
		},
		"algorithm_instance_modified_time": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Time the algorithm instance was last modified, in RFC3339 format. Null if never modified.",
		},
		"certificate_status": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Status of the certificate for this key type: `CSR_READY`, `CERT_UPLOAD_PROCESSING`, `READY_FOR_USE`, or `ABANDONED`.",
		},
		"csr_expiration_date": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Date when the CSR expires, in RFC3339 format.",
		},
		"csr_pem": resourceschema.StringAttribute{
			Computed:    true,
			Description: "PEM-encoded certificate signing request.",
		},
		"signed_certificate_issuer": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Issuer field of the signed certificate. Null until a certificate is uploaded.",
		},
		"signed_certificate_not_valid_after_date": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Expiration date of the signed certificate, in RFC3339 format. Null until a certificate is uploaded.",
		},
		"signed_certificate_not_valid_before_date": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Start of validity of the signed certificate, in RFC3339 format. Null until a certificate is uploaded.",
		},
		"signed_certificate_pem": resourceschema.StringAttribute{
			Computed:    true,
			Description: "PEM-encoded signed certificate. Null until a certificate is uploaded.",
		},
		"signed_certificate_serial_number": resourceschema.StringAttribute{
			Computed:    true,
			Description: "Serial number of the signed certificate in hex format. Null until a certificate is uploaded.",
		},
		"signed_certificate_sha256_fingerprint": resourceschema.StringAttribute{
			Computed:    true,
			Description: "SHA-256 fingerprint of the signed certificate. Null until a certificate is uploaded.",
		},
		"trust_chain_pem": resourceschema.StringAttribute{
			Computed:    true,
			Description: "PEM-encoded trust chain uploaded alongside the signed certificate. Null if none was uploaded.",
		},
	}
}

// mapActivationCore maps the activation fields shared by activationDetailModel and networkActivationModel.
func mapActivationCore(a cloudcertificates.GetActivationStatusResponse) commonActivationModel {
	return commonActivationModel{
		ActivationType:          types.StringValue(a.ActivationType),
		ActivationStatus:        types.StringValue(a.ActivationStatus),
		GenerationID:            types.Int64Value(a.GenerationID),
		ActivationCreatedTime:   fwdate.TimeRFC3339Value(a.ActivationCreatedTime),
		ActivationModifiedTime:  fwdate.TimeRFC3339Value(a.ActivationModifiedTime),
		CreatedBy:               types.StringValue(a.CreatedBy),
		ModifiedBy:              types.StringValue(a.ModifiedBy),
		TotalHostnameCount:      types.Int64PointerValue(a.TotalHostnameCount),
		InProgressHostnameCount: types.Int64PointerValue(a.InProgressHostnameCount),
		PreEmptedBy:             types.Int64PointerValue(a.PreEmptedBy),
		ErrorTypes:              tf.StringValueOrNullIfEmpty(a.ErrorTypes),
	}
}

// mapActivationDetail maps the computed detail fields of a GetActivationStatusResponse.
func mapActivationDetail(a cloudcertificates.GetActivationStatusResponse) activationDetailModel {
	return activationDetailModel{
		commonActivationModel: mapActivationCore(a),
		TargetEnvironment:     types.StringValue(a.TargetEnvironment),
	}
}

// mapActivation maps a GetActivationStatusResponse to an activationModel.
func mapActivation(a cloudcertificates.GetActivationStatusResponse) activationModel {
	return activationModel{
		ActivationID:          types.Int64Value(a.ActivationID),
		LineageID:             types.Int64Value(a.LineageID),
		activationDetailModel: mapActivationDetail(a),
	}
}

// mergeResourceAttributes is mergeAttributes' resource-schema counterpart. Keys from extra override base on conflict.
func mergeResourceAttributes(base, extra map[string]resourceschema.Attribute) map[string]resourceschema.Attribute {
	merged := make(map[string]resourceschema.Attribute, len(base)+len(extra))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range extra {
		merged[k] = v
	}
	return merged
}
