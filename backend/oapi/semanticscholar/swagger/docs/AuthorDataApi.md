# {{classname}}

All URIs are relative to */graph/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetGraphGetAuthor**](AuthorDataApi.md#GetGraphGetAuthor) | **Get** /author/{author_id} | Details about an author
[**GetGraphGetAuthorPapers**](AuthorDataApi.md#GetGraphGetAuthorPapers) | **Get** /author/{author_id}/papers | Details about an author&#x27;s papers
[**GetGraphGetAuthorSearch**](AuthorDataApi.md#GetGraphGetAuthorSearch) | **Get** /author/search | Search for authors by name
[**PostGraphGetAuthors**](AuthorDataApi.md#PostGraphGetAuthors) | **Post** /author/batch | Get details for multiple authors at once

# **GetGraphGetAuthor**
> AuthorWithPapers GetGraphGetAuthor(ctx, authorId, optional)
Details about an author

Examples: <ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/1741101</code></li>     <ul>         <li>Returns the author's authorId and name.</li>     </ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/1741101?fields=aliases,papers</code></li>     <ul>         <li>Returns the author's authorId, aliases, and list of papers.  </li>         <li>Each paper has its paperId plus its title.</li>     </ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/1741101?fields=url,papers.abstract,papers.authors</code></li>     <ul>         <li>Returns the author's authorId, url, and list of papers.  </li>         <li>Each paper has its paperId, abstract, and list of authors.</li>         <li>In that list of authors, each author has their authorId and name.</li>     </ul>     <br>     Limitations:     <ul>         <li>Can only return up to 10 MB of data at a time.</li>     </ul> </ul>

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **authorId** | **string**|  | 
 **optional** | ***AuthorDataApiGetGraphGetAuthorOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AuthorDataApiGetGraphGetAuthorOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **fields** | **optional.String**| A comma-separated list of the fields to be returned.&lt;br&gt;&lt;br&gt;  The following case-sensitive author fields are recognized: &lt;ul&gt;     &lt;li&gt;&lt;code&gt;authorId&lt;/code&gt; - S2 unique ID for this author&lt;/li&gt; &lt;li&gt;&lt;code&gt;externalIds&lt;/code&gt; - ORCID/DBLP IDs for this author, if known&lt;/li&gt; &lt;li&gt;&lt;code&gt;url&lt;/code&gt; - URL on the Semantic Scholar website&lt;/li&gt; &lt;li&gt;&lt;code&gt;name&lt;/code&gt; - Author&#x27;s name&lt;/li&gt; &lt;li&gt;&lt;code&gt;aliases&lt;/code&gt; - List of names the author has used on publications over time, not intended to be displayed     to users. WARNING: this list may be out of date or contain deadnames of authors who have     changed their name. (see https://en.wikipedia.org/wiki/Deadnaming)&lt;/li&gt; &lt;li&gt;&lt;code&gt;affiliations&lt;/code&gt; - Author&#x27;s affiliations - sourced from claimed authors who have set affiliation on their S2 author page.&lt;/li&gt; &lt;li&gt;&lt;code&gt;homepage&lt;/code&gt; - Author&#x27;s homepage&lt;/li&gt; &lt;li&gt;&lt;code&gt;paperCount&lt;/code&gt; - Author&#x27;s total publications count&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationCount&lt;/code&gt; - Author&#x27;s total citations count&lt;/li&gt; &lt;li&gt;&lt;code&gt;hIndex&lt;/code&gt; - See the S2 &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#h-index\&quot;&gt;FAQ&lt;/a&gt; on h-index&lt;/li&gt;     &lt;li&gt;&lt;code&gt;papers&lt;/code&gt;         &lt;ul&gt;             &lt;li&gt;&lt;code&gt;paperId&lt;/code&gt; - Always included. A unique (string) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;corpusId&lt;/code&gt; - A second unique (numeric) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;url&lt;/code&gt; - URL on the Semantic Scholar website&lt;/li&gt; &lt;li&gt;&lt;code&gt;title&lt;/code&gt; - Included if no fields are specified&lt;/li&gt; &lt;li&gt;&lt;code&gt;venue&lt;/code&gt; - Normalized venue name&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationVenue&lt;/code&gt; - Publication venue meta-data for the paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;year&lt;/code&gt; - Year of publication&lt;/li&gt; &lt;li&gt;&lt;code&gt;authors&lt;/code&gt; - Up to 500 will be returned.  Will include: &lt;code&gt;authorId&lt;/code&gt; &amp; &lt;code&gt;name&lt;/code&gt;&lt;/li&gt; &lt;li&gt;To get more detailed information about an author&#x27;s papers, use the &lt;code&gt;/author/{author_id}/papers&lt;/code&gt; endpoint&lt;/li&gt; &lt;li&gt;&lt;code&gt;externalIds&lt;/code&gt;IDs from external sources - Supports ArXiv, MAG, ACL, PubMed, Medline, PubMedCentral, DBLP, DOI&lt;/li&gt; &lt;li&gt;&lt;code&gt;abstract&lt;/code&gt; - The paper&#x27;s abstract. Note that due to legal reasons, this may be missing even if we display an abstract on the website&lt;/li&gt; &lt;li&gt;&lt;code&gt;referenceCount&lt;/code&gt; - Total number of papers referenced by this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationCount&lt;/code&gt; - Total number of citations S2 has found for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;influentialCitationCount&lt;/code&gt; - More information &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#influential-citations\&quot;&gt;here&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;isOpenAccess&lt;/code&gt; - More information &lt;a href&#x3D;\&quot;https://www.openaccess.nl/en/what-is-open-access\&quot;&gt;here&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;openAccessPdf&lt;/code&gt; - A link to the paper if it is open access, and we have a direct link to the pdf&lt;/li&gt; &lt;li&gt;&lt;code&gt;fieldsOfStudy&lt;/code&gt; - A list of high-level academic categories from external sources&lt;/li&gt; &lt;li&gt;&lt;code&gt;s2FieldsOfStudy&lt;/code&gt; - A list of academic categories, sourced from either external sources or our internally developed &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#how-does-semantic-scholar-determine-a-papers-field-of-study\&quot;&gt;classifier&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationTypes&lt;/code&gt; - Journal Article, Conference, Review, etc&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationDate&lt;/code&gt; - YYYY-MM-DD, if available&lt;/li&gt; &lt;li&gt;&lt;code&gt;journal&lt;/code&gt; - Journal name, volume, and pages, if available&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationStyles&lt;/code&gt; - Generates bibliographical citation of paper. Currently supported styles: BibTeX&lt;/li&gt;         &lt;/ul&gt;     &lt;/li&gt; &lt;/ul&gt; | 

### Return type

[**AuthorWithPapers**](AuthorWithPapers.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetGraphGetAuthorPapers**
> PaperBatch GetGraphGetAuthorPapers(ctx, authorId, optional)
Details about an author's papers

Fetch the papers of an author in batches.<br> Only retrieves the most recent 10,000 citations/references for papers belonging to the batch.<br> To retrieve the full set of citations for a paper, use the /paper/{paper_id}/citations endpoint <br><br> Examples: <ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/1741101/papers</code></li>     <ul>         <li>Return with offset=0, and data is a list of the first 100 papers.</li>         <li>Each paper has its paperId and title.</li>     </ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/1741101/papers?fields=url,year,authors&limit=2</code></li>     <ul>         <li>Returns with offset=0, next=2, and data is a list of 2 papers.</li>         <li>Each paper has its paperId, url, year, and list of authors.</li>         <li>Each author has their authorId and name.</li>     </ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/1741101/papers?fields=citations.authors&offset=260</code></li>     <ul>         <li>Returns with offset=260, and data is a list of the last 4 papers.</li>         <li>Each paper has its paperId and a list of citations.</li>         <li>Each citation has its paperId and a list of authors.</li>         <li>Each author has their authorId and name.</li>     </ul> </ul>

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **authorId** | **string**|  | 
 **optional** | ***AuthorDataApiGetGraphGetAuthorPapersOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AuthorDataApiGetGraphGetAuthorPapersOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **offset** | **optional.Int32**| When returning a list of results, start with the element at this position in the list. | [default to 0]
 **limit** | **optional.Int32**| The maximum number of results to return.&lt;br&gt; Must be &lt;&#x3D; 1000 | [default to 100]
 **fields** | **optional.String**| A comma-separated list of the fields to be returned.&lt;br&gt;&lt;br&gt;  The following case-sensitive paper fields are recognized: &lt;ul&gt;     &lt;li&gt;&lt;code&gt;paperId&lt;/code&gt; - Always included. A unique (string) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;corpusId&lt;/code&gt; - A second unique (numeric) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;url&lt;/code&gt; - URL on the Semantic Scholar website&lt;/li&gt; &lt;li&gt;&lt;code&gt;title&lt;/code&gt; - Included if no fields are specified&lt;/li&gt; &lt;li&gt;&lt;code&gt;venue&lt;/code&gt; - Normalized venue name&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationVenue&lt;/code&gt; - Publication venue meta-data for the paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;year&lt;/code&gt; - Year of publication&lt;/li&gt; &lt;li&gt;&lt;code&gt;authors&lt;/code&gt; - Up to 500 will be returned.  Will include: &lt;code&gt;authorId&lt;/code&gt; &amp; &lt;code&gt;name&lt;/code&gt;&lt;/li&gt; &lt;li&gt;To get more detailed information about an author&#x27;s papers, use the &lt;code&gt;/author/{author_id}/papers&lt;/code&gt; endpoint&lt;/li&gt; &lt;li&gt;&lt;code&gt;externalIds&lt;/code&gt;IDs from external sources - Supports ArXiv, MAG, ACL, PubMed, Medline, PubMedCentral, DBLP, DOI&lt;/li&gt; &lt;li&gt;&lt;code&gt;abstract&lt;/code&gt; - The paper&#x27;s abstract. Note that due to legal reasons, this may be missing even if we display an abstract on the website&lt;/li&gt; &lt;li&gt;&lt;code&gt;referenceCount&lt;/code&gt; - Total number of papers referenced by this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationCount&lt;/code&gt; - Total number of citations S2 has found for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;influentialCitationCount&lt;/code&gt; - More information &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#influential-citations\&quot;&gt;here&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;isOpenAccess&lt;/code&gt; - More information &lt;a href&#x3D;\&quot;https://www.openaccess.nl/en/what-is-open-access\&quot;&gt;here&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;openAccessPdf&lt;/code&gt; - A link to the paper if it is open access, and we have a direct link to the pdf&lt;/li&gt; &lt;li&gt;&lt;code&gt;fieldsOfStudy&lt;/code&gt; - A list of high-level academic categories from external sources&lt;/li&gt; &lt;li&gt;&lt;code&gt;s2FieldsOfStudy&lt;/code&gt; - A list of academic categories, sourced from either external sources or our internally developed &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#how-does-semantic-scholar-determine-a-papers-field-of-study\&quot;&gt;classifier&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationTypes&lt;/code&gt; - Journal Article, Conference, Review, etc&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationDate&lt;/code&gt; - YYYY-MM-DD, if available&lt;/li&gt; &lt;li&gt;&lt;code&gt;journal&lt;/code&gt; - Journal name, volume, and pages, if available&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationStyles&lt;/code&gt; - Generates bibliographical citation of paper. Currently supported styles: BibTeX&lt;/li&gt;     &lt;li&gt;&lt;code&gt;citations&lt;/code&gt;&lt;/li&gt;     &lt;ul&gt;         &lt;li&gt;&lt;code&gt;paperId&lt;/code&gt; - Always included. A unique (string) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;corpusId&lt;/code&gt; - A second unique (numeric) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;url&lt;/code&gt; - URL on the Semantic Scholar website&lt;/li&gt; &lt;li&gt;&lt;code&gt;title&lt;/code&gt; - Included if no fields are specified&lt;/li&gt; &lt;li&gt;&lt;code&gt;venue&lt;/code&gt; - Normalized venue name&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationVenue&lt;/code&gt; - Publication venue meta-data for the paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;year&lt;/code&gt; - Year of publication&lt;/li&gt; &lt;li&gt;&lt;code&gt;authors&lt;/code&gt; - Up to 500 will be returned.  Will include: &lt;code&gt;authorId&lt;/code&gt; &amp; &lt;code&gt;name&lt;/code&gt;&lt;/li&gt; &lt;li&gt;To get more detailed information about an author&#x27;s papers, use the &lt;code&gt;/author/{author_id}/papers&lt;/code&gt; endpoint&lt;/li&gt;         &lt;li&gt;Total number of citations will be truncated at 10,000 for the entire batch.&lt;/li&gt;         &lt;li&gt;To fetch more citations per paper, reduce the number of papers in the batch with &lt;code&gt;limit&#x3D;&lt;/code&gt; or use the &lt;code&gt;/paper/{paper_id}/citations&lt;/code&gt; endpoint.&lt;/li&gt;     &lt;/ul&gt;     &lt;li&gt;&lt;code&gt;references&lt;/code&gt;&lt;/li&gt;     &lt;ul&gt;         &lt;li&gt;&lt;code&gt;paperId&lt;/code&gt; - Always included. A unique (string) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;corpusId&lt;/code&gt; - A second unique (numeric) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;url&lt;/code&gt; - URL on the Semantic Scholar website&lt;/li&gt; &lt;li&gt;&lt;code&gt;title&lt;/code&gt; - Included if no fields are specified&lt;/li&gt; &lt;li&gt;&lt;code&gt;venue&lt;/code&gt; - Normalized venue name&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationVenue&lt;/code&gt; - Publication venue meta-data for the paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;year&lt;/code&gt; - Year of publication&lt;/li&gt; &lt;li&gt;&lt;code&gt;authors&lt;/code&gt; - Up to 500 will be returned.  Will include: &lt;code&gt;authorId&lt;/code&gt; &amp; &lt;code&gt;name&lt;/code&gt;&lt;/li&gt; &lt;li&gt;To get more detailed information about an author&#x27;s papers, use the &lt;code&gt;/author/{author_id}/papers&lt;/code&gt; endpoint&lt;/li&gt;         &lt;li&gt;Same fields supported as for papers above&lt;/li&gt;         &lt;li&gt;Total number of references will be truncated at 10,000 for the entire batch.&lt;/li&gt;         &lt;li&gt;To fetch more references per paper, reduce the number of papers in the batch with &lt;code&gt;limit&#x3D;&lt;/code&gt; or use the &lt;code&gt;/paper/{paper_id}/references&lt;/code&gt; endpoint.&lt;/li&gt;     &lt;/ul&gt; &lt;/ul&gt;    | 

### Return type

[**PaperBatch**](PaperBatch.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetGraphGetAuthorSearch**
> AuthorSearchBatch GetGraphGetAuthorSearch(ctx, query, optional)
Search for authors by name

Examples: <ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/search?query=adam+smith</code></li>     <ul>         <li>Returns with total=490, offset=0, next=100, and data is a list of 100 authors.</li>         <li>Each author has their authorId and name. </li>     </ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/search?query=adam+smith&fields=name,aliases,url,papers.title,papers.year&limit=5</code></li>     <ul>         <li>Returns with total=490, offset=0, next=5, and data is a list of 5 authors.</li>         <li>Each author has authorId, name, aliases, url, and a list of their papers title and year.</li>     </ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/search?query=totalGarbageNonsense</code></li>     <ul>         <li>Returns with total = 0, offset=0, and data is a list of 0 author.</li>     </ul>     <br>     Limitations: <ul>     <li>Can only return up to 10 MB of data at a time.</li> </ul>

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **query** | **string**| A plain-text search query string. * No special query syntax is supported. * Hyphenated query terms yield no matches (replace it with space to find matches)  Specifying &lt;code&gt;papers&lt;/code&gt; fields in the request will return all papers linked to each author in the results, set a &lt;code&gt;limit&lt;/code&gt; on the search results to reduce output size and latency. | 
 **optional** | ***AuthorDataApiGetGraphGetAuthorSearchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AuthorDataApiGetGraphGetAuthorSearchOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **offset** | **optional.Int32**| When returning a list of results, start with the element at this position in the list. | [default to 0]
 **limit** | **optional.Int32**| The maximum number of results to return.&lt;br&gt; Must be &lt;&#x3D; 1000 | [default to 100]
 **fields** | **optional.String**| A comma-separated list of the fields to be returned.&lt;br&gt;&lt;br&gt;  The following case-sensitive author fields are recognized: &lt;ul&gt;     &lt;li&gt;&lt;code&gt;authorId&lt;/code&gt; - S2 unique ID for this author&lt;/li&gt; &lt;li&gt;&lt;code&gt;externalIds&lt;/code&gt; - ORCID/DBLP IDs for this author, if known&lt;/li&gt; &lt;li&gt;&lt;code&gt;url&lt;/code&gt; - URL on the Semantic Scholar website&lt;/li&gt; &lt;li&gt;&lt;code&gt;name&lt;/code&gt; - Author&#x27;s name&lt;/li&gt; &lt;li&gt;&lt;code&gt;aliases&lt;/code&gt; - List of names the author has used on publications over time, not intended to be displayed     to users. WARNING: this list may be out of date or contain deadnames of authors who have     changed their name. (see https://en.wikipedia.org/wiki/Deadnaming)&lt;/li&gt; &lt;li&gt;&lt;code&gt;affiliations&lt;/code&gt; - Author&#x27;s affiliations - sourced from claimed authors who have set affiliation on their S2 author page.&lt;/li&gt; &lt;li&gt;&lt;code&gt;homepage&lt;/code&gt; - Author&#x27;s homepage&lt;/li&gt; &lt;li&gt;&lt;code&gt;paperCount&lt;/code&gt; - Author&#x27;s total publications count&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationCount&lt;/code&gt; - Author&#x27;s total citations count&lt;/li&gt; &lt;li&gt;&lt;code&gt;hIndex&lt;/code&gt; - See the S2 &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#h-index\&quot;&gt;FAQ&lt;/a&gt; on h-index&lt;/li&gt;     &lt;li&gt;&lt;code&gt;papers&lt;/code&gt;         &lt;ul&gt;             &lt;li&gt;&lt;code&gt;paperId&lt;/code&gt; - Always included. A unique (string) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;corpusId&lt;/code&gt; - A second unique (numeric) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;url&lt;/code&gt; - URL on the Semantic Scholar website&lt;/li&gt; &lt;li&gt;&lt;code&gt;title&lt;/code&gt; - Included if no fields are specified&lt;/li&gt; &lt;li&gt;&lt;code&gt;venue&lt;/code&gt; - Normalized venue name&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationVenue&lt;/code&gt; - Publication venue meta-data for the paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;year&lt;/code&gt; - Year of publication&lt;/li&gt; &lt;li&gt;&lt;code&gt;authors&lt;/code&gt; - Up to 500 will be returned.  Will include: &lt;code&gt;authorId&lt;/code&gt; &amp; &lt;code&gt;name&lt;/code&gt;&lt;/li&gt; &lt;li&gt;To get more detailed information about an author&#x27;s papers, use the &lt;code&gt;/author/{author_id}/papers&lt;/code&gt; endpoint&lt;/li&gt; &lt;li&gt;&lt;code&gt;externalIds&lt;/code&gt;IDs from external sources - Supports ArXiv, MAG, ACL, PubMed, Medline, PubMedCentral, DBLP, DOI&lt;/li&gt; &lt;li&gt;&lt;code&gt;abstract&lt;/code&gt; - The paper&#x27;s abstract. Note that due to legal reasons, this may be missing even if we display an abstract on the website&lt;/li&gt; &lt;li&gt;&lt;code&gt;referenceCount&lt;/code&gt; - Total number of papers referenced by this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationCount&lt;/code&gt; - Total number of citations S2 has found for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;influentialCitationCount&lt;/code&gt; - More information &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#influential-citations\&quot;&gt;here&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;isOpenAccess&lt;/code&gt; - More information &lt;a href&#x3D;\&quot;https://www.openaccess.nl/en/what-is-open-access\&quot;&gt;here&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;openAccessPdf&lt;/code&gt; - A link to the paper if it is open access, and we have a direct link to the pdf&lt;/li&gt; &lt;li&gt;&lt;code&gt;fieldsOfStudy&lt;/code&gt; - A list of high-level academic categories from external sources&lt;/li&gt; &lt;li&gt;&lt;code&gt;s2FieldsOfStudy&lt;/code&gt; - A list of academic categories, sourced from either external sources or our internally developed &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#how-does-semantic-scholar-determine-a-papers-field-of-study\&quot;&gt;classifier&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationTypes&lt;/code&gt; - Journal Article, Conference, Review, etc&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationDate&lt;/code&gt; - YYYY-MM-DD, if available&lt;/li&gt; &lt;li&gt;&lt;code&gt;journal&lt;/code&gt; - Journal name, volume, and pages, if available&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationStyles&lt;/code&gt; - Generates bibliographical citation of paper. Currently supported styles: BibTeX&lt;/li&gt;         &lt;/ul&gt;     &lt;/li&gt; &lt;/ul&gt; | 

### Return type

[**AuthorSearchBatch**](AuthorSearchBatch.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostGraphGetAuthors**
> AuthorWithPapers PostGraphGetAuthors(ctx, body, optional)
Get details for multiple authors at once

* Fields is a single-value string parameter, not a multi-value one. * It is a query parameter, not to be submitted in the POST request's body.  In python:      r = requests.post(         'https://api.semanticscholar.org/graph/v1/author/batch',         params={'fields': 'name,hIndex,citationCount'},         json={\"ids\":[\"1741101\", \"1780531\"]}     )     print(json.dumps(r.json(), indent=2))      [       {         \"authorId\": \"1741101\",         \"name\": \"Oren Etzioni\",         \"citationCount\": 34803,         \"hIndex\": 86       },       {         \"authorId\": \"1780531\",         \"name\": \"Daniel S. Weld\",         \"citationCount\": 35526,         \"hIndex\": 89       }     ]  Other Examples: <ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/batch</code></li>     <ul>         <li><code>{\"ids\":[\"1741101\", \"1780531\", \"48323507\"]}</code></li>         <li>Returns details for 3 authors.</li>         <li>Each author returns the field authorId and name if no other fields are specified.</li>     </ul>     <li><code>https://api.semanticscholar.org/graph/v1/author/batch?fields=url,name,paperCount,papers,papers.title,papers.openAccessPdf</code></li>     <ul>         <li><code>{\"ids\":[\"1741101\", \"1780531\", \"48323507\"]}</code></li>         <li>Returns authorID, url, name, paperCount, and list of papers for 3 authors.</li>         <li>Each paper has its paperID, title, and link if available.</li>     </ul> </ul> <br> Limitations: <ul>     <li>Can only process 1,000 author ids at a time.</li>     <li>Can only return up to 10 MB of data at a time.</li> </ul>

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AuthorBatchRequest**](AuthorBatchRequest.md)|  | 
 **optional** | ***AuthorDataApiPostGraphGetAuthorsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AuthorDataApiPostGraphGetAuthorsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **fields** | **optional.**| A comma-separated list of the fields to be returned.&lt;br&gt;&lt;br&gt;  The following case-sensitive author fields are recognized: &lt;ul&gt;     &lt;li&gt;&lt;code&gt;authorId&lt;/code&gt; - S2 unique ID for this author&lt;/li&gt; &lt;li&gt;&lt;code&gt;externalIds&lt;/code&gt; - ORCID/DBLP IDs for this author, if known&lt;/li&gt; &lt;li&gt;&lt;code&gt;url&lt;/code&gt; - URL on the Semantic Scholar website&lt;/li&gt; &lt;li&gt;&lt;code&gt;name&lt;/code&gt; - Author&#x27;s name&lt;/li&gt; &lt;li&gt;&lt;code&gt;aliases&lt;/code&gt; - List of names the author has used on publications over time, not intended to be displayed     to users. WARNING: this list may be out of date or contain deadnames of authors who have     changed their name. (see https://en.wikipedia.org/wiki/Deadnaming)&lt;/li&gt; &lt;li&gt;&lt;code&gt;affiliations&lt;/code&gt; - Author&#x27;s affiliations - sourced from claimed authors who have set affiliation on their S2 author page.&lt;/li&gt; &lt;li&gt;&lt;code&gt;homepage&lt;/code&gt; - Author&#x27;s homepage&lt;/li&gt; &lt;li&gt;&lt;code&gt;paperCount&lt;/code&gt; - Author&#x27;s total publications count&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationCount&lt;/code&gt; - Author&#x27;s total citations count&lt;/li&gt; &lt;li&gt;&lt;code&gt;hIndex&lt;/code&gt; - See the S2 &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#h-index\&quot;&gt;FAQ&lt;/a&gt; on h-index&lt;/li&gt;     &lt;li&gt;&lt;code&gt;papers&lt;/code&gt;         &lt;ul&gt;             &lt;li&gt;&lt;code&gt;paperId&lt;/code&gt; - Always included. A unique (string) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;corpusId&lt;/code&gt; - A second unique (numeric) identifier for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;url&lt;/code&gt; - URL on the Semantic Scholar website&lt;/li&gt; &lt;li&gt;&lt;code&gt;title&lt;/code&gt; - Included if no fields are specified&lt;/li&gt; &lt;li&gt;&lt;code&gt;venue&lt;/code&gt; - Normalized venue name&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationVenue&lt;/code&gt; - Publication venue meta-data for the paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;year&lt;/code&gt; - Year of publication&lt;/li&gt; &lt;li&gt;&lt;code&gt;authors&lt;/code&gt; - Up to 500 will be returned.  Will include: &lt;code&gt;authorId&lt;/code&gt; &amp; &lt;code&gt;name&lt;/code&gt;&lt;/li&gt; &lt;li&gt;To get more detailed information about an author&#x27;s papers, use the &lt;code&gt;/author/{author_id}/papers&lt;/code&gt; endpoint&lt;/li&gt; &lt;li&gt;&lt;code&gt;externalIds&lt;/code&gt;IDs from external sources - Supports ArXiv, MAG, ACL, PubMed, Medline, PubMedCentral, DBLP, DOI&lt;/li&gt; &lt;li&gt;&lt;code&gt;abstract&lt;/code&gt; - The paper&#x27;s abstract. Note that due to legal reasons, this may be missing even if we display an abstract on the website&lt;/li&gt; &lt;li&gt;&lt;code&gt;referenceCount&lt;/code&gt; - Total number of papers referenced by this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationCount&lt;/code&gt; - Total number of citations S2 has found for this paper&lt;/li&gt; &lt;li&gt;&lt;code&gt;influentialCitationCount&lt;/code&gt; - More information &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#influential-citations\&quot;&gt;here&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;isOpenAccess&lt;/code&gt; - More information &lt;a href&#x3D;\&quot;https://www.openaccess.nl/en/what-is-open-access\&quot;&gt;here&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;openAccessPdf&lt;/code&gt; - A link to the paper if it is open access, and we have a direct link to the pdf&lt;/li&gt; &lt;li&gt;&lt;code&gt;fieldsOfStudy&lt;/code&gt; - A list of high-level academic categories from external sources&lt;/li&gt; &lt;li&gt;&lt;code&gt;s2FieldsOfStudy&lt;/code&gt; - A list of academic categories, sourced from either external sources or our internally developed &lt;a href&#x3D;\&quot;https://www.semanticscholar.org/faq#how-does-semantic-scholar-determine-a-papers-field-of-study\&quot;&gt;classifier&lt;/a&gt;&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationTypes&lt;/code&gt; - Journal Article, Conference, Review, etc&lt;/li&gt; &lt;li&gt;&lt;code&gt;publicationDate&lt;/code&gt; - YYYY-MM-DD, if available&lt;/li&gt; &lt;li&gt;&lt;code&gt;journal&lt;/code&gt; - Journal name, volume, and pages, if available&lt;/li&gt; &lt;li&gt;&lt;code&gt;citationStyles&lt;/code&gt; - Generates bibliographical citation of paper. Currently supported styles: BibTeX&lt;/li&gt;         &lt;/ul&gt;     &lt;/li&gt; &lt;/ul&gt; | 

### Return type

[**AuthorWithPapers**](AuthorWithPapers.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

