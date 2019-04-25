$(function () {

    $(document).ready(function () {
        let $successPopupEl = $('#successPopup');
        let displayPopup = $successPopupEl.data("display");
        if (displayPopup) {
            $successPopupEl.modal('show');
        }

        let appName = getUrlParameter('waitingforapp');
        if (appName) {
            let status = getApplicationStatus(appName);
            if (status === STATUS.CREATED || status === STATUS.FAILED) {
                let uri = window.location.href;
                window.history.replaceState({}, document.title, uri.substring(0, uri.indexOf("?")));
            } else {
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }
        }
    });

});

let STATUS = {
    CREATED: 'created',
    FAILED: 'failed'
};
let delayTime = 10000;

function getApplicationStatus(appName) {
    let status;
    $.each($('.edp-table tr td.app-name a'), function () {
        if ($(this).text().trim() === appName) {
            status = $(this).parents('tr').find('.app-status').data('status').trim();
            return false;
        }
    });
    return status;
}