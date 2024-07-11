# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [priv/secrets/v1alpha1/secrets.proto](#priv_secrets_v1alpha1_secrets-proto)
    - [Action](#priv-secrets-v1alpha1-Action)
    - [BatchRequest](#priv-secrets-v1alpha1-BatchRequest)
    - [BatchResponse](#priv-secrets-v1alpha1-BatchResponse)
    - [DeleteSecretRequest](#priv-secrets-v1alpha1-DeleteSecretRequest)
    - [DeleteSecretResponse](#priv-secrets-v1alpha1-DeleteSecretResponse)
    - [ListSecretsRequest](#priv-secrets-v1alpha1-ListSecretsRequest)
    - [ListSecretsResponse](#priv-secrets-v1alpha1-ListSecretsResponse)
    - [PatchSecretRequest](#priv-secrets-v1alpha1-PatchSecretRequest)
    - [PatchSecretResponse](#priv-secrets-v1alpha1-PatchSecretResponse)
    - [Secret](#priv-secrets-v1alpha1-Secret)
    - [Secret.EnvironmentValuesEntry](#priv-secrets-v1alpha1-Secret-EnvironmentValuesEntry)
  
    - [SecretsService](#priv-secrets-v1alpha1-SecretsService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="priv_secrets_v1alpha1_secrets-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## priv/secrets/v1alpha1/secrets.proto



<a name="priv-secrets-v1alpha1-Action"></a>

### Action
Action is designed to represent a CRUD operation on a single Secret


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| patch_secret | [PatchSecretRequest](#priv-secrets-v1alpha1-PatchSecretRequest) |  | Reserving for future CRUD operations we may introduce CreateSecretRequest create_secret = 1; UpdateSecretRequest create_secret = 2; |
| delete_secret | [DeleteSecretRequest](#priv-secrets-v1alpha1-DeleteSecretRequest) |  |  |






<a name="priv-secrets-v1alpha1-BatchRequest"></a>

### BatchRequest
BatchRequest composes Actions to apply multiple CRUD methods in a single request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| actions | [Action](#priv-secrets-v1alpha1-Action) | repeated |  |






<a name="priv-secrets-v1alpha1-BatchResponse"></a>

### BatchResponse







<a name="priv-secrets-v1alpha1-DeleteSecretRequest"></a>

### DeleteSecretRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [string](#string) |  | project_id is the Project that this secret is namespaced within |
| secret_name | [string](#string) |  | secret_name is the name of the Secret to delete |
| environments | [string](#string) | repeated | environments must be one of: &#39;dev&#39;, &#39;preview&#39; or &#39;prod&#39; |






<a name="priv-secrets-v1alpha1-DeleteSecretResponse"></a>

### DeleteSecretResponse







<a name="priv-secrets-v1alpha1-ListSecretsRequest"></a>

### ListSecretsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [string](#string) |  | project_id is the Project that this secret is namespaced within E.g. &#34;projects/proj_sdfasdlfkj&#34; In the future, we could have Org-level secrets, which are shared across Projects, for which a new field may be introduced. |






<a name="priv-secrets-v1alpha1-ListSecretsResponse"></a>

### ListSecretsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| secrets | [Secret](#priv-secrets-v1alpha1-Secret) | repeated |  |






<a name="priv-secrets-v1alpha1-PatchSecretRequest"></a>

### PatchSecretRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [string](#string) |  | project_id is the Project that this secret is namespaced within |
| secret | [Secret](#priv-secrets-v1alpha1-Secret) |  | secret is the Secret to patch. If the Secret does not exist, it will be created. It can also be somewhat partial, in that it must have name and value but the environment_values map can be selectively filled. |






<a name="priv-secrets-v1alpha1-PatchSecretResponse"></a>

### PatchSecretResponse







<a name="priv-secrets-v1alpha1-Secret"></a>

### Secret
Secret is a resource that represents a Jetify Secret inside a Project.

NOTE: in this Secrets API, `org_id` is implicitly part of this API. We assume
that `org_id` will be retrieved from the Authorization JWT token.

id is a unique identifier of this Secret within its Project.
Reserving for now to stay compliant with this guideline, but not implementing.
https://cloud.google.com/apis/design/resource_names#resource_name_as_string
string id = 1;


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name for the secret |
| environment_values | [Secret.EnvironmentValuesEntry](#priv-secrets-v1alpha1-Secret-EnvironmentValuesEntry) | repeated | environment_values is a map of environment name to value. The environment name *must* be one of: &#39;dev&#39;, &#39;preview&#39; or &#39;prod&#39;, and the name must be all lowercase. In the future, this constraint may be relaxed to allow dynamically defined environments. |






<a name="priv-secrets-v1alpha1-Secret-EnvironmentValuesEntry"></a>

### Secret.EnvironmentValuesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [bytes](#bytes) |  |  |





 

 

 


<a name="priv-secrets-v1alpha1-SecretsService"></a>

### SecretsService
SecretsService provides CRUD methods for Secrets, although
we omit CreateSecret for now. It also has helper group methods for listing
secrets and batch operations.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Batch | [BatchRequest](#priv-secrets-v1alpha1-BatchRequest) | [BatchResponse](#priv-secrets-v1alpha1-BatchResponse) | Batch composes multiple CRUD requests into a single request. |
| DeleteSecret | [DeleteSecretRequest](#priv-secrets-v1alpha1-DeleteSecretRequest) | [DeleteSecretResponse](#priv-secrets-v1alpha1-DeleteSecretResponse) | DeleteSecret deletes an existing Secret. |
| ListSecrets | [ListSecretsRequest](#priv-secrets-v1alpha1-ListSecretsRequest) | [ListSecretsResponse](#priv-secrets-v1alpha1-ListSecretsResponse) | ListSecrets returns a list of Secrets for a given Project. |
| PatchSecret | [PatchSecretRequest](#priv-secrets-v1alpha1-PatchSecretRequest) | [PatchSecretResponse](#priv-secrets-v1alpha1-PatchSecretResponse) | PatchSecret partially updates an existing Secret, or creates it if it doesn&#39;t exist. |

 



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

