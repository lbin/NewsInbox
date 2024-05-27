# Citation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Contexts** | [***interface{}**](interface{}.md) | Snippets of text where the reference is mentioned | [optional] [default to null]
**Intents** | [***interface{}**](interface{}.md) | https://www.semanticscholar.org/faq#citation-intent | [optional] [default to null]
**ContextsWithIntent** | [***interface{}**](interface{}.md) | Similar to \&quot;contexts\&quot; but each context comes its associated intents. Each object returned in this list has two keys: &#x27;context&#x27; - the text snippet, and &#x27;intents&#x27; - associated intents for this context. | [optional] [default to null]
**IsInfluential** | **bool** | https://www.semanticscholar.org/faq#influential-citations | [optional] [default to null]
**CitingPaper** | [***AllOfCitationCitingPaper**](AllOfCitationCitingPaper.md) | Details about the citing paper | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

