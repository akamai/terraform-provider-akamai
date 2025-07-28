package mtlskeystore

import (
	"context"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clientCertificateDataSource{}
	_ datasource.DataSourceWithConfigure = &clientCertificateDataSource{}
)

type (
	clientCertificateDataSource struct {
		meta meta.Meta
	}

	clientCertificateDataSourceModel struct {
		CertificateID               types.Int64    `tfsdk:"certificate_id"`
		IncludeAssociatedProperties types.Bool     `tfsdk:"include_associated_properties"`
		CertificateName             types.String   `tfsdk:"certificate_name"`
		CreatedBy                   types.String   `tfsdk:"created_by"`
		CreatedDate                 types.String   `tfsdk:"created_date"`
		Geography                   types.String   `tfsdk:"geography"`
		KeyAlgorithm                types.String   `tfsdk:"key_algorithm"`
		NotificationEmails          types.Set      `tfsdk:"notification_emails"`
		SecureNetwork               types.String   `tfsdk:"secure_network"`
		Signer                      types.String   `tfsdk:"signer"`
		Subject                     types.String   `tfsdk:"subject"`
		Versions                    []versionModel `tfsdk:"versions"`
		Current                     *versionModel  `tfsdk:"current"`
		Previous                    *versionModel  `tfsdk:"previous"`
	}

	versionModel struct {
		Version                  types.Int64            `tfsdk:"version"`
		VersionGUID              types.String           `tfsdk:"version_guid"`
		Status                   types.String           `tfsdk:"status"`
		CreatedBy                types.String           `tfsdk:"created_by"`
		CreatedDate              types.String           `tfsdk:"created_date"`
		ExpiryDate               types.String           `tfsdk:"expiry_date"`
		Issuer                   types.String           `tfsdk:"issuer"`
		KeyAlgorithm             types.String           `tfsdk:"key_algorithm"`
		EllipticCurve            types.String           `tfsdk:"elliptic_curve"`
		KeySizeInBytes           types.String           `tfsdk:"key_size_in_bytes"`
		SignatureAlgorithm       types.String           `tfsdk:"signature_algorithm"`
		Subject                  types.String           `tfsdk:"subject"`
		CertificateBlock         *certificateBlockModel `tfsdk:"certificate_block"`
		CertificateSubmittedBy   types.String           `tfsdk:"certificate_submitted_by"`
		CertificateSubmittedDate types.String           `tfsdk:"certificate_submitted_date"`
		CSRBlock                 *csrBlockModel         `tfsdk:"csr_block"`
		DeleteRequestedDate      types.String           `tfsdk:"delete_requested_date"`
		ScheduledDeleteDate      types.String           `tfsdk:"scheduled_delete_date"`
		IssuedDate               types.String           `tfsdk:"issued_date"`
		Properties               []propertyModel        `tfsdk:"properties"`
		Validation               validationModel        `tfsdk:"validation"`
	}

	certificateBlockModel struct {
		Certificate  types.String `tfsdk:"certificate"`
		KeyAlgorithm types.String `tfsdk:"key_algorithm"`
		TrustChain   types.String `tfsdk:"trust_chain"`
	}

	csrBlockModel struct {
		CSR          types.String `tfsdk:"csr"`
		KeyAlgorithm types.String `tfsdk:"key_algorithm"`
	}

	propertyModel struct {
		AssetID         types.Int64  `tfsdk:"asset_id"`
		GroupID         types.Int64  `tfsdk:"group_id"`
		PropertyName    types.String `tfsdk:"property_name"`
		PropertyVersion types.Int64  `tfsdk:"property_version"`
	}

	validationModel struct {
		Errors   []validationErrorModel `tfsdk:"errors"`
		Warnings []validationErrorModel `tfsdk:"warnings"`
	}

	validationErrorModel struct {
		Message string `tfsdk:"message"`
		Reason  string `tfsdk:"reason"`
		Type    string `tfsdk:"type"`
	}
)

// NewClientCertificateDataSource returns a new mtls keystore client certificate data source.
func NewClientCertificateDataSource() datasource.DataSource {
	return &clientCertificateDataSource{}
}

// Metadata configures data source's meta information.
func (d *clientCertificateDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlskeystore_client_certificate"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *clientCertificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
					req.ProviderData))
		}
	}()
	d.meta = meta.Must(req.ProviderData)
}

// Schema is used to define data source's terraform schema.
func (d *clientCertificateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	versionSchema := map[string]schema.Attribute{
		"version": schema.Int64Attribute{
			Description: "The unique identifier of the client certificate version.",
			Computed:    true,
		},
		"version_guid": schema.StringAttribute{
			Description: "Unique identifier for the client certificate version. Use it to configure mutual authentication (mTLS) sessions between the origin and edge servers in Property Manager's Mutual TLS Origin Keystore behavior.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The client certificate version status. Possible values: `AWAITING_SIGNED_CERTIFICATE`, `DEPLOYMENT_PENDING`, `DEPLOYED`, or `DELETE_PENDING`.",
			Computed:    true,
		},
		"created_by": schema.StringAttribute{
			Description: "The user who created the client certificate version.",
			Computed:    true,
		},
		"created_date": schema.StringAttribute{
			Description: "An ISO 8601 timestamp indicating the client certificate version's creation.",
			Computed:    true,
		},
		"expiry_date": schema.StringAttribute{
			Description: "An ISO 8601 timestamp indicating when the client certificate version expires.",
			Computed:    true,
		},
		"issuer": schema.StringAttribute{
			Description: "The signing entity of the client certificate version.",
			Computed:    true,
		},
		"key_algorithm": schema.StringAttribute{
			Description: "Identifies the client certificate version's encryption algorithm. Supported values are `RSA` and `ECDSA`.",
			Computed:    true,
		},
		"elliptic_curve": schema.StringAttribute{
			Description: "Specifies the key elliptic curve when the key algorithm `ECDSA` is used.",
			Computed:    true,
		},
		"key_size_in_bytes": schema.StringAttribute{
			Description: "The private key length of the client certificate version when the key algorithm `RSA` is used.",
			Computed:    true,
		},
		"signature_algorithm": schema.StringAttribute{
			Description: "Specifies the algorithm that secures the data exchange between the edge server and origin.",
			Computed:    true,
		},
		"subject": schema.StringAttribute{
			Description: "The public key's entity stored in the client certificate version's subject public key field.",
			Computed:    true,
		},
		"certificate_block": schema.SingleNestedAttribute{
			Description: "Details of the certificate block for the client certificate version.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"certificate": schema.StringAttribute{
					Description: "A text representation of the client certificate in PEM format.",
					Computed:    true,
				},
				"key_algorithm": schema.StringAttribute{
					Description: "Identifies the CA certificate's encryption algorithm. Possible values: `RSA` or `ECDSA`.",
					Computed:    true,
				},
				"trust_chain": schema.StringAttribute{
					Description: "A text representation of the trust chain in PEM format.",
					Computed:    true,
				},
			},
		},
		"certificate_submitted_by": schema.StringAttribute{
			Description: "The user who uploaded the `THIRD_PARTY` client certificate version.",
			Computed:    true,
		},
		"certificate_submitted_date": schema.StringAttribute{
			Description: "An ISO 8601 timestamp indicating when the `THIRD_PARTY` signer client certificate version was uploaded.",
			Computed:    true,
		},
		"csr_block": schema.SingleNestedAttribute{
			Description: "Details of the Certificate Signing Request (CSR) for the client certificate version.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"csr": schema.StringAttribute{
					Description: "Text of the certificate signing request.",
					Computed:    true,
				},
				"key_algorithm": schema.StringAttribute{
					Description: "Identifies the CA certificate's encryption algorithm. Possible values: `RSA` or `ECDSA`.",
					Computed:    true,
				},
			},
		},
		"delete_requested_date": schema.StringAttribute{
			Description: "An ISO 8601 timestamp indicating the client certificate version's deletion request.",
			Computed:    true,
		},
		"scheduled_delete_date": schema.StringAttribute{
			Description: "An ISO 8601 timestamp indicating the client certificate version's scheduled deletion.",
			Computed:    true,
		},
		"issued_date": schema.StringAttribute{
			Description: "An ISO 8601 timestamp indicating the client certificate version's availability.",
			Computed:    true,
		},
		"properties": schema.ListNestedAttribute{
			Description: "A list of properties associated with the client certificate.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"asset_id": schema.Int64Attribute{
						Description: "The unique identifier of the asset.",
						Required:    true,
					},
					"group_id": schema.Int64Attribute{
						Description: "The unique identifier of the group.",
						Required:    true,
					},
					"property_name": schema.StringAttribute{
						Description: "The name of the property.",
						Required:    true,
					},
					"property_version": schema.Int64Attribute{
						Description: "The version of the property.",
						Required:    true,
					},
				},
			},
		},
		"validation": schema.SingleNestedAttribute{
			Description: "Validation results for the client certificate version.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"errors": schema.ListNestedAttribute{
					Description: "Validation errors that need to be resolved for the request to succeed.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"message": schema.StringAttribute{
								Description: "Specifies the error details.",
								Computed:    true,
							},
							"reason": schema.StringAttribute{
								Description: "Specifies the error root cause.",
								Computed:    true,
							},
							"type": schema.StringAttribute{
								Description: "Specifies the error category.",
								Computed:    true,
							},
						},
					},
				},
				"warnings": schema.ListNestedAttribute{
					Description: "Validation warnings that can be resolved.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"message": schema.StringAttribute{
								Description: "Specifies the warning details.",
								Computed:    true,
							},
							"reason": schema.StringAttribute{
								Description: "Specifies the warning root cause.",
								Computed:    true,
							},
							"type": schema.StringAttribute{
								Description: "Specifies the warning category.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Retrieve client certificate with its versions.",
		Attributes: map[string]schema.Attribute{
			"certificate_id": schema.Int64Attribute{
				Required:    true,
				Description: "Identifies each client certificate.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"include_associated_properties": schema.BoolAttribute{
				Optional:    true,
				Description: "If set to true will list associated properties to that certificate version.",
			},
			"certificate_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the client certificate.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the client certificate.",
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "An ISO 8601 timestamp indicating the client certificate's creation.",
			},
			"geography": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies the type of network to deploy the client certificate. Possible values: `CORE`, `RUSSIA_AND_CORE`, or `CHINA_AND_CORE`.",
			},
			"key_algorithm": schema.StringAttribute{
				Computed:    true,
				Description: "The cryptographic algorithm used for key generation. Possible values: `RSA` or `ECDSA`.",
			},
			"notification_emails": schema.SetAttribute{
				Computed:    true,
				Description: "The email addresses to notify for client certificate-related issues.",
				ElementType: types.StringType,
			},
			"secure_network": schema.StringAttribute{
				Computed:    true,
				Description: "Identifies the network deployment type. Possible values: `STANDARD_TLS` or `ENHANCED_TLS`.",
			},
			"signer": schema.StringAttribute{
				Computed:    true,
				Description: "The signing entity of the client certificate. Possible values: `AKAMAI` or `THIRD_PARTY`.",
			},
			"subject": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies the client certificate. The `CN` attribute is required and is included in the subject.",
			},
			"previous": schema.SingleNestedAttribute{
				Description: "Details of the previous client certificate version.",
				Computed:    true,
				Attributes:  versionSchema,
			},
			"current": schema.SingleNestedAttribute{
				Description: "Details of the current client certificate version.",
				Computed:    true,
				Attributes:  versionSchema,
			},
			"versions": schema.ListNestedAttribute{
				Description: "A list of client certificate versions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: versionSchema,
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *clientCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificate Data Source Read")

	var data clientCertificateDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := Client(d.meta)

	clientCertificate, err := client.GetClientCertificate(ctx, mtlskeystore.GetClientCertificateRequest{
		CertificateID: data.CertificateID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Read Client Certificate failed", err.Error())
		return
	}

	versionsResponse, err := client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID:               data.CertificateID.ValueInt64(),
		IncludeAssociatedProperties: data.IncludeAssociatedProperties.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Read Client Certificate failed", err.Error())
		return
	}
	if versionsResponse == nil {
		resp.Diagnostics.AddError("Read Client Certificate failed", "Unexpected nil response for client certificate versions.")
		return
	}

	versionsData := extractVersions(versionsResponse.Versions)
	diags := data.parseClientCertificate(ctx, clientCertificate, versionsData)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

type versionsData struct {
	otherVersions []mtlskeystore.ClientCertificateVersion
	current       *mtlskeystore.ClientCertificateVersion
	previous      *mtlskeystore.ClientCertificateVersion
}

func extractVersions(versions []mtlskeystore.ClientCertificateVersion) versionsData {
	var otherVersions []mtlskeystore.ClientCertificateVersion
	var current, previous *mtlskeystore.ClientCertificateVersion
	for _, version := range versions {
		if version.VersionAlias == nil {
			otherVersions = append(otherVersions, version)
		} else {
			if *version.VersionAlias == "CURRENT" {
				current = &version
			} else {
				previous = &version
			}
		}
	}

	return versionsData{
		otherVersions: otherVersions,
		current:       current,
		previous:      previous,
	}
}

// parseClientCertificate maps the client certificate response to the data source model.
func (d *clientCertificateDataSourceModel) parseClientCertificate(ctx context.Context, clientCertificate *mtlskeystore.GetClientCertificateResponse, versionsData versionsData) diag.Diagnostics {
	var diags diag.Diagnostics

	d.CertificateName = types.StringValue(clientCertificate.CertificateName)
	d.CreatedBy = types.StringValue(clientCertificate.CreatedBy)
	d.CreatedDate = types.StringValue(clientCertificate.CreatedDate.Format(time.RFC3339))
	d.Geography = types.StringValue(clientCertificate.Geography)
	d.KeyAlgorithm = types.StringValue(clientCertificate.KeyAlgorithm)
	d.SecureNetwork = types.StringValue(clientCertificate.SecureNetwork)
	d.Signer = types.StringValue(clientCertificate.Signer)
	d.Subject = types.StringValue(clientCertificate.Subject)

	diags.Append(d.setNotificationEmails(ctx, clientCertificate.NotificationEmails)...)
	if diags.HasError() {
		return diags
	}

	for _, version := range versionsData.otherVersions {
		d.Versions = append(d.Versions, *convertDataToVersionModel(version))
	}

	if versionsData.current != nil {
		d.Current = convertDataToVersionModel(*versionsData.current)
	}
	if versionsData.previous != nil {
		d.Previous = convertDataToVersionModel(*versionsData.previous)
	}

	return diags
}

func convertDataToVersionModel(v mtlskeystore.ClientCertificateVersion) *versionModel {
	versionModel := versionModel{
		Version:                  types.Int64Value(v.Version),
		VersionGUID:              types.StringValue(v.VersionGUID),
		Status:                   types.StringValue(v.Status),
		CreatedBy:                types.StringValue(v.CreatedBy),
		CreatedDate:              types.StringValue(v.CreatedDate.Format(time.RFC3339)),
		ExpiryDate:               types.StringPointerValue(formatOptionalRFC3339(v.ExpiryDate)),
		Issuer:                   types.StringPointerValue(v.Issuer),
		KeyAlgorithm:             types.StringValue(v.KeyAlgorithm),
		EllipticCurve:            types.StringPointerValue(v.EllipticCurve),
		KeySizeInBytes:           types.StringPointerValue(v.KeySizeInBytes),
		SignatureAlgorithm:       types.StringPointerValue(v.SignatureAlgorithm),
		Subject:                  types.StringPointerValue(v.Subject),
		IssuedDate:               types.StringPointerValue(formatOptionalRFC3339(v.IssuedDate)),
		CertificateSubmittedBy:   types.StringPointerValue(v.CertificateSubmittedBy),
		CertificateSubmittedDate: types.StringPointerValue(formatOptionalRFC3339(v.CertificateSubmittedDate)),
		DeleteRequestedDate:      types.StringPointerValue(formatOptionalRFC3339(v.DeleteRequestedDate)),
		ScheduledDeleteDate:      types.StringPointerValue(formatOptionalRFC3339(v.ScheduledDeleteDate)),
		Properties:               parseProperties(v.AssociatedProperties),
		Validation:               parseValidation(v.Validation),
	}

	if v.CertificateBlock != nil {
		versionModel.CertificateBlock = &certificateBlockModel{
			Certificate:  types.StringValue(v.CertificateBlock.Certificate),
			KeyAlgorithm: types.StringValue(v.CertificateBlock.KeyAlgorithm),
			TrustChain:   types.StringValue(v.CertificateBlock.TrustChain),
		}
	}
	if v.CSRBlock != nil {
		versionModel.CSRBlock = &csrBlockModel{
			CSR:          types.StringValue(v.CSRBlock.CSR),
			KeyAlgorithm: types.StringValue(v.CSRBlock.KeyAlgorithm),
		}
	}

	return &versionModel
}

func (d *clientCertificateDataSourceModel) setNotificationEmails(ctx context.Context, emails []string) diag.Diagnostics {
	notificationEmails, diags := types.SetValueFrom(ctx, types.StringType, emails)
	if diags.HasError() {
		return diags
	}
	d.NotificationEmails = notificationEmails
	return nil
}

func parseProperties(properties []mtlskeystore.AssociatedProperty) []propertyModel {
	var props []propertyModel
	if properties == nil {
		return props
	}
	for _, prop := range properties {
		props = append(props, propertyModel{
			AssetID:         types.Int64Value(prop.AssetID),
			GroupID:         types.Int64Value(prop.GroupID),
			PropertyName:    types.StringValue(prop.PropertyName),
			PropertyVersion: types.Int64Value(prop.PropertyVersion),
		})
	}
	return props
}

func parseValidation(validation mtlskeystore.ValidationResult) validationModel {
	var errors []validationErrorModel
	for _, err := range validation.Errors {
		errors = append(errors, validationErrorModel{
			Message: err.Message,
			Reason:  err.Reason,
			Type:    err.Type,
		})
	}

	var warnings []validationErrorModel
	for _, warning := range validation.Warnings {
		warnings = append(warnings, validationErrorModel{
			Message: warning.Message,
			Reason:  warning.Reason,
			Type:    warning.Type,
		})
	}

	return validationModel{
		Errors:   errors,
		Warnings: warnings,
	}
}

func formatOptionalRFC3339(t *time.Time) *string {
	if t == nil {
		return nil
	}
	parsed := (*t).Format(time.RFC3339)
	return &parsed
}
