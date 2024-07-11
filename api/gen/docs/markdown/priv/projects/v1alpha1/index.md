# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [priv/projects/v1alpha1/projects.proto](#priv_projects_v1alpha1_projects-proto)
    - [CountProjectsWithDeploymentRequest](#priv-projects-v1alpha1-CountProjectsWithDeploymentRequest)
    - [CountProjectsWithDeploymentResponse](#priv-projects-v1alpha1-CountProjectsWithDeploymentResponse)
    - [CreateProjectRequest](#priv-projects-v1alpha1-CreateProjectRequest)
    - [CreateProjectResponse](#priv-projects-v1alpha1-CreateProjectResponse)
    - [DeleteProjectRequest](#priv-projects-v1alpha1-DeleteProjectRequest)
    - [DeleteProjectResponse](#priv-projects-v1alpha1-DeleteProjectResponse)
    - [GetProjectRequest](#priv-projects-v1alpha1-GetProjectRequest)
    - [GetProjectResponse](#priv-projects-v1alpha1-GetProjectResponse)
    - [ListProjectsRequest](#priv-projects-v1alpha1-ListProjectsRequest)
    - [ListProjectsResponse](#priv-projects-v1alpha1-ListProjectsResponse)
    - [PatchProjectRequest](#priv-projects-v1alpha1-PatchProjectRequest)
    - [PatchProjectResponse](#priv-projects-v1alpha1-PatchProjectResponse)
    - [Project](#priv-projects-v1alpha1-Project)
    - [SearchProjectsRequest](#priv-projects-v1alpha1-SearchProjectsRequest)
    - [SearchProjectsResponse](#priv-projects-v1alpha1-SearchProjectsResponse)
    - [UpdateProjectRequest](#priv-projects-v1alpha1-UpdateProjectRequest)
    - [UpdateProjectResponse](#priv-projects-v1alpha1-UpdateProjectResponse)
  
    - [ProjectsService](#priv-projects-v1alpha1-ProjectsService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="priv_projects_v1alpha1_projects-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## priv/projects/v1alpha1/projects.proto
API to manage projects


<a name="priv-projects-v1alpha1-CountProjectsWithDeploymentRequest"></a>

### CountProjectsWithDeploymentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| org_id | [string](#string) |  |  |






<a name="priv-projects-v1alpha1-CountProjectsWithDeploymentResponse"></a>

### CountProjectsWithDeploymentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| count | [int32](#int32) |  |  |






<a name="priv-projects-v1alpha1-CreateProjectRequest"></a>

### CreateProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| org_id | [string](#string) |  | The id of the organization under which to create the project. |
| project | [Project](#priv-projects-v1alpha1-Project) |  | The project object that you want to create. ID must be left empty since it will be assigned by there server. |






<a name="priv-projects-v1alpha1-CreateProjectResponse"></a>

### CreateProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#priv-projects-v1alpha1-Project) |  | The created project object. |






<a name="priv-projects-v1alpha1-DeleteProjectRequest"></a>

### DeleteProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [string](#string) |  | The unique id of the project you want to delete |






<a name="priv-projects-v1alpha1-DeleteProjectResponse"></a>

### DeleteProjectResponse







<a name="priv-projects-v1alpha1-GetProjectRequest"></a>

### GetProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The unique id used to identify the project. It&#39;s the same id returned by a project creation request or a project list request. |






<a name="priv-projects-v1alpha1-GetProjectResponse"></a>

### GetProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#priv-projects-v1alpha1-Project) |  | The requested project object. |






<a name="priv-projects-v1alpha1-ListProjectsRequest"></a>

### ListProjectsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| org_id | [string](#string) |  | The id of the organization under which you want to list projects. |
| page_size | [int32](#int32) |  | The maximum number of objects returned in a response. This number can range between 1 and 100, and it defaults to 10. |
| page_token | [string](#string) |  | A cursor for use in pagination. |






<a name="priv-projects-v1alpha1-ListProjectsResponse"></a>

### ListProjectsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| projects | [Project](#priv-projects-v1alpha1-Project) | repeated | A list containing the retrieved project objects. The list can be empty. |
| next_page_token | [string](#string) |  | A cursor to fetch the next page. |






<a name="priv-projects-v1alpha1-PatchProjectRequest"></a>

### PatchProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#priv-projects-v1alpha1-Project) |  | The project object you want to patch. |
| patch_mask | [google.protobuf.FieldMask](#google-protobuf-FieldMask) |  | A field mask containing the list of fields you would like to update. Only the fields listed in this field mask will be modified. |






<a name="priv-projects-v1alpha1-PatchProjectResponse"></a>

### PatchProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#priv-projects-v1alpha1-Project) |  | The updated project object. |






<a name="priv-projects-v1alpha1-Project"></a>

### Project
The project object

Projects describe a specific folder in a code repository that you develop on,
build and deploy as a unit. For example a backend that exposes your internal API,
and a nodejs server that serves your webapp, would each be a project.

In a multi-repo world, each project might live at the root of it&#39;s own git repository.
In a monorepo world, several projects might be in the same git repository, but have
different root directories.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | Unique identifier for the ojbect |
| repo | [string](#string) |  | The source control repository where the project&#39;s code lives Usually it&#39;s a git repository.

TODO: document format of string. Is this &lt;owner&gt;/&lt;repo&gt; (assumes GitHub), or a url? Or either? |
| directory | [string](#string) |  | The directory within the repository where the project&#39;s code lives. It should be a path within the repo, and if left empty, it defaults to the root of the repo.

RFC: Vercel calls this root_directory. Do we prefer that? I think it&#39;s because they also have output_directory ... but maybe we&#39;ll need to add another directory in the future too? |
| name | [string](#string) |  |  |






<a name="priv-projects-v1alpha1-SearchProjectsRequest"></a>

### SearchProjectsRequest
TBD






<a name="priv-projects-v1alpha1-SearchProjectsResponse"></a>

### SearchProjectsResponse
TBD






<a name="priv-projects-v1alpha1-UpdateProjectRequest"></a>

### UpdateProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#priv-projects-v1alpha1-Project) |  | The project object you want to update. |






<a name="priv-projects-v1alpha1-UpdateProjectResponse"></a>

### UpdateProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#priv-projects-v1alpha1-Project) |  | The updated project object. |





 

 

 


<a name="priv-projects-v1alpha1-ProjectsService"></a>

### ProjectsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetProject | [GetProjectRequest](#priv-projects-v1alpha1-GetProjectRequest) | [GetProjectResponse](#priv-projects-v1alpha1-GetProjectResponse) | Get a project

Retrieves the details of an existing project identified by its unique project id. |
| ListProjects | [ListProjectsRequest](#priv-projects-v1alpha1-ListProjectsRequest) | [ListProjectsResponse](#priv-projects-v1alpha1-ListProjectsResponse) | List the projects in an organization

Lists the projects belonging to the given organization. The projects are sorted by creation date, with the most recently created projects appearing first. |
| CountProjectsWithDeployment | [CountProjectsWithDeploymentRequest](#priv-projects-v1alpha1-CountProjectsWithDeploymentRequest) | [CountProjectsWithDeploymentResponse](#priv-projects-v1alpha1-CountProjectsWithDeploymentResponse) | Count the number of projects with Deployments

Given an org_id, counts the number of projects in an organization that have enabled deployments. |
| SearchProjects | [SearchProjectsRequest](#priv-projects-v1alpha1-SearchProjectsRequest) | [SearchProjectsResponse](#priv-projects-v1alpha1-SearchProjectsResponse) | Search for projects in an organization

Searches for products previously created in the given organization. Don&#39;t use search in read-after-write flows where strict consistency is necessary. |
| CreateProject | [CreateProjectRequest](#priv-projects-v1alpha1-CreateProjectRequest) | [CreateProjectResponse](#priv-projects-v1alpha1-CreateProjectResponse) | Create a new project

Creates a new project in the specified organization. The authenticated user must be a member of the organization. |
| DeleteProject | [DeleteProjectRequest](#priv-projects-v1alpha1-DeleteProjectRequest) | [DeleteProjectResponse](#priv-projects-v1alpha1-DeleteProjectResponse) | Delete a project

Deletes the project specified by the given id. |
| PatchProject | [PatchProjectRequest](#priv-projects-v1alpha1-PatchProjectRequest) | [PatchProjectResponse](#priv-projects-v1alpha1-PatchProjectResponse) | Patch a project

Patches the specified project with the provided fields. Any fields that are not provided, will be left unchanged. |
| UpdateProject | [UpdateProjectRequest](#priv-projects-v1alpha1-UpdateProjectRequest) | [UpdateProjectResponse](#priv-projects-v1alpha1-UpdateProjectResponse) | Update a project

Updates the specified project by setting the values of the provided fields. All fields will be updates. If you&#39;d like to partially update some fields, use Patch instead. |

 



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

