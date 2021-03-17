const maxJiraFields = 3;

function getTemplate() {
    return $('.jiraIssueMetadata .jira-issue-metadata-row').length === 0
        ? $('.full-template.jira-issue-metadata-row').clone()
        : $('.partial-template.jira-issue-metadata-row').clone();
}

function generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substring(2);
}

function enableButton() {
    $('.add-jira-field.circle.plus').removeClass('disable').removeAttr('disabled');
}

function toggleSelectOption() {
    toggleSelectOptions.bind(this)();
}

function toggleSelectOptions() {
    let oldData;
    if (typeof $(this).attr('data-old') !== 'undefined') {
        oldData = $.parseJSON($(this).attr('data-old'));
    }

    let selectedOption = $(this).find('option:selected').val();
    $.each($('section .jira-issue-metadata-row select'), function () {
        if ($(this).find('option:selected').val() === selectedOption) {
            return
        }
        this.selectize.removeOption(selectedOption);

        if (typeof oldData !== 'undefined') {
            this.selectize.addOption({value: oldData.key, text: oldData.value});
        }
    });

    $(this).attr('data-old', JSON.stringify({
        key: $(this).val(),
        value: $(this).text()
    }));
}

function tryToDisableButton() {
    if ($('.jiraIssueMetadata .jira-issue-metadata-row').length === maxJiraFields) {
        $('.add-jira-field.circle.plus').addClass('disable').attr('disabled');
        return
    }
}

function removeJiraIssueRow() {
    let $row = $(this).parents('.jira-issue-metadata-row'),
        $rows = $('section .jira-issue-metadata-row');
    if ($row.is(':first-child') && $rows.length >= 2) {
        let $fieldLabel = $('.jiraFieldNameLabel.template').clone(),
            $patternLabel = $('.jiraPatternLabel.template').clone();
        $fieldLabel.insertBefore($($rows[1]).find('select.jiraIssueFields'));
        $fieldLabel.removeClass('template').removeClass('hide-element');
        $patternLabel.insertBefore($($rows[1]).find('input.jiraPattern'));
        $patternLabel.removeClass('template').removeClass('hide-element');
    }

    toggleSelectOptions.bind($row.find('select.jiraIssueFields'))();
    $row.remove();
    enableButton();
}

function toggleJiraView() {
    let $jiraEl = $('.jiraServerBlock'),
        $commitMessagePatternBlockEl = $('.commitMessagePatternBlock'),
        $ticketNamePatternBlockEl = $('.ticketNamePatternBlock'),
        $jiraIssueMetadata = $('.jiraIssueMetadata'),
        $jiraIssueSelects = $('section .jiraIssueMetadata .jira-issue-metadata-row select.jiraIssueFields'),
        $jiraPatternInputs = $('section .jiraIssueMetadata .jira-issue-metadata-row input.jiraPattern');
    if ($(this).is(':checked')) {
        $jiraEl.removeClass('hide-element')
            .find('select[name="jiraServer"]')
            .prop('disabled', false);

        $commitMessagePatternBlockEl.removeClass('hide-element')
            .find('input[id="commitMessagePattern"]')
            .prop('disabled', false);

        $ticketNamePatternBlockEl.removeClass('hide-element')
            .find('input[id="ticketNamePattern"]')
            .prop('disabled', false);

        $jiraIssueSelects.prop('disabled', false);
        $jiraPatternInputs.prop('disabled', false);

        $jiraIssueMetadata.removeClass('hide-element');
        return;
    }
    $jiraEl.addClass('hide-element')
        .find('select[name="jiraServer"]')
        .prop('disabled', true);

    $commitMessagePatternBlockEl.addClass('hide-element')
        .find('input[id="commitMessagePattern"]')
        .prop('disabled', true);

    $ticketNamePatternBlockEl.addClass('hide-element')
        .find('input[id="ticketNamePattern"]')
        .prop('disabled', true);

    $jiraIssueSelects.prop('disabled', true);
    $jiraPatternInputs.prop('disabled', true);

    $jiraIssueMetadata.addClass('hide-element');
}