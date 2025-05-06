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
});