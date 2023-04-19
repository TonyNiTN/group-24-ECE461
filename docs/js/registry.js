// Populate RegistryTable
// This should use /packages endpoint
function renderTable() {
    $.ajax({
        'url': "/packages",
        'method': "POST",
        'contentType': 'application/json'
    }).done( function(data) {
        $('#registryTable').dataTable( {
            "aaData": data,
            "columns": [
                { "data": "ID" },
                { "data": "Name" },
                { "data": "Version" },
                { "data": "" },
            ]
        })
    })
    
    // button onClick=""
}

function deleteEntry() {

}

function editEntry() {

}

function openEntry() {

}

function searchRegistry() {
    
}