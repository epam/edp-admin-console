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
                // $('#successPopup').modal('show');
                $.notify({
                        icon: 'glyphicon glyphicon-ok-circle alert-icon',
                        message: "Provisioning has been started."
                    },
                    {
                        type: 'success',
                        delay: 10000,
                        animate: {
                            enter:'animated fadeInRight',
                            exit: 'animated fadeOutRight'
                        }
                    });
            }
            location.hash = '';
        }

        let uri = window.location.pathname;
        let codebaseName = getUrlParameter('waitingforcodebase');
        if (codebaseName) {
            let status = $("tr[data-codebase-name='" + codebaseName + "']").attr("data-codebase-status");

            if (status == STATUS.CREATED)
            {
                $.notify({
                        icon: 'glyphicon glyphicon-ok-circle alert-icon',
                        message: "Provisioning has been finished."
                    },
                    {
                        type: 'success',
                        delay: 5000,
                        animate: {
                            enter:'animated fadeInRight',
                            exit: 'animated fadeOutRight'
                        }
                    });
            }
            else
            if (status == STATUS.FAILED) {
                $.notify({
                        icon: 'glyphicon gglyphicon-warning-sign alert-icon',
                        message: "Provisioning has been failed."
                    },
                    {
                        type: 'error',
                        delay: 5000,
                        animate: {
                            enter:'animated fadeInRight',
                            exit: 'animated fadeOutRight'
                        }
                    });

            }
            else {
                uri += "?waitingforcodebase=" + codebaseName;
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }
        }
        window.history.replaceState({}, document.title, uri);
    });
});