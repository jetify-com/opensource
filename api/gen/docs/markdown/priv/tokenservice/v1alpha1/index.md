# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [priv/tokenservice/v1alpha1/tokenservice.proto](#priv_tokenservice_v1alpha1_tokenservice-proto)
    - [CreateTokenRequest](#priv-tokenservice-v1alpha1-CreateTokenRequest)
    - [CreateTokenResponse](#priv-tokenservice-v1alpha1-CreateTokenResponse)
    - [GetAccessTokenRequest](#priv-tokenservice-v1alpha1-GetAccessTokenRequest)
    - [GetAccessTokenResponse](#priv-tokenservice-v1alpha1-GetAccessTokenResponse)
    - [Token](#priv-tokenservice-v1alpha1-Token)
  
    - [TokenService](#priv-tokenservice-v1alpha1-TokenService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="priv_tokenservice_v1alpha1_tokenservice-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## priv/tokenservice/v1alpha1/tokenservice.proto
API to manage token service


<a name="priv-tokenservice-v1alpha1-CreateTokenRequest"></a>

### CreateTokenRequest







<a name="priv-tokenservice-v1alpha1-CreateTokenResponse"></a>

### CreateTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [Token](#priv-tokenservice-v1alpha1-Token) |  |  |






<a name="priv-tokenservice-v1alpha1-GetAccessTokenRequest"></a>

### GetAccessTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  |  |






<a name="priv-tokenservice-v1alpha1-GetAccessTokenResponse"></a>

### GetAccessTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |






<a name="priv-tokenservice-v1alpha1-Token"></a>

### Token



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | ID is public and represents the token in the DB. |
| secret | [string](#string) |  | Secret is private. We do not store secret in db (only hash) |
| name | [string](#string) |  | if needed. Not adding yet, because not needed. org_id ? Could be &#34;subject&#34; instead since token may belong to org, user, project, etc scopes ? Not sure we want to use scopes. I think rbac is probably better. |





 

 

 


<a name="priv-tokenservice-v1alpha1-TokenService"></a>

### TokenService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetAccessToken | [GetAccessTokenRequest](#priv-tokenservice-v1alpha1-GetAccessTokenRequest) | [GetAccessTokenResponse](#priv-tokenservice-v1alpha1-GetAccessTokenResponse) |  |
| CreateToken | [CreateTokenRequest](#priv-tokenservice-v1alpha1-CreateTokenRequest) | [CreateTokenResponse](#priv-tokenservice-v1alpha1-CreateTokenResponse) |  |

 



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

