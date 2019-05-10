$(function () {

    function isStageAdded() {
        let currentStageName = $('#stageName').val();
        let stageNames = $('input[name="stageName"]').map(function () {
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
        }

        if (validateFields()) {
            let stageData = collectStageData();
            let $stageCreationModal = $('#stage-creation');
            let stageToEdit = $('.stage-info.' + $stageCreationModal.attr('old-name'));
            replaceOldValues(stageToEdit, stageData);
            resetFields();
            toggleAdding();

            $stageCreationModal.modal('hide');
        }
    });

    $('.stage-modal-close').click(function () {
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

});

function validateFields() {
    return handleStageNameValidation() & handleStageDescriptionValidation() & handleStepNameValidation();
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
        '<input id="stageNameForm" name="stageName" type="hidden" value="' + stageData.stageName + '">' +
        '<input id="stageDescForm" name="' + stageData.stageName + '-stageDesc" type="hidden" value="' + stageData.stageDesc + '">' +
        '<input id="nameOfStepForm" name="' + stageData.stageName + '-nameOfStep" type="hidden" value="' + stageData.nameOfStep + '">' +
        '<input id="qualityGateTypeForm" name="' + stageData.stageName + '-qualityGateType" type="hidden" value="' + stageData.qualityGateType + '">' +
        '<input id="triggerTypeForm" name="' + stageData.stageName + '-triggerType" type="hidden" value="' + stageData.triggerType + '">' +
        '    </div>').appendTo($('.stages-list'));
}

function resetFields() {
    $('#qualityGateType option:first').prop('selected', true);
    $('#triggerType option:first').prop('selected', true);
    $("#stage-creation input").val("");

    $('input.non-valid-input').removeClass('non-valid-input');
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
    fillFields(stageName);
}

function fillFields(stageName) {
    let $stageEl = $('.stage-info.' + stageName);
    $('#stageName').val($stageEl.find('#stageNameForm').val());
    $('#stageDesc').val($stageEl.find('#stageDescForm').val());
    $('#nameOfStep').val($stageEl.find('#nameOfStepForm').val());
    $("#qualityGateType").val($stageEl.find('#qualityGateTypeForm').val());
    $("#triggerType").val($stageEl.find('#triggerTypeForm').val());
}

function toggleAdding() {
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
    return {
        stageName: $('#stageName').val(),
        stageDesc: $('#stageDesc').val(),
        nameOfStep: $('#nameOfStep').val(),
        qualityGateType: $('#qualityGateType').val(),
        triggerType: $('#triggerType').val(),
    };
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

function isStageNameValid() {
    let checkStageName = function (stageName) {
        return /^[a-z0-9][a-z0-9-._]*[a-z0-9]$/.test(stageName);
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
        return /^[a-z0-9][a-z0-9-._]*[a-z0-9]$/.test(stepName);
    };

    let $stepNameEl = $('#nameOfStep');
    return !(!$stepNameEl.val() || !checkStepName($stepNameEl.val()));
}