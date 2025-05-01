$(document).ready(function () {
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
                // Redirect to dashboard
                window.location.href = '/dashboard';
            },
            error: function (xhr, status, error) {
                console.error('Could not start scan:', error);
            }
        });
    });
});