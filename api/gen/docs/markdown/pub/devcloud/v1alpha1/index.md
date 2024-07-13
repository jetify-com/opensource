# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [pub/devcloud/v1alpha1/sandbox.proto](#pub_devcloud_v1alpha1_sandbox-proto)
    - [CreateSandboxRequest](#pub-devcloud-v1alpha1-CreateSandboxRequest)
    - [CreateSandboxRequest.EnvironmentVariablesEntry](#pub-devcloud-v1alpha1-CreateSandboxRequest-EnvironmentVariablesEntry)
    - [CreateSandboxResponse](#pub-devcloud-v1alpha1-CreateSandboxResponse)
    - [GetSandboxRequest](#pub-devcloud-v1alpha1-GetSandboxRequest)
    - [GetSandboxResponse](#pub-devcloud-v1alpha1-GetSandboxResponse)
    - [Sandbox](#pub-devcloud-v1alpha1-Sandbox)
  
    - [SandboxState](#pub-devcloud-v1alpha1-SandboxState)
  
    - [SandboxService](#pub-devcloud-v1alpha1-SandboxService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="pub_devcloud_v1alpha1_sandbox-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pub/devcloud/v1alpha1/sandbox.proto
API to manage Jetify devcloud Sandbox environments


<a name="pub-devcloud-v1alpha1-CreateSandboxRequest"></a>

### CreateSandboxRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| external_billing_tag | [string](#string) |  | Optional, user provided. Used for billing. |
| repo | [string](#string) |  |  |
| subdir | [string](#string) |  |  |
| ref | [string](#string) |  |  |
| environment_variables | [CreateSandboxRequest.EnvironmentVariablesEntry](#pub-devcloud-v1alpha1-CreateSandboxRequest-EnvironmentVariablesEntry) | repeated |  |






<a name="pub-devcloud-v1alpha1-CreateSandboxRequest-EnvironmentVariablesEntry"></a>

### CreateSandboxRequest.EnvironmentVariablesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="pub-devcloud-v1alpha1-CreateSandboxResponse"></a>

### CreateSandboxResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sandbox | [Sandbox](#pub-devcloud-v1alpha1-Sandbox) |  |  |






<a name="pub-devcloud-v1alpha1-GetSandboxRequest"></a>

### GetSandboxRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="pub-devcloud-v1alpha1-GetSandboxResponse"></a>

### GetSandboxResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sandbox | [Sandbox](#pub-devcloud-v1alpha1-Sandbox) |  |  |






<a name="pub-devcloud-v1alpha1-Sandbox"></a>

### Sandbox



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| external_billing_tag | [string](#string) |  |  |
| repo | [string](#string) |  |  |
| subdir | [string](#string) |  |  |
| ref | [string](#string) |  |  |
| url | [string](#string) |  | possibly empty while creating |
| state | [SandboxState](#pub-devcloud-v1alpha1-SandboxState) |  | enum |
| access_token | [string](#string) |  |  |





 


<a name="pub-devcloud-v1alpha1-SandboxState"></a>

### SandboxState
SandboxState represents the state of a sandbox and maps to workstationspb.Workstation_State
in the GCP Workstations API.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SANDBOX_STATE_UNSPECIFIED | 0 | Do not use. |
| SANDBOX_STATE_STARTING | 1 | The workstation is not yet ready to accept requests from users but will be soon. |
| SANDBOX_STATE_RUNNING | 2 | The workstation is ready to accept requests from users. |
| SANDBOX_STATE_STOPPING | 3 | The workstation is being stopped. |
| SANDBOX_STATE_STOPPED | 4 | The workstation is stopped and will not be able to receive requests until it is started. |


 

 


<a name="pub-devcloud-v1alpha1-SandboxService"></a>

### SandboxService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateSandbox | [CreateSandboxRequest](#pub-devcloud-v1alpha1-CreateSandboxRequest) | [CreateSandboxResponse](#pub-devcloud-v1alpha1-CreateSandboxResponse) |  |
| GetSandbox | [GetSandboxRequest](#pub-devcloud-v1alpha1-GetSandboxRequest) | [GetSandboxResponse](#pub-devcloud-v1alpha1-GetSandboxResponse) |  |

 



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

