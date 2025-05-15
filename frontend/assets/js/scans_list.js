$(document).ready(function () {
    console.log("scans_list.js was loaded");
    $('a.delete-scan').on('click', function (e) {
        e.preventDefault(); // Prevent the link from navigating
        console.log("click event was triggered");
        const href = $(this).attr('href');

        $.ajax({
            url: href,
            type: 'DELETE',
            success: function () {

                // delete a row from table
                $(e.target).closest('tr').remove();
            },
            error: function (xhr, status, error) {
                console.error('Logout failed:', error);
            }
        });
    });

    // register click event for the new scan button
    $('#new-scan').on('click', function (e) {
        // disable the button to prevent multiple clicks
        $(this).prop('disabled', true);

        $('span#hint').removeClass('d-none');
        // show user feedback
        $(this).html('Starting Scan...');
        $.ajax({
            url: '/scans',
            type: 'POST',
            success: function () {
                // Redirect to scans list
                window.location.href = '/scans';
            },
            error: function (xhr, status, error) {
                console.error('Could not start scan:', error);
            }
        });
    });

    setInterval(function () {
        let shouldRefresh = false;
    
        // Iterate over each .scan-status cell
        $('.scan-status').each(function () {
            const statusText = $(this).text().trim(); // Get the text content and trim whitespace
            if (statusText === 'Pending' || statusText === 'Running') {
                shouldRefresh = true; // Set flag to true if condition is met
                return false; // Break out of the loop early
            }
        });
    
        // Refresh the page if any status is 'Pending' or 'Running'
        if (shouldRefresh) {
            location.reload();
        }
    }, 10000);
});