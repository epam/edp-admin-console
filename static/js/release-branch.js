$(function () {

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#branchExistsModal') {
                $('.branch-exists-modal').show();
                $('#releaseBranchModal').modal('show');
            }
            if (anchor === '#branchSuccessModal') {
                $('#successPopup').modal('show');
            }
            location.hash = '';
        }
    });

    $('#btn-modal-close, #btn-cross-close').click(function () {
        $('.branch-exists-modal').hide();
        $('#branchName,#commitNumber').removeClass('non-valid-input');
    });

    $('.modal-release-branch').click(function () {
        $('#releaseBranchModal').modal('show');
    });

    $('#create-release-branch').click(function () {
        $('.branch-exists-modal').hide();
        let isBranchValid = isBranchNameValid();
        if (!isBranchValid) {
            $('#branchName').addClass('non-valid-input');
        } else {
            $('#branchName').removeClass('non-valid-input');
        }

        let isCommitValid = isHashCommitValid();
        if (!isCommitValid) {
            $('#commitNumber').addClass('non-valid-input');
        } else {
            $('#commitNumber').removeClass('non-valid-input');
        }

        if (isBranchValid && isCommitValid) {
            $('#create-branch-action').submit();
        }
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
    return /^[a-z\d](?:[a-z\d]|-(?=[a-z\d])){0,38}$/i.test(branchName);
}

function checkHashCommit(hashCommit) {
    return /\b([a-f0-9]{40})\b/.test(hashCommit);
}
