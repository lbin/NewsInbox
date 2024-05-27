# Error400

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Error_** | **string** | Depending on the case, error message may be any of these: &lt;ul&gt;     &lt;li&gt;&lt;code&gt;\&quot;Unrecognized or unsupported fields: [bad1, bad2, etc...]\&quot;&lt;/code&gt;&lt;/li&gt;     &lt;li&gt;&lt;code&gt;\&quot;Unacceptable query params: [badK1&#x3D;badV1, badK2&#x3D;badV2, etc...}]\&quot;&lt;/code&gt;&lt;/li&gt;     &lt;li&gt;&lt;code&gt;\&quot;Response would exceed maximum size....\&quot;&lt;/code&gt;&lt;/li&gt;         &lt;ul&gt;&lt;li&gt;This error will occur when the response exceeds 10 MB. Suggestions to either break the request into smaller batches, or make use of the limit and offset features will be presented.&lt;/li&gt;&lt;/ul&gt;     &lt;li&gt;A custom message string&lt;/li&gt;&lt;/ul&gt; | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

