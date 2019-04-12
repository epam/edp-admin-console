$(function () {

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#branchExistsModal') {
                let errorMessage = 'Release branch with ' + getUrlParameter('errorExistingBranch') + ' name is already exists.';
                $('.branch-exists-modal').text(errorMessage).show();
                $('#releaseBranchModal').modal('show');
            }
            if (anchor === '#branchSuccessModal') {
                $('#successPopup').modal('show');
            }
            window.history.pushState(null, null, window.location.pathname);
        }
    });

    $('#btn-modal-close, #btn-cross-close').click(function () {
        $('.branch-exists-modal').hide();
        $('#branchName,#commitNumber').val('').removeClass('non-valid-input');
    });

    $('.modal-release-branch').click(function () {
        $('#releaseBranchModal').modal('show');
    });

    $('#create-release-branch').click(function () {
        $('.branch-exists-modal').hide();
        let isBranchValid = handleBranchNameValidation();
        let isCommitValid = handleCommitHashvalidation();

        if (isBranchValid && isCommitValid) {
            $('#create-branch-action').submit();
        }
    });

    $('#branchName').focusout(function () {
        handleBranchNameValidation();
    });

    $('#commitNumber').focusout(function () {
        handleCommitHashvalidation();
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
    return /^[a-z][a-z0-9-.]*[a-z0-9]$/.test(branchName);
}

function checkHashCommit(hashCommit) {
    return /\b([a-f0-9]{40})\b/.test(hashCommit);
}

function handleBranchNameValidation() {
    let isBranchValid = isBranchNameValid();
    if (!isBranchValid) {
        $('#branchName').addClass('non-valid-input');
    } else {
        $('#branchName').removeClass('non-valid-input');
    }
    return isBranchValid;
}

function handleCommitHashvalidation() {
    let isCommitValid = isHashCommitValid();
    if (!isCommitValid) {
        $('#commitNumber').addClass('non-valid-input');
    } else {
        $('#commitNumber').removeClass('non-valid-input');
    }
    return isCommitValid;
}

function getUrlParameter(sParam) {
    let sPageURL = window.location.search.substring(1),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : decodeURIComponent(sParameterName[1]);
        }
    }
};