package vSummary

import (
    "errors"
    "fmt"
    "net/url"
    "regexp"
    "strconv"
    "strings"
)

// parseParts takes the split out parts of the field string, verifies that they are
// syntactically valid and then parses them out
//  for example columns[i][search][regex] would come in as
//       field:  'columns[2][search][regex]'
//       nameparts[0]  'columns'
//       nameparts[1]  '2]'
//       nameparts[2]  'search]'
//       nameparts[3]  'regex]'
func parseParts(field string, nameparts []string) (index int, elem1 string, elem2 string, err error) {
    defaultErr := fmt.Errorf("Invalid order[] element %v", field)
    numRegex, err := regexp.Compile("^[0-9]+]$")
    if err != nil {
        return
    }
    elemRegex, err := regexp.Compile("^[a-z]+]$")
    if err != nil {
        return
    }
    if len(nameparts) != 3 && len(nameparts) != 4 {
        err = defaultErr
        return
    }
    // Make sure it is a number followed by the closing ]
    if !numRegex.MatchString(nameparts[1]) {
        err = defaultErr
        return
    }
    // And parse it as a number to make sure
    numstr := strings.TrimSuffix(nameparts[1], "]")
    index, err = strconv.Atoi(numstr)
    if err != nil {
        return
    }
    // Check that the next index is a name token followed by a ]
    if !elemRegex.MatchString(nameparts[2]) {
        err = defaultErr
        return
    }
    // Strip off the trailing ]
    elem1 = strings.TrimSuffix(nameparts[2], "]")
    // If we had a third element, check to make sure it is also close by a ]
    if len(nameparts) == 4 {
        if !elemRegex.MatchString(nameparts[3]) {
            err = defaultErr
            return
        }
        // And trim off the ]
        elem2 = strings.TrimSuffix(nameparts[3], "]")
    }
    // Let's sanity check and make sure they aren't returning an index that is way out of range.
    // We shall assume that no more than 200 columns are being returned
    if index > 200 || index < 0 {
        err = defaultErr
    }
    return
}

// ParseDatatablesRequest checks the HTTP request to see if it corresponds
// to a datatables AJAX data request and parses the request data into
// the DataTablesInfo structure.
//
// This structure can be used by MySQLFilter and MySQLOrderby to generate a
// MySQL query to run against a database.
//
// For example assuming you are going to fill in a response structure to DataTables
// such as:
//
//   type QueryResponse struct {
//       DateAdded   time.Time
//       Status      string
//       Email       struct {
//           Name      string
//           Email     string
//       }
//   }
//   var emailQueueFields = map[string]string{
//       "DateAdded":          "t1.dateadded",
//       "Status":             "t1.status",
//       "Email.Name":         "t2.Name",
//       "Email.Email":        "t2.Email",
//   }
//
//   const baseQuery = `
//       SELECT t1.dateadded
//             ,t1.status
//             ,t2.Name
//             ,t2.Email
//       FROM infotable t1
//       LEFT JOIN usertable t2
//         ON t1.key = t2.key`
//
//       // See if we have a where clause to add to the base query
//       query := baseQuery
//       sqlPart, err := di.MySQLFilter(sqlFields)
//       // If we did have a where filter, append it.  Note that it doesn't put the " WHERE "
//       // in front because we might be doing a boolean operation.
//       if sqlPart != "" {
//           query += " WHERE " + sqlPart
//       }
//       sqlPart, err = di.MySQLOrderby(sqlFields)
//       query += sqlPart
//
// At that point you have a query that you can send straight to mySQL
//
func ParseDatatablesRequest(form url.Values) (res *DataTablesInfo, err error) {
    var index int
    var elem string
    var elem2 string
    foundDraw := false
    res = &DataTablesInfo{}
    // Let the request parse the post values into the r.Form structure

    for field, value := range form {
        // Remember that HTML sends us an array of values, but for datatables we only have one entry so we
        // we can shortcut and take the first element (which will be the only element) of the field.
        val0 := value[0]
        // Split out on the [ into pieces so we can see what the name is.  Note that we will have another
        // routine split out remainder of the string.
        nameparts := strings.Split(field, "[")
        switch nameparts[0] {
        case "draw":
            foundDraw = true
            res.Draw, err = strconv.Atoi(val0)
        case "start":
            res.Start, err = strconv.Atoi(val0)
        case "length":
            res.Length, err = strconv.Atoi(val0)
        case "search":
            if len(nameparts) != 2 {
                err = fmt.Errorf("Invalid search[] element %v", field)
            } else if nameparts[1] == "value]" {
                res.Searchval = val0
            } else if nameparts[1] == "regex]" {
                res.UseRegex = (val0 == "true")
            } else {
                err = fmt.Errorf("Invalid search[] element %v", field)
            }
        case "order":
            index, elem, _, err = parseParts(field, nameparts)
            if err == nil {
                // Make sure there is a spot to store this one.  Note that we may see
                // order[3][column] before we see order[0][dir]
                for len(res.Order) <= index {
                    res.Order = append(res.Order, OrderInfo{})
                }
                switch elem {
                case "column":
                    res.Order[index].ColNum, err = strconv.Atoi(val0)
                case "dir":
                    res.Order[index].Direction = Asc
                    if val0 == "desc" {
                        res.Order[index].Direction = Desc
                    }
                }
            }
        case "columns":
            index, elem, elem2, err = parseParts(field, nameparts)
            // First make sure we have a valid column number to work against
            if err == nil {
                // Fill up the slice to get to the spot where it is going
                // because the columns may come out of order.. I.e. we may see
                // columns[4][search][value] before we see columns[0][data]
                for len(res.Columns) <= index {
                    res.Columns = append(res.Columns, ColData{})
                }
            }
            // Now fill in the field in the column slice
            switch elem {
            case "data":
                res.Columns[index].Data = val0
            case "name":
                res.Columns[index].Name = val0
            case "searchable":
                res.Columns[index].Searchable = (val0 != "false")
            case "orderable":
                res.Columns[index].Orderable = (val0 != "false")
            case "search":
                switch elem2 {
                case "value":
                    res.Columns[index].Searchval = val0
                case "regex":
                    res.Columns[index].UseRegex = (val0 != "false")
                }
            }
        }
        // Any errors along the way and we get out.
        if err != nil {
            return
        }
    }
    // If no Draw was specified in the request, then this isn't a datatables request and we can safely ignore it
    if !foundDraw {
        res = nil
        err = errors.New("Not a DataTables request")
    } else {
        // We have a valid datatables request.  See if we actually have any filtering
        res.HasFilter = false
        // Check the global search value to see if it has anything on it
        if res.Searchval != "" {
            // We do have a filter so note that for later
            res.HasFilter = true
            // If they ask for a regex but don't use any regular expressions, then turn off regex for efficiency
            if res.UseRegex && !strings.ContainsAny(res.Searchval, "^$.*+|[]?") {
                res.UseRegex = false
            }
            // Escape the single quotes and any escape characters and then quote the string
            res.Searchval = strings.Replace(res.Searchval, "\\", "\\\\", -1)
            res.Searchval = "'" + strings.Replace(res.Searchval, "'", "\\'", -1) + "'"
        }
        // Now we check all of the columns to see if they have search expressions
        for _, colData := range res.Columns {
            if colData.Searchval != "" {
                // We have a search expression so we remember we have a filter
                res.HasFilter = true
                // CHeck for any regular expression characters and turn off regex if not
                if colData.UseRegex && !strings.ContainsAny(colData.Searchval, "[]^$.*?+") {
                    colData.UseRegex = false
                }
                // Escape the single quotes and any escape characters and then quote the string
                colData.Searchval = strings.Replace(colData.Searchval, "\\", "\\\\", -1)
                colData.Searchval = "'" + strings.Replace(colData.Searchval, "'", "\\'", -1) + "'"
            }
        }
    }
    return
}
