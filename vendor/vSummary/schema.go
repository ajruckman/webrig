package vSummary

type DataTablesInfo struct {
    // HasFilter Indicates there is a filter on the data to apply.  It is used to optimize generating
    // the query filters
    HasFilter bool
    // Draw counter. This is used by DataTables to ensure that the Ajax returns
    // from server-side processing requests are drawn in sequence by DataTables
    // (Ajax requests are asynchronous and thus can return out of sequence).
    // This is used as part of the draw return parameter (see below).
    Draw int
    // Start is the paging first record indicator.
    // This is the start point in the current data set (0 index based - i.e. 0 is the first record).
    Start int
    // Length is the number of records that the table can display in the current draw.
    // It is expected that the number of records returned will be equal to this number, unless the server has fewer records to return.
    //  Note that this can be -1 to indicate that all records should be returned (although that negates any benefits of server-side processing!)
    Length int
    // Searchval holds the global search value. To be applied to all columns which have searchable as true.
    Searchval string
    // UseRegex is true if the global filter should be treated as a regular expression for advanced searching.
    //  Note that normally server-side processing scripts will not perform regular expression
    //  searching for performance reasons on large data sets, but it is technically possible and at the discretion of your script.
    UseRegex bool
    // Order provides information about what columns are to be ordered in the results and which direction
    Order []OrderInfo
    // Columns provides a mapping of what fields are to be searched
    Columns []ColData
}

type ColData struct {
    // columns[i][name] Column's name, as defined by columns.name.
    Name string
    // columns[i][data] Column's data source, as defined by columns.data.
    // It is poss
    Data string
    // columns[i][searchable]	boolean	Flag to indicate if this column is searchable (true) or not (false).
    // This is controlled by columns.searchable.
    Searchable bool
    // columns[i][orderable] Flag to indicate if this column is orderable (true) or not (false).
    // This is controlled by columns.orderable.
    Orderable bool
    // columns[i][search][value] Search value to apply to this specific column.
    Searchval string
    // columns[i][search][regex]
    // Flag to indicate if the search term for this column should be treated as regular expression (true) or not (false).
    // As with global search, normally server-side processing scripts will not perform regular expression searching
    // for performance reasons on large data sets, but it is technically possible and at the discretion of your script.
    UseRegex bool
}

// OrderInfo tracks the list of columns to sort by and in which direction to sort them.
type OrderInfo struct {
    // ColNum indicates Which column to apply sorting to (zero based index to the Columns data)
    ColNum int
    // Direction tells us which way to sort
    Direction SortDir
}

type DataTablesResponse struct {
    // Draw counter. This is used by DataTables to ensure that the Ajax returns
    // from server-side processing requests are drawn in sequence by DataTables
    // (Ajax requests are asynchronous and thus can return out of sequence).
    // This is used as part of the draw return parameter (see below).
    Draw int `json:"draw"`
    // Total records, before filtering (i.e. the total number of records in the database)
    RecordsTotal int `json:"recordsTotal"`
    // Total records, after filtering
    // (i.e. the total number of records after filtering has been applied - not just the number of records being returned for this page of data).
    RecordsFiltered int `json:"recordsFiltered"`
    // The data to be displayed in the table. This is an array of data source objects, one for each row, which will be used by DataTables.
    // Note that this parameter's name can be changed using the ajax option's dataSrc property.
    Data []map[string]interface{} `json:"data"`
}

type SortDir int
