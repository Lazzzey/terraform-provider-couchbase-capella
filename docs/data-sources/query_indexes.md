---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_query_indexes Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  
---

# couchbase-capella_query_indexes (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bucket_name` (String)
- `cluster_id` (String)
- `organization_id` (String)
- `project_id` (String)

### Optional

- `collection_name` (String)
- `scope_name` (String)

### Read-Only

- `data` (Attributes Set) (see [below for nested schema](#nestedatt--data))

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `definition` (String)
- `index_name` (String)
