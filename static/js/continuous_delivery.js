$(function () {

    $(document).ready(function () {
        let uri = window.location.pathname;
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#cdPipelineSuccessModal') {
                showNotification(true);
            }
            location.hash = '';
        }

        let pipelineName = getUrlParameter('waitingforcdpipeline');
        if (pipelineName) {
            let status = getCDPipelineStatus(pipelineName);
            if (status === STATUS.IN_PROGRESS) {
                uri += "?waitingforcdpipeline="+pipelineName;
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }

            window.history.replaceState({}, document.title, uri);
        }
    });
});

let STATUS = {
    IN_PROGRESS: 'inactive'
};
let delayTime = 3000;

function getCDPipelineStatus(pipelineName) {
    let status;
    $.each($('.cd-pipeline-name'), function () {
        if ($(this).text().trim() === pipelineName) {
            status = $(this).parents('tr').find('.cd-pipeline-status').data('status').trim();
            return false;
        }
    });
    return status;
}

function showNotification(ok, delay) {
    $.notify({
            icon: ok ? 'glyphicon glyphicon-ok-circle alert-icon' : 'glyphicon gglyphicon-warning-sign alert-icon',
            message: ok ? 'Provisioning has been started.' : 'Provisioning has been failed.'
        },
        {
            type: ok ? 'success' : 'error',
            delay: delay ? delay : 5000,
            animate: {
                enter: 'animated fadeInRight',
                exit: 'animated fadeOutRight'
            }
        });
}