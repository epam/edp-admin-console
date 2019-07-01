$(function () {

    function isStageAdded() {
        let currentStageName = $('#stageName').val();
        let stageNames = $('input[name="stageName"][data-mode="add"]').map(function () {
            return $(this).val();
        }).get();
        return !!($.inArray(currentStageName, stageNames) !== -1);
    }

    $('.add-stage').click(function () {
        let $stageAddedMsgEl = $('.invalid-feedback.stage-added-msg');
        if (validateFields()) {
            if (isStageAdded()) {
                $stageAddedMsgEl.show();
                return;
            }
            $stageAddedMsgEl.hide();
            appendStage(collectStageData());
            resetFields();

            $('#stage-creation').modal('hide');
        }
    });

    $('.confirm-edit-stage').click(function () {
        function replaceOldValues(stageToEdit, stageData) {
            stageToEdit.removeClass($('#stage-creation').attr('old-name'));
            stageToEdit.find('.stage-name a').text(stageData.stageName);
            stageToEdit.find('.edit-stage').attr('name', stageData.stageName);
            stageToEdit.find('.remove-stage').attr('name', stageData.stageName);
            stageToEdit.addClass(stageData.stageName);
            stageToEdit.find('#stageNameForm').val(stageData.stageName);
            stageToEdit.find('#stageDescForm').val(stageData.stageDesc).attr('name', stageData.stageName + '-stageDesc');
            stageToEdit.find('#nameOfStepForm').val(stageData.nameOfStep).attr('name', stageData.stageName + '-nameOfStep');
            stageToEdit.find('#qualityGateTypeForm').val(stageData.qualityGateType).attr('name', stageData.stageName + '-qualityGateType');
            stageToEdit.find('#triggerTypeForm').val(stageData.triggerType).attr('name', stageData.stageName + '-triggerType');

            $('[name="' + stageData.stageName + '-autotests"]').remove();
            $.map(stageData.autotests, function (v, i) {
                $('[name="' + i + '-' + stageData.stageName + '-autotestBranch"]').remove();
            });

            if (stageData.autotests && stageData.qualityGateType === 'autotests') {
                $('.autotests-checkbox-info input:checked').each(function () {
                    let $stageBlockEl = $('.stage-info.' + stageData.stageName);
                    let branch = stageData.autotests[$(this).attr('value')];

                    $('<input data-target-select="' + branch + '" data-target="#' + $(this).attr('value') + '-checkbox' + '" id="' + $(this).attr('value') + '" name="' + stageData.stageName + '-autotests' + '" type="hidden" value="' + $(this).attr('value') + '">').appendTo($stageBlockEl);

                    let autotest = $(this).attr('value');
                    $('<input id="' + autotest + '-' + branch + '" name="' + autotest + '-' + stageData.stageName + '-autotestBranch" type="hidden" value="' + branch + '">').appendTo($stageBlockEl);
                });
            }
        }

        let $stageAddedMsgEl = $('.invalid-feedback.stage-added-msg');
        if (validateFields()) {
            if (isStageAdded()) {
                $stageAddedMsgEl.show();
                return;
            }
            $stageAddedMsgEl.hide();
            let stageData = collectStageData();
            let $stageCreationModal = $('#stage-creation');
            let stageToEdit = $('.stage-info.' + $stageCreationModal.attr('old-name'));
            replaceOldValues(stageToEdit, stageData);
            resetFields();
            toggleAdding();

            $stageCreationModal.modal('hide');
        }
    });

    $('.stage-modal-close, .cancel-edit-stage').click(function () {
        resetFields();
        toggleAdding();
    });

    $('#stageName').focusout(function () {
        let isValid = handleStageNameValidation();
        let $stageAddedMsgEl = $('.invalid-feedback.stage-added-msg');
        if (isValid) {
            if (isStageAdded()) {
                $stageAddedMsgEl.show();
                return;
            }
            $stageAddedMsgEl.hide();
        }
    });

    $('#stageDesc').focusout(function () {
        handleStageDescriptionValidation();
    });

    $('#nameOfStep').focusout(function () {
        handleStepNameValidation();
    });

    $('.tooltip-icon').tooltip();

    $('#qualityGateType').change(function () {
        let $autotestsEl = $('.autotests');
        if (this.value === 'autotests') {
            $autotestsEl.show();
        } else {
            $autotestsEl.hide();
            $('#qualityGateType').removeClass('non-valid-input');
            $('.autotests-validation-msg').hide();
        }
    });

    $('.autotests-checkbox').change(function () {
        if (this.checked) {
            $('#qualityGateType').removeClass('non-valid-input');
            $('.autotests-validation-msg').hide();
        }
    });
});

function validateFields() {
    return handleStageNameValidation() & handleStageDescriptionValidation() & handleStepNameValidation() & handleQualityGateTypeValidation();
}

function appendStage(stageData) {
    $('<div class="d-flex stage-info ' + stageData.stageName + '">\n' +
        '        <div class="form-group w-50 mb-2 stage-name">\n' +
        '<a class="edit-stage" onclick="editStage(this.name)" name="' + stageData.stageName + '" href="#">\n' + stageData.stageName + '</a>\n' +
        '        </div>\n' +
        '        <div class="d-flex flex-column justify-content-end mb-2">\n' +
        '            <button type="button" onclick="removeStage(this.name)" class="delete remove-stage" name="' + stageData.stageName + '" data-toggle="modal" data-target="#exampleModal">\n' +
        '                <i class="icon-trashcan"></i>\n' +
        '            </button>\n' +
        '        </div>\n' +
        '<input data-mode="add" id="stageNameForm" name="stageName" type="hidden" value="' + stageData.stageName + '">' +
        '<input id="stageDescForm" name="' + stageData.stageName + '-stageDesc" type="hidden" value="' + stageData.stageDesc + '">' +
        '<input id="nameOfStepForm" name="' + stageData.stageName + '-nameOfStep" type="hidden" value="' + stageData.nameOfStep + '">' +
        '<input id="qualityGateTypeForm" name="' + stageData.stageName + '-qualityGateType" type="hidden" value="' + stageData.qualityGateType + '">' +
        '<input id="triggerTypeForm" name="' + stageData.stageName + '-triggerType" type="hidden" value="' + stageData.triggerType + '">' +
        '    </div>').appendTo($('.stages-list'));

    if (stageData.autotests && stageData.qualityGateType === 'autotests') {
        $.map(stageData.autotests, function (v, i) {
            let $stageBlockEl = $('.stage-info.' + stageData.stageName);
            $('<input data-target-select="' + v + '" data-target="#' + i + '-checkbox' + '" id="' + i + '" name="' + stageData.stageName + '-autotests' + '" type="hidden" value="' + i + '">').appendTo($stageBlockEl);
            $('<input id="' + i + '-' + v + '" name="' + i + '-' + stageData.stageName + '-autotestBranch' + '" type="hidden" value="' + v + '">').appendTo($stageBlockEl);
        });
    }
}

function resetFields() {
    $('#qualityGateType option:first').prop('selected', true);
    $('#triggerType option:first').prop('selected', true);
    $('#stage-creation input[type="text"]').val("");
    $('.autotests').hide();
    $('.autotests-checkbox').prop('checked', false);

    $('input.non-valid-input, select.non-valid-input').removeClass('non-valid-input');
    $('div.invalid-feedback').hide();
}

function removeStage(stageName) {
    $('.stage-info.' + stageName).remove();
}

function editStage(stageName) {
    toggleEditing();
    let $stageCreationModal = $('#stage-creation');
    $stageCreationModal.attr('old-name', stageName);
    $stageCreationModal.modal('show');
    $('#stageNameForm[value="' + stageName + '"]').attr('data-mode', 'edit');
    fillFields(stageName);
}

function fillFields(stageName) {
    let $stageEl = $('.stage-info.' + stageName);
    let qualityGateTypeVal = $stageEl.find('#qualityGateTypeForm').val();
    $('#stageName').val($stageEl.find('#stageNameForm').val());
    $('#stageDesc').val($stageEl.find('#stageDescForm').val());
    $('#nameOfStep').val($stageEl.find('#nameOfStepForm').val());
    $("#qualityGateType").val(qualityGateTypeVal);
    $("#triggerType").val($stageEl.find('#triggerTypeForm').val());

    if (qualityGateTypeVal === 'autotests') {
        $('.autotests').show();

        $('input[name="' + stageName + '-autotests"]').each(function () {
            $($(this).attr('data-target')).prop('checked', true);

            let autotest = $(this).attr('value');
            let branch = $(this).attr('data-target-select');
            $('.' + autotest + '-branch').val(branch);
        });
    }
}

function toggleAdding() {
    $('.stage-info input[data-mode="edit"]').attr('data-mode', 'add');
    $('#add-header').show();
    $('#edit-header').hide();
    $('button.add-stage').show();
    $('button.confirm-edit-stage').hide();
}

function toggleEditing() {
    $('#add-header').hide();
    $('#edit-header').show();
    $('button.add-stage').hide();
    $('button.confirm-edit-stage').show();
}

function collectStageData() {
    let stageData = {
        stageName: $('#stageName').val(),
        stageDesc: $('#stageDesc').val(),
        nameOfStep: $('#nameOfStep').val(),
        qualityGateType: $('#qualityGateType').val(),
        autotests: undefined,
        triggerType: $('#triggerType').val(),
    };

    if (stageData.qualityGateType === 'autotests') {
        let autotests = {};
        $('.autotests-checkbox-info input:checked').each(function () {
            let autName = $(this).attr('value');
            autotests[autName] = $('.' + autName + '-branch').val();
        });
        stageData.autotests = autotests;
    }

    return stageData;
}

function handleStageNameValidation() {
    let $stageNameEl = $('#stageName');
    let valid = isStageNameValid();
    if (!valid) {
        $stageNameEl.addClass('non-valid-input');
        $stageNameEl.parents('div.form-group').find('.invalid-feedback.stage-name-msg').show();
    } else {
        $stageNameEl.removeClass('non-valid-input');
        $stageNameEl.parents('div.form-group').find('.invalid-feedback.stage-name-msg').hide();
    }
    return valid;
}

function handleStageDescriptionValidation() {
    let $stageDescEl = $('#stageDesc');
    let valid = isStageDescriptionValid();
    if (!valid) {
        $stageDescEl.addClass('non-valid-input');
        $stageDescEl.parents('div.form-group').find('.invalid-feedback').show();
    } else {
        $stageDescEl.removeClass('non-valid-input');
        $stageDescEl.parents('div.form-group').find('.invalid-feedback').hide();
    }
    return valid;
}

function handleStepNameValidation() {
    let $nameOfStep = $('#nameOfStep');
    let $validationMsgEl = $('.step-name-validation-msg');
    let valid = isStepNameValid();
    if (!valid) {
        $nameOfStep.addClass('non-valid-input');
        $validationMsgEl.show();
    } else {
        $nameOfStep.removeClass('non-valid-input');
        $validationMsgEl.hide();
    }
    return valid;
}

function handleQualityGateTypeValidation() {
    let $qualityGateTypeEl = $('#qualityGateType');
    let $autotestsValidationMsgEl = $('.autotests-validation-msg');
    let qualityTypeVal = $('#qualityGateType').children("option:selected").val();
    if (qualityTypeVal === 'autotests') {
        let isChecked = $('.autotests-checkbox').is(':checked');
        if (!isChecked) {
            $qualityGateTypeEl.addClass('non-valid-input');
            $autotestsValidationMsgEl.show();
        } else {
            $qualityGateTypeEl.removeClass('non-valid-input');
            $autotestsValidationMsgEl.hide();
        }
        return isChecked;
    }
    $qualityGateTypeEl.removeClass('non-valid-input');
    $autotestsValidationMsgEl.hide();
    return true;
}

function isStageNameValid() {
    let checkStageName = function (stageName) {
        return /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/.test(stageName);
    };

    let $stageNameEl = $('#stageName');
    return !(!$stageNameEl.val() || !checkStageName($stageNameEl.val()));
}

function isStageDescriptionValid() {
    let $stageDescriptionEl = $('#stageDesc');
    return $stageDescriptionEl.val().length !== 0;
}

function isStepNameValid() {
    let checkStepName = function (stepName) {
        return /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/.test(stepName);
    };

    let $stepNameEl = $('#nameOfStep');
    return !(!$stepNameEl.val() || !checkStepName($stepNameEl.val()));
}