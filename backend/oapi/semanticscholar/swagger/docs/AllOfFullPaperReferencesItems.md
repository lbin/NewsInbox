# AllOfFullPaperReferencesItems

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PaperId** | **string** | A unique (string) identifier for this paper.&lt; | [default to null]
**CorpusId** | **string** | A second unique (numeric) identifier for this paper. | [optional] [default to null]
**ExternalIds** | [***interface{}**](interface{}.md) | Other catalog IDs for this paper, if known. Supports ArXiv, MAG, ACL, PubMed, Medline, PubMedCentral, DBLP, DOI. | [optional] [default to null]
**Url** | **string** | URL on the Semantic Scholar website | [optional] [default to null]
**Title** | **string** |  | [optional] [default to null]
**Abstract** | **string** | The paper&#x27;s abstract. Note that due to legal reasons, this may be missing even if we display an abstract on the website. | [optional] [default to null]
**Venue** | **string** | normalized venue name | [optional] [default to null]
**PublicationVenue** | **string** | Details about the journal or conference in which this paper was published | [optional] [default to null]
**Year** | **int32** | year of publication | [optional] [default to null]
**ReferenceCount** | **int32** | Total number of papers referenced by this paper | [optional] [default to null]
**CitationCount** | **int32** | Total number of citations S2 has found for this paper | [optional] [default to null]
**InfluentialCitationCount** | **int32** | https://www.semanticscholar.org/faq#influential-citations | [optional] [default to null]
**IsOpenAccess** | **bool** | https://www.openaccess.nl/en/what-is-open-access | [optional] [default to null]
**OpenAccessPdf** | **string** | A link to the paper if it is open access, and we have a direct link to the pdf. As well as the paper&#x27;s status. More info on status here: https://en.wikipedia.org/wiki/Open_access#Colour_naming_system | [optional] [default to null]
**FieldsOfStudy** | [***interface{}**](interface{}.md) | A list of high-level academic categories from external sources. | [optional] [default to null]
**S2FieldsOfStudy** | [***interface{}**](interface{}.md) | This field returns a list of objects, where each has two keys: &#x27;category&#x27; and &#x27;source&#x27;. There are two sources: &#x27;external&#x27; - same as &#x27;fieldsOfStudy&#x27;, &#x27;s2-fos-model&#x27; - an internally developed classifier, see https://www.semanticscholar.org/faq#how-does-semantic-scholar-determine-a-papers-field-of-study | [optional] [default to null]
**PublicationTypes** | **[]string** | The type of this publication | [optional] [default to null]
**PublicationDate** | **string** | Year-month-day when this paper was published | [optional] [default to null]
**Journal** | [***interface{}**](interface{}.md) | Journal name, volume, and pages | [optional] [default to null]
**CitationStyles** | [***interface{}**](interface{}.md) | Bibliographic citations for paper, currently supported styles: BibTeX | [optional] [default to null]
**Authors** | [**[]AuthorInfo**](AuthorInfo.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

