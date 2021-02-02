$(function () {

    let REGEX = {
        PIPELINE_NAME: /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/
    };

    !function () {
        _sendGetRequest(true, `${$('input[id="basepath"]').val()}/api/v1/edp/cd-pipeline`,
            function (pipes) {
                $('.pipeline-block').attr('data-pipes', JSON.stringify(pipes));
            }, function (resp) {
                console.error(resp)
            });
    }();

    function validatePipelineInfo(event) {
        let $pipelineBlockEl = $('.pipeline-block');

        resetErrors($pipelineBlockEl);

        let isValid = isPipelineInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($pipelineBlockEl);
            return isValid;
        }
        blockIsValid($pipelineBlockEl);

        return isValid;
    }

    function isPipelineInfoValid() {
        let $pipelineNameInputEl = $('#pipelineName'),
            isPipelineNameValid = isFieldValid($pipelineNameInputEl, REGEX.PIPELINE_NAME),
            isNameInUse = false;

        if (!isPipelineNameValid) {
            $('.invalid-feedback.pipeline-name-validation').show();
            $pipelineNameInputEl.addClass('is-invalid');
        }

        let $pipeExistsErrorBlock = $('.pipe-exists');
        if (_isPipeNameInUse()) {
            isNameInUse = true;
            $pipeExistsErrorBlock.show();
            $pipelineNameInputEl.addClass('pipe-exists-invalid');
        } else {
            isNameInUse = false;
            $pipeExistsErrorBlock.hide();
            $pipelineNameInputEl.removeClass('pipe-exists-invalid');
        }

        return isPipelineNameValid && !isNameInUse;
    }

    function validateApplicationInfo(event) {
        let $applicationBlockEl = $('.application-block');

        resetErrors($applicationBlockEl);

        let isValid = isApplicationInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($applicationBlockEl);
            return isValid;
        }
        blockIsValid($applicationBlockEl);

        return isValid;
    }

    function isApplicationInfoValid() {
        let isApplicationBlockValid = $('.app-checkbox[name="app"]').is(':checked');

        if (!isApplicationBlockValid) {
            $('.app-checkbox-error').show();
        }

        return isApplicationBlockValid;
    }

    function validateStageInfo(event) {
        let $stageBlockEl = $('.stage-block');

        resetErrors($stageBlockEl);

        let isValid = isStageInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($stageBlockEl);
            return isValid;
        }
        blockIsValid($stageBlockEl);

        return isValid;

    }

    function isStageInfoValid() {
        let isStageValid = $('.stages-list .stage-info').length > 0;

        if (!isStageValid) {
            $('.stage-error').show();
        }

        return isStageValid;
    }

    function resetErrors($el) {
        $el.find('input.is-invalid').removeClass('is-invalid');
        $el.find('.invalid-feedback:not(.pipe-exists)').hide();
    }

    $('.application-checkbox :checkbox').change(function () {
        let $selectEl = $('.select-' + $(this).attr('id')),
            $checkboxEl = $('.checkbox-' + $(this).attr('id'));
        if ($(this).is(':checked')) {
            $selectEl.prop('disabled', false);
            $checkboxEl.prop('disabled', false);
            $('.app-checkbox-error').hide();
            blockIsValid($('.application-block'));
        } else {
            $selectEl.prop('disabled', true);
            $checkboxEl.prop('disabled', true);
        }
    });

    $('#pipelineName').focusout(function () {
        let isPipelineValueValid = isPipelineValid();
        if (!isPipelineValueValid) {
            $('#pipelineName').addClass('non-valid-input');
            $('.invalid-feedback.pipeline-name-validation').show();
        } else {
            $('#pipelineName').removeClass('non-valid-input');
            $('.invalid-feedback.pipeline-name-validation').hide();
        }

        let $pipeExistsErrorBlock = $('.pipe-exists');
        if (_isPipeNameInUse()) {
            $pipeExistsErrorBlock.show();
            $(this).addClass('pipe-exists-invalid');
        } else {
            $pipeExistsErrorBlock.hide();
            $(this).removeClass('pipe-exists-invalid');
        }
    });

    function _isPipeNameInUse() {
        let pipes = JSON.parse($('.pipeline-block').attr('data-pipes')),
            pipe = $.grep(pipes, function (c) {
                return c.name === $('#pipelineName').val()
            });
        return pipe.length === 1;
    }


    function isPipelineValid() {
        let $pipelineNameEl = $('#pipelineName');
        return !(!$pipelineNameEl.val() || !checkPipelineName($pipelineNameEl.val()));
    }

    function checkPipelineName(pipelineName) {
        return /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/.test(pipelineName);
    }

    $('.pipeline-info-button').click(function (event) {
        validatePipelineInfo(event);
    });

    $('.application-info-button').click(function (event) {
        validateApplicationInfo(event);
    });

    $('.stage-info-button').click(function (event) {
        validateStageInfo(event);
    });

    $('.create-cd-pipeline').click(function (event) {
        event.preventDefault();

        let canCreateCDPipeline = validatePipelineInfo(event) &
            validateApplicationInfo(event) & validateStageInfo(event);

        if (canCreateCDPipeline) {
            $('#createCDCR').submit();
        }
    });

    $('.add-stage-modal').click(function () {
        $('#stage-creation').modal('show');
        disableSelectElems();
    });

});

function disableSelectElems() {
    $.each($('.autotests-checkbox'), function () {
        let $selectEl = $('select[name="' + $(this).attr('value') + '-autotestBranch"]');
        if (!$(this).is(':checked')) {
            $selectEl.attr('disabled', 'disabled');
        } else {
            $selectEl.removeAttr('disabled');
        }
    });
}