$(function () {

    let STATUS = {
        IN_PROGRESS: 'In Progress'
    };
    let delayTime = 10000;

    $(document).ready(function () {

        let uri = window.location.pathname;
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#branchExistsModal') {
                let errorMessage = 'Release branch with ' + getUrlParameter('errorExistingBranch') + ' name is already exists.';
                $('.branch-exists-modal').text(errorMessage).show();
                $('#releaseBranchModal').modal('show');
            }
            if (anchor === '#branchSuccessModal') {
               // $('#successPopup').modal('show');
                showNotification(true);
            }
        }

        let branchName = getUrlParameter('waitingforbranch');
        if (branchName) {
            let branchStatus = $("tr[data-branch-name='"+branchName+"']").attr("data-branch-status");
            if (status != STATUS.IN_PROGRESS) {
                uri += "?waitingforbranch="+branchName;
                setTimeout(function () {
                    location.reload();
                }, delayTime);
            }
            else {
                showNotification(true);
            }
        }
        window.history.replaceState({}, document.title, uri);
    });

    $('#btn-modal-close, #btn-cross-close').click(function () {
        $('.branch-exists-modal').hide();
        $('#branchName,#commitNumber').val('').removeClass('non-valid-input');
        $('.invalid-feedback.branch-name').hide();
        $('.invalid-feedback.commit-message').hide();
    });

    $('.modal-release-branch').click(function () {
        $('#releaseBranchModal').modal('show');
    });

    $('#create-release-branch').click(function () {
        $('.branch-exists-modal').hide();
        let isBranchValid = handleBranchNameValidation();
        let isCommitValid = handleCommitHashValidation();

        if (isBranchValid && isCommitValid) {
            $('#create-branch-action').submit();
        }
    });

    $('#branchName').focusout(function () {
        handleBranchNameValidation();
    });

    $('#commitNumber').focusout(function () {
        handleCommitHashValidation();
    });

});

function isBranchNameValid() {
    let $branchName = $('#branchName');
    return !(!$branchName.val() || !checkBranchName($branchName.val()));
}

function isHashCommitValid() {
    let $commitNumber = $('#commitNumber');

    if ($commitNumber.val().length === 0) {
        return true;
    } else {
        return !(!$commitNumber.val() || !checkHashCommit($commitNumber.val()));
    }
}

function checkBranchName(branchName) {
    return /^[a-z0-9][a-z0-9-._]*[a-z0-9]$/.test(branchName);
}

function checkHashCommit(hashCommit) {
    return /\b([a-f0-9]{40})\b/.test(hashCommit);
}

function handleBranchNameValidation() {
    let isBranchValid = isBranchNameValid();
    if (!isBranchValid) {
        $('#branchName').addClass('non-valid-input');
        $('.invalid-feedback.branch-name').show();
    } else {
        $('#branchName').removeClass('non-valid-input');
        $('.invalid-feedback.branch-name').hide();
    }
    return isBranchValid;
}

function handleCommitHashValidation() {
    let isCommitValid = isHashCommitValid();
    if (!isCommitValid) {
        $('#commitNumber').addClass('non-valid-input');
        $('.invalid-feedback.commit-message').show();
    } else {
        $('#commitNumber').removeClass('non-valid-input');
        $('.invalid-feedback.commit-message').hide();
    }
    return isCommitValid;
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