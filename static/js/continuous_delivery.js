$(function () {

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#cdPipelineSuccessModal') {
               // $('#successPopup').modal('show');
                showNotification(true);
            }
            location.hash = '';
        }

        let pipelineName = getUrlParameter('waitingforcdpipeline');
        if (pipelineName) {
            let status = getCDPipelineStatus(pipelineName);
            if (status == STATUS.CREATED || status == STATUS.FAILED)
            {
                showNotification(status == STATUS.CREATED);
            }
            else {
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }
            window.history.replaceState({}, document.title, window.location.pathname);
        }
    });
});

let STATUS = {
    CREATED: 'created',
    FAILED: 'failed'
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
            delay: delay ? delay: 5000,
            animate: {
                enter: 'animated fadeInRight',
                exit: 'animated fadeOutRight'
            }
        });
}