# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [priv/organizations/v1alpha1/organizations.proto](#priv_organizations_v1alpha1_organizations-proto)
    - [CreateOrganizationRequest](#priv-organizations-v1alpha1-CreateOrganizationRequest)
    - [CreateOrganizationResponse](#priv-organizations-v1alpha1-CreateOrganizationResponse)
    - [DeleteOrganizationRequest](#priv-organizations-v1alpha1-DeleteOrganizationRequest)
    - [DeleteOrganizationResponse](#priv-organizations-v1alpha1-DeleteOrganizationResponse)
    - [GetOrganizationRequest](#priv-organizations-v1alpha1-GetOrganizationRequest)
    - [GetOrganizationResponse](#priv-organizations-v1alpha1-GetOrganizationResponse)
    - [Organization](#priv-organizations-v1alpha1-Organization)
  
    - [OrganizationsService](#priv-organizations-v1alpha1-OrganizationsService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="priv_organizations_v1alpha1_organizations-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## priv/organizations/v1alpha1/organizations.proto
API to manage  an organization.


<a name="priv-organizations-v1alpha1-CreateOrganizationRequest"></a>

### CreateOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization | [Organization](#priv-organizations-v1alpha1-Organization) |  |  |






<a name="priv-organizations-v1alpha1-CreateOrganizationResponse"></a>

### CreateOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization | [Organization](#priv-organizations-v1alpha1-Organization) |  |  |






<a name="priv-organizations-v1alpha1-DeleteOrganizationRequest"></a>

### DeleteOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="priv-organizations-v1alpha1-DeleteOrganizationResponse"></a>

### DeleteOrganizationResponse







<a name="priv-organizations-v1alpha1-GetOrganizationRequest"></a>

### GetOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="priv-organizations-v1alpha1-GetOrganizationResponse"></a>

### GetOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization | [Organization](#priv-organizations-v1alpha1-Organization) |  |  |






<a name="priv-organizations-v1alpha1-Organization"></a>

### Organization
The organization object


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | Unique identifier for the object |
| name | [string](#string) |  |  |
| slug | [string](#string) |  |  |
| email | [string](#string) |  |  |
| stripe_id | [string](#string) |  |  |





 

 

 


<a name="priv-organizations-v1alpha1-OrganizationsService"></a>

### OrganizationsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetOrganization | [GetOrganizationRequest](#priv-organizations-v1alpha1-GetOrganizationRequest) | [GetOrganizationResponse](#priv-organizations-v1alpha1-GetOrganizationResponse) |  |
| CreateOrganization | [CreateOrganizationRequest](#priv-organizations-v1alpha1-CreateOrganizationRequest) | [CreateOrganizationResponse](#priv-organizations-v1alpha1-CreateOrganizationResponse) |  |
| DeleteOrganization | [DeleteOrganizationRequest](#priv-organizations-v1alpha1-DeleteOrganizationRequest) | [DeleteOrganizationResponse](#priv-organizations-v1alpha1-DeleteOrganizationResponse) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

