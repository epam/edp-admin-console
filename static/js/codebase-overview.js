$(function () {

    let STATUS = {
        IN_PROGRESS: 'inactive'
    };
    let delayTime = 10000;

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#codebaseSuccessModal') {
                showNotification(true);
            }
            location.hash = '';
        }

        let uri = window.location.pathname;
        let codebaseName = getUrlParameter('waitingforcodebase');
        if (codebaseName) {
            let status = $("tr[data-codebase-name='" + codebaseName + "']").attr("data-codebase-status");
            if (status === STATUS.IN_PROGRESS) {
                uri += "?waitingforcodebase=" + codebaseName;
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }
        }
        window.history.replaceState({}, document.title, uri);
    });
});


function showNotification(ok, delay) {
    $.notify({
            icon: ok ? 'glyphicon glyphicon-ok-circle alert-icon' : 'glyphicon gglyphicon-warning-sign alert-icon',
            message: ok ? 'Provisioning has been started.' : 'Provisioning has been failed.'
        },
        {
            type: ok ? 'success' : 'error',
            delay: delay ? delay: 5000,
            animate: {
                enter: 'animated fadeInRight',
                exit: 'animated fadeOutRight'
            }
        });
}