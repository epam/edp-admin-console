$(function () {

    let STATUS = {
        CREATED: 'created',
        FAILED: 'failed',
        IN_PROGRESS: 'in_progress'
    };
    let delayTime = 10000;

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#codebaseSuccessModal') {
                $('#successPopup').modal('show');
            }
            location.hash = '';
        }

        let uri = window.location.pathname;
        let codebaseName = getUrlParameter('waitingforcodebase');
        if (codebaseName) {
            let status = $("tr[data-codebase-name='" + codebaseName + "']").attr("data-codebase-status");
            if (status !== STATUS.CREATED && status !== STATUS.FAILED) {
                uri += "?waitingforcodebase=" + codebaseName;
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }
        }
        window.history.replaceState({}, document.title, uri);
    });
});