$(function () {

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#cdPipelineSuccessModal') {
                $('#successPopup').modal('show');
            }
            location.hash = '';
        }

        let pipelineName = getUrlParameter('waitingforcdpipeline');
        if (pipelineName) {
            let status = getCDPipelineStatus(pipelineName);
            if (status === STATUS.CREATED || status === STATUS.FAILED) {
                window.history.replaceState({}, document.title, window.location.pathname);
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