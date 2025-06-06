# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [pub/sandbox/v1alpha1/sandbox.proto](#pub_sandbox_v1alpha1_sandbox-proto)
    - [CreateSandboxRequest](#pub-sandbox-v1alpha1-CreateSandboxRequest)
    - [CreateSandboxRequest.EnvironmentVariablesEntry](#pub-sandbox-v1alpha1-CreateSandboxRequest-EnvironmentVariablesEntry)
    - [CreateSandboxResponse](#pub-sandbox-v1alpha1-CreateSandboxResponse)
    - [DeleteSandboxRequest](#pub-sandbox-v1alpha1-DeleteSandboxRequest)
    - [DeleteSandboxResponse](#pub-sandbox-v1alpha1-DeleteSandboxResponse)
    - [GetSandboxRequest](#pub-sandbox-v1alpha1-GetSandboxRequest)
    - [GetSandboxResponse](#pub-sandbox-v1alpha1-GetSandboxResponse)
    - [ListSandboxesFilter](#pub-sandbox-v1alpha1-ListSandboxesFilter)
    - [ListSandboxesRequest](#pub-sandbox-v1alpha1-ListSandboxesRequest)
    - [ListSandboxesResponse](#pub-sandbox-v1alpha1-ListSandboxesResponse)
    - [Sandbox](#pub-sandbox-v1alpha1-Sandbox)
    - [StartSandboxRequest](#pub-sandbox-v1alpha1-StartSandboxRequest)
    - [StartSandboxResponse](#pub-sandbox-v1alpha1-StartSandboxResponse)
    - [StopSandboxRequest](#pub-sandbox-v1alpha1-StopSandboxRequest)
    - [StopSandboxResponse](#pub-sandbox-v1alpha1-StopSandboxResponse)
  
    - [MachineType](#pub-sandbox-v1alpha1-MachineType)
    - [Mode](#pub-sandbox-v1alpha1-Mode)
    - [SandboxState](#pub-sandbox-v1alpha1-SandboxState)
  
    - [SandboxService](#pub-sandbox-v1alpha1-SandboxService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="pub_sandbox_v1alpha1_sandbox-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pub/sandbox/v1alpha1/sandbox.proto
API to manage Jetify Sandbox environments


<a name="pub-sandbox-v1alpha1-CreateSandboxRequest"></a>

### CreateSandboxRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| external_billing_tag | [string](#string) |  | Optional, user provided. Used for billing. |
| repo | [string](#string) |  |  |
| subdir | [string](#string) |  |  |
| ref | [string](#string) |  |  |
| environment_variables | [CreateSandboxRequest.EnvironmentVariablesEntry](#pub-sandbox-v1alpha1-CreateSandboxRequest-EnvironmentVariablesEntry) | repeated |  |
| name | [string](#string) |  | Optional, human readable name for the sandbox. |
| project_id | [string](#string) |  | Optional. Project must belong to organization. |
| machine_type | [MachineType](#pub-sandbox-v1alpha1-MachineType) |  | Optional. The machine type to use for the sandbox. Defaults to MACHINE_TYPE_4CPU_16MEM. |
| github_token | [string](#string) |  |  |
| mode | [Mode](#pub-sandbox-v1alpha1-Mode) |  |  |






<a name="pub-sandbox-v1alpha1-CreateSandboxRequest-EnvironmentVariablesEntry"></a>

### CreateSandboxRequest.EnvironmentVariablesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="pub-sandbox-v1alpha1-CreateSandboxResponse"></a>

### CreateSandboxResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sandbox | [Sandbox](#pub-sandbox-v1alpha1-Sandbox) |  |  |
| long_running_operation_id | [string](#string) |  | This field is experimental and may be removed in the future. |






<a name="pub-sandbox-v1alpha1-DeleteSandboxRequest"></a>

### DeleteSandboxRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="pub-sandbox-v1alpha1-DeleteSandboxResponse"></a>

### DeleteSandboxResponse







<a name="pub-sandbox-v1alpha1-GetSandboxRequest"></a>

### GetSandboxRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| wait_for_long_running_operation_id | [string](#string) |  | This field is experimental and may be removed in the future. |






<a name="pub-sandbox-v1alpha1-GetSandboxResponse"></a>

### GetSandboxResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sandbox | [Sandbox](#pub-sandbox-v1alpha1-Sandbox) |  |  |






<a name="pub-sandbox-v1alpha1-ListSandboxesFilter"></a>

### ListSandboxesFilter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [string](#string) |  | Optional. If passed in, only sandboxes for this project are returned. Project must belong to organization |
| member_id | [string](#string) |  | Optional. If passed in, only sandboxes for this member are returned. |






<a name="pub-sandbox-v1alpha1-ListSandboxesRequest"></a>

### ListSandboxesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fetch_status_and_url | [bool](#bool) |  | fetch_status_and_url could be renamed to `fetch_workstation` to populate fields dependent on the cloud resource (i.e. Cloud Workstation) |
| filter | [ListSandboxesFilter](#pub-sandbox-v1alpha1-ListSandboxesFilter) |  |  |






<a name="pub-sandbox-v1alpha1-ListSandboxesResponse"></a>

### ListSandboxesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sandboxes | [Sandbox](#pub-sandbox-v1alpha1-Sandbox) | repeated |  |






<a name="pub-sandbox-v1alpha1-Sandbox"></a>

### Sandbox



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| external_billing_tag | [string](#string) |  |  |
| repo | [string](#string) |  | repo is the github repository that is cloned. This may be user owned, or the template repo. |
| subdir | [string](#string) |  | The subdirectory within the repo to checkout. Defaults to the root of the repo. |
| ref | [string](#string) |  | The git ref to checkout. This can be a branch, tag, or commit hash. Defaults to the default branch. |
| url | [string](#string) |  | Will be empty if the sandbox is not running. If present, it will contain access token. |
| state | [SandboxState](#pub-sandbox-v1alpha1-SandboxState) |  |  |
| access_token | [string](#string) |  | Token used to make requests to the sandbox. Use in the Authorization header as a Bearer token. |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | The time the sandbox was created. |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | The time the sandbox was last updated. |
| name | [string](#string) |  |  |
| project_id | [string](#string) |  | Optional. For sandboxes created with a project_id. |
| is_trial | [bool](#bool) |  | True if the sandbox is a trial sandbox. Trial sandboxes are deleted after 14 days and have other restrictions. |
| mode | [Mode](#pub-sandbox-v1alpha1-Mode) |  |  |
| has_gpu | [bool](#bool) |  |  |






<a name="pub-sandbox-v1alpha1-StartSandboxRequest"></a>

### StartSandboxRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="pub-sandbox-v1alpha1-StartSandboxResponse"></a>

### StartSandboxResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sandbox | [Sandbox](#pub-sandbox-v1alpha1-Sandbox) |  |  |
| long_running_operation_id | [string](#string) |  | This field is experimental and may be removed in the future. |






<a name="pub-sandbox-v1alpha1-StopSandboxRequest"></a>

### StopSandboxRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="pub-sandbox-v1alpha1-StopSandboxResponse"></a>

### StopSandboxResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sandbox | [Sandbox](#pub-sandbox-v1alpha1-Sandbox) |  |  |





 


<a name="pub-sandbox-v1alpha1-MachineType"></a>

### MachineType


| Name | Number | Description |
| ---- | ------ | ----------- |
| MACHINE_TYPE_UNSPECIFIED | 0 | Do not use. |
| MACHINE_TYPE_2CPU_8MEM | 1 |  |
| MACHINE_TYPE_4CPU_16MEM | 2 |  |
| MACHINE_TYPE_4CPU_15MEM_1T4 | 3 |  |



<a name="pub-sandbox-v1alpha1-Mode"></a>

### Mode


| Name | Number | Description |
| ---- | ------ | ----------- |
| MODE_UNSPECIFIED | 0 |  |
| MODE_GIT_REPO | 1 | MODE_GIT_REPO sets up a sandbox with the provided git repo (at specified `ref`), and the working directory set to the `subdir` in the repo. |
| MODE_TEMPLATE | 2 | MODE_TEMPLATE sets up a sandbox to use a template. The code from the git repo (at specified `ref`) is checked out and the root of the sandbox code folder is set to the `subdir` in the repo. The original git repo is discarded. A new git repo is initialized. |



<a name="pub-sandbox-v1alpha1-SandboxState"></a>

### SandboxState
SandboxState represents the state of a sandbox.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SANDBOX_STATE_UNSPECIFIED | 0 | Do not use. |
| SANDBOX_STATE_STARTING | 1 | The workstation is not yet ready to accept requests from users but will be soon. |
| SANDBOX_STATE_RUNNING | 2 | The workstation is ready to accept requests from users. |
| SANDBOX_STATE_STOPPING | 3 | The workstation is being stopped. |
| SANDBOX_STATE_STOPPED | 4 | The workstation is stopped and will not be able to receive requests until it is started. |
| SANDBOX_STATE_DELETED | 5 | Reserved for future use. Currently deleted sandboxes are not returned by api. |


 

 


<a name="pub-sandbox-v1alpha1-SandboxService"></a>

### SandboxService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateSandbox | [CreateSandboxRequest](#pub-sandbox-v1alpha1-CreateSandboxRequest) | [CreateSandboxResponse](#pub-sandbox-v1alpha1-CreateSandboxResponse) |  |
| GetSandbox | [GetSandboxRequest](#pub-sandbox-v1alpha1-GetSandboxRequest) | [GetSandboxResponse](#pub-sandbox-v1alpha1-GetSandboxResponse) |  |
| DeleteSandbox | [DeleteSandboxRequest](#pub-sandbox-v1alpha1-DeleteSandboxRequest) | [DeleteSandboxResponse](#pub-sandbox-v1alpha1-DeleteSandboxResponse) |  |
| ListSandboxes | [ListSandboxesRequest](#pub-sandbox-v1alpha1-ListSandboxesRequest) | [ListSandboxesResponse](#pub-sandbox-v1alpha1-ListSandboxesResponse) |  |
| StartSandbox | [StartSandboxRequest](#pub-sandbox-v1alpha1-StartSandboxRequest) | [StartSandboxResponse](#pub-sandbox-v1alpha1-StartSandboxResponse) | Starts sandbox. If sandbox is already starting or running, it does nothing. May return error if sandbox is stopping. |
| StopSandbox | [StopSandboxRequest](#pub-sandbox-v1alpha1-StopSandboxRequest) | [StopSandboxResponse](#pub-sandbox-v1alpha1-StopSandboxResponse) | Stops sandbox. If sandbox is already stopped, it does nothing. May return error if sandbox is starting. |

 



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

