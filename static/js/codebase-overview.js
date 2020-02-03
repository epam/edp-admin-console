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
            } else if (anchor === '#codebaseIsUsed') {
                let codebase = getUrlParameter('codebase'),
                    pipeline = getUrlParameter('pipeline'),
                    $modal = $("#delete-confirmation");
                $modal.find('.invalid-feedback.server-error').show()
                    .text(`Codebase ${codebase} is used by CD Pipeline(s) ${pipeline}`);
                $("#delete-confirmation").data('codebase', codebase).modal('show');
            } else if (anchor === '#codebaseIsDeleted') {
                let codebase = getUrlParameter('codebase');
                showNotification(true, `Codebase ${codebase} was marked for deletion.`);
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

    $('.delete-codebase').click(function () {
        let codebase = $(this).data('codebase'),
            $modal = $("#delete-confirmation");
        $modal.data('codebase', codebase).modal('show');
    });

    $('.delete-confirmation').click(function () {
        let $modal = $("#delete-confirmation"),
            targetName = $modal.data('codebase'),
            confirmationName = $modal.find('#entity-name').val(),
            $errName = $modal.find('.invalid-feedback.different-name'),
            $errUsed = $modal.find('.invalid-feedback.server-error');
        if (targetName !== confirmationName) {
            $errName.show();
            return
        }
        $errName.hide();
        $errUsed.hide();
        $errUsed.text('');
        $("#delete-action").submit();
    });

    $('.close,.cancel-delete').click(function () {
        let $modal = $("#delete-confirmation");
        $modal.find('.invalid-feedback.different-name').hide();
        $modal.find('.invalid-feedback.server-error').text('').hide();
        $modal.find('#entity-name').val('');
    });

});

function showNotification(ok, msg, delay) {
    $.notify({
            icon: ok ? 'glyphicon glyphicon-ok-circle alert-icon' : 'glyphicon gglyphicon-warning-sign alert-icon',
            message: msg ? msg : (ok ? 'Provisioning has been started.' : 'Provisioning has been failed.')
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
