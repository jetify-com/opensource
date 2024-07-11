# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [priv/members/v1alpha1/members.proto](#priv_members_v1alpha1_members-proto)
    - [CreateMemberRequest](#priv-members-v1alpha1-CreateMemberRequest)
    - [CreateMemberResponse](#priv-members-v1alpha1-CreateMemberResponse)
    - [DeleteMemberRequest](#priv-members-v1alpha1-DeleteMemberRequest)
    - [DeleteMemberResponse](#priv-members-v1alpha1-DeleteMemberResponse)
    - [GetMemberRequest](#priv-members-v1alpha1-GetMemberRequest)
    - [GetMemberResponse](#priv-members-v1alpha1-GetMemberResponse)
    - [Member](#priv-members-v1alpha1-Member)
  
    - [MembersService](#priv-members-v1alpha1-MembersService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="priv_members_v1alpha1_members-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## priv/members/v1alpha1/members.proto
API to manage members of an organization.


<a name="priv-members-v1alpha1-CreateMemberRequest"></a>

### CreateMemberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| member | [Member](#priv-members-v1alpha1-Member) |  |  |






<a name="priv-members-v1alpha1-CreateMemberResponse"></a>

### CreateMemberResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| member | [Member](#priv-members-v1alpha1-Member) |  |  |






<a name="priv-members-v1alpha1-DeleteMemberRequest"></a>

### DeleteMemberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="priv-members-v1alpha1-DeleteMemberResponse"></a>

### DeleteMemberResponse







<a name="priv-members-v1alpha1-GetMemberRequest"></a>

### GetMemberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The unique identifier of the member to retrieve. |






<a name="priv-members-v1alpha1-GetMemberResponse"></a>

### GetMemberResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| member | [Member](#priv-members-v1alpha1-Member) |  | The requested member object. |






<a name="priv-members-v1alpha1-Member"></a>

### Member
The member object


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | Unique identifier for the object |
| organization | [priv.organizations.v1alpha1.Organization](#priv-organizations-v1alpha1-Organization) |  |  |





 

 

 


<a name="priv-members-v1alpha1-MembersService"></a>

### MembersService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetMember | [GetMemberRequest](#priv-members-v1alpha1-GetMemberRequest) | [GetMemberResponse](#priv-members-v1alpha1-GetMemberResponse) | Get logged-in member

Retrieves the details of an existing member identified by its unique member id. |
| CreateMember | [CreateMemberRequest](#priv-members-v1alpha1-CreateMemberRequest) | [CreateMemberResponse](#priv-members-v1alpha1-CreateMemberResponse) |  |
| DeleteMember | [DeleteMemberRequest](#priv-members-v1alpha1-DeleteMemberRequest) | [DeleteMemberResponse](#priv-members-v1alpha1-DeleteMemberResponse) |  |

 



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

