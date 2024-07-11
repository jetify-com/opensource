# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [priv/nix/v1alpha1/nix.proto](#priv_nix_v1alpha1_nix-proto)
    - [AWSCredentials](#priv-nix-v1alpha1-AWSCredentials)
    - [GetAWSCredentialsRequest](#priv-nix-v1alpha1-GetAWSCredentialsRequest)
    - [GetAWSCredentialsResponse](#priv-nix-v1alpha1-GetAWSCredentialsResponse)
    - [GetBinCacheRequest](#priv-nix-v1alpha1-GetBinCacheRequest)
    - [GetBinCacheResponse](#priv-nix-v1alpha1-GetBinCacheResponse)
    - [NixBinCache](#priv-nix-v1alpha1-NixBinCache)
  
    - [Permission](#priv-nix-v1alpha1-Permission)
  
    - [NixService](#priv-nix-v1alpha1-NixService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="priv_nix_v1alpha1_nix-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## priv/nix/v1alpha1/nix.proto
API to manage nix features


<a name="priv-nix-v1alpha1-AWSCredentials"></a>

### AWSCredentials



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_key_id | [string](#string) |  |  |
| secret_key | [string](#string) |  |  |
| session_token | [string](#string) |  |  |
| expiration | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="priv-nix-v1alpha1-GetAWSCredentialsRequest"></a>

### GetAWSCredentialsRequest







<a name="priv-nix-v1alpha1-GetAWSCredentialsResponse"></a>

### GetAWSCredentialsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| credentials | [AWSCredentials](#priv-nix-v1alpha1-AWSCredentials) |  |  |






<a name="priv-nix-v1alpha1-GetBinCacheRequest"></a>

### GetBinCacheRequest







<a name="priv-nix-v1alpha1-GetBinCacheResponse"></a>

### GetBinCacheResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nix_bin_cache_uri | [string](#string) |  | **Deprecated.** nix_bin_cache_uri is used by devbox 0.10.7 and below |
| caches | [NixBinCache](#priv-nix-v1alpha1-NixBinCache) | repeated |  |






<a name="priv-nix-v1alpha1-NixBinCache"></a>

### NixBinCache



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uri | [string](#string) |  |  |
| permissions | [Permission](#priv-nix-v1alpha1-Permission) | repeated |  |





 


<a name="priv-nix-v1alpha1-Permission"></a>

### Permission


| Name | Number | Description |
| ---- | ------ | ----------- |
| PERMISSION_UNSPECIFIED | 0 |  |
| PERMISSION_READ | 1 |  |
| PERMISSION_WRITE | 2 |  |


 

 


<a name="priv-nix-v1alpha1-NixService"></a>

### NixService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetBinCache | [GetBinCacheRequest](#priv-nix-v1alpha1-GetBinCacheRequest) | [GetBinCacheResponse](#priv-nix-v1alpha1-GetBinCacheResponse) |  |
| GetAWSCredentials | [GetAWSCredentialsRequest](#priv-nix-v1alpha1-GetAWSCredentialsRequest) | [GetAWSCredentialsResponse](#priv-nix-v1alpha1-GetAWSCredentialsResponse) |  |

 



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

