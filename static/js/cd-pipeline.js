$(function () {

    $(document).ready(function () {
        let anchor = $(location).attr('hash');
        if (anchor) {
            if (anchor === '#cdPipelineSuccessModal') {
                $('#successPopup').modal('show');
            }
            location.hash = '';
        }
    });

    $('.application-checkbox :checkbox').change(function () {
        let $selectEl = $('.select-' + $(this).attr('id'));
        if ($(this).is(':checked')) {
            $selectEl.prop('disabled', false);
            $('.app-checkbox-error').hide();
        } else {
            $selectEl.prop('disabled', true);
        }
    });

    $('#pipelineName').focusout(function () {
        let isPipelineValueValid = isPipelineValid();
        if (!isPipelineValueValid) {
            $('#pipelineName').addClass('non-valid-input');
            $('.invalid-feedback.pipeline-name').show();
        } else {
            $('#pipelineName').removeClass('non-valid-input');
            $('.invalid-feedback.pipeline-name').hide();
        }
    });

    $('.create-cd-pipeline').click(function (e) {
        e.preventDefault();
        let isPipelineValueValid = isPipelineValid();
        if (!isPipelineValueValid) {
            $('#pipelineName').addClass('non-valid-input');
            $('.invalid-feedback.pipeline-name').show();
        } else {
            $('#pipelineName').removeClass('non-valid-input');
            $('.invalid-feedback.pipeline-name').hide();
        }

        let areCheckboxesChecked = isAppCheckboxesValid();
        if (!areCheckboxesChecked) {
            $('.app-checkbox-error').show();
        } else {
            $('.app-checkbox-error').hide();
        }

        if (isPipelineValueValid && areCheckboxesChecked) {
            $('#createCDCR').submit();
        }
    });

    $('.add-stage-modal').click(function () {
        $('#stage-creation').modal('show');
    });

});

function isAppCheckboxesValid() {
    return $('.app-checkbox').is(':checked');
}

function isPipelineValid() {
    let $pipelineNameEl = $('#pipelineName');
    return !(!$pipelineNameEl.val() || !checkPipelineName($pipelineNameEl.val()));
}

function checkPipelineName(pipelineName) {
    return /^[a-z][a-z0-9-.]*[a-z0-9]$/.test(pipelineName);
}