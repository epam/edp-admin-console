$(function () {

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

        let isStageValid = isStageAdded();
        if (!isStageValid) {
            $('.stage-error').show();
        } else {
            $('.stage-error').hide();
        }

        if (isPipelineValueValid && areCheckboxesChecked && isStageValid) {
            $('#createCDCR').submit();
        }
    });

    $('.add-stage-modal').click(function () {
        $('#stage-creation').modal('show');
        disableSelectElems();
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
    return /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/.test(pipelineName);
}

function isStageAdded() {
    return $('.stages-list .stage-info').length > 0;
}

function disableSelectElems() {
    $.each($('.autotests-checkbox'), function() {
        let $selectEl = $('select[name="' + $(this).attr('value') + '-autotestBranch"]');
        if (!$(this).is(':checked')) {
            $selectEl.attr('disabled', 'disabled');
        } else {
            $selectEl.removeAttr('disabled');
        }
    });
}