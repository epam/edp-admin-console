$(function () {

    let STATUS = {
        IN_PROGRESS: 'in_progress'
    };
    let delayTime = 10000;

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#autotestSuccessModal') {
                $('#successPopup').modal('show');
            }
            location.hash = '';
        }

        let uri = window.location.pathname;
        let autotestName = getUrlParameter('waitingforautotest');
        if (autotestName) {
            let status = $("tr[data-autotest-name='" + autotestName + "']").attr("data-autotest-status");
            if (status === STATUS.IN_PROGRESS) {
                uri += "?waitingforautotest=" + autotestName;
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }
        }
        window.history.replaceState({}, document.title, uri);
    });
});